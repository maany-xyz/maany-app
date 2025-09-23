package keeper

import (
    "context"
    "fmt"

    "cosmossdk.io/errors"
    "cosmossdk.io/log"
    "cosmossdk.io/math"

    authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

    storetypes "cosmossdk.io/store/types"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    consumertypes "github.com/cosmos/interchain-security/v5/x/ccv/consumer/types"

    "github.com/maany-xyz/maany-app/osmomath"
    "github.com/maany-xyz/maany-app/x/feeburner/types"
)

type (
    Keeper struct {
        cdc      codec.BinaryCodec
        storeKey storetypes.StoreKey
        memKey   storetypes.StoreKey

        accountKeeper types.AccountKeeper
        bankKeeper    types.BankKeeper
        authority     string

        // Optional routing to x/incentives. If IncentivesKeeper is non-nil and
        // IncentivesGaugeID > 0 and StakeFraction > 0, a portion of native fees
        // will be added to the configured gauge each time BurnAndDistribute runs.
        incentivesKeeper  types.IncentivesKeeper
        incentivesGaugeID uint64
        stakeFraction     osmomath.Dec
    }
)

var KeyBurnedFees = []byte("BurnedFees")

func NewKeeper(
    cdc codec.BinaryCodec,
    storeKey,
    memKey storetypes.StoreKey,
    accountKeeper types.AccountKeeper,
    bankKeeper types.BankKeeper,
    authority string,
) *Keeper {
    return &Keeper{
        cdc:           cdc,
        storeKey:      storeKey,
        memKey:        memKey,
        accountKeeper: accountKeeper,
        bankKeeper:    bankKeeper,
        authority:     authority,
        stakeFraction: osmomath.ZeroDec(),
    }
}

func (k Keeper) GetAuthority() string {
    return k.authority
}

// SetIncentivesRouting configures optional routing of a portion of native fees
// into an incentives gauge. Passing a nil incentives keeper or zero gaugeID or
// zero fraction disables the routing.
func (k *Keeper) SetIncentivesRouting(ik types.IncentivesKeeper, gaugeID uint64, fraction osmomath.Dec) {
    k.incentivesKeeper = ik
    k.incentivesGaugeID = gaugeID
    k.stakeFraction = fraction
}

// RecordBurnedFees adds `amount` to the total amount of burned NTRN tokens
func (k Keeper) RecordBurnedFees(ctx sdk.Context, amount sdk.Coin) {
	store := ctx.KVStore(k.storeKey)

	totalBurnedNeutronsAmount := k.GetTotalBurnedNeutronsAmount(ctx)
	totalBurnedNeutronsAmount.Coin = totalBurnedNeutronsAmount.Coin.Add(amount)

	store.Set(KeyBurnedFees, k.cdc.MustMarshal(&totalBurnedNeutronsAmount))
}

// GetTotalBurnedNeutronsAmount gets the total burned amount of NTRN tokens
func (k Keeper) GetTotalBurnedNeutronsAmount(ctx sdk.Context) types.TotalBurnedNeutronsAmount {
	store := ctx.KVStore(k.storeKey)

	var totalBurnedNeutronsAmount types.TotalBurnedNeutronsAmount
	bzTotalBurnedNeutronsAmount := store.Get(KeyBurnedFees)
	if bzTotalBurnedNeutronsAmount != nil {
		k.cdc.MustUnmarshal(bzTotalBurnedNeutronsAmount, &totalBurnedNeutronsAmount)
	}

	if totalBurnedNeutronsAmount.Coin.Denom == "" {
		totalBurnedNeutronsAmount.Coin = sdk.NewCoin(k.GetParams(ctx).NeutronDenom, math.NewInt(0))
	}

	return totalBurnedNeutronsAmount
}

// SetTotalBurnedNeutronsAmount sets the total burned amount of NTRN tokens
func (k Keeper) SetTotalBurnedNeutronsAmount(ctx sdk.Context, totalBurnedNeutronsAmount types.TotalBurnedNeutronsAmount) {
	store := ctx.KVStore(k.storeKey)

	store.Set(KeyBurnedFees, k.cdc.MustMarshal(&totalBurnedNeutronsAmount))
}

// BurnAndDistribute is an important part of tokenomics. It does few things:
// 1. Burns NTRN fee coins distributed to consumertypes.ConsumerRedistributeName in ICS (https://github.com/cosmos/interchain-security/blob/86046926502f7b0ba795bebcdd1fdc97ac776573/x/ccv/consumer/keeper/distribution.go#L67)
// 2. Updates total amount of burned NTRN coins
// 3. Sends non-NTRN fee tokens to reserve contract address
// Panics if no `consumertypes.ConsumerRedistributeName` module found OR could not burn NTRN tokens
func (k Keeper) BurnAndDistribute(ctx sdk.Context) {
    moduleAddr := k.accountKeeper.GetModuleAddress(consumertypes.ConsumerRedistributeName)
    if moduleAddr == nil {
        panic("ConsumerRedistributeName must have module address")
    }

    params := k.GetParams(ctx)
    balances := k.bankKeeper.GetAllBalances(ctx, moduleAddr)
    fundsForReserve := make(sdk.Coins, 0, len(balances))

    for _, balance := range balances {
        if !balance.IsZero() {
            if balance.Denom == params.NeutronDenom {
                // Do not burn native consumer fees anymore.
                // If incentives routing is configured, route ALL native fees to the gauge.
                if k.incentivesKeeper != nil && k.incentivesGaugeID > 0 {
                    fullStake := sdk.NewCoins(balance)
                    if err := k.incentivesKeeper.AddToGaugeRewards(ctx, moduleAddr, fullStake, k.incentivesGaugeID); err != nil {
                        panic(errors.Wrapf(err, "failed to add fees to incentives gauge"))
                    }
                    k.Logger(ctx).Info(
                        "routed all native consumer fees to incentives gauge",
                        "gauge_id", k.incentivesGaugeID,
                        "amount", fullStake.String(),
                    )
                } else {
                    // If no incentives routing configured, send to treasury if set; otherwise keep funds in module.
                    if params.TreasuryAddress != "" {
                        addr, err := sdk.AccAddressFromBech32(params.TreasuryAddress)
                        if err == nil {
                            if err := k.bankKeeper.SendCoins(ctx, moduleAddr, addr, sdk.NewCoins(balance)); err != nil {
                                panic(errors.Wrapf(err, "failed sending native consumer fees to Treasury"))
                            }
                            k.Logger(ctx).Info("sent native consumer fees to treasury", "treasury", params.TreasuryAddress, "amount", balance.String())
                        } else {
                            k.Logger(ctx).Info("treasury address invalid; keeping native consumer fees in module", "treasury", params.TreasuryAddress, "amount", balance.String())
                        }
                    } else {
                        k.Logger(ctx).Info("no incentives routing or treasury configured; keeping native consumer fees in module", "amount", balance.String())
                    }
                }
            } else {
                fundsForReserve = append(fundsForReserve, balance)
            }
        }
    }

    if len(fundsForReserve) > 0 {
        addr, err := sdk.AccAddressFromBech32(params.TreasuryAddress)
        if err != nil {
			// there's no way we face this kind of situation in production, since it means the chain is misconfigured
			// still, in test environments it might be the case when the chain is started without Reserve
			// in such case we just burn the tokens
            err := k.bankKeeper.BurnCoins(ctx, consumertypes.ConsumerRedistributeName, fundsForReserve)
            if err != nil {
                panic(errors.Wrapf(err, "failed to burn tokens during fee processing"))
            }
            k.Logger(ctx).Info("burned non-native consumer fees (no treasury configured)", "amount", fundsForReserve.String())
        } else {
            err = k.bankKeeper.SendCoins(
                ctx,
                moduleAddr, addr,
                fundsForReserve,
            )
            if err != nil {
                panic(errors.Wrapf(err, "failed sending funds to Reserve"))
            }
            k.Logger(ctx).Info("sent non-native consumer fees to treasury", "treasury", params.TreasuryAddress, "amount", fundsForReserve.String())
        }
    }
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// FundCommunityPool is method to satisfy DistributionKeeper interface for packet-forward-middleware Keeper.
// The original method sends coins to a community pool of a chain.
// The current method sends coins to a Fee Collector module which collects fee on consumer chain.
func (k Keeper) FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error {
    return k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, authtypes.FeeCollectorName, amount)
}
