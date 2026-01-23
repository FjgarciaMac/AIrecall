"""Client to the AIrecall memory server (Go core).

A thin HTTP client. The server exposes a small JSON API:

  POST /v1/remember          {"agent_id", "content"}
  POST /v1/facts             {"agent_id", "key", "value"}
  POST /v1/recall            {"agent_id", "query", "top_k"}
  GET  /v1/facts/{agent}/{key}
  POST /v1/summarize         {"agent_id"}
"""

from __future__ import annotations

import json
from typing import Optional

import requests


class MemoryClient:
    def __init__(self, base_url: str, agent_id: str, timeout: float = 10.0):
        self.base_url = base_url.rstrip("/")
        self.agent_id = agent_id
        self.timeout = timeout

    def _post(self, path: str, payload: dict) -> dict:
        r = requests.post(f"{self.base_url}{path}", json=payload,
                          timeout=self.timeout)
        r.raise_for_status()
        return r.json()

    def remember(self, content: str) -> None:
        self._post("/v1/remember", {"agent_id": self.agent_id, "content": content})

    def remember_fact(self, key: str, value: str) -> None:
        self._post("/v1/facts", {"agent_id": self.agent_id, "key": key,
                                 "value": value})

    def recall(self, query: str, top_k: int = 5) -> list[str]:
        data = self._post("/v1/recall", {"agent_id": self.agent_id,
                                         "query": query, "top_k": top_k})
        return data.get("memories", [])

    def recall_fact(self, key: str) -> Optional[str]:
        try:
            r = requests.get(
                f"{self.base_url}/v1/facts/{self.agent_id}/{key}",
                timeout=self.timeout)
            if r.status_code == 200:
                return r.json().get("value")
        except requests.RequestException:
            pass
        return None

    def summarize(self) -> dict:
        return self._post("/v1/summarize", {"agent_id": self.agent_id})
// draft note 1396
