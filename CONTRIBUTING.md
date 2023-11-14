# Contributing to AIrecall

Thanks for considering a contribution. AIrecall is deliberately small:
a Python SDK, a Go memory server, and a shared SQLite schema.

## Layout

- `sdk-python/airecall` - the SDK (memory, client, adapters, cli)
- `server/` - the Go memory server (`airecalld`) and its tests
- `docs/` - architecture and usage notes

## Ground rules

1. One logical change per pull request.
2. The SQLite schema is the contract between SDK and server - changes need
   a note in `docs/architecture.md` and a CHANGELOG entry.
3. Keep the SDK dependency-free apart from `requests`; adapters live in
   `airecall/adapters/` and must import optional frameworks lazily.

## Running the tests

Python (3.9+):

```
cd sdk-python
pip install -e ".[dev]"
pytest
```

Go (1.22+):

```
cd server
go test ./... -count=1
```

Both suites must be green before a PR is merged. `make test` runs both.
