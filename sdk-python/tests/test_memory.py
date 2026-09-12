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


def test_agents_isolated(tmp_path):
    a = Memory(agent_id="a", db_path=str(tmp_path / "x.db"))
    b = Memory(agent_id="b", db_path=str(tmp_path / "x.db"))
    a.remember("secret for agent a")
    assert b.recall("secret") == []
    a.close()
    b.close()


def test_top_k_limits(memory):
    for i in range(10):
        memory.remember(f"keyword unique-{i} memory")
    hits = memory.recall("keyword", top_k=3)
    assert len(hits) == 3


def test_negative_top_k_is_empty(memory):
    memory.remember("alpha beta")
    assert memory.recall("alpha", top_k=-1) == []


def test_summarize_noop_dev(memory):
    memory.remember("one episode")
    result = memory.summarize()
    assert result["episodes"] == 1


def test_context_manager(tmp_path):
    with Memory(agent_id="x", db_path=str(tmp_path / "y.db")) as m:
        m.remember("inside context")
        assert m.recall("inside")
    assert os.path.exists(str(tmp_path / "y.db"))