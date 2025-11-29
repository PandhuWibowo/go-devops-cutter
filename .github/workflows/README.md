# GitHub Actions Workflows

This directory contains the CI/CD workflows for go-devops-cutter.

## Workflows

### 1. PR Tests (`pr-tests.yml`)

**Triggers:**
- When a pull request is opened
- When commits are pushed to an open PR
- When a closed PR is reopened

**Target Branches:**
- `main`
- `develop`

**Jobs:**

#### Test Job
- Runs all unit tests
- Generates test coverage report
- Displays coverage summary in PR
- Warns if coverage is below 50%
- Runs on Go 1.24

#### Build Job
- Builds CLI binary (`cutter`)
- Builds API binary (`devops-cutter-api`)
- Verifies binaries work
- Uploads artifacts for 7 days

#### Lint Job
- Checks code formatting with `go fmt`
- Runs `go vet` for static analysis
- Verifies `go.mod` and `go.sum` are tidy

**Usage:**
This workflow runs automatically on every PR. No manual intervention needed.

### 2. CI (`ci.yml`)

**Triggers:**
- Push to `main` branch
- Push to `develop` branch
- Push of version tags (e.g., `v1.0.0`)

**Jobs:**

#### Test Job
- Runs all unit tests
- Generates and displays coverage report

#### Build Job
- Depends on test job passing
- Builds all binaries (CLI + API)
- Verifies binaries work
- Uploads artifacts for 30 days

**Usage:**
This workflow runs automatically on pushes to main/develop branches and version tags.

## Requirements

All workflows require:
- Go 1.24
- Make
- Tests to pass
- Code to build successfully

## Coverage Threshold

The PR test workflow warns if coverage drops below 50%. This is a soft warning and won't fail the build, but reviewers should pay attention.

## Artifacts

Built binaries are uploaded as artifacts:
- PR builds: 7 days retention
- Main/develop builds: 30 days retention

Download artifacts from the Actions tab in GitHub.

## Local Testing

Before pushing, run these commands locally to ensure CI will pass:

```bash
# Run tests
make test

# Check coverage
make test-coverage-func

# Check formatting
gofmt -s -l .

# Run vet
go vet ./...

# Verify go.mod is tidy
go mod tidy

# Build everything
make build
```

## Troubleshooting

**Tests failing in CI but passing locally?**
- Ensure you've committed all files
- Check if you're on the correct Go version (1.24)
- Run `go mod tidy` and commit changes

**Lint job failing?**
- Run `gofmt -s -w .` to auto-format
- Run `go mod tidy` to clean dependencies
- Commit the changes

**Build job failing?**
- Ensure `make build` works locally
- Check for missing dependencies
