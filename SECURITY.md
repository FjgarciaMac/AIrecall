# Security

AIrecall stores agent memories locally. The security model is simple: the
store file is only as safe as the machine it lives on.

## Supported versions

| Version | Supported |
|---------|-----------|
| 0.4.x   | yes       |
| < 0.4   | no        |

## Reporting a vulnerability

Use GitHub's private vulnerability reporting (Security -> Report a
vulnerability) or email the maintainers directly. Include the version,
a minimal reproduction, and the impact you believe it has. We aim to
respond within 7 days.

## Scope

- SQL injection through the retrieval or fact APIs.
- Path traversal via `--db` / `storage.backend` values.
- Anything that lets one agent read another agent's memories.
- Prompt-injection paths that reach tool execution in the adapters.
