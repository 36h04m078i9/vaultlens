# vaultlens

A read-only Vault secret explorer with fuzzy search and audit-log integration for ops teams.

---

## Installation

```bash
go install github.com/yourorg/vaultlens@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/vaultlens.git && cd vaultlens && go build -o vaultlens .
```

---

## Usage

Set your Vault address and token, then run `vaultlens` to start exploring:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

vaultlens
```

Use the interactive fuzzy search to navigate secret paths:

```bash
# Search a specific mount path
vaultlens --mount secret/ops

# Output matching secret paths as JSON (no values exposed)
vaultlens --mount secret/ops --output json
```

> **Note:** `vaultlens` is strictly read-only. It never writes, deletes, or modifies secrets. All access is logged via Vault's audit backend.

---

## Features

- 🔍 Fuzzy search across secret paths
- 📋 Audit-log integration — every lookup is traceable
- 🔒 Read-only by design — safe for shared ops environments
- 📦 Single binary, no dependencies

---

## Configuration

| Flag | Env Var | Description |
|------|---------|-------------|
| `--addr` | `VAULT_ADDR` | Vault server address |
| `--token` | `VAULT_TOKEN` | Vault authentication token |
| `--mount` | `VAULTLENS_MOUNT` | Secret engine mount path |
| `--output` | — | Output format: `table` (default) or `json` |

---

## License

MIT © yourorg