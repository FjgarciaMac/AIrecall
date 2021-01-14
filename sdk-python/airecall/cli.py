"""airecall CLI - manage memories from the terminal.

Talks to the in-process store by default; use --server to target the
Go memory server.
"""

from __future__ import annotations

import argparse
import sys
from typing import Optional

from airecall.core.memory import Memory


def cmd_init(args) -> int:
    Memory(agent_id=args.agent, db_path=args.db)
    print(f"initialized memory store at {args.db}")
    return 0


