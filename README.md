# taskmanager

A tiny Task Manager backend written in Go, structured with a Clean Architecture-ish layering.

- Task fields: `priority`, `assignedUserId`, `status`, `description`
- Operations: create / update / delete / get / list-by-user
- Storage: in-memory repository (no DB yet)

## Overview

This app exposes a JSON HTTP API.

**Layering (high level):**

- **`internal/entity`**: core entities + domain validation
- **`internal/usecase`**: application logic (calls repository interface only)
- **`internal/repository`**: output ports (interfaces) used by usecases
- **`internal/infra`**: concrete implementations (currently `memory`)
- **`internal/framework/http`**: HTTP server adapter (JSON <-> usecase)
- **`cmd/taskmanager`**: composition root (DI/wiring + server boot)

## Directory structure

```text
taskmanager/
  cmd/
    taskmanager/
      main.go                    # server entrypoint (wiring + ListenAndServe)

  internal/
    entity/
      task.go                    # Task entity + Validate()

    repository/
      task_repository.go         # TaskRepository interface
      errors.go                  # repository.ErrNotFound (infra-side)

    usecase/
      task_usecase.go            # Task usecases (Create/Update/Delete/Get/ListByUser)
      errors.go                  # usecase.ErrNotFound (framework-facing)
      task_usecase_test.go       # minimal unit test

    infra/
      memory/
        task_repository.go       # in-memory TaskRepository implementation

    framework/
      http/
        handler/
          task_handler.go        # HTTP handlers
        router/
          router.go              # chi router + route registration

  go.mod
  go.sum
  README.md
```

## How to run

### Requirements

- Go **1.20+** (this repo pins `chi` to a Go 1.20 compatible version)

### Start server

```zsh
cd taskmanager

go run ./cmd/taskmanager
```

Environment variables:

- `TASKMANAGER_ADDR` (default: `:8080`)

## HTTP API

Base URL: `http://localhost:8080`

### Common notes

- All request/response bodies are JSON.
- Error response format:

```json
{ "error": "...message..." }
```

- `usecase.ErrNotFound` is mapped to HTTP **404**.
- Other errors are mapped to HTTP **400** (for now).

### Health

#### `GET /healthz`

**200**

Body (plain text):

```text
ok
```

### Tasks

#### Task JSON shape

```json
{
  "id": "<uuid>",
  "priority": 1,
  "assignedUserId": "user-1",
  "status": "todo",
  "description": "write README",
  "createdAt": "2026-03-15T00:00:00Z",
  "updatedAt": "2026-03-15T00:00:00Z"
}
```

- `status` is one of: `todo`, `in_progress`, `done`
- `createdAt/updatedAt` are RFC3339 strings

#### `POST /tasks`

Create a task.

Request body:

```json
{
  "priority": 1,
  "assignedUserId": "user-1",
  "status": "todo",
  "description": "buy milk"
}
```

Responses:

- **201**: created task

Curl example:

```zsh
curl -s -X POST http://localhost:8080/tasks \
  -H 'Content-Type: application/json' \
  -d '{"priority":1,"assignedUserId":"user-1","status":"todo","description":"buy milk"}' | jq
```

Example output:

```json
{
  "id": "2b0b2f0f-1d3b-4d6e-9f55-2b2b4a5a7890",
  "priority": 1,
  "assignedUserId": "user-1",
  "status": "todo",
  "description": "buy milk",
  "createdAt": "2026-03-15T12:34:56Z",
  "updatedAt": "2026-03-15T12:34:56Z"
}
```

#### `GET /tasks/{id}`

Get a single task.

Responses:

- **200**: task
- **404**: not found

Curl:

```zsh
curl -s http://localhost:8080/tasks/<taskId> | jq
```

#### `PUT /tasks/{id}`

Update a task.

Request body:

```json
{
  "priority": 2,
  "assignedUserId": "user-2",
  "status": "in_progress",
  "description": "buy milk and bread"
}
```

Responses:

- **200**: updated task
- **404**: not found

Curl (note: replace `<taskId>`):

```zsh
curl -s -X PUT http://localhost:8080/tasks/<taskId> \
  -H 'Content-Type: application/json' \
  -d '{"priority":2,"assignedUserId":"user-2","status":"in_progress","description":"buy milk and bread"}' | jq
```

#### `DELETE /tasks/{id}`

Delete a task.

Responses:

- **204**: deleted
- **404**: not found

Curl:

```zsh
curl -i -X DELETE http://localhost:8080/tasks/<taskId>
```

### List tasks by user

#### `GET /users/{userId}/tasks`

List tasks assigned to a user.

Responses:

- **200**: array of tasks

Curl:

```zsh
curl -s http://localhost:8080/users/user-2/tasks | jq
```

Example output:

```json
[
  {
    "id": "2b0b2f0f-1d3b-4d6e-9f55-2b2b4a5a7890",
    "priority": 2,
    "assignedUserId": "user-2",
    "status": "in_progress",
    "description": "buy milk and bread",
    "createdAt": "2026-03-15T12:34:56Z",
    "updatedAt": "2026-03-15T12:35:10Z"
  }
]
```

## Development

Run unit tests:

```zsh
cd taskmanager

go test ./...
```

## Notes / limitations

- Persistence is in-memory only; data is lost on restart.
- Error mapping is minimal: not-found -> 404, everything else -> 400.
- No authentication/authorization yet.
# CleanArch.taskmanager
