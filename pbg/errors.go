package pbg

import "errors"

var (
    // Server building
    ErrUnspecifiedDataM  = errors.New("Data mechanism not set in the builder")
    ErrUnspecifiedSessM  = errors.New("Session mechanism not set in the builder")
    ErrUnspecifiedAuthM  = errors.New("Authorization mechanism not set in the builder")

    ErrInvalidConfiguration = errors.New("Invalid Configuration value used in the builder")
    ErrInvalidAPIResponser  = errors.New("Invalid APIResponser value used in the builder")

    ErrMoveNotFound    = errors.New("Move not found")
    ErrPokèmonNotFound = errors.New("Pokèmon not found")
    ErrTrainerNotFound = errors.New("Trainer not found")
    ErrSessionNotFound = errors.New("Session not found")

    ErrTrainerAlreadyExists = errors.New("Trainer already exists")
)
