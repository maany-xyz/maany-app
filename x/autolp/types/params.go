package types

import (
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk/types"
    paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// DefaultParams returns default autolp params (empty allowlist).
func DefaultParams() Params { return Params{AllowedSubmitters: []string{}} }

// Validate validates params.
func (p Params) Validate() error {
    for _, a := range p.AllowedSubmitters {
        if _, err := sdk.AccAddressFromBech32(a); err != nil {
            return fmt.Errorf("invalid allowed submitter: %s", a)
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
