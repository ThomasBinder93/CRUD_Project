# Pushing to GitHub

This guide helps you push the CRUD project to GitHub.

## Prerequisites

- Git installed ([download](https://git-scm.com/download))
- GitHub account ([create one](https://github.com/join))
- GitHub SSH key configured (optional but recommended)

## Steps

### 1. Create a New Repository on GitHub

1. Go to [github.com/new](https://github.com/new)
2. Repository name: `CRUD_Project` (or your preferred name)
3. Description: "Production-Grade CRUD Backend API with Go"
4. Choose Public or Private
5. **Do NOT** initialize with README, .gitignore, or license (we have them)
6. Click "Create repository"

### 2. Initialize Git Locally

```bash
cd ~/path/to/CRUD_Project

# Initialize Git
git init

# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: Production-grade CRUD API"
```

### 3. Add Remote and Push

Replace `USERNAME` with your GitHub username and `REPOSITORY` with your repo name:

```bash
# Add remote (HTTPS)
git remote add origin https://github.com/ThomasBinder93/CRUD_Project.git

# Or use SSH (if configured)
# git remote add origin git@github.com:USERNAME/REPOSITORY.git

# Rename main branch if needed (Git uses 'master' by default)
git branch -M main

# Push to GitHub
git push -u origin main
```

### 4. Verify

1. Go to your GitHub repository URL
2. You should see all files including:
   - README.md
   - Makefile
   - Docker configuration
   - Go source code
   - GitHub workflows

## Configuring GitHub

### Set up branch protection rules

1. Go to Settings → Branches
2. Add rule for `main` branch
3. Enable:
   - "Require a pull request before merging"
   - "Require status checks to pass before merging"
   - "Require branches to be up to date"

### Configure CI/CD Secrets (Optional)

For Docker Hub integration, add secrets in Settings → Secrets and variables → Actions:

- `DOCKER_USERNAME`: Your Docker Hub username
- `DOCKER_PASSWORD`: Your Docker Hub token

## Future Commits

```bash
# Make changes
git add .
git commit -m "Your commit message"
git push origin main
```

## Creating Releases

```bash
# Tag a release
git tag -a v1.0.0 -m "Version 1.0.0"
git push origin v1.0.0
```

## Additional Resources

- [GitHub Docs: Creating a repository from a template](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template)
- [GitHub Docs: Adding locally hosted code to GitHub](https://docs.github.com/en/migrations/importing-source-code/using-the-command-line-to-import-source-code/adding-locally-hosted-code-to-github)
- [GitHub Docs: Managing remote repositories](https://docs.github.com/en/get-started/getting-started-with-git/managing-remote-repositories)
