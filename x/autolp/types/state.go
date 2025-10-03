package types

// LPState represents the lifecycle state of autolp for a given (connection_id, ica_id).
// Stored as a small string for simplicity and readability.
type LPState string

const (
    LPStateUninitialized LPState = "UNINITIALIZED"
    LPStateFunding       LPState = "FUNDING"
    LPStateReady         LPState = "READY"
    LPStateCreatingPool  LPState = "CREATING_POOL"
    LPStateActive        LPState = "ACTIVE"
    LPStatePaused        LPState = "PAUSED"
)

