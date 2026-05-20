# Kalicoforge rename checklist (v0.4.0)

Planned **in-repo rebrand** at `v0.4.0`. Execute on the existing `synthread/proj/Klippventurer` project—no repository fork or history split.

## Preconditions

- [ ] `v0.3.x` USB installer preview milestones documented in [docs/preview-v0.3.md](../../docs/preview-v0.3.md) are met or explicitly deferred with release notes
- [ ] Public preview CI/CD publishes stable and preview docs channels (see [ci/docs-repo.yml](../../ci/docs-repo.yml))
- [ ] N32 validation evidence captured for supported Adventurer 3 targets
- [ ] SSH dev-deploy runbook exercised on at least one hardware profile

## GitLab project (manual, maintainer)

- [ ] Rename GitLab project/path to Kalicoforge (or `synthread/proj/Kalicoforge`) per org naming policy
- [ ] Update default branch protection and Pages URL
- [ ] Redirect old project URL if GitLab rename does not preserve redirects
- [ ] Update `glab` remotes and CI variables (`CI_PAGES_URL`, package registry paths)

## Repository content

- [ ] README title, logos, and links ([README.md](../../README.md), [images/](../../images/))
- [ ] Docs site name in [mkdocs.yml](../../mkdocs.yml)
- [ ] `.lens` references and runbooks
- [ ] OTA/preview artifact namespace (see [ota-preview-channel.md](../schemas/ota-preview-channel.md))
- [ ] Support matrix and compatibility strings where user-facing

## Out of scope for rename

- Splitting git history into a new repo
- Broad eMMC deployment (ADR 0002 / 0005 still apply until separately decided)
- Claiming production OTA before safestrap health gates are proven

## After rename

- [ ] Publish `v0.4.0` release notes describing Klippventurer → Kalicoforge continuity
- [ ] Update Discord/docs links and issue templates
