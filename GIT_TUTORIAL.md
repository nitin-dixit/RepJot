# Mastering Professional Git Workflows & Repository History

This comprehensive tutorial documents the precise methodologies, internal mechanics, and commands used to rebuild this repository into a professional, production-grade project. It is designed as an educational asset so you can implement, manage, and automate these workflows on your own.

---

## 1. The Core Philosophy: Milestone-Driven History

In professional software development, git history is **documentation**. A messy commit history with giant diffs or messages like `fix`, `work in progress`, or `changes` makes code reviews difficult, degrades repository maintainability, and hides the architectural evolution of the system.

### Semantic Branch Separation
To maintain code quality, we separate development into logical, isolated blocks called **Milestones**. Each milestone:
1. Solves a specific problem or implements a single cohesive feature.
2. Has its own dedicated feature branch (`feature/<name>`).
3. Is integrated into the main stream (`develop` or `main`) through a formal **Pull Request (PR)**.
4. Preserves compile-ability and passes all tests at the end of each commit.

---

## 2. Deep Dive: Git Internals (Author vs. Committer Dates)

Every commit in Git contains **two** distinct timestamps:
* **Author Date**: When the change was originally written.
* **Committer Date**: When the change was applied/committed to the Git tree.

Under normal circumstances, these are identical. However, when we perform rebasing, cherry-picking, or simulated chronological reconstructions, they can diverge. 

### Overriding Timestamps
Git allows you to explicitly override these dates using environment variables before executing the `git commit` command:

```bash
GIT_AUTHOR_DATE="YYYY-MM-DDTHH:MM:SS" GIT_COMMITTER_DATE="YYYY-MM-DDTHH:MM:SS" git commit -m "feat: description"
```

* Format: ISO 8601 (e.g., `2026-05-12T10:00:00`) or standard RFC 2822.
* To check both timestamps in your repository, run:
  ```bash
  git log --pretty=fuller
  ```

---

## 3. Step-by-Step Repository Professionalization Guide

This is the exact sequence of steps you can follow to rebuild or professionalize any project repository from a completed "flat" working directory.

### Phase 1: Establish a Safe Backup
Always create a complete backup branch containing all your work (including uncommitted and untracked files) before rewriting history:
```bash
# Create and switch to a backup branch
git checkout -b backup-full-workout

# Stage and commit all files (bypassing pre-commit hooks temporarily)
git add .
git commit -m "backup: initial full workout project state" --no-verify
```

### Phase 2: Rebuilding Chronological Milestones
We reset the target branch (`develop` or `main`) to the base state and rebuild features sequentially.

#### Milestone 1: HTTP Scaffolding
Check out only the scaffolding files from your backup (or base branch):
```bash
# Switch to develop and reset to the initial commit
git checkout develop
git reset --hard <initial-commit-hash>

# Create feature branch
git checkout -b feature/bootstrap-server

# Checkout the base server files from the backup branch
git checkout backup-full-workout -- main.go go.mod go.sum internal/app/app.go internal/routes/routes.go internal/api/workout_handler.go

# FIX THE EOF ERROR: In Go, placeholder files must have a package header.
# Write a syntactically correct shell inside internal/store/workout_store.go:
echo "package store" > internal/store/workout_store.go
```
Commit it with your backdated timestamp:
```bash
GIT_AUTHOR_DATE="2026-05-12T10:00:00" GIT_COMMITTER_DATE="2026-05-12T10:00:00" git commit -am "feat: initialize the project with chi router and basic handlers"
```

#### Milestone 2: Database Setup & Migrations
```bash
# Create feature branch off develop
git checkout develop
git checkout -b feature/database-setup

# Check out database files and migrations
git checkout backup-full-workout -- docker-compose.yml internal/store/database.go migrations/

# Run go get to populate dependencies in go.mod
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/pressly/goose/v3

# Commit with backdated timestamp
GIT_AUTHOR_DATE="2026-05-14T11:00:00" GIT_COMMITTER_DATE="2026-05-14T11:00:00" git commit -am "feat: setup PostgreSQL with goose migrations and docker-compose"
```

---

## 4. GitHub CLI (`gh`) Automation & Local Merge Integration

To simulate a real-world team workflow, we create Pull Requests on GitHub for each milestone and merge them. To keep merge timestamps matching our backdated commit dates, we use a hybrid **CLI-PR + Local Merge** workflow.

### Why not merge on GitHub directly?
If you merge the PR using the GitHub Web UI or `gh pr merge` on GitHub's servers, GitHub will stamp the merge commit with the *current real-time date*. Doing the merge **locally** allows us to override the merge commit's date and then push the merged history!

### The Automated Script Workflow
1. **Push Branch**: Push the local backdated feature branch:
   ```bash
   git push origin feature/database-setup --force
   ```
2. **Create Pull Request**: Create the PR programmatically:
   ```bash
   gh pr create --title "feat: database migrations & setup" --body "Sets up standard PostgreSQL connection pool with goose embedding." --head feature/database-setup --base develop
   ```
3. **Merge Locally**: Switch to `develop` and perform a non-fast-forward merge with backdated dates:
   ```bash
   git checkout develop
   GIT_AUTHOR_DATE="2026-05-15T13:00:00" GIT_COMMITTER_DATE="2026-05-15T13:00:00" git merge feature/database-setup --no-ff -m "Merge pull request #2 from feature/database-setup"
   ```
4. **Push develop**: Pushing `develop` to GitHub will cause GitHub to automatically match the commits and mark the PR as **Merged**!
   ```bash
   git push origin develop
   ```

---

## 5. Gatekeeping: Husky & Commitlint Integration

A professional project uses git hooks to prevent syntax errors, failing tests, or sloppy commit messages from ever reaching the repository.

### What is Husky & Commitlint?
* **Husky**: A lightweight wrapper around native Git hooks (like `pre-commit` or `commit-msg`).
* **Commitlint**: A linter that checks if your commit messages conform to the **Conventional Commits** standard (e.g., `feat: ...`, `fix: ...`, `chore: ...`).

### Configuration & Setup
1. **Initialize package.json**:
   ```json
   {
     "name": "goproject",
     "version": "1.0.0",
     "scripts": {
       "test": "go test ./..."
     },
     "devDependencies": {
       "@commitlint/cli": "^19.3.0",
       "@commitlint/config-conventional": "^19.2.2",
       "husky": "^9.0.11"
     }
   }
   ```
2. **Commitlint Config (`commitlint.config.js`)**:
   ```javascript
   module.exports = { extends: ['@commitlint/config-conventional'] };
   ```
3. **Husky Pre-Commit Hook (`.husky/pre-commit`)**:
   Runs your Go tests before allowing a commit:
   ```bash
   #!/bin/sh
   . "$(dirname "$0")/_/husky.sh"
   
   npm test
   ```
4. **Husky Commit-Msg Hook (`.husky/commit-msg`)**:
   Lints your commit messages:
   ```bash
   #!/bin/sh
   . "$(dirname "$0")/_/husky.sh"
   
   npx --no -- commitlint --edit "$1"
   ```

---

## 6. Resolving Integration Errors

### Satisfying Foreign Key Constraints in Unit Tests
If a table has a foreign key constraint (e.g., `workouts.user_id` referencing `users.id` with `NOT NULL`), writing a unit test that creates a workout directly will fail if the database starts empty or truncated:
```
ERROR: insert or update on table "workouts" violates foreign key constraint (SQLSTATE 23503)
```

**Solution**:
1. Your test setup function must first create and insert a dummy user record.
2. Retrieve the inserted user's ID.
3. Assign that ID to the workout structure being tested:
   ```go
   var userID int
   err = db.QueryRow(`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`, "test_user", "test@example.com", "hash").Scan(&userID)
   // ...
   workout := &Workout{
       UserID: userID,
       Title: "pushup day",
       // ...
   }
   ```

---

## Summary Command Cheat Sheet

Use these quick commands in your daily workflows:

| Action | Command |
| :--- | :--- |
| **Stage & Commit** | `git add . && git commit -m "feat: description"` |
| **Commit (bypass hooks)** | `git commit -m "chore: bypass" --no-verify` |
| **Backdated Commit** | `GIT_AUTHOR_DATE="YYYY-MM-DDTHH:MM:SS" GIT_COMMITTER_DATE="YYYY-MM-DDTHH:MM:SS" git commit -m "feat: description"` |
| **Audit Logs** | `git log --pretty=fuller` |
| **Simulated PR Merge** | `GIT_AUTHOR_DATE="..." GIT_COMMITTER_DATE="..." git merge branch --no-ff -m "Merge pull request #X from branch"` |
| **Tag Release** | `git tag -a v1.0.0 -m "Release v1.0.0" && git push origin v1.0.0` |

*Happy Coding!*
