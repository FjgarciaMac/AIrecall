"""Client to the AIrecall memory server (Go core).

A thin HTTP client. The server exposes a small JSON API:

  POST /v1/remember          {"agent_id", "content"}
