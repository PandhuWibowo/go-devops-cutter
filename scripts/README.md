# Git Hooks for go-devops-cutter

This directory contains git hooks that automatically run quality checks before commits and pushes.

## Quick Start

Install the hooks by running:

```bash
./scripts/install-hooks.sh
```

## Available Hooks

### Pre-commit Hook

Runs **before** every `git commit` to ensure code quality.

**Checks performed:**
1. ✅ Code formatting (`gofmt`)
2. ✅ Static analysis (`go vet`)
3. ✅ Dependencies are tidy (`go.mod`/`go.sum`)
4. ✅ All unit tests pass
5. ✅ Project builds successfully

**Example output:**
```
🔍 Running pre-commit checks...
📝 Checking code formatting...
✓ Code formatting passed
🔍 Running go vet...
✓ go vet passed
📦 Checking go.mod and go.sum...
✓ go.mod and go.sum are tidy
🧪 Running unit tests...
✓ All tests passed
🔨 Building project...
✓ Build successful
✅ All pre-commit checks passed!
```

**If any check fails, the commit is blocked.**

### Pre-push Hook

Runs **before** every `git push` to ensure code quality before sharing.

**Checks performed:**
1. ✅ All unit tests pass with coverage
2. ✅ Coverage report generated
3. ✅ Project builds successfully
4. ⚠️  Warns about TODO/FIXME comments
5. ⚠️  Warns about debug statements
6. ⚠️  Warns if coverage < 40%
7. ⚠️  Confirms push to main/master branches

**Example output:**
```
🚀 Running pre-push checks...
Branch: feat/my-feature
Remote: origin

🧪 Running unit tests with coverage...
✓ All tests passed
📊 Checking test coverage...
Total coverage: 40.8%
🔨 Building all binaries...
✓ Build successful
🔍 Checking for common issues...
✅ All pre-push checks passed! Proceeding with push...
```

## Skipping Hooks

Sometimes you need to skip hooks (e.g., WIP commits). Use `--no-verify`:

```bash
# Skip pre-commit hook
git commit --no-verify -m "WIP: work in progress"

# Skip pre-push hook
git push --no-verify
```

**⚠️ Use sparingly!** Skipping hooks means bypassing quality checks.

## Uninstalling Hooks

To remove the hooks:

```bash
rm .git/hooks/pre-commit
rm .git/hooks/pre-push
```

Or simply delete them from your `.git/hooks/` directory.

## Customizing Hooks

The hook source files are in `scripts/hooks/`. To modify:

1. Edit the hook file in `scripts/hooks/`
2. Run `./scripts/install-hooks.sh` to reinstall
3. The updated hook will be copied to `.git/hooks/`

**Note:** `.git/hooks/` is not tracked by git, so each developer must install hooks locally.

## Troubleshooting

### Hook not running?

Check if the hook is executable:
```bash
ls -la .git/hooks/pre-commit
ls -la .git/hooks/pre-push
```

If not, make them executable:
```bash
chmod +x .git/hooks/pre-commit
chmod +x .git/hooks/pre-push
```

### Hook failing incorrectly?

1. Try running the checks manually:
   ```bash
   make test
   make build
   gofmt -s -l .
   go vet ./...
   go mod tidy
   ```

2. Check the hook script output for specific errors

3. If needed, skip with `--no-verify` and investigate

### Hook not finding commands?

Ensure these are installed and in your PATH:
- `go` (version 1.24+)
- `make`
- `bc` (for coverage calculations)

## Best Practices

### Do:
- ✅ Install hooks immediately after cloning
- ✅ Keep hooks updated (`git pull` + `./scripts/install-hooks.sh`)
- ✅ Fix issues highlighted by hooks
- ✅ Add tests before committing new code

### Don't:
- ❌ Habitually skip hooks with `--no-verify`
- ❌ Commit failing tests
- ❌ Push untested code
- ❌ Ignore warnings about coverage or debug statements

## Performance

**Pre-commit:** ~5-15 seconds (depending on test suite size)
**Pre-push:** ~10-20 seconds (includes coverage generation)

If hooks become too slow:
- Consider running only fast unit tests in pre-commit
- Move slower integration tests to pre-push
- Use `--no-verify` for WIP branches, but run hooks before PR

## Integration with CI/CD

These hooks mirror the checks run by GitHub Actions:
- Pre-commit ≈ GitHub Actions lint + test jobs
- Pre-push ≈ GitHub Actions full CI pipeline

Running hooks locally **reduces CI failures** and saves time.
