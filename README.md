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

You can download the latest version of `dtctl` from the [Releases](https://github.com/francislance/dtctl/releases/latest) page.

### macOS and Linux

1. **Download and unzip the binary**:

   #### For macOS:
   ```bash
   curl -L -o dtctl.zip https://github.com/francislance/dtctl/releases/latest/download/dtctl-macos-amd64.zip
   ```
   #### For Linux:
   ```bash
   curl -L -o dtctl.zip https://github.com/francislance/dtctl/releases/latest/download/dtctl-linux-amd64.zip
   ```

2. **Download and unzip the binary**:

   ```bash
   unzip dtctl.zip
   ```

3. Make the binary executable:

   ```bash 
   chmod +x dtctl
   ```

4. Move the binary to a directory in your PATH:

   ```bash
   sudo mv dtctl /usr/local/bin/
   ```

5. Verify the installation:

   ```bash
   dtctl --version
   ```

### Windows
1. Download and unzip the binary:
   - Download dtctl-vX.X.X-windows.zip from the Releases page. 
   - Unzip the file to extract dtctl.exe. 
2. Add the binary to your PATH:
   - Move dtctl.exe to a directory that's in your PATH, or add the directory containing dtctl.exe to your PATH environment variable.
3. Verify the installation:

   ```cmd
   dtctl --version
   ```

---

## Configuration

Before using `dtctl`, you need to configure it with your Dependency-Track server details. This configuration allows you to manage multiple server contexts and switch between them as needed.

The configuration is stored in a file located at `~/.dtctl/config.json`.

### Adding a Context

To add a new context, use the `add-context` command with a unique name, the Dependency-Track server URL, and your API token.

```bash
dtctl config add-context mycontext --url="https://dependency-track.example.com" --token="your-api-key"
```

### Switching Contexts
Set the current context to use for operations. This allows you to switch between different Dependency-Track server configurations seamlessly.
```bash
dtctl config use-context production
```

---

## Usage

Once you have configured your contexts, you can use `dtctl` to interact with your Dependency-Track server. Below are the primary commands and their usage.

### Projects

Retrieve and display all projects from the current context's Dependency-Track server.

```bash
dtctl get projects
```

```bash
dtctl get projects --tag="springboot"
```

### Policies
```bash
dtctl get policies
```

### Components
```bash
# get all components
dtctl get components
```
```bash
# get all components under a project with specific tag
dtctl get components --tag="container"
```
```bash
# get all components under a project with fields
# (available: projectname, projectuuid, sha256, sha1, md5)
dtctl get components --show-fields="projectname,projectuuid,sha256,sha1,md5" --tag="container"
```

### Hash Policy Condition

Sample updating of hash policy condition:
```bash
dtctl set hashpolicycondition --uuid="1cf6c518-149a-43a6-991d-276d163c5852" --operator="IS_NOT" --subject="COMPONENT_HASH" --algorithm="SHA-256" --algorithm-value="928b2691494882b361bbe4f70fcf3fa9fbcb5a2bbe88f2b42f7e93f2c8cc726b"
```
---

### Evaluate a Policy

```bash
# set the uuid of the policy
dtctl eval policy --uuid="c4583613-1e43-4346-ac2d-db3d4e19491a"
```