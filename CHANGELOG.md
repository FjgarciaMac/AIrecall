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

## [0.3.0] - 2023-10-09

### Added
- Summarization pass with fact promotion
- `airecall summarize` CLI command

## [0.2.0] - 2020-11-05

### Added
- Python SDK with the `Memory` class
- Semantic facts API (`remember_fact` / `recall_fact`)
- LangChain retriever adapter
- `airecall fact` CLI command

## [0.1.0] - 2018-06-21

### Added
- Episodic memory store (SQLite)
- Hybrid retrieval groundwork (keyword scoring)
- First recall API

## [0.0.2] - 2016-09-14

### Added
- Schema draft for episodes and facts tables

## [0.0.1] - 2015-05-20

### Added
- Project scaffold and design notes
