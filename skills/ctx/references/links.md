# ctx links

Extract all links from a webpage.

## Usage

```bash
ctx links <url>
ctx links <url> --internal-only
ctx links <url> --visible-only
```

| Flag | Default | Description |
|---|---|---|
| (positional) | optional | URL (can also be in `-d` body) |
| `--visible-only` | false | Only visible links |
| `--internal-only` | false | Exclude external domain links |
| `--no-cache` | false | Bypass cache |
| `-d` | | Full API request body (JSON5, `@file`, or `-` for stdin) |

Output: one URL per line on stdout.

## When to use

- Explore a documentation site's structure before selectively reading pages
- Find all pages linked from an index/hub page

## Typical workflow

```bash
# 1. Get links from a docs index
ctx links https://docs.example.com/

# 2. Read the relevant ones
ctx read https://docs.example.com/api/auth
```
