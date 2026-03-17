# ctx Self-Improvement: Filing Feedback

When you encounter friction while using ctx, you should file an issue. This helps improve the tool for all future agents.

## Three categories of problems

### 1. Skill guidance issue
The instructions in SKILL.md or a reference file misled you. You followed the documented workflow but got an unexpected result, or the documentation was missing a critical detail.

**Label**: `skill-guidance`

### 2. CLI functionality issue
A ctx command doesn't behave as expected — wrong output, missing flags, broken error messages, parameter design that creates confusion.

**Label**: `cli-bug`

### 3. System error (unrecoverable)
External service failure (Cloudflare API down, GitHub rate limit, network error). These are not ctx bugs — do NOT file an issue. Just report the error to the user.

## How to file

Use `gh` CLI to create an issue on the ctx repository:

```bash
gh issue create \
  --repo ethan-huo/ctx \
  --title "<concise title>" \
  --label "<label>" \
  --body "$(cat <<'EOF'
## What happened
<describe what you tried to do and what went wrong>

## Expected behavior
<what you expected to happen>

## Reproduction
<the exact command(s) you ran>

## Context
- Skill instruction that misled (if applicable): <quote the specific line>
- Command output / error message: <paste relevant output>
- Suggestion for improvement (if any): <your idea>
EOF
)"
```

## When to file

- A command produced confusing output that wasted your time
- You had to work around a skill instruction that was wrong or incomplete
- A parameter name or default value was misleading
- Error messages didn't help you recover

## When NOT to file

- External service is down (Cloudflare, GitHub) — not a ctx issue
- User gave you wrong input — not a ctx issue
- You're unsure if it's a problem — ask the user first before filing
