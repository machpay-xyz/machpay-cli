# Contributing to MachPay CLI

Thank you for your interest in contributing to MachPay CLI! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful, inclusive, and constructive. We're all here to build something great together.

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git
- Make (optional, for convenience)

### Setup

```bash
# Clone the repository
git clone https://github.com/machpay/machpay-cli.git
cd machpay-cli

# Install dependencies
go mod download

# Build
go build -o machpay ./cmd/machpay

# Run tests
go test ./...
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 2. Make Changes

- Follow the existing code style
- Add tests for new functionality
- Update documentation as needed

### 3. Test

```bash
# Run all tests
go test -v ./...

# Run with race detection
go test -race ./...

# Run with coverage
go test -cover ./...
```

### 4. Commit

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new command
fix: resolve issue with config loading
docs: update README
test: add tests for wallet module
refactor: simplify process manager
```

### 5. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then open a Pull Request on GitHub.

## Project Structure

```
machpay-cli/
├── cmd/machpay/          # Main entry point
├── internal/
│   ├── auth/             # Authentication
│   ├── cmd/              # Cobra commands
│   ├── config/           # Configuration management
│   ├── gateway/          # Gateway download & process
│   ├── tui/              # Terminal UI components
│   └── wallet/           # Wallet/key management
├── scripts/              # Install scripts
└── .github/workflows/    # CI/CD
```

## Testing Guidelines

- Write table-driven tests where appropriate
- Mock external dependencies (HTTP, filesystem)
- Aim for >80% coverage on new code
- Test edge cases and error conditions

Example:

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "foo", "bar", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Something(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("got = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Pull Request Guidelines

- Keep PRs focused and reasonably sized
- Include tests for new functionality
- Update documentation if needed
- Ensure CI passes
- Request review from maintainers

## Release Process

Releases are automated via GoReleaser when a version tag is pushed:

```bash
git tag v1.2.3
git push origin v1.2.3
```

## Questions?

- Open a [Discussion](https://github.com/machpay/machpay-cli/discussions)
- Join our [Discord](https://discord.gg/machpay)
- Email: hello@machpay.xyz

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

