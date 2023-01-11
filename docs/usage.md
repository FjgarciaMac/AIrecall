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
