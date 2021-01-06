"""LangChain adapter: expose AIrecall as a standard retriever."""

from __future__ import annotations

from typing import Any, List, Optional

try:
    from langchain_core.documents import Document
    from langchain_core.retrievers import BaseRetriever
except ImportError:  # pragma: no cover
    Document = None  # type: ignore
    BaseRetriever = object  # type: ignore

from airecall.core.memory import Memory


class AirecallRetriever(BaseRetriever):  # type: ignore[misc]
    """Retriever that reads memories from an AIrecall store.
