package types

import (
    "context"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// IncentivesKeeper defines the minimal interface from x/incentives needed by feeburner.
// This allows feeburner to route a portion of fees into an incentives gauge when available.
type IncentivesKeeper interface {
    AddToGaugeRewards(ctx sdk.Context, owner sdk.AccAddress, coins sdk.Coins, gaugeID uint64) error
}

// AccountKeeper defines the minimal subset needed by feeburner.
type AccountKeeper interface {
    GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the minimal subset needed by feeburner.
// Uses context.Context to match SDK v0.50 keeper signatures.
type BankKeeper interface {
    GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
    BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
    SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
    SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}
