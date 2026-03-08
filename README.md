# GitHub Issues Manager

Interactive Go CLI to create and search GitHub issues.

## Project Structure

- `cmd/issues`: CLI entrypoint
- `internal/app`: interactive flow and action routing
- `internal/cli`: action parsing and menu/usage rendering
- `internal/config`: environment configuration helpers
- `internal/githubapi`: GitHub REST API client and models

## Required Environment Variables

- `GITHUB_API_TOKEN`
- `GITHUB_OWNER`
- `GITHUB_REPO`

Use `.env.example` as a template.

## Run

```bash
go run ./cmd/issues
```
