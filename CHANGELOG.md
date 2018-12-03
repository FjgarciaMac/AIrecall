# Changelog

All notable changes to AIrecall are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project
adheres to [Semantic Versioning](https://semver.org/).

## [0.4.0] - 2026-03-02

### Added
- Go memory server (`airecalld`) with JSON API
- MCP adapter scaffold in the server

### Changed
- Hybrid retrieval scoring tuned (0.65 keyword / 0.35 vector)

## [0.3.0] - 2025-11-14

### Added
- Summarization pass with fact promotion
- `airecall summarize` CLI command

## [0.2.0] - 2025-06-09

### Added
- Semantic facts API (`remember_fact` / `recall_fact`)
- LangChain `MemoryRetriever` adapter
- `airecall fact` CLI command

## [0.1.0] - 2024-11-18

### Added
- Episodic memory store (SQLite)
- Hybrid retrieval (keyword + vector)
- Python SDK with `Memory` class
- `airecall init | store | recall` CLI

## [0.0.2] - 2024-07-25

### Added
- Schema draft for episodes and facts tables

## [0.0.1] - 2024-03-12

### Added
- Project scaffold and design notes

<!-- draft note 308 -->
