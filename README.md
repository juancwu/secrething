# 🔐 Secrething: Secure Secret Management Made Simple

Secrething is your go-to solution for securely storing, managing, and sharing sensitive information within your organization. It keeps your secrets neatly compartmentalized and protected.

## ✨ Features

- **End-to-End Encryption**: Secrets are encrypted on the client side - Secrething never sees plaintext data
- **Team Sharing**: Securely share credentials with team members through the groups system
- **Fine-grained Permissions**: Control who can access, view, and modify your secrets
- **Two-Factor Authentication**: Enhanced security with TOTP (Time-based One-Time Password)
- **Intuitive CLI**: Command-line interface with TUI support for easy management
- **API Access**: RESTful API for integration with your existing tools
- **Audit Logs**: Track who accessed what and when

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- SQLite database (or Turso for production)
- Resend.com account (for email verification)
- [Goose](https://github.com/pressly/goose) for database migrations

### Installation

1. Clone the repository

```bash
git clone https://github.com/juancwu/secrething.git
cd secrething
```

2. Install dependencies

```bash
go mod download
```

3. Create a `.env` file in the project root with the following variables:

```
PORT=8080
DB_URL=file:secrething.db
JWT_SECRET=your-secret-key
RESEND_API_KEY=your-resend-api-key
APP_URL=http://localhost:8080
```

4. Build the project

```bash
# Build the server
go build -o bin/secrething cmd/server/main.go

# Build the CLI
go build -o bin/secrething-cli cmd/cli/main.go
```

## 📚 Usage

### Development Scripts

The project includes several utility scripts to make development easier:

#### Setup Script

```bash
./setup.sh
```

Sets up the local development environment by:
- Creating required directories
- Installing dependencies (Go, Air, Goose, sqlc, Turso CLI)
- Downloading Go dependencies
- Generating SQL code
- Setting up a local database
- Running database migrations

#### Development Server

```bash
./dev.sh
```

Runs the development server with hot reload using Air.

#### CLI Shortcut

```bash
./cli.sh
```

Runs the CLI application in interactive TUI mode.

#### Database Migrations

```bash
./migrate.sh [command]
```

Database migration utility with the following commands:
- `up`: Apply all migrations
- `up-by-one`: Apply next migration
- `down`: Roll back all migrations
- `down-by-one`: Roll back most recent migration
- `status`: Show migration status
- `create [name]`: Create a new migration

#### Run Tests

```bash
./run_tests.sh
```

Runs tests with coverage by:
- Starting a local Turso database instance
- Running migrations
- Executing tests and generating coverage report
- Terminating the database instance

### Start the Server

```bash
./bin/secrething
```

### Using the CLI

The CLI can be run in two modes:

#### Interactive TUI Mode

```bash
./bin/secrething-cli
```

This launches an interactive terminal user interface where you can:

- Register/login to your account
- Set up 2FA with TOTP
- Manage your secrets (secret containers)
- Create and join groups
- Invite team members

#### Command Mode

```bash
# Login to your account
./bin/secrething-cli login

# Create a new secret container
./bin/secrething-cli secret new my-api-keys

# Add a secret to a container
./bin/secrething-cli secret add my-api-keys AWS_SECRET_KEY=abcdefg

# List all secret containers
./bin/secrething-cli secret list

# Share a secret with a group
./bin/secrething-cli group invite DevTeam john@example.com
```

## 🛡️ Security

Secrething is designed with security at its core:

- Client-side encryption ensures your secrets never leave your machine in plaintext
- Two-factor authentication (TOTP) protects your account
- Fine-grained permission system prevents unauthorized access
- No plaintext storage of sensitive data
- Email verification for new accounts

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request