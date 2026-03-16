---
name: find-docs
description: >-
  Retrieves authoritative, up-to-date technical documentation, API references,
  configuration details, and code examples for any developer technology.

  Use this skill whenever answering technical questions or writing code that
  interacts with external technologies. This includes libraries, frameworks,
  programming languages, SDKs, APIs, CLI tools, cloud services, infrastructure
  tools, and developer platforms.

  Common scenarios:
  - looking up API endpoints, classes, functions, or method parameters
  - checking configuration options or CLI commands
  - answering "how do I" technical questions
  - generating code that uses a specific library or service
  - debugging issues related to frameworks, SDKs, or APIs
  - retrieving setup instructions, examples, or migration guides
  - verifying version-specific behavior or breaking changes

  Prefer this skill whenever documentation accuracy matters or when model
  knowledge may be outdated.
---

# docs7 — Library Documentation Finder

Find library documentation, then read the full source documents. Two-step workflow: **search → read**.

Binary: `docs7`

## Workflow

### Step 1: Find documentation sources

```bash
docs7 docs <library-name> "<query>"
```

This returns a list of relevant documents with descriptions, plus their URLs.

```bash
# Examples
docs7 docs mlx-swift "GPU stream thread safety"
docs7 docs sparkle "appcast auto update configuration"
docs7 docs convex "Swift client authentication"
```

If you already know the library ID (format `/owner/repo`), pass it directly:

```bash
docs7 docs /ml-explore/mlx-swift "lazy evaluation"
```

### Step 2: Read the full documents

Pick the most relevant URL(s) from Step 1 and read them:

```bash
docs7 read <url>
```

The `read` command auto-detects the URL type:

| URL pattern | Strategy |
|---|---|
| `github://owner/repo/path` | GitHub API (authenticated via `gh auth`) |
| `https://github.com/.../blob/...` | GitHub API (auto-converted) |
| `https://...` (serves markdown) | Direct fetch with `Accept: text/markdown` |
| `https://...` (serves HTML) | Jina Reader fallback → clean markdown |

### Putting it together

```bash
# 1. Find
docs7 docs react "useEffect cleanup async"
# output:
#   1. **React Hooks Reference**
#      - ...
#   ---
#   - github://facebook/react/docs/hooks-reference.md

# 2. Read the most relevant one
docs7 read github://facebook/react/docs/hooks-reference.md
```

## Writing Good Queries

The query directly affects result quality. Be specific.

| Quality | Example |
|---|---|
| Good | `"SwiftUI NavigationStack path binding programmatic navigation"` |
| Good | `"Express.js middleware error handling async"` |
| Bad | `"navigation"` |
| Bad | `"middleware"` |

Include the programming language or framework name when ambiguous.

## When to use `search` instead of `docs`

Use `docs7 search` when you need to **find the right library first**, before querying its docs:

```bash
# "Which library is this?"
docs7 search swift-testing
docs7 search convex "mobile client"
```

This returns a ranked list of matching libraries with IDs you can feed to `docs`.

## Limits

- Do not call `docs7 docs` more than 3 times per question. Use the best result you have.
- Coverage is best for libraries with GitHub repos or documentation sites. Apple-native C APIs (CoreAudio, CGEventTap) are not indexed — use the `apple-docs` skill for those.
