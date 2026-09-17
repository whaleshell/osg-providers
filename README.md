<h1 align="center">osg-providers</h1>

<p align="center">
  <strong>Provider catalog & compose</strong><br>
  Builtin provider profiles and effective-policy composition for sandboxes.
</p>
<p align="center">
  <a href="https://github.com/zorneth/osg-providers/actions/workflows/ci.yml"><img src="https://github.com/zorneth/osg-providers/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://pkg.go.dev/github.com/zorneth/osg-providers"><img src="https://pkg.go.dev/badge/github.com/zorneth/osg-providers.svg" alt="Go Reference"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License"></a>
  <a href="https://github.com/zorneth/osg-providers"><img src="https://img.shields.io/badge/Go-1.27+-00ADD8?logo=go" alt="Go Version"></a>
</p>
<p align="center">
  <sub>Part of the <a href="https://github.com/zorneth">zorneth / osg</a> ecosystem</sub>
</p>

---

## Overview

**osg-providers** ships YAML provider profiles (Cursor, GitHub, NVIDIA, …) and composes them onto a base policy to produce the effective network/credential set a sandbox runs with.

### Key Features

| Category | Capabilities |
|----------|--------------|
| **Catalog** | Builtin profiles under `profiles/` |
| **Compose** | Merge provider endpoints + env keys into base policy |
| **Custom** | Import/override profiles via gateway API |
| **Parity** | OpenShell-style provider attach model |

---

## Installation

```bash
go get github.com/zorneth/osg-providers@latest
```

**Requirements:** Go 1.27+

---

## Quick Start

```go
import (
    "github.com/zorneth/osg-core/policy"
    "github.com/zorneth/osg-providers/provider"
)

base, _ := policy.Load("base.yaml")
prof, _ := provider.LoadFile("profiles/cursor.yaml")
effective := provider.EffectivePolicy(base, []provider.Layer{{
    InstanceName: "cursor",
    Profile:      prof,
}}, false)
_ = effective
```

Profiles: [`profiles/`](./profiles/).

---

## Package Structure

| Path | Purpose |
|------|---------|
| `provider/` | Load, validate, compose |
| `profiles/` | Builtin YAML catalogs |


---

## Related

| Resource | Link |
|----------|------|
| Roadmap | [ROADMAP.md](./ROADMAP.md) |
| Organization | [https://github.com/zorneth](https://github.com/zorneth) |
| Organization overview | [github.com/zorneth](https://github.com/zorneth) |
| pkg.go.dev | [`github.com/zorneth/osg-providers`](https://pkg.go.dev/github.com/zorneth/osg-providers) |

## License

[MIT](./LICENSE) © zorneth
