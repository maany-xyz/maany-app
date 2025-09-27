package types

import (
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk/types"
    paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultParams returns default autolp params (empty allowlist).
func DefaultParams() Params {
    return Params{
        AllowedSubmitters:   []string{},
        EnsureIcaEnabled:    false,
        ConnectionId:        "",
        IcaId:               "",
        EnsureIcaInterval:   50, // default every 50 blocks
    }
}

// Validate validates params.
func (p Params) Validate() error {
    for _, a := range p.AllowedSubmitters {
        if _, err := sdk.AccAddressFromBech32(a); err != nil {
            return fmt.Errorf("invalid allowed submitter: %s", a)
        }
    }
    if p.EnsureIcaEnabled {
        if p.ConnectionId == "" {
            return fmt.Errorf("connection_id must be set when ensure_ica_enabled is true")
        }
        if p.IcaId == "" {
            return fmt.Errorf("ica_id must be set when ensure_ica_enabled is true")
        }
        if p.EnsureIcaInterval == 0 {
            return fmt.Errorf("ensure_ica_interval must be > 0 when ensure_ica_enabled is true")
        }
    }
    return nil
}

func ParamKeyTable() paramtypes.KeyTable {
    return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements ParamSet for x/params.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
    return paramtypes.ParamSetPairs{
        paramtypes.NewParamSetPair(KeyAllowedSubmitters, &p.AllowedSubmitters, validateAllowedSubmitters),
        paramtypes.NewParamSetPair([]byte("EnsureIcaEnabled"), &p.EnsureIcaEnabled, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair([]byte("ConnectionId"), &p.ConnectionId, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair([]byte("IcaId"), &p.IcaId, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair([]byte("EnsureIcaInterval"), &p.EnsureIcaInterval, func(i interface{}) error { return nil }),
    }
}

func validateAllowedSubmitters(i interface{}) error {
    lst, ok := i.([]string)
    if !ok {
        return fmt.Errorf("invalid type for AllowedSubmitters: %T", i)
    }
    for _, a := range lst {
        if _, err := sdk.AccAddressFromBech32(a); err != nil {
            return fmt.Errorf("invalid allowed submitter: %s", a)
        }
    }
    return nil
}
