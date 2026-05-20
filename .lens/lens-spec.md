# `.lens` Spec Sheet

## Purpose

`.lens/` is the repo-carried durable memory layer.

It exists to hold approved context that should move with the repository and remain reviewable in git.

## Core Role

`.lens/` is for durable project memory, not full operational history.

It should complement a server-side memory system rather than replace one.

## What Belongs In `.lens/`

- prompts and working conventions
- runbooks
- MCP setup notes
- CI/CD notes
- approved project memory
- stable schemas and reference documents
- portable references to server-backed file-attached chats or room artifacts
- small curated helper surfaces such as:
  - `.lens/runbooks/`
  - `.lens/inbox/`
  - `.lens/prompts/`
  - `.lens/schemas/`
  - `.lens/tools/`

## Inbox Role

`.lens/inbox/` is for small, human-authored, reviewable intake items that have not yet been normalized.

Examples:

- agent handoff notes
- candidate runbook or standards changes
- distilled observations that still need routing or confirmation

Inbox items are not automatically durable memory. They should be either promoted, rewritten, moved elsewhere, or deleted once processed.

## What Does Not Belong In `.lens/`

- raw full chat logs
- heavy conversational history
- transient operational logs
- autonomous-run journals by default
- private agent memory
- per-agent preference or dissent history
- secrets or credentials
- large machine-generated caches
- raw firmware extraction trees
- scratch notes that have not been curated
- bulky local analysis artifacts better kept in ignored paths

## State Boundary

- `.lens/`: approved, portable, repo-specific durable context
- server memory/database: heavy chat history, agent context, approvals, journals, attached threads, operational history
- client UI: review and interaction surface, not source of truth

## Design Rules

1. Keep `.lens/` concise and reviewable.
2. Prefer references over mirroring large server-side state.
3. Promote content into `.lens/` only when it is durable and worth carrying in git.
4. Keep `.lens/` human-readable.
5. Keep changes attributable and easy to diff.
6. Treat inbox items as intake, not as a final destination.

## Future Public-Docs Direction

This file is the seed for the public `.lens` documentation.

When promoted, the public doc should also define:

- recommended directory layout
- reference formats for server-backed attached chats
- versioning and migration rules
- review/approval policy for `.lens` updates
