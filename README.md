# simplhttp 🚀

**simplhttp** is a modern, mobile-friendly HTTP file server and web server designed as a powerful replacement for `python -m http.server`.

Share files instantly, upload from any device, browse directories, download entire folders as ZIP archives, and secure access with HTTPS and authentication—all from a single static Go binary.

Built with Go and a modern React-powered interface, simplhttp makes file sharing and management effortless on desktop and mobile browsers alike.

[![Build Status](https://github.com/naskopw/simplhttp/actions/workflows/build.yml/badge.svg)](https://github.com/naskopw/simplhttp/actions)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25.4-blue.svg)](go.mod)

---

## ✨ Why simplhttp?

* 🚀 Better alternative to `python -m http.server`
* 🌐 Modern HTTP/HTTPS file server
* 📱 Mobile-friendly responsive web interface
* 📤 Upload files directly from your browser
* 📂 Browse and manage directories easily
* 📦 Download entire folders as ZIP archives
* 🔐 Built-in authentication and HTTPS support
* 🧩 Single static binary with no external dependencies

---

## ✨ Features

* 🌐 **Modern HTTP File Server**

  * Serve files and directories over HTTP or HTTPS.
  * Simple setup with sensible defaults.

* 🎨 **Responsive Web UI**

  * Clean, modern interface built with React and Mantine.
  * Optimized for desktop, tablet, and mobile devices.

* 📤 **File Uploads**

  * Drag-and-drop support.
  * Upload files directly from your browser.

* 📂 **Directory Navigation**

  * Fast and intuitive file browsing.
  * Navigate directories with ease.

* 📦 **ZIP Downloads**

  * Download entire directories as ZIP archives on the fly.

* 🔐 **Security Features**

  * Basic Authentication support (`user:pass`).
  * Read-only mode to prevent modifications.
  * Automatic HTTPS with self-signed certificates.
  * Support for custom TLS certificates.

* 📦 **Single Binary Distribution**

  * No runtime dependencies.
  * Everything embedded into one executable.

* 🚀 **High Performance**

  * Built with Go and Echo.
  * Lightweight and efficient.

---

## 🚀 Quick Start

### Installation

Download the latest binary for your platform from the Releases page.

Or install directly with Go:

```bash
go install github.com/naskopw/simplhttp@latest
```

### Start a File Server

Serve the current directory:

```bash
simplhttp
```

Specify a port, directory, and a custom base route:

```bash
simplhttp -p 8080 -d /path/to/my/files -R /my-files
```

### Serve a Specific Directory

```bash
simplhttp -d /path/to/files
```

### Change the Port

```bash
simplhttp -p 3000
```

### Enable Authentication

```bash
simplhttp --auth "admin:password"
```

### Enable Read-Only Mode

```bash
simplhttp --auth "admin:password" --readonly
```

### Enable HTTPS

```bash
simplhttp --https
```

### Use Custom TLS Certificates

```bash
simplhttp \
  --cert server.crt \
  --key server.key
```

---

## 🛠️ Command Line Options

| Flag         | Shorthand | Default | Description                                        |
| ------------ | --------- | ------- | -------------------------------------------------- |
| `--port`     | `-p`      | `8080`  | Port to run the server on                          |
| `--dir`      | `-d`      | `.`     | Root directory to serve                            |
| `--route`    | `-R`      | `/`     | Base route to serve the UI and files               |
| `--auth`     | `-a`      | `""`    | Basic auth credentials (`user:pass`)               |
| `--readonly` | `-r`      | `false` | Disable uploads and file modifications             |
| `--max-size` | `-m`      | `100MB` | Maximum upload size (e.g. `100MB`, `1GB`)          |
| `--https`    |           | `false` | Enable HTTPS (self-signed if no cert/key provided) |
| `--cert`     |           | `""`    | Path to TLS certificate file                       |
| `--key`      |           | `""`    | Path to TLS private key file                       |

---

## 📱 Mobile Friendly

simplhttp is fully responsive and works great on:

* 📱 Smartphones
* 📲 Tablets
* 💻 Laptops
* 🖥️ Desktop computers

Manage files, upload content, and browse directories from any modern browser without installing additional software.

---

## 🔒 Security

simplhttp includes several built-in security features:

* Basic Authentication
* HTTPS with automatic self-signed certificates
* Custom TLS certificate support
* Read-only mode
* Configurable upload size limits

For internet-facing deployments, it is recommended to:

* Enable HTTPS
* Configure authentication
* Use strong passwords
* Run behind a reverse proxy when appropriate

---

> **Note:** Internal API routes are served under `/simpl/api/` to avoid conflicts with your files.


## 🏗️ Development

### Prerequisites

* Go 1.25+
* Node.js
* npm

### Build From Source

Clone the repository:

```bash
git clone https://github.com/naskopw/simplhttp.git
cd simplhttp
```

Build the frontend and Go binary:

```bash
make build-production
```

The resulting `simplhttp` executable will be generated in the project root.

---

## 🤝 Contributing

Contributions, bug reports, feature requests, and pull requests are welcome.

If you find a bug or have an idea for an improvement, please open an issue.

---

## 📄 License

This project is licensed under the MIT License.

See the [LICENSE](LICENSE) file for details.
