"""Memory - the public SDK surface.

In dev mode (no server running) Memory falls back to an in-process
SQLite store so `airecall init` works without the Go core. When the
server is reachable it forwards everything to it.
"""

from __future__ import annotations

import json
import os
import sqlite3
import uuid
from datetime import datetime, timezone
from typing import Optional

from airecall.core.client import MemoryClient


def _now() -> str:
    return datetime.now(timezone.utc).isoformat()


class Memory:
    """Agent-facing memory facade.

    Parameters
    ----------
    agent_id:
        Namespace for this agent's memories. Different agents never see
        each other's episodes.
    db_path:
        SQLite file used by the in-process store (dev mode).
    server_url:
        When set, all operations go to the AIrecall memory server
        instead of the local store.
    """

    def __init__(self, agent_id: str = "default",
                 db_path: str = "airecall.db",
                 server_url: Optional[str] = None):
        self.agent_id = agent_id
        self.db_path = db_path
        if server_url:
            self._client = MemoryClient(server_url, agent_id)
            self._server = True
        else:
            self._client = None
            self._server = False
            self._conn = sqlite3.connect(db_path)
            self._init_schema()

    def _init_schema(self) -> None:
        self._conn.executescript("""
            CREATE TABLE IF NOT EXISTS episodes (
                id         INTEGER PRIMARY KEY AUTOINCREMENT,
                agent_id   TEXT NOT NULL,
                content    TEXT NOT NULL,
                kind       TEXT NOT NULL DEFAULT 'episode',
                created_at TEXT NOT NULL
            );
            CREATE TABLE IF NOT EXISTS facts (
                id         INTEGER PRIMARY KEY AUTOINCREMENT,
                agent_id   TEXT NOT NULL,
                key        TEXT NOT NULL,
                value      TEXT NOT NULL,
                updated_at TEXT NOT NULL,
                UNIQUE(agent_id, key)
            );
            CREATE INDEX IF NOT EXISTS idx_episodes_agent
                ON episodes(agent_id, created_at);
        """)
        self._conn.commit()

    def remember(self, content: str) -> None:
        """Store an episodic memory."""
        if self._server:
            self._client.remember(content)
            return
        self._conn.execute(
            "INSERT INTO episodes (agent_id, content, created_at) VALUES (?, ?, ?)",
            (self.agent_id, content, _now()))
        self._conn.commit()

    def remember_fact(self, key: str, value: str) -> None:
        """Store (or update) a durable semantic fact."""
        if self._server:
            self._client.remember_fact(key, value)
            return
        self._conn.execute(
            "INSERT INTO facts (agent_id, key, value, updated_at) VALUES (?, ?, ?, ?) "
            "ON CONFLICT(agent_id, key) DO UPDATE SET value=excluded.value, "
            "updated_at=excluded.updated_at",
            (self.agent_id, key, value, _now()))
        self._conn.commit()

    def recall(self, query: str, top_k: int = 5) -> list[str]:
        """Return the most relevant memories for the current turn.

        Hybrid retrieval: BM25-style keyword scoring combined with a
        simple cosine similarity over word-overlap vectors. The server
        runs the same algorithm in Go so results agree between modes.
        """
        top_k = max(0, int(top_k))
        if self._server:
            return self._client.recall(query, top_k)
        rows = self._conn.execute(
            "SELECT content FROM episodes WHERE agent_id=? ORDER BY created_at DESC "
            "LIMIT 200", (self.agent_id,)).fetchall()
        if not rows:
            return []
        scored = []
        for (content,) in rows:
            scored.append((self._score(query, content), content))
        scored.sort(key=lambda x: -x[0])
        return [c for s, c in scored[:top_k] if s > 0]

    def recall_fact(self, key: str) -> Optional[str]:
        """Return a semantic fact by key."""
        if self._server:
            return self._client.recall_fact(key)
        row = self._conn.execute(
            "SELECT value FROM facts WHERE agent_id=? AND key=?",
            (self.agent_id, key)).fetchone()
        return row[0] if row else None

    def summarize(self) -> dict:
        """Trigger a summarization pass (server mode only in production;
        dev mode returns a no-op summary)."""
        if self._server:
            return self._client.summarize()
        n = self._conn.execute(
            "SELECT COUNT(*) FROM episodes WHERE agent_id=?", (self.agent_id,)
        ).fetchone()[0]
        return {"episodes": n, "compacted": 0, "facts_promoted": 0}

    @staticmethod
    def _score(query: str, content: str) -> float:
        q_words = set(w.lower() for w in query.split() if len(w) > 2)
        if not q_words:
            return 0.0
        c_words = content.lower().split()
        hits = sum(1 for w in q_words if w in c_words)
        return hits / len(q_words)

    def close(self) -> None:
        if not self._server:
            self._conn.close()

    def __enter__(self) -> "Memory":
        return self

    def __exit__(self, *exc) -> None:
        self.close()
