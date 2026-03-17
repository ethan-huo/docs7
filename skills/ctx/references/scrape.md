# ctx scrape

Extract specific elements from a webpage by CSS selector.

## Usage

```bash
ctx scrape <url> -s "h1" -s "p.description"
ctx scrape <url> -s "table.api-params" --text-only
```

| Flag | Short | Default | Description |
|---|---|---|---|
| (positional) | | optional | URL (can also be in `-d` body) |
| `--selector` | `-s` | | CSS selectors (repeatable: `-s "h1" -s "table"`) |
| `--text-only` | | false | Output plain text instead of JSON |
| `--data` | `-d` | | Full API request body (JSON5, `@file`, or `-` for stdin) |

When using `-d`, selectors go in the body as `elements` array — no `-s` flags needed.

## When to use

- Extract specific parts of a page (API tables, code blocks, pricing) without full-page markdown
- `ctx read` returns too much — scrape is surgical

## Output format

Default JSON:
```json
[
  {
    "selector": "h1",
    "results": [
      {"text": "API Reference", "html": "<h1>API Reference</h1>", "width": 800, "height": 40}
    ]
  }
]
```

With `--text-only`: one text value per line.

## Full API control via -d

```bash
ctx scrape -d '{
  url: "https://example.com",
  elements: [{selector: "h1"}, {selector: "table.pricing"}]
}'
```
