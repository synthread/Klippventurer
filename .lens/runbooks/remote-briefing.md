# Remote Briefing Runbook

## Purpose

Use this runbook when answering questions like "what's new on the remote" for Klippventurer.

This exists because commit-only summaries are not complete if merge request review activity has not been checked.

## Required Coverage

A complete remote briefing should include:

1. Branch/ref changes and divergence.
2. Relevant merge requests for the active branch.
3. Review activity, especially unresolved comments from named reviewers.
4. A clear warning when GitLab auth blocks MR visibility.

## Minimal Flow

1. `git fetch --all --prune`
2. `git status --short --branch`
3. `git log` / `git diff` against the tracked upstream branch
4. `glab auth status`
5. `glab mr list --all --source-branch <branch>`
6. `glab mr view <id> --notes`

## Response Structure

Prefer this structure:

- `Code changes`
- `Review activity`
- `Action needed`

## Failure Rule

Do not present a commit-only summary as a complete remote update when MR comments have not been checked.
