# CI and docs publishing

## Layout

- Root include: [.gitlab-ci.yml](../../.gitlab-ci.yml) → [ci/docs-repo.yml](../../ci/docs-repo.yml)
- Site config: [mkdocs.yml](../../mkdocs.yml)
- Hooks (preview banner, external link `rel`): [ci/mkdocs_hooks.py](../../ci/mkdocs_hooks.py)

Replaces the former `synthread/docs-viewer` include (inaccessible to CI; see issue #2).

## Channels

| `DOCS_CHANNEL` | Trigger branches | Behavior |
| --- | --- | --- |
| `stable` | `main`, tags | Production Pages site |
| `preview` | `next/firmware-package`, `preview/*`, MRs | Preview banner; **build only** (no Pages deploy, avoids overwriting stable site) |

## Local build

```sh
pip install mkdocs mkdocs-material pymdown-extensions
DOCS_CHANNEL=preview mkdocs serve
```

## Verify on GitLab

```sh
glab ci list -R synthread/proj/Klippventurer
glab ci trace -R synthread/proj/Klippventurer <pipeline-id>
```

## N32 probe artifacts

Preview branches run `build:n32probe` and publish `build/n32probe_usb/` as a pipeline artifact (see [ota-preview-channel.md](../schemas/ota-preview-channel.md)).
