"""LangChain adapter: expose AIrecall as a standard retriever."""

from __future__ import annotations

from typing import Any, List, Optional

try:
    from langchain_core.documents import Document
