# Single-repo policy

## Decision

Klippventurer uses **one GitLab repository** for all product, preview, and tooling work:

- Canonical project: `synthread/proj/Klippventurer`
- Canonical remote: `https://gitlab.com/synthread/proj/Klippventurer.git`

There is **no parallel Kalicoforge repository**. The local workspace directory may be named `kalicoforge` for convenience; that is not a second product repo.

## Branch model

| Branch / pattern | Role |
| --- | --- |
| `main` | Stable `v0.2.x` manual-install line and published stable docs |
| `next/firmware-package` | Active `v0.3.x` USB installer preview integration |
| `preview/*` | Short-lived preview packaging (e.g. `preview/kalicoforge-0.4` for v0.4 staging) |

Preview branches land into `next/firmware-package` or `main` via merge requests; they are not long-lived forks.

## Kalicoforge rename

The product rename to **Kalicoforge** happens at the planned **`v0.4.0`** release by renaming/rebranding this same GitLab project and updating docs/assets. See [kalicoforge-rename-v0.4.md](kalicoforge-rename-v0.4.md).

## Agent rule

Do not create, recommend, or plan work against a second Kalicoforge GitLab repo unless this policy is explicitly revised.
