"""Tests for the in-process Memory store and retrieval scoring."""

import os
import tempfile

import pytest

from airecall.core.memory import Memory


@pytest.fixture()
def memory(tmp_path):
    m = Memory(agent_id="test-agent", db_path=str(tmp_path / "m.db"))
