# ctx screenshot

Capture a webpage as a PNG image.

## Usage

```bash
ctx screenshot <url>
ctx screenshot <url> --full-page
ctx screenshot <url> --selector ".main-content" -o output.png
```

| Flag | Short | Default | Description |
|---|---|---|---|
| (positional) | | optional | URL (can also be in `-d` body) |
| `--output` | `-o` | auto-generated in cache dir | Output file path |
| `--full-page` | | false | Capture full scrollable page |
| `--selector` | | | Screenshot only the element matching this CSS selector |
| `--no-cache` | | false | Bypass cache |
| `--data` | `-d` | | Full API request body (JSON5, `@file`, or `-` for stdin) |

Output: file path on stdout. Read the image with your file-read tool to view it.

## When to use

- Page has visual information that markdown can't capture (UI layouts, charts, diagrams)
- Need to verify how a rendered page looks

## Full API control

```bash
ctx screenshot -d '{
  url: "https://example.com",
  viewport: {width: 390, height: 844},
  screenshotOptions: {type: "jpeg", quality: 80}
}'
```

Flags (`--full-page`, `--selector`, `-o`) override corresponding fields in `-d` body.
