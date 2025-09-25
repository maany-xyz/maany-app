package autolp

import (
    "encoding/json"

    abci "github.com/cometbft/cometbft/abci/types"
    "cosmossdk.io/core/appmodule"
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/codec"
    cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/types/module"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"

    "github.com/maany-xyz/maany-app/x/autolp/keeper"
    "github.com/maany-xyz/maany-app/x/autolp/types"
)

type AppModuleBasic struct{ cdc codec.BinaryCodec }

func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic { return AppModuleBasic{cdc: cdc} }

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) { types.RegisterInterfaces(reg) }

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
    _ = cdc // not used; use std json to avoid protobuf msgs
    bz, _ := json.Marshal(types.DefaultGenesis())
    return bz
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
    _ = cdc
    var gs types.GenesisState
    if err := json.Unmarshal(bz, &gs); err != nil {
        return err
    }
    return gs.Validate()
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
    // no-op for now; we defined QueryServer in Go
    _ = clientCtx
    _ = mux
}

type AppModule struct {
    AppModuleBasic
    keeper keeper.Keeper
}

func NewAppModule(cdc codec.Codec, k keeper.Keeper) AppModule {
    return AppModule{AppModuleBasic: NewAppModuleBasic(cdc), keeper: k}
}

// Interface assertions
var (
    _ appmodule.AppModule  = AppModule{}
    _ module.AppModule     = AppModule{}
    _ module.AppModuleBasic = AppModuleBasic{}
)

// markers for appmodule compile-time checks
func (am AppModule) IsOnePerModuleType() {}
func (am AppModule) IsAppModule() {}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (am AppModule) RegisterServices(cfg module.Configurator) {
    types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
    types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
    var gs types.GenesisState
    // don't use codec here to avoid proto dependency
    _ = json.Unmarshal(data, &gs)
    InitGenesis(ctx, am.keeper, gs)
    return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
    gs := ExportGenesis(ctx, am.keeper)
    bz, _ := json.Marshal(gs)
    return bz
}

func (am AppModule) ConsensusVersion() uint64 { return types.ConsensusVersion }

func (AppModule) QuerierRoute() string { return types.RouterKey }

// Ensure it satisfies interfaces
var _ module.AppModuleBasic = AppModuleBasic{}
