#  Grpcframe

<div align="center">
<img src="https://img.shields.io/badge/Go-1.24-00ADD8?style=for-the-badge&logo=go&logoColor=white"/>
</div>

**grpcframe** is a robust CLI tool for scaffolding and managing gRPC-based Go applications. It utilizes the Cobra CLI framework to organize and streamline operations such as module creation, protobuf generation, and database migrations.

---

## 🚀 Technologies Used

- [Cobra](https://github.com/spf13/cobra) — CLI framework for Go
- [gRPC](https://grpc.io/) — High-performance RPC framework
- [Protocol Buffers (Protobuf)](https://developers.google.com/protocol-buffers) — Interface definition language
- [Buf](https://buf.build/) — Tooling for Protobuf
- [sqlc](https://sqlc.dev/) — Generate type-safe Go from SQL
- [golang-migrate](https://github.com/golang-migrate/migrate) — Database migrations

---

## 📦 Usage

```bash
grpcframe [command]
```

## 📚 Available Commands

### 🔧 Project Initialization

- `init [dir] [module-name]`  
  Initializes a new project in the specified directory.

### 🧬 Module Management

- `module`  
  Entry point for module-related commands.

- `module add [module-name] [target-module]`  
  Adds a new gRPC module with handlers.

### 📄 Protobuf Generation

- `protogen`  
  Generates Go code from Protobuf definitions using `buf`.

### 🛠 SQLc Generation

- `sqlc`  
  Runs SQLc code generation.

### 🗃 Database Migration

- `migrate`  
  Root command for migrations.

- `migrate up`  
  Applies all pending migrations.

- `migrate down [number]`  
  Rolls back the given number of migrations (default: 1).

- `migrate force [version]`  
  Forces database to a specific version.

- `migrate version`  
  Displays the current database migration version.

### 🌐 Gateway Registration

- `gateway`  
  Registers all gRPC gateway endpoints automatically.

### 🧾 Help & Autocompletion

- `help`  
  Displays help for a specific command.

- `completion`  
  Generates autocompletion scripts for various shells.

---

## 🆘 Flags

```bash
  -h, --help     help for grpcframe
  -t, --toggle   Help message for toggle
```

---

## 💡 Example

```bash
grpcframe init myproject github.com/swan/grpcframe
grpcframe module add user service
grpcframe migrate up
grpcframe protogen
grpcframe sqlc
```
Here's how to document the download instructions for users in your `README.md`, with clear options for different platforms:

## 📥 Download and Installation

### Option 1: Download Pre-built Binaries (Recommended)

1. Go to the [Latest Release](https://github.com/SwanHtetAungPhyo/grpcframe/releases/latest) page
2. Download the appropriate binary for your system:

  - **Windows**:
    - 64-bit: `grpcframe-windows-amd64.exe`
    - 32-bit: `grpcframe-windows-386.exe`

  - **macOS**:
    - Intel: `grpcframe-darwin-amd64`
    - Apple Silicon: `grpcframe-darwin-arm64`

  - **Linux**:
    - 64-bit: `grpcframe-linux-amd64`
    - 32-bit: `grpcframe-linux-386`
    - ARM: `grpcframe-linux-arm`
    - ARM64: `grpcframe-linux-arm64`

3. Make the binary executable (Linux/macOS):
   ```bash
   chmod +x grpcframe-*
   ```
4. Move to your PATH (optional but recommended):
   ```bash
   sudo mv grpcframe-* /usr/local/bin/grpcframe
   ```

### Option 2: Install via cURL (Linux/macOS)

```bash
# Download and install directly
curl -L https://github.com/SwanHtetAungPhyo/grpcframe/releases/latest/download/grpcframe-$(uname -s)-$(uname -m) -o grpcframe
chmod +x grpcframe
sudo mv grpcframe /usr/local/bin/
```

### Option 3: Build from Source

```bash
git clone https://github.com/SwanHtetAungPhyo/grpcframe.git
cd grpcframe
go build -o grpcframe .
```

### Verification (Optional)

Verify the checksum matches what's in the release's `checksums.txt` file:

```bash
sha256sum grpcframe-*  # Linux/macOS
certutil -hashfile grpcframe.exe SHA256  # Windows
```

### Shell Completions

To enable tab completion:

```bash
# Bash
grpcframe completion bash > /etc/bash_completion.d/grpcframe

# Zsh
grpcframe completion zsh > "${fpath[1]}/_grpcframe"

# Fish
grpcframe completion fish > ~/.config/fish/completions/grpcframe.fish
```

---
## 📜 License

MIT License — © Swan Htet Aung Phyo