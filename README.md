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

# Set up 

```zsh
    
git clone https://github.com/SwanHtetAungPhyo/grpcframe.git
     
cd grpcframe
    
    # Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o grpcframe .

# Linux (32-bit)
GOOS=linux GOARCH=386 go build -o grpcframe .

# Linux ARM
GOOS=linux GOARCH=arm go build -o grpcframe .
GOOS=linux GOARCH=arm64 go build -o grpcframe .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o grpcframe .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o grpcframe .

# Windows
GOOS=windows GOARCH=amd64 go build -o grpcframe .

sudo mv grpcframe /usr/local/bin

# command completion
grpcframe completion

# use this after set up
grpcframe --help     
```
---

## 📜 License

MIT License — © Swan Htet Aung Phyo