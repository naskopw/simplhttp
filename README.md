# simplhttp 🚀

**simplhttp** is a modern, bidirectional HTTP file server designed as a powerful replacement for `python -m http.server`. It features a beautiful React-based UI, supports file uploads, directory downloads (as ZIP), and comes with built-in security features like Basic Auth and HTTPS.

All assets are embedded into a single, static binary, making it extremely easy to distribute and run anywhere.

[![Build Status](https://github.com/naskopw/simplhttp/actions/workflows/build.yml/badge.svg)](https://github.com/naskopw/simplhttp/actions)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25.4-blue.svg)](go.mod)

## ✨ Features

- 🎨 **Modern UI:** Clean, responsive interface built with React and Mantine.
- 📤 **File Uploads:** Drag-and-drop support for uploading files directly from your browser.
- 📂 **Directory Navigation:** Fast and intuitive browsing of your file system.
- 🤐 **ZIP Downloads:** Download entire directories as a single ZIP archive on the fly.
- 🔐 **Security:**
  - Basic Authentication support (`user:pass`).
  - Read-only mode to prevent modifications.
  - Automatic HTTPS (self-signed) or custom certificates.
- 📦 **Single Binary:** No dependencies, everything is packed into one file.
- 🚀 **Performant:** Built with Go and Echo for high efficiency.

## 🚀 Quick Start

### Installation

Download the latest binary for your platform from the [Releases](https://github.com/naskopw/simplhttp/releases) page.

Alternatively, install via Go:

```bash
go install github.com/naskopw/simplhttp@latest
```

### Usage

Run the server in your current directory:

```bash
simplhttp
```

Specify a port and directory:

```bash
simplhttp -p 8080 -d /path/to/my/files
```

Enable Basic Auth and Read-Only mode:

```bash
simplhttp --auth "admin:password" --readonly
```

Enable HTTPS:

```bash
simplhttp --https
```

## 🛠️ Command Line Options

| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--port` | `-p` | `8080` | Port to run the server on |
| `--dir` | `-d` | `.` | Root directory to serve |
| `--auth` | `-a` | `""` | Basic auth credentials (`user:pass`) |
| `--readonly`| `-r` | `false` | Disable uploads and file modifications |
| `--max-size`| `-m` | `100MB` | Max upload size (e.g., 100MB, 1GB) |
| `--https` | | `false` | Enable HTTPS (self-signed if no cert/key) |
| `--cert` | | `""` | Path to TLS certificate file |
| `--key` | | `""` | Path to TLS private key file |

## 🏗️ Development

### Prerequisites

- [Go](https://go.dev/doc/install) (1.25+)
- [Node.js](https://nodejs.org/) & [npm](https://www.npmjs.com/)

### Build from source

1. Clone the repository:
   ```bash
   git clone https://github.com/naskopw/simplhttp.git
   cd simplhttp
   ```

2. Build the UI and the Go binary:
   ```bash
   make build-production
   ```

The resulting `simplhttp` binary will be in the root directory.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
