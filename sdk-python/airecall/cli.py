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


def cmd_store(args) -> int:
    memory = Memory(agent_id=args.agent, db_path=args.db,
                    server_url=args.server)
    memory.remember(" ".join(args.text))
    print("stored")
    return 0


def cmd_recall(args) -> int:
    memory = Memory(agent_id=args.agent, db_path=args.db,
                    server_url=args.server)
    hits = memory.recall(" ".join(args.query), top_k=args.top_k)
    for hit in hits:
        print(f"- {hit}")
    return 0


def cmd_fact(args) -> int:
    memory = Memory(agent_id=args.agent, db_path=args.db,
                    server_url=args.server)
    if args.value is None:
        value = memory.recall_fact(args.key)
        if value is None:
            print("(no such fact)")
            return 1
        print(value)
        return 0
    memory.remember_fact(args.key, args.value)
    print("stored")
    return 0


def cmd_summarize(args) -> int:
    memory = Memory(agent_id=args.agent, db_path=args.db,
                    server_url=args.server)
    result = memory.summarize()
    print(f"scanned {result.get('episodes', 0)} episodes")
    return 0


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="airecall",
