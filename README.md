# AIrecall

**A drop-in long-term memory layer for AI agents.**
Episodic + semantic memory, hybrid retrieval, and auto-summarization — so
your agent remembers what matters across sessions.

```
pip install airecall-sdk
airecall init
```

| | |
|---|---|
| Languages | Python SDK + Go memory server |
| Storage | SQLite + in-process vector index |
| License | MIT |
| Dependencies | SDK stdlib-only · server one pure-Go dep |

---

## Why This Exists

Most agent frameworks give you a context window, not a memory. Close the
session and everything the agent learned — user preferences, past
decisions, corrections it was given — is gone. Stuff it all into the
prompt instead, and you're paying for and diluting your context with old
information, most of which is irrelevant to the current turn.

AIrecall sits between your agent and a persistent store, and gives it an
actual memory:

| | | |
|---|---|---|
| 🧠 **Episodic memory** | What happened, in order (conversations, actions, outcomes) |
| 📌 **Semantic memory** | Durable facts and preferences distilled out of those episodes |
| 🔎 **Hybrid retrieval** | Keyword + vector search pulls back what matters for *this* turn |
| 🗜️ **Auto-summarization** | Old episodes are compressed, not deleted — long-term memory stays cheap |

The goal is a memory layer that's boring to integrate and hard to notice —
until you turn it off and the agent forgets your name.

---

## How It Works

```
        +---------------------+        local call / gRPC        +----------------------+
        |      Agent code      | -------------------------------> |     AIrecall Core     |
        |    (Python SDK)      | <------------------------------- |   (memory server)     |
        +---------------------+                                  +-----------+----------+
                                                                             |
                                                                 +-----------v----------+
                                                                 |    Storage Engine      |
                                                                 |  SQLite + vector idx   |
                                                                 +-----------+----------+
                                                                             |
                                                                 +-----------v----------+
                                                                 |   Optional MCP         |
                                                                 |   adapter              |
                                                                 +-----------------------+
```

On every turn, the SDK sends the current query to the core, which does a
**hybrid retrieval** pass — keyword + vector similarity — over both
episodic and semantic memory, and returns the top-k relevant memories to
inject into your prompt. In the background, a summarizer periodically
walks older episodes, extracts durable facts into semantic memory, and
compacts the rest.

---

## Install

```bash
# Python SDK - what your agent code imports
pip install airecall-sdk

# Core memory server - runs locally or as a sidecar
pip install airecall-sdk[server]
# or run it standalone:
airecall serve
```

`airecall init` scaffolds a local SQLite-backed store so you can start
storing and recalling memories immediately, no separate service required.
