## Overview

Do you primarily develop in a single top-level directory like `/dev`, while juggling multiple languages and constantly configuring project and workspace settings in VS Code? If so, this CLI is for you. It serves as a VS Code workspace and project scaffolding tool, helping streamline the setup of workspaces and project configurations across different languages. I built this tool to solve my own frustrations of having to frequently reconfigure VS Code extensions and workspace settings. The goal is to make it easier to create isolated projects where workspace and project settings may differ, avoiding repetitive configuration tasks and allowing you to dive into coding faster.

## Features

- Create and configure multiple project types under a unified VS Code workspace.
- Seamlessly switch between project types with pre-configured settings.
- Cross-platform support for Windows, macOS, and Linux.

## Installation

You can download the entire latest release folder, which includes the binary and the `global-config`. The following instructions assume where your `/bin` folder might be, but you can customize this as needed.

### Quick Install

#### For macOS

```bash
curl -L https://github.com/zackiles/unified-workspace-generator/releases/latest/download/unified-workspace-generator.zip -o unified-workspace-generator.zip
unzip unified-workspace-generator.zip
chmod +x unified-workspace-generator/create-project
sudo mv unified-workspace-generator/create-project /usr/local/bin/create-project
```

#### For Windows (Using PowerShell)

```powershell
Invoke-WebRequest -Uri https://github.com/zackiles/unified-workspace-generator/releases/latest/download/unified-workspace-generator.zip -OutFile unified-workspace-generator.zip
Expand-Archive -Path unified-workspace-generator.zip -DestinationPath . -Force
Move-Item -Path unified-workspace-generator\create-project.exe -Destination C:\Windows\System32\
```

#### For Linux

```bash
curl -L https://github.com/zackiles/unified-workspace-generator/releases/latest/download/unified-workspace-generator.zip -o unified-workspace-generator.zip
unzip unified-workspace-generator.zip
chmod +x unified-workspace-generator/create-project
sudo mv unified-workspace-generator/create-project /usr/local/bin/create-project
```

### Adding `create-project` to PATH for Cross-Platform Usage

To ensure that `create-project` can be run from anywhere in your terminal, you need to add the installation path to your system's PATH environment variable. Here are the instructions for macOS, Windows, and Linux.

#### For macOS and Linux

1. Open your terminal and edit your shell profile (e.g., `~/.bashrc`, `~/.zshrc`, or `~/.bash_profile`).

   ```bash
   nano ~/.bashrc  # or ~/.zshrc depending on your shell
   ```

2. Add the following line to include the directory where `create-project` is installed in your PATH:

   ```bash
   export PATH=$PATH:/usr/local/bin
   ```

3. Save the file and then source it to apply the changes:

   ```bash
   source ~/.bashrc  # or ~/.zshrc or ~/.bash_profile
   ```

4. Verify that the binary is accessible from any directory:

   ```bash
   create-project --help
   ```

#### For Windows (Using PowerShell)

1. Open PowerShell as Administrator.
2. Use the following command to add `C:\Windows\System32` to the system's PATH if it's not already there (in case the binary was moved to this directory):

   ```powershell
   [System.Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Windows\System32", [System.EnvironmentVariableTarget]::Machine)
   ```

3. Restart PowerShell or your Command Prompt for the changes to take effect.

4. Verify that the binary is accessible from any directory:

   ```powershell
   create-project --help
   ```

## Prerequisites

This project uses [Mage](https://github.com/magefile/mage) for task automation. Please make sure you have it installed globally before running any Mage tasks.

To install Mage globally, run:

```bash
go install github.com/magefile/mage@latest
```
