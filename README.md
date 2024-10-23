# Ubuntu to Fedora Omakub Script Converter

A Go-based tool that converts Omakub shell scripts to be compatible with Fedora systems. This tool features a Terminal User Interface (TUI) that makes it easy to convert installation scripts and other shell scripts from Ubuntu's package management commands to their Fedora equivalents.

## Features

- Interactive Terminal User Interface using Bubble Tea
- Automatic conversion of common Ubuntu commands to Fedora equivalents:
  - `apt` → `dnf`
  - `apt-get` → `dnf`
  - `add-apt-repository` → `dnf config-manager --add-repo`
  - And more package management command conversions
- Integration with the omakub repository for script management
- Batch processing of multiple shell scripts
- Preserves original script functionality while ensuring Fedora compatibility

## Installation

1. Ensure you have Go 1.23.2 or later installed
2. Clone this repository:
```bash
git clone https://github.com/yourusername/ubuntu-to-fedora.git
cd ubuntu-to-fedora
```
3. Build the project:
```bash
go build
```

## Usage

Run the application:
```bash
./ubuntu-to-fedora
```

The TUI will guide you through the process of:
1. Selecting source scripts to convert
2. Converting Ubuntu commands to their Fedora equivalents
3. Saving the converted scripts

## Command Conversions

The tool automatically converts the following Ubuntu commands to their Fedora equivalents:

| Ubuntu Command | Fedora Equivalent |
|---------------|-------------------|
| `sudo apt update` | `sudo dnf update` |
| `sudo apt upgrade` | `sudo dnf upgrade` |
| `sudo apt install` | `sudo dnf install` |
| `add-apt-repository` | `sudo dnf config-manager --add-repo` |
| `sudo apt autoremove` | `sudo dnf autoremove` |
| `sudo apt-get` | `sudo dnf` |
| `apt-get` | `dnf` |
| `sudo apt` | `sudo dnf` |
| `apt` | `dnf` |

## Dependencies

- Go 1.23.2 or later
- External packages:
  - github.com/charmbracelet/bubbletea - Terminal UI framework
  - github.com/go-git/go-git/v5 - Git operations
  - github.com/charmbracelet/lipgloss - Terminal styling
  - github.com/stretchr/testify - Testing framework

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
