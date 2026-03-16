# docs7

Library documentation finder — uses [Context7](https://context7.com) index to locate documentation sources, then reads full documents instead of RAG chunks.

## Why

Tools like ctx7 and ref both solve "find docs for AI agents" but with trade-offs:

- **ctx7** has great search but returns fragmented RAG chunks (60-200 tokens each)
- **ref** returns full documents but search accuracy is inconsistent

docs7 takes ctx7's search index and discards the chunks, keeping only the source URLs. Then it reads the full original documents via GitHub API or HTTP with markdown content negotiation.

## Install

```bash
make install
```

Requires Go 1.24+. Installs to `bin/docs7`, symlinks to `/usr/local/bin/docs7` if `bin/` isn't in PATH.

## Usage

```bash
# Find documentation sources for a library
docs7 docs mlx-swift "GPU stream thread safety"
docs7 docs sparkle "appcast auto update"
docs7 docs convex "Swift client authentication"

# Read a full document
docs7 read github://ml-explore/mlx-swift/Source/MLX/Documentation.docc/MLXArray.md
docs7 read https://sparkle-project.org/documentation/index

# Find a library by name
docs7 search react-native
```

## Authentication

docs7 shares credentials with ctx7 (`~/.context7/credentials.json`).

```bash
# Login (opens browser, OAuth PKCE)
docs7 login

# Check status
docs7 whoami

# Logout
docs7 logout
```

GitHub reads use your `gh auth` token automatically.

## How `read` works

| URL | Strategy |
|---|---|
| `github://owner/repo/path` | GitHub Contents API |
| `https://github.com/.../blob/...` | Auto-converted to GitHub API |
| Any `https://` | `Accept: text/markdown` negotiation → Jina Reader fallback |

## AI Agent Integration

The `skills/find-docs/` directory contains a SKILL.md for Claude Code / Cursor / similar tools. Install it with your agent's skill mechanism.
