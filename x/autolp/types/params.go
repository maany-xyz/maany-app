package types

// DefaultParams returns default autolp params (none currently).
func DefaultParams() Params { return Params{} }

// Validate validates params; no-op for now.
func (p Params) Validate() error { return nil }

