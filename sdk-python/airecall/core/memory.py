"""Memory - the public SDK surface.

In dev mode (no server running) Memory falls back to an in-process
SQLite store so `airecall init` works without the Go core. When the
server is reachable it forwards everything to it.
"""

from __future__ import annotations

import json
import os
import sqlite3
import uuid
from datetime import datetime, timezone
from typing import Optional
