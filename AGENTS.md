# BIT303_shop Project Instructions

## Development Scope

- This project is developed from 0 to 1.
- Keep the current stage focused on the feature being requested.
- Do not add unrelated mall features early.
- For each backend feature, follow the user's GoFrame workflow:
  1. Design table structure and update initialization SQL.
  2. Run `gf gen dao` to generate dao/do/entity when tables are ready.
  3. Write API request/response structs.
  4. Write internal model input/output structs.
  5. Write logic and generate service code.
  6. Register the service in the corresponding logic module.
  7. Write controller code.
  8. Register routes in `internal/cmd/cmd.go`.
  9. Verify the feature.

## Git Workflow For Each Feature

- Before developing a feature, create a new branch from `main`.
- Default branch naming:
  - `feature/employee-login`
  - `feature/points-management`
  - `feature/goods-management`
- Do not develop feature work directly on `main`.
- Recommended start sequence:

```bash
git checkout main
git pull
git checkout -b feature/<feature-name>
```

## Commit And Review Gate

- After implementing and verifying a feature, it is allowed to run:

```bash
git add .
git commit -m "<feature summary>"
```

- After the local commit, stop and let the user review before any of these actions:
  - Push the feature branch to remote.
  - Merge or sync the feature branch into `main`.
  - Push `main` to remote.

## Prohibited Without Explicit User Confirmation

- Do not push any branch to remote.
- Do not merge into `main`.
- Do not push `main`.
- Do not use force push.
- Do not rewrite Git history.
- Do not run destructive Git commands such as `git reset --hard` unless the user explicitly requests them.

## Push Failure Handling

- If remote push fails because of network, credential, or GitHub authentication issues, keep the local commit and tell the user the exact manual push command to run.
