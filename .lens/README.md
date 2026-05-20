# .lens

Repo-carried project memory for Klippventurer.

Use `.lens/` for approved, durable, reviewable project context that should travel with the repository.

`.lens/inbox/` is for small curated intake items that still need processing before promotion into stable `.lens` locations such as runbooks, schemas, or policies.

Raw research dumps, machine-generated artifacts, transient logs, and bulky local analysis outputs should not be stored in `.lens/`; keep them out of git unless they have been curated into reviewable documents.

Do not use `.lens/` as:

- the primary chat history store
- a dump for transient operational logs
- private agent memory
- a shadow runtime or undocumented automation layer

## Start Here

- [`.lens` spec sheet](lens-spec.md)
- [`.lens` version table](lens-version-table.md)
- [Single-repo policy](policies/single-repo.md)
- [Kalicoforge rename checklist (v0.4.0)](policies/kalicoforge-rename-v0.4.md)

## Relationship To Server Memory

- `.lens/` stores portable references, prompts, runbooks, policies, schemas, and other approved repo-specific durable context.
- A server-side memory store should hold heavier conversational state, file-attached chat threads, operational journals, approvals, and private agent memory.
- `.lens/` may point at server-backed context, but it should not mirror the whole server-side chat log into git.

## Promotion Path

This local `.lens/` guidance is intended to become public-facing documentation once the spec stabilizes.
