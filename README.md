# dtctl

**dtctl** is a command-line interface (CLI) tool for interacting with [Dependency-Track](https://dependencytrack.org/), allowing you to manage contexts and perform various operations such as fetching projects. Inspired by tools like `kubectl`, `dtctl` aims to simplify Dependency-Track operations directly from your terminal.

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)

---

## Features

- **Manage Multiple Contexts**: Easily add, switch, and view different Dependency-Track server configurations.
- **Seamless Context Switching**: Quickly switch between different server contexts without reconfiguring each time.
- **Fetch Projects**: Retrieve and display all projects from the current Dependency-Track server.
- **Cross-Platform Support**: Available for Linux, macOS, and Windows, ensuring broad usability.
- **Extensible Architecture**: Designed to allow the addition of more commands and functionalities in the future.

---

## Installation

### Download Pre-built Binaries

Pre-built binaries are available for Linux, macOS, and Windows on the [Releases](https://github.com/yourusername/dtctl/releases) page.

1. **Download the binary** for your operating system from the [Releases](https://github.com/yourusername/dtctl/releases) page.

2. **Make the binary executable** (if required):

    ```bash
    chmod +x dtctl
    ```

3. **Move the binary to a directory in your `PATH`**:

    ```bash
    sudo mv dtctl /usr/local/bin/
    ```

   *Alternatively, you can move it to any directory that's included in your system's `PATH` environment variable.*

### Build from Source

#### Prerequisites

- [Go 1.20+](https://golang.org/dl/)

#### Steps

1. **Clone the repository**:

    ```bash
    git clone https://github.com/yourusername/dtctl.git
    cd dtctl
    ```

2. **Build the binary**:

    ```bash
    go build -o dtctl
    ```

3. **Move the binary to a directory in your `PATH`**:

    ```bash
    sudo mv dtctl /usr/local/bin/
    ```

   *Ensure that `/usr/local/bin/` is in your `PATH`. You can verify this by running:*

    ```bash
    echo $PATH
    ```

   *If it's not included, you can add it by editing your shell profile (e.g., `.bashrc`, `.zshrc`).*

---

## Configuration

Before using `dtctl`, you need to configure it with your Dependency-Track server details. This configuration allows you to manage multiple server contexts and switch between them as needed.

The configuration is stored in a file located at `~/.dtctl/config.json`.

### Adding a Context

To add a new context, use the `add-context` command with a unique name, the Dependency-Track server URL, and your API token.

```bash
dtctl config add-context mycontext --url="https://dependency-track.example.com" --token="your-api-key"
dtctl config add-context production --url="https://dt.example.com" --token="abcd1234efgh5678ijkl"
```

### Switching Contexts
Set the current context to use for operations. This allows you to switch between different Dependency-Track server configurations seamlessly.
```bash
dtctl config use-context production
```

---

## Usage

Once you have configured your contexts, you can use `dtctl` to interact with your Dependency-Track server. Below are the primary commands and their usage.

### Fetching Projects

Retrieve and display all projects from the current context's Dependency-Track server.

```bash
dtctl get projects
dtctl get policies
dtctl get components
```

---

