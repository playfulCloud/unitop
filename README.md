# Unitop

Unitop is a terminal UI for monitoring and managing systemd services.

It gives you a live service table, name filtering, common `systemctl` actions,
and quick access to recent `journalctl` logs without leaving the terminal.

![Unitop demo](docs/demo.gif)

## Features

- Live service state refresh from `systemctl show`
- Service discovery from `systemctl list-unit-files`
- Manual service lists for focused monitoring
- Filter services by name
- Start, stop, restart, enable, and disable services
- Open recent logs for the selected service with `journalctl -u`
- Configurable refresh interval and discovery rules
- Thread-safe in-memory service state store

## Requirements

- Go 1.25 or newer
- Linux with systemd
- `systemctl` available in `PATH`
- `journalctl` available in `PATH`

Some service actions may require elevated permissions depending on your system
policy. Unitop does not currently manage `sudo` for you.

## Install

Install the latest tagged version with Go:

```sh
go install github.com/playfulCloud/unitop/cmd/unitop@latest
```

Or download a prebuilt Linux binary from the
[GitHub releases](https://github.com/playfulCloud/unitop/releases) page.

The project can also be run directly from source:

```sh
go run ./cmd/unitop
```

Or built as a local binary:

```sh
make build
./bin/unitop
```

By default, Unitop reads its config from:

```text
~/.config/unitop/unitop.yaml
```

If the file does not exist, Unitop creates it with a standard configuration on
first run.

Use a custom config path with:

```sh
unitop --config ./configs/unitop.yaml
```

Print the current version with:

```sh
unitop --version
```

## Configuration

Example:

```yaml
mode: all # selected | all
refresh_interval: 5s

services:
  - docker.service
  - NetworkManager.service
  - ssh.service

discovery:
  include:
    - "*.service"
  exclude:
    - "systemd-*"
    - "user@*.service"
    - "getty@*.service"
    - "autovt@*.service"
  states:
    - disabled
    - enabled
    - enabled-runtime
    - linked
    - linked-runtime
```

### Modes

`mode: selected` uses the explicit `services` list.

`mode: all` discovers services with `systemctl list-unit-files` and applies the
`discovery` filters.

### Discovery

`include` and `exclude` use shell-style patterns such as `*.service` and
`systemd-*`.

If no discovery states are configured, Unitop defaults to:

```yaml
states:
  - enabled
  - enabled-runtime
  - linked
  - linked-runtime
```

## Keybindings

| Key | Action |
| --- | --- |
| `up` / `k` | Move selection up |
| `down` / `j` | Move selection down |
| `/` | Start a new filter |
| `enter` | Apply filter |
| `esc` | Close filter or clear active filter |
| `r` | Restart selected service |
| `s` | Start selected service |
| `x` | Stop selected service |
| `e` | Enable selected service |
| `d` | Disable selected service |
| `l` | Open recent logs for selected service |
| `q` / `ctrl+c` | Quit |

## Development

Run the test suite:

```sh
go test ./...
```

Run race tests:

```sh
go test -race ./...
```

Run vet:

```sh
go vet ./...
```

The Makefile also provides:

```sh
make test
make test-race
make vet
make build
```

## Project Status

Unitop is usable as a source-built or release-built tool. It creates a default
config on first run, supports custom config paths, blocks duplicate service
actions, uses longer timeouts for service actions, starts monitoring
asynchronously, and adapts the table height to the terminal size.

Planned improvements:

- Homebrew tap and AUR packaging
- Shell completions
- Changelog

## License

Unitop is distributed under the MIT License. See [LICENSE](LICENSE).
