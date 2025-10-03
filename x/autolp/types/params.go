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
        BaseHostDenom:       "",
        QuoteHostDenom:      "",
        CreatePoolRequired:  sdk.Coins{},
        AutoCreatePoolEnabled: false,
        IcqUpdatePeriod:     20,
        IcaTxTimeoutSeconds: 600,
        TransferTimeoutSeconds: 600,
        HostAccountHrp:      "",
        BaseLocalDenom:      "uapp",
        BaseDexChannel:      "",
        QuoteLocalDenom:     "uusdc",
        QuoteDexChannel:     "",
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

    // If pool auto-creation is enabled, validate host denoms and thresholds
    if p.AutoCreatePoolEnabled {
        if p.BaseHostDenom == "" {
            return fmt.Errorf("base_host_denom must be set when auto_create_pool_enabled is true")
        }
        if p.QuoteHostDenom == "" {
            return fmt.Errorf("quote_host_denom must be set when auto_create_pool_enabled is true")
        }
        if len(p.CreatePoolRequired) == 0 {
            return fmt.Errorf("create_pool_required must have at least one entry when auto_create_pool_enabled is true")
        }
        for _, c := range p.CreatePoolRequired {
            if !c.IsValid() || !c.IsPositive() {
                return fmt.Errorf("invalid create_pool_required coin: %s", c)
            }
        }
        if p.IcqUpdatePeriod == 0 {
            return fmt.Errorf("icq_update_period must be > 0 when auto_create_pool_enabled is true")
        }
        if p.IcaTxTimeoutSeconds == 0 {
            return fmt.Errorf("ica_tx_timeout_seconds must be > 0 when auto_create_pool_enabled is true")
        }
        if p.TransferTimeoutSeconds == 0 {
            return fmt.Errorf("transfer_timeout_seconds must be > 0 when auto_create_pool_enabled is true")
        }
        if p.HostAccountHrp == "" {
            return fmt.Errorf("host_account_hrp must be set when auto_create_pool_enabled is true")
        }
        // Optional runtime derivation support for base_local_denom
        if p.BaseLocalDenom != "" && p.BaseHostDenom == "" && p.BaseDexChannel == "" {
            // If thresholds are specified in base_local_denom, either BaseHostDenom or BaseDexChannel should be set
            // We don't know if thresholds include base_local_denom here; warn via validation only if both are empty.
            // Keeping it soft to allow USDC-only pools.
        }
        // Optional runtime derivation support for quote_local_denom
        if p.QuoteLocalDenom != "" && p.QuoteHostDenom == "" && p.QuoteDexChannel == "" {
            // Same soft guidance for quote.
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
        paramtypes.NewParamSetPair(KeyEnsureIcaEnabled, &p.EnsureIcaEnabled, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyConnectionId, &p.ConnectionId, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyIcaId, &p.IcaId, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyEnsureIcaInterval, &p.EnsureIcaInterval, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyBaseHostDenom, &p.BaseHostDenom, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyQuoteHostDenom, &p.QuoteHostDenom, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyCreatePoolRequired, &p.CreatePoolRequired, validateCreatePoolRequired),
        paramtypes.NewParamSetPair(KeyAutoCreatePoolEnabled, &p.AutoCreatePoolEnabled, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyICQUpdatePeriod, &p.IcqUpdatePeriod, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyICATxTimeoutSeconds, &p.IcaTxTimeoutSeconds, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyTransferTimeoutSeconds, &p.TransferTimeoutSeconds, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyHostAccountHRP, &p.HostAccountHrp, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyBaseLocalDenom, &p.BaseLocalDenom, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyBaseDexChannel, &p.BaseDexChannel, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyQuoteLocalDenom, &p.QuoteLocalDenom, func(i interface{}) error { return nil }),
        paramtypes.NewParamSetPair(KeyQuoteDexChannel, &p.QuoteDexChannel, func(i interface{}) error { return nil }),
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

func validateCreatePoolRequired(i interface{}) error {
    switch v := i.(type) {
    case sdk.Coins:
        for _, c := range v {
            if !c.IsValid() || !c.IsPositive() {
                return fmt.Errorf("invalid create_pool_required coin: %s", c)
            }
        }
        return nil
    case []sdk.Coin:
        for _, c := range v {
            if !c.IsValid() || !c.IsPositive() {
                return fmt.Errorf("invalid create_pool_required coin: %s", c)
            }
        }
        return nil
    default:
        return fmt.Errorf("invalid type for CreatePoolRequired: %T", i)
    }
}
