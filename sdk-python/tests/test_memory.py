"""Tests for the in-process Memory store and retrieval scoring."""

import os
import tempfile

import pytest

from airecall.core.memory import Memory


@pytest.fixture()
def memory(tmp_path):
    m = Memory(agent_id="test-agent", db_path=str(tmp_path / "m.db"))
    yield m
    m.close()


def test_remember_and_recall(memory):
    memory.remember("User asked about refund policy for ORD-9921")
    memory.remember("User prefers email over phone")
    hits = memory.recall("refund policy")
    assert any("refund" in h for h in hits)


def test_recall_empty(memory):
    assert memory.recall("anything") == []


def test_facts_upsert(memory):
    memory.remember_fact("preferred_contact", "email")
    assert memory.recall_fact("preferred_contact") == "email"
    memory.remember_fact("preferred_contact", "phone")
    assert memory.recall_fact("preferred_contact") == "phone"


