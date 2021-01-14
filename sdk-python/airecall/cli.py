"""airecall CLI - manage memories from the terminal.

Talks to the in-process store by default; use --server to target the
Go memory server.
"""

from __future__ import annotations

import argparse
import sys
