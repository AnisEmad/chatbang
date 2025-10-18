# GitHub Actions Workflows Documentation

This project uses two GitHub Actions workflows to automate continuous integration, testing, and release processes.

## Table of Contents

- [CICD Workflow](#cicd-workflow)
- [Release Workflow](#release-workflow)
- [Setup Instructions](#setup-instructions)
- [Triggering Workflows](#triggering-workflows)

---

## CICD Workflow

**File:** `.github/workflows/main.yml`

### Purpose

The CICD workflow automatically builds and tests the application whenever code is pushed to the main branch or when a pull request targets the main branch. This ensures code quality and catches issues early in the development process.

### Trigger Events

This workflow runs on:
- **Push events** to the `main` branch
- **Pull request events** targeting the `main` branch

### Jobs

#### `build`

Runs on: `ubuntu-latest`

**Steps:**

1. **Checkout Code**
   - Action: `actions/checkout@v5.0.0`
   - Clones the repository code to the runner

2. **Setup Go Environment**
   - Action: `actions/setup-go@v6.0.0`
   - Installs and configures the Go programming language environment

3. **Install Dependencies**
   - Command: `go mod tidy`
   - Downloads and organizes all required Go modules

4. **Build**
   - Command: `go build main.go`
   - Compiles the application to verify it builds successfully

5. **Test**
   - Command: `go test -v`
   - Runs all tests with verbose output to ensure code correctness

### Example Usage

The workflow triggers automatically when you:

```bash
# Push to main branch
git push origin main

# Or create a pull request to main
gh pr create --base main
```

---

## Release Workflow

**File:** `.github/workflows/release.yml`

### Purpose

The Release workflow automates the creation of GitHub releases with compiled binaries. When you tag a version, it builds the application, runs tests, and publishes a release with the binary artifact.

### Trigger Events

This workflow runs on:
- **Push events** with tags matching the pattern `v*` (e.g., `v1.0.0`, `v2.1.3`)

### Permissions

- `contents: write` - Required to create releases and upload assets

### Jobs

#### `build-and-release`

Runs on: `ubuntu-latest`

**Steps:**

1. **Checkout Code**
   - Action: `actions/checkout@v5.0.0`
   - Clones the repository code to the runner

2. **Setup Go**
   - Action: `actions/setup-go@v6.0.0`
   - Installs and configures the Go programming language environment

3. **Install Dependencies**
   - Command: `go mod tidy`
   - Downloads and organizes all required Go modules

4. **Run Tests**
   - Command: `go test -v`
   - Executes all tests to ensure the release is stable

5. **Build Binary**
   - Creates a `dist` directory
   - Compiles the application into a binary named `chatbang`
   - Output: `dist/chatbang`

6. **Create GitHub Release**
   - Action: `softprops/action-gh-release@v2`
   - Creates a new GitHub release for the tag
   - Uploads the `chatbang` binary as a release asset
   - Uses the repository's `GITHUB_TOKEN` for authentication

### Example Usage

To trigger a release:

```bash
# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0

# Or use GitHub CLI
gh release create v1.0.0
```

---

## Setup Instructions

### Prerequisites

1. Both workflow files should be placed in `.github/workflows/` directory:
   ```
   .github/
   └── workflows/
       ├── main.yml
       └── release.yml
   ```

2. Ensure your repository has:
   - A `main.go` file in the root directory
   - A `go.mod` file for dependency management
   - Go test files (optional but recommended)

### Repository Configuration

No additional secrets are required. The workflows use the automatically provided `GITHUB_TOKEN` which has sufficient permissions for these operations.

---

## Triggering Workflows

### CICD Workflow

**Automatically triggers when:**
- You push commits to the `main` branch
- Someone opens a pull request targeting `main`
- Someone pushes updates to an existing PR targeting `main`

### Release Workflow

**Manually trigger a release:**

```bash
# Create a new version tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# Push the tag to GitHub
git push origin v1.0.0
```

After pushing the tag, the workflow will:
1. Run all tests
2. Build the binary
3. Create a GitHub release
4. Attach the `chatbang` binary to the release

---

## Workflow Status

You can monitor workflow runs in your repository:
- Navigate to the **Actions** tab in your GitHub repository
- View running, completed, or failed workflows
- Click on individual runs to see detailed logs

---

## Troubleshooting

### CICD Workflow Fails

- **Build fails:** Check that all dependencies are listed in `go.mod`
- **Tests fail:** Review test output in the workflow logs

### Release Workflow Fails

- **Binary not created:** Ensure `main.go` exists and compiles successfully
- **Release creation fails:** Verify the tag follows the `v*` pattern (e.g., `v1.0.0`)
- **Permission denied:** Ensure repository settings allow Actions to create releases

---

## Best Practices

1. **Always test before tagging:** The CICD workflow will catch issues on pull requests
2. **Use semantic versioning:** Follow the `vMAJOR.MINOR.PATCH` format (e.g., `v1.2.3`)
3. **Write descriptive commit messages:** They appear in release notes by default
4. **Keep workflows updated:** Periodically update action versions for security patches
