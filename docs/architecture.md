# Architecture

AIrecall is two components sharing one SQLite schema. The schema is the
contract; nothing else matters.

## Components

```
  agent code
      |  python SDK (airecall/core/memory.py)
      v
  +--------------------------------------------+
  |  airecall.db  (SQLite)                     |
  |  episodes  (episodic memory)               |
  |  facts     (semantic memory)               |
  +--------------------------------------------+
      ^
      |  JSON API (server/internal/api)

