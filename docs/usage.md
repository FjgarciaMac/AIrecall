# Usage

Practical walkthroughs for AIrecall.

## Remember and recall

```python
from airecall import Memory

memory = Memory(agent_id="support-bot")
memory.remember("User asked about refund policy for ORD-9921")
memory.remember_fact("preferred_contact", "email")

hits = memory.recall("what did we tell this user about refunds?")
for h in hits:
    print("-", h)
```

## Use the CLI

```bash
airecall init
airecall store "user prefers email over phone" --agent support-bot
## Pinning facts across sessions

When a user states a durable preference, promote it immediately instead of waiting for the summarizer:

```python
memory.remember_fact("preferred_contact", "email")
memory.remember_fact("timezone", "utc")
```
Facts are upserted per agent, so a later correction overwrites cleanly - no stale duplicates.


## Handling long sessions

Sessions that exceed `max_turns` are compacted automatically. The summary keeps entity references intact, so follow-up questions about a user mentioned early in the session still resolve correctly. Disable with `compact=False` if you pipe the raw transcript elsewhere.
