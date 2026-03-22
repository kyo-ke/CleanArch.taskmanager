package usecase

import "errors"

// ErrNotFound is returned when a requested resource doesn't exist.
// Framework layers (HTTP, gRPC, etc.) should depend on this error instead of repository-layer errors.
var ErrNotFound = errors.New("not found")
