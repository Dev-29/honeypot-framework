# honeypot-framework

A small Go-based honeypot framework for capturing and logging attacker interactions across multiple services.

## Overview

This project provides a lightweight framework to run simple service emulators (SSH, HTTP, etc.), collect interaction data, and store logs either to a file or a local SQLite database. It's intended for research, learning, and controlled lab use only.

## Features

- Structured logging (JSON) and file output
- Optional SQLite storage for attack logs
- Simple service modules under `pkg/services`

## Requirements

- Go 1.18+ (module-enabled)

## Quick start

1. Clone the repo (if you haven't already):

    git clone https://github.com/Dev-29/honeypot-framework.git
    cd honeypot-framework

2. Build:

    go build ./...

3. Configure

Edit `config/config.json` to adjust services, database path, or logging options. A sample config is included in `config/config.json`.

4. Run

    ./honeypot-framework

Or run the main package with `go run`:

    go run ./...

## Project layout

- `pkg/` — core packages (database, logger, services)
- `config/` — configuration files
- `web/` — optional web UI or static assets
- `logs/` — default location for file logs

## Logging & Database

- File logging is implemented in `pkg/logger`.
- SQLite support is available in `pkg/database`; set the path in `config/config.json`.

## License

See the `LICENSE` file in the repository root.

## Contributing

Contributions are welcome. Open issues or PRs for bugs and improvements.

## Security & Ethics

Use this software responsibly. Do not deploy it on networks where you don't have explicit permission to monitor traffic or capture data. Ensure collected data is handled lawfully and ethically.
