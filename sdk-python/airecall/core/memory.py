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

