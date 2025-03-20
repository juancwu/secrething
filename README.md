# 🍱 Konbini: Secure Secret Management Made Simple

Konbini (Japanese for "convenience store") is your go-to solution for securely storing, managing, and sharing sensitive information within your organization. Like a well-organized bento box, Konbini keeps your secrets neatly compartmentalized and protected.

## ✨ Features

- **End-to-End Encryption**: Secrets are encrypted on the client side - Konbini never sees plaintext data
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
git clone https://github.com/juancwu/konbini.git
cd konbini
```

2. Install dependencies

```bash
go mod download
```

3. Create a `.env` file in the project root with the following variables:

```
PORT=8080
DB_URL=file:konbini.db
JWT_SECRET=your-secret-key
RESEND_API_KEY=your-resend-api-key
APP_URL=http://localhost:8080
```

4. Build the project

```bash
# Build the server
go build -o bin/konbini cmd/server/main.go

# Build the CLI
go build -o bin/konbini-cli cmd/cli/main.go
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
./bin/konbini
```

### Using the CLI

The CLI can be run in two modes:

#### Interactive TUI Mode

```bash
./bin/konbini-cli
```

This launches an interactive terminal user interface where you can:

- Register/login to your account
- Set up 2FA with TOTP
- Manage your bentos (secret containers)
- Create and join groups
- Invite team members

#### Command Mode

```bash
# Login to your account
./bin/konbini-cli login

# Create a new bento
./bin/konbini-cli bento new my-api-keys

# Add a secret to a bento
./bin/konbini-cli bento add my-api-keys AWS_SECRET_KEY=abcdefg

# List all bentos
./bin/konbini-cli bento list

# Share a bento with a group
./bin/konbini-cli group invite DevTeam john@example.com
```

## 🎯 What is a Bento?

In Konbini, a "bento" is a container for your secrets:

- Each bento has a unique name and can contain multiple "ingredients" (key-value pairs)
- Bentos can be shared with other users through groups
- Permissions control who can view or modify each bento
- All bento contents are encrypted on the client side before being sent to the server

## 🛡️ Security

Konbini is designed with security at its core:

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

