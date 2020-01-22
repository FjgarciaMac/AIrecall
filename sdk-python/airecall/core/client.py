"""Client to the AIrecall memory server (Go core).

A thin HTTP client. The server exposes a small JSON API:

  POST /v1/remember          {"agent_id", "content"}
  POST /v1/facts             {"agent_id", "key", "value"}
  POST /v1/recall            {"agent_id", "query", "top_k"}
  GET  /v1/facts/{agent}/{key}
  POST /v1/summarize         {"agent_id"}
"""

from __future__ import annotations
