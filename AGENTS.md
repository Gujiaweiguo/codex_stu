<INSTRUCTIONS>
# Repository Collaboration Guide (Flexible)

## Repository Purpose
- This repository is mainly a set of study notes/tutorials for using Codex and doing project development (mostly Markdown).
- Priorities: clear explanations, reproducible steps, and commands/code that can be copy-pasted and used directly.

## Default Workflow
1) Read the relevant Markdown first (usually `README.md` and/or `docs/codex_cource.md`).
2) Only ask 1–3 clarifying questions when necessary; otherwise make the “smallest reasonable” change.
3) When modifying code snippets, try to compile/run-verify in a temporary directory; keep dependencies minimal.

## Scope and Safety
- Unless explicitly requested, prefer modifying only Markdown and assets under `image/`.
- Avoid large rewrites; keep changes focused, reversible, and easy to review.
- Unless explicitly requested, do not introduce new toolchains/formatters.
- Unless explicitly requested, do not automatically `git commit`/`git push`.
- Do not write sensitive information into the repo (passwords, tokens, private keys, DSNs, etc.); use placeholders in examples or recommend environment variables/`.env` (and make sure it is included in `.gitignore`).

## Language and Style
- Default to communicating and writing docs in Chinese, unless the user explicitly requests English/bilingual.
- Show commands and code in fenced code blocks (```).
- In shell examples, prefer using `python3` for simple JSON field extraction (assuming `jq` may not be installed), while optionally mentioning `jq` as an alternative.
- Keep writing concise: provide executable steps/conclusions first, then only the necessary background; avoid long, generic explanations.
- When referencing files, include paths (and line numbers when needed) to make navigation easier.

## Conversation and Change Conventions (Suggested)
- When the user says “from scratch”: prefer organizing steps as “create new directory → `git init` → initialize dependencies → implement in small steps → verify”.
- When the user asks to “carefully review” documentation: focus on copy-pastable commands, dependency versions, prerequisites (e.g., whether `python3/jq/mysql` is needed), and whether examples can actually run.
- When a command needs to be executed: prefer reading/verifying it first; for write/install commands, explain what will change and the impact before running them.
- Default to a “plan → confirmation → step-by-step execution” workflow: propose a 3–6 step plan (each step: which files to change + how to verify + expected result), then wait for the user to confirm “start step N” before editing code/docs.
- During step-by-step execution, follow “scope constraints”: if the user specifies a file list, modify only those files; otherwise propose the smallest reasonable set of files and explain why.
- If the plan needs adjustment, pause implementation first: explain why, provide an updated plan, and wait for confirmation.
- Each step’s delivery includes: change summary, review points (diff/risk), minimal verification commands (with expected output), and next-step suggestions.

## Verification and Delivery
- If you touch commands/code snippets, provide a minimal verification path (e.g., `go test ./...`, `go run .` + a `curl` check).
- When delivering, briefly list changed files and key points; don’t paste large blocks of repeated content.
</INSTRUCTIONS>
