# ctx crawl

Crawl a website, collecting markdown content from multiple pages.

## Usage

```bash
ctx crawl <url> --limit 20
ctx crawl <url> --limit 50 --depth 2 --include "*/api/*" --exclude "*/changelog/*"
```

| Flag | Default | Description |
|---|---|---|
| (positional) | required | URL to crawl, or job ID (UUID) to resume |
| `--limit` | 10 | Max pages to crawl |
| `--depth` | 0 | Max link depth |
| `--include` | (none) | URL include patterns (glob, repeatable) |
| `--exclude` | (none) | URL exclude patterns (glob, repeatable) |
| `--no-wait` | false | Start and return job ID without waiting |
| `--cancel` | false | Cancel a running crawl job |
| `-d` | | Full API request body (JSON5, `@file`, or stdin) |

## Output format

Each crawled page outputs as:
```
## https://docs.example.com/page1
<markdown content>

---
## https://docs.example.com/page2
<markdown content>
```

Progress goes to stderr. Data goes to stdout.

## When to use

- Need content from many pages under a documentation site
- Single `ctx read` isn't enough — you need breadth across pages

## Async mode

For large crawls, use `--no-wait` to get a job ID immediately:

```bash
ctx crawl <url> --no-wait          # prints job ID
ctx crawl <job-id>                 # resume / get results
ctx crawl <job-id> --cancel        # cancel
```

The command auto-detects whether the argument is a URL or UUID.

## Timeout

Polling times out after 10 minutes. Resume with the job ID:
```bash
ctx crawl <job-id>
```
