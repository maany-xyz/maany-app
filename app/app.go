package app

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	v502 "github.com/maany-xyz/maany-app/app/upgrades/v5.0.2"
	v504 "github.com/maany-xyz/maany-app/app/upgrades/v5.0.4"
	v505 "github.com/maany-xyz/maany-app/app/upgrades/v5.0.5"








	ibcratelimit "github.com/maany-xyz/maany-app/x/ibc-rate-limit"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"

	appconfig "github.com/maany-xyz/maany-app/app/config"











	v500 "github.com/maany-xyz/maany-app/app/upgrades/v5.0.0"
	"github.com/maany-xyz/maany-app/x/globalfee"
	globalfeetypes "github.com/maany-xyz/maany-app/x/globalfee/types"

	"cosmossdk.io/log"
	db "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec/address"


	"github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward"
	ibctestingtypes "github.com/cosmos/ibc-go/v8/testing/types"
	"github.com/cosmos/interchain-security/v5/testutil/integration"
	ccv "github.com/cosmos/interchain-security/v5/x/ccv/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	tendermint "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"

	"github.com/maany-xyz/maany-app/docs"

	"github.com/maany-xyz/maany-app/app/upgrades"

	"github.com/maany-xyz/maany-app/x/cron"

	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"

	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"


	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"

	ibcratelimitkeeper "github.com/maany-xyz/maany-app/x/ibc-rate-limit/keeper"
	ibcratelimittypes "github.com/maany-xyz/maany-app/x/ibc-rate-limit/types"


	ibcporttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	"github.com/spf13/cast"

	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	cronkeeper "github.com/maany-xyz/maany-app/x/cron/keeper"
	crontypes "github.com/maany-xyz/maany-app/x/cron/types"


	"github.com/cosmos/admin-module/v2/x/adminmodule"
	adminmodulecli "github.com/cosmos/admin-module/v2/x/adminmodule/client/cli"
	adminmodulekeeper "github.com/cosmos/admin-module/v2/x/adminmodule/keeper"
	adminmoduletypes "github.com/cosmos/admin-module/v2/x/adminmodule/types"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	appparams "github.com/maany-xyz/maany-app/app/params"
	"github.com/maany-xyz/maany-app/wasmbinding"
	"github.com/maany-xyz/maany-app/x/contractmanager"
	contractmanagermodulekeeper "github.com/maany-xyz/maany-app/x/contractmanager/keeper"
	contractmanagermoduletypes "github.com/maany-xyz/maany-app/x/contractmanager/types"

	"github.com/maany-xyz/maany-app/x/autolp"
	autolpkeeper "github.com/maany-xyz/maany-app/x/autolp/keeper"
	autolptypes "github.com/maany-xyz/maany-app/x/autolp/types"


	"github.com/maany-xyz/maany-app/x/feeburner"
	feeburnerkeeper "github.com/maany-xyz/maany-app/x/feeburner/keeper"
	feeburnertypes "github.com/maany-xyz/maany-app/x/feeburner/types"
	"github.com/maany-xyz/maany-app/x/feerefunder"
	feekeeper "github.com/maany-xyz/maany-app/x/feerefunder/keeper"
	ibchooks "github.com/maany-xyz/maany-app/x/ibc-hooks"
	ibchookstypes "github.com/maany-xyz/maany-app/x/ibc-hooks/types"
	"github.com/maany-xyz/maany-app/x/incentives"
	incentiveskeeper "github.com/maany-xyz/maany-app/x/incentives/keeper"
	incentivestypes "github.com/maany-xyz/maany-app/x/incentives/types"
	"github.com/maany-xyz/maany-app/x/interchainqueries"
	interchainqueriesmodulekeeper "github.com/maany-xyz/maany-app/x/interchainqueries/keeper"
	interchainqueriesmoduletypes "github.com/maany-xyz/maany-app/x/interchainqueries/types"
	"github.com/maany-xyz/maany-app/x/interchaintxs"
	interchaintxskeeper "github.com/maany-xyz/maany-app/x/interchaintxs/keeper"
	interchaintxstypes "github.com/maany-xyz/maany-app/x/interchaintxs/types"
	"github.com/maany-xyz/maany-app/x/lockup"
	lockupkeeper "github.com/maany-xyz/maany-app/x/lockup/keeper"
	lockuptypes "github.com/maany-xyz/maany-app/x/lockup/types"
	transferSudo "github.com/maany-xyz/maany-app/x/transfer"
	wrapkeeper "github.com/maany-xyz/maany-app/x/transfer/keeper"

	feetypes "github.com/maany-xyz/maany-app/x/feerefunder/types"

	ccvconsumer "github.com/cosmos/interchain-security/v5/x/ccv/consumer"
	ccvconsumerkeeper "github.com/cosmos/interchain-security/v5/x/ccv/consumer/keeper"
	ccvconsumertypes "github.com/cosmos/interchain-security/v5/x/ccv/consumer/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	pfmkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/keeper"
	pfmtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"





	globalfeekeeper "github.com/maany-xyz/maany-app/x/globalfee/keeper"
	gmpmiddleware "github.com/maany-xyz/maany-app/x/gmp"



	blocksdk "github.com/skip-mev/block-sdk/v2/block"

	// epochs module (reused from Osmosis)
	epochs "github.com/maany-xyz/maany-app/x/epochs"
	epochskeeper "github.com/maany-xyz/maany-app/x/epochs/keeper"
	epochstypes "github.com/maany-xyz/maany-app/x/epochs/types"






	"github.com/skip-mev/block-sdk/v2/abci/checktx"








	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"




)

// noOpProtorevKeeper is a minimal implementation that satisfies
// incentives' ProtorevKeeper interface by always reporting that a route exists.
// This allows distributing rewards in the native denom without AMM pools.
type noOpProtorevKeeper struct{}

func (noOpProtorevKeeper) GetPoolForDenomPairNoOrder(ctx sdk.Context, denom1, denom2 string) (uint64, error) {
    return 1, nil
}

// simpleTxFeesKeeper returns the chain base denom for incentives' fee checks.
type simpleTxFeesKeeper struct{}

func (simpleTxFeesKeeper) GetBaseDenom(ctx sdk.Context) (string, error) {
    return appparams.DefaultDenom, nil
}

const (
    Name = "maanyappd"
)

var (
	Upgrades = []upgrades.Upgrade{
		v500.Upgrade,
		v502.Upgrade,
		v504.Upgrade,
		v505.Upgrade,
	}


	DefaultNodeHome string




	ModuleBasics = module.NewBasicManager(



		auth.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		vesting.AppModuleBasic{},
		tendermint.AppModuleBasic{},
		consensus.AppModuleBasic{},
		epochs.AppModuleBasic{},




		ibc.AppModuleBasic{},
		ccvconsumer.AppModuleBasic{},
		ica.AppModuleBasic{},
		transferSudo.AppModuleBasic{},



		interchainqueries.AppModuleBasic{},
		interchaintxs.AppModuleBasic{},

		ibcratelimit.AppModuleBasic{},
		packetforward.AppModuleBasic{},
		feerefunder.AppModuleBasic{},
		feeburner.AppModuleBasic{},
		contractmanager.AppModuleBasic{},
		cron.AppModuleBasic{},
		autolp.AppModuleBasic{},

		globalfee.AppModule{},



		wasm.AppModuleBasic{},
		ibchooks.AppModuleBasic{},
		lockup.AppModuleBasic{},
		incentives.AppModuleBasic{},




		adminmodule.NewAppModuleBasic(
			govclient.NewProposalHandler(
				adminmodulecli.NewSubmitParamChangeProposalTxCmd,
			),
			govclient.NewProposalHandler(
				adminmodulecli.NewCmdSubmitUpgradeProposal,
			),
			govclient.NewProposalHandler(
				adminmodulecli.NewCmdSubmitCancelUpgradeProposal,
			),
		),





	)


	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:                    nil,

		ibctransfertypes.ModuleName:                   {authtypes.Minter, authtypes.Burner},
		icatypes.ModuleName:                           nil,
		wasmtypes.ModuleName:                          {},
		interchainqueriesmoduletypes.ModuleName:       nil,
		feetypes.ModuleName:                           nil,
		feeburnertypes.ModuleName:                     nil,
		lockuptypes.ModuleName:                        nil,
		incentivestypes.ModuleName:                    nil,
		ccvconsumertypes.ConsumerRedistributeName:     {authtypes.Burner},
		ccvconsumertypes.ConsumerToSendToProviderName: nil,
		crontypes.ModuleName:                          nil,





	}
)

var (
	_ runtime.AppI            = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
	_ ibctesting.TestingApp   = (*App)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+Name)

	appconfig.GetDefaultConfig()
}




type App struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	configurator module.Configurator

	encodingConfig appparams.EncodingConfig

	invCheckPeriod uint


	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey


	AccountKeeper     authkeeper.AccountKeeper
	AdminmoduleKeeper adminmodulekeeper.Keeper
	AuthzKeeper       authzkeeper.Keeper
	BankKeeper        bankkeeper.BaseKeeper


	CapabilityKeeper    *capabilitykeeper.Keeper
	SlashingKeeper      slashingkeeper.Keeper
	CrisisKeeper        crisiskeeper.Keeper
	UpgradeKeeper       upgradekeeper.Keeper
	ParamsKeeper        paramskeeper.Keeper
	IBCKeeper           *ibckeeper.Keeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper
	EvidenceKeeper      evidencekeeper.Keeper
	TransferKeeper      wrapkeeper.KeeperTransferWrapper
	FeeGrantKeeper      feegrantkeeper.Keeper

	FeeKeeper           *feekeeper.Keeper
	FeeBurnerKeeper     *feeburnerkeeper.Keeper
	ConsumerKeeper      ccvconsumerkeeper.Keeper
	CronKeeper          cronkeeper.Keeper
	PFMKeeper           *pfmkeeper.Keeper

	GlobalFeeKeeper     globalfeekeeper.Keeper
	EpochsKeeper        *epochskeeper.Keeper
	LockupKeeper        *lockupkeeper.Keeper
	IncentivesKeeper    *incentiveskeeper.Keeper
	AutolpKeeper        autolpkeeper.Keeper

	PFMModule packetforward.AppModule

	TransferStack           *ibchooks.IBCMiddleware
	Ics20WasmHooks          *ibchooks.WasmHooks
	RateLimitingICS4Wrapper *ibcratelimit.ICS4Wrapper
	HooksICS4Wrapper        ibchooks.ICS4Middleware


	ScopedIBCKeeper         capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper    capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper        capabilitykeeper.ScopedKeeper
	ScopedInterTxKeeper     capabilitykeeper.ScopedKeeper
	ScopedCCVConsumerKeeper capabilitykeeper.ScopedKeeper

	InterchainQueriesKeeper interchainqueriesmodulekeeper.Keeper
	InterchainTxsKeeper     interchaintxskeeper.Keeper
	ContractManagerKeeper   contractmanagermodulekeeper.Keeper

	ConsensusParamsKeeper consensusparamkeeper.Keeper

	WasmKeeper     wasmkeeper.Keeper
	ContractKeeper *wasmkeeper.PermissionedKeeper










	mm *module.Manager


	sm *module.SimulationManager



	checkTxHandler checktx.CheckTx






}


func (app *App) AutoCLIOpts(initClientCtx client.Context) autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule)
	for _, m := range app.mm.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(app.mm.Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
		ClientCtx:             initClientCtx,
	}
}

func (app *App) GetTestBankKeeper() integration.TestBankKeeper {
	return app.BankKeeper
}

func (app *App) GetTestAccountKeeper() integration.TestAccountKeeper {
	return app.AccountKeeper
}

func (app *App) GetTestSlashingKeeper() integration.TestSlashingKeeper {
	return app.SlashingKeeper
}

func (app *App) GetTestEvidenceKeeper() evidencekeeper.Keeper {
	return app.EvidenceKeeper
}


func New(
	logger log.Logger,
	db db.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig appparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	wasmOpts []wasmkeeper.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	overrideWasmVariables()

	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(Name, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := storetypes.NewKVStoreKeys(
		authzkeeper.StoreKey, authtypes.StoreKey, banktypes.StoreKey, slashingtypes.StoreKey,
		paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey, feegrant.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, icacontrollertypes.StoreKey,
		icahosttypes.StoreKey, capabilitytypes.StoreKey,
		interchainqueriesmoduletypes.StoreKey, contractmanagermoduletypes.StoreKey, interchaintxstypes.StoreKey, wasmtypes.StoreKey, feetypes.StoreKey,
        feeburnertypes.StoreKey, adminmoduletypes.StoreKey, ccvconsumertypes.StoreKey,
		pfmtypes.StoreKey,
		crontypes.StoreKey, ibcratelimittypes.ModuleName, ibchookstypes.StoreKey, consensusparamtypes.StoreKey, crisistypes.StoreKey,
	epochstypes.StoreKey,
	lockuptypes.StoreKey,
		incentivestypes.StoreKey,
		autolptypes.StoreKey,


		globalfeetypes.StoreKey,
	)
	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey, feetypes.MemStoreKey)

	app := &App{
		BaseApp:           bApp,
		cdc:               legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
		encodingConfig:    encodingConfig,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// Epochs keeper (hooks are set later once all keepers are initialized)
	app.EpochsKeeper = epochskeeper.NewKeeper(keys[epochstypes.StoreKey])


	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]), authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(), runtime.EventService{})
	bApp.SetParamStore(&app.ConsensusParamsKeeper.ParamsStore)


	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])


	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedICAControllerKeeper := app.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	app.ScopedTransferKeeper = scopedTransferKeeper
	scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)
	scopedInterTxKeeper := app.CapabilityKeeper.ScopeToModule(interchaintxstypes.ModuleName)
	scopedCCVConsumerKeeper := app.CapabilityKeeper.ScopeToModule(ccvconsumertypes.ModuleName)


	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		address.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[authz.ModuleName]), appCodec, app.MsgServiceRouter(), app.AccountKeeper,
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.AccountKeeper,
		app.BlockedAddrs(),
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
		logger,
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		&app.ConsumerKeeper,
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)
	app.CrisisKeeper = *crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
		invCheckPeriod,
		&app.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
		address.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, runtime.NewKVStoreService(keys[feegrant.StoreKey]), app.AccountKeeper)
	app.UpgradeKeeper = *upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		app.BaseApp,
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)















	app.ConsumerKeeper = ccvconsumerkeeper.NewNonZeroKeeper(
		appCodec,
		keys[ccvconsumertypes.StoreKey],
		app.GetSubspace(ccvconsumertypes.ModuleName),
	)


	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), &app.ConsumerKeeper, app.UpgradeKeeper, scopedIBCKeeper, authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)




	app.FeeKeeper = feekeeper.NewKeeper(
		appCodec,
		keys[feetypes.StoreKey],
		memKeys[feetypes.MemStoreKey],
		app.IBCKeeper.ChannelKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)
	feeModule := feerefunder.NewAppModule(appCodec, *app.FeeKeeper, app.AccountKeeper, app.BankKeeper)

	app.ContractManagerKeeper = *contractmanagermodulekeeper.NewKeeper(
		appCodec,
		keys[contractmanagermoduletypes.StoreKey],
		keys[contractmanagermoduletypes.MemStoreKey],
		&app.WasmKeeper,
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)










	app.WireICS20PreWasmKeeper(appCodec)
	app.PFMModule = packetforward.NewAppModule(app.PFMKeeper, app.GetSubspace(pfmtypes.ModuleName))

	// autolp keeper (uses transfer wrapper and interchaintxs query)
	app.AutolpKeeper = autolpkeeper.NewKeeper(
		appCodec,
		keys[autolptypes.StoreKey],
		app.TransferKeeper,
		app.InterchainTxsKeeper,
		app.ICAControllerKeeper,
		icacontrollerkeeper.NewMsgServerImpl(&app.ICAControllerKeeper),
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)

	app.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec, keys[icacontrollertypes.StoreKey], app.GetSubspace(icacontrollertypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper, app.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper, app.MsgServiceRouter(),
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)

	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, keys[icahosttypes.StoreKey], app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper, app.IBCKeeper.PortKeeper,
		app.AccountKeeper, scopedICAHostKeeper, app.MsgServiceRouter(),
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)
	app.ICAHostKeeper.WithQueryRouter(app.GRPCQueryRouter())


	app.FeeBurnerKeeper = feeburnerkeeper.NewKeeper(
		appCodec,
		keys[feeburnertypes.StoreKey],
		keys[feeburnertypes.MemStoreKey],
		app.AccountKeeper,
		&app.BankKeeper,
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)
	feeBurnerModule := feeburner.NewAppModule(appCodec, *app.FeeBurnerKeeper)

	app.GlobalFeeKeeper = globalfeekeeper.NewKeeper(appCodec, keys[globalfeetypes.StoreKey], authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String())

	// Lockup and Incentives keepers
	app.LockupKeeper = lockupkeeper.NewKeeper(
		keys[lockuptypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.FeeBurnerKeeper, // implements FundCommunityPool
		app.GetSubspace(lockuptypes.ModuleName),
	)

	app.IncentivesKeeper = incentiveskeeper.NewKeeper(
		keys[incentivestypes.StoreKey],
		app.GetSubspace(incentivestypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		*app.LockupKeeper,
		*app.EpochsKeeper,
		app.FeeBurnerKeeper, // CommunityPoolKeeper
		simpleTxFeesKeeper{}, // TxFeesKeeper
		nil, // PoolIncentiveKeeper
		noOpProtorevKeeper{}, // ProtorevKeeper
	)

	// Register incentives epoch hooks now that keepers exist
	app.EpochsKeeper = app.EpochsKeeper.SetHooks(epochstypes.NewMultiEpochHooks(app.IncentivesKeeper.Hooks()))


	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(keys[evidencetypes.StoreKey]), &app.ConsumerKeeper, app.SlashingKeeper,
		address.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()), runtime.ProvideCometInfoService(),
	)

	app.EvidenceKeeper = *evidenceKeeper

	app.ConsumerKeeper = ccvconsumerkeeper.NewKeeper(
		appCodec,
		keys[ccvconsumertypes.StoreKey],
		app.GetSubspace(ccvconsumertypes.ModuleName),
		scopedCCVConsumerKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		app.IBCKeeper.ConnectionKeeper,
		app.IBCKeeper.ClientKeeper,
		app.SlashingKeeper,
		&app.BankKeeper,
		app.AccountKeeper,
		app.TransferKeeper.Keeper,

		app.IBCKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
		address.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		address.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
	app.ConsumerKeeper = *app.ConsumerKeeper.SetHooks(app.SlashingKeeper.Hooks())
	consumerModule := ccvconsumer.NewAppModule(app.ConsumerKeeper, app.GetSubspace(ccvconsumertypes.ModuleName))



























	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm cfg: %s", err))
	}






	supportedFeatures := []string{"iterator", "stargate", "staking", "neutron", "cosmwasm_1_1", "cosmwasm_1_2", "cosmwasm_1_3", "cosmwasm_1_4", "cosmwasm_2_0", "cosmwasm_2_1"}


	adminRouterLegacy := govv1beta1.NewRouter()
	adminRouterLegacy.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper))

	app.AdminmoduleKeeper = *adminmodulekeeper.NewKeeper(
		appCodec,
		keys[adminmoduletypes.StoreKey],
		keys[adminmoduletypes.MemStoreKey],
		adminRouterLegacy,
		app.MsgServiceRouter(),
		IsConsumerProposalAllowlisted,
		isSdkMessageWhitelisted,
	)
	adminModule := adminmodule.NewAppModule(appCodec, app.AdminmoduleKeeper)

	app.InterchainQueriesKeeper = *interchainqueriesmodulekeeper.NewKeeper(
		appCodec,
		keys[interchainqueriesmoduletypes.StoreKey],
		keys[interchainqueriesmoduletypes.MemStoreKey],
		app.IBCKeeper,
		&app.BankKeeper,
		app.ContractManagerKeeper,
		interchainqueriesmodulekeeper.Verifier{},
		interchainqueriesmodulekeeper.TransactionVerifier{},
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)
	app.InterchainTxsKeeper = *interchaintxskeeper.NewKeeper(
		appCodec,
		keys[interchaintxstypes.StoreKey],
		memKeys[interchaintxstypes.MemStoreKey],
		app.IBCKeeper.ChannelKeeper,
		app.ICAControllerKeeper,
		icacontrollerkeeper.NewMsgServerImpl(&app.ICAControllerKeeper),
		contractmanager.NewSudoLimitWrapper(app.ContractManagerKeeper, &app.WasmKeeper),
		app.FeeKeeper,
		app.BankKeeper,
		func(ctx sdk.Context) string { return app.FeeBurnerKeeper.GetParams(ctx).TreasuryAddress },
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)

















	app.CronKeeper = *cronkeeper.NewKeeper(
		appCodec,
		keys[crontypes.StoreKey],
		keys[crontypes.MemStoreKey],
		app.AccountKeeper,
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)
    wasmOpts = append(wasmbinding.RegisterCustomPlugins(
        &app.InterchainTxsKeeper,
        &app.InterchainQueriesKeeper,

        &app.AdminmoduleKeeper,
        app.FeeBurnerKeeper,
        app.FeeKeeper,
        &app.BankKeeper,
        &app.CronKeeper,
        &app.ContractManagerKeeper,
        nil,


        nil,
    ), wasmOpts...)

	queryPlugins := wasmkeeper.WithQueryPlugins(
		&wasmkeeper.QueryPlugins{
			Stargate: wasmkeeper.AcceptListStargateQuerier(wasmbinding.AcceptedStargateQueries(), app.GRPCQueryRouter(), appCodec),
			Grpc:     wasmkeeper.AcceptListGrpcQuerier(wasmbinding.AcceptedStargateQueries(), app.GRPCQueryRouter(), appCodec),
		})
	wasmOpts = append(wasmOpts, queryPlugins)

	app.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[wasmtypes.StoreKey]),
		app.AccountKeeper,
		&app.BankKeeper,
		nil,
		nil,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
		wasmOpts...,
	)

	app.CronKeeper.WasmMsgServer = wasmkeeper.NewMsgServerImpl(&app.WasmKeeper)
	cronModule := cron.NewAppModule(appCodec, app.CronKeeper)


	ibcRouter := ibcporttypes.NewRouter()

    icaModule := ica.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper)

    // Build ICA controller stack: combine interchaintxs and autolp IBC modules, then wrap with ICA controller middleware.
    var icaControllerStack ibcporttypes.IBCModule
    {
        base := interchaintxs.NewIBCModule(app.InterchainTxsKeeper)
        // autolp IBC for event hooks
        autolpIBC := autolp.NewIBCModule(app.AutolpKeeper)
        combined := autolp.CombineIBCModules(base, autolpIBC)
        icaControllerStack = icacontroller.NewIBCMiddleware(combined, app.ICAControllerKeeper)
    }

	icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)

	interchainQueriesModule := interchainqueries.NewAppModule(
		appCodec,
		keys[interchainqueriesmoduletypes.StoreKey],
		app.InterchainQueriesKeeper,
		app.AccountKeeper,
		app.BankKeeper,
	)
	interchainTxsModule := interchaintxs.NewAppModule(appCodec, app.InterchainTxsKeeper, app.AccountKeeper, app.BankKeeper)
	contractManagerModule := contractmanager.NewAppModule(appCodec, app.ContractManagerKeeper)
	ibcRateLimitmodule := ibcratelimit.NewAppModule(appCodec, app.RateLimitingICS4Wrapper.IbcratelimitKeeper, app.RateLimitingICS4Wrapper)
	ibcHooksModule := ibchooks.NewAppModule(app.AccountKeeper)

	transferModule := transferSudo.NewAppModule(app.TransferKeeper)

	app.ContractKeeper = wasmkeeper.NewDefaultPermissionKeeper(app.WasmKeeper)

	app.RateLimitingICS4Wrapper.ContractKeeper = app.ContractKeeper
	app.Ics20WasmHooks.ContractKeeper = &app.WasmKeeper

    ibcRouter.AddRoute(icacontrollertypes.SubModuleName, icaControllerStack).
        AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
        AddRoute(ibctransfertypes.ModuleName, app.TransferStack).
        AddRoute(interchaintxstypes.ModuleName, icaControllerStack).
        AddRoute(wasmtypes.ModuleName, wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper, app.IBCKeeper.ChannelKeeper)).
        AddRoute(ccvconsumertypes.ModuleName, consumerModule)
    app.IBCKeeper.SetRouter(ibcRouter)





	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))




	app.setupUpgradeStoreLoaders()

    app.mm = module.NewManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.ConsumerKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		upgrade.NewAppModule(&app.UpgradeKeeper, address.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,
		consumerModule,
		icaModule,
		app.PFMModule,
		interchainQueriesModule,
		interchainTxsModule,
		feeModule,
		feeBurnerModule,
		contractManagerModule,
		adminModule,
		ibcRateLimitmodule,
		ibcHooksModule,
		cronModule,
		autolp.NewAppModule(appCodec, app.AutolpKeeper),
        globalfee.NewAppModule(app.GlobalFeeKeeper, app.GetSubspace(globalfee.ModuleName), app.AppCodec(), app.keys[globalfee.ModuleName]),
        lockup.NewAppModule(*app.LockupKeeper, app.AccountKeeper, app.BankKeeper),
        incentives.NewAppModule(*app.IncentivesKeeper, app.AccountKeeper, app.EpochsKeeper),
        epochs.NewAppModule(*app.EpochsKeeper),





		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),

        crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),


	)

	app.mm.SetOrderPreBlockers(
		upgradetypes.ModuleName,
	)





    app.mm.SetOrderBeginBlockers(

		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		vestingtypes.ModuleName,
        ibchost.ModuleName,
        epochstypes.ModuleName,
        lockuptypes.ModuleName,
        incentivestypes.ModuleName,
		ibctransfertypes.ModuleName,
		authtypes.ModuleName,
		authz.ModuleName,
		banktypes.ModuleName,
		crisistypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		ccvconsumertypes.ModuleName,
		icatypes.ModuleName,
		interchainqueriesmoduletypes.ModuleName,
		interchaintxstypes.ModuleName,

		contractmanagermoduletypes.ModuleName,
		wasmtypes.ModuleName,
		feetypes.ModuleName,
		feeburnertypes.ModuleName,
		adminmoduletypes.ModuleName,
		ibcratelimittypes.ModuleName,
		ibchookstypes.ModuleName,
		pfmtypes.ModuleName,
		crontypes.ModuleName,


		globalfee.ModuleName,


		consensusparamtypes.ModuleName,


	)

    app.mm.SetOrderEndBlockers(

		crisistypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		authz.ModuleName,
		banktypes.ModuleName,
		slashingtypes.ModuleName,
		vestingtypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		ibchost.ModuleName,
		ibctransfertypes.ModuleName,
		ccvconsumertypes.ModuleName,
		icatypes.ModuleName,
		interchainqueriesmoduletypes.ModuleName,
		interchaintxstypes.ModuleName,
		contractmanagermoduletypes.ModuleName,
		wasmtypes.ModuleName,
		feetypes.ModuleName,
		feeburnertypes.ModuleName,
		adminmoduletypes.ModuleName,
		ibcratelimittypes.ModuleName,
		ibchookstypes.ModuleName,
		pfmtypes.ModuleName,
		crontypes.ModuleName,


		globalfee.ModuleName,


        consensusparamtypes.ModuleName,
        epochstypes.ModuleName,
        lockuptypes.ModuleName,
		incentivestypes.ModuleName,
		autolptypes.ModuleName,


	)






    app.mm.SetOrderInitGenesis(

		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		ibctransfertypes.ModuleName,
		authz.ModuleName,
		banktypes.ModuleName,
		vestingtypes.ModuleName,
		slashingtypes.ModuleName,
		crisistypes.ModuleName,
		ibchost.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		feegrant.ModuleName,
		ccvconsumertypes.ModuleName,
		icatypes.ModuleName,
		interchainqueriesmoduletypes.ModuleName,
		interchaintxstypes.ModuleName,
		contractmanagermoduletypes.ModuleName,
		wasmtypes.ModuleName,
		feetypes.ModuleName,
		feeburnertypes.ModuleName,
		adminmoduletypes.ModuleName,
		ibcratelimittypes.ModuleName,
		ibchookstypes.ModuleName,
		pfmtypes.ModuleName,
		crontypes.ModuleName,
		globalfee.ModuleName,




        consensusparamtypes.ModuleName,
        epochstypes.ModuleName,
        lockuptypes.ModuleName,
		incentivestypes.ModuleName,
		autolptypes.ModuleName,

	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	err = app.mm.RegisterServices(app.configurator)
	if err != nil {
		panic(fmt.Sprintf("failed to register services: %s", err))
	}

	app.setupUpgradeHandlers()


	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, nil, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		wasm.NewAppModule(appCodec, &app.WasmKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasmtypes.ModuleName)),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,
		ibcRateLimitmodule,
		consumerModule,
		icaModule,
		app.PFMModule,
		interchainQueriesModule,
		interchainTxsModule,
		feeBurnerModule,
		cronModule,

	)
	app.sm.RegisterStoreDecoders()


	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)


	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)


	baseLane := app.CreateLanes()
	mempool, err := blocksdk.NewLanedMempool(app.Logger(), []blocksdk.Lane{baseLane})
	if err != nil {
		panic(err)
	}


	app.SetMempool(mempool)


	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				FeegrantKeeper:  app.FeeGrantKeeper,
				SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			BankKeeper:            app.BankKeeper,
			AccountKeeper:         app.AccountKeeper,
			IBCKeeper:             app.IBCKeeper,
			GlobalFeeKeeper:       app.GlobalFeeKeeper,
			WasmConfig:            &wasmConfig,
			TXCounterStoreService: runtime.NewKVStoreService(keys[wasmtypes.StoreKey]),
			ConsumerKeeper:        app.ConsumerKeeper,

		},
		app.Logger(),
	)
	if err != nil {
		panic(err)
	}
	app.SetAnteHandler(anteHandler)

















































































	parityCheckTx := checktx.NewMempoolParityCheckTx(
		app.Logger(),
		mempool,
		app.GetTxConfig().TxDecoder(),
		app.BaseApp.CheckTx,
		app.BaseApp,
	)

	app.SetCheckTx(parityCheckTx.CheckTx())



































































	if manager := app.SnapshotManager(); manager != nil {
		err := manager.RegisterExtensions(
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
		)
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %s", err))
		}
	}

	if loadLatest {
		app.LoadLatest()
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	app.ScopedWasmKeeper = scopedWasmKeeper
	app.ScopedInterTxKeeper = scopedInterTxKeeper
	app.ScopedCCVConsumerKeeper = scopedCCVConsumerKeeper

	return app
}

func (app *App) LoadLatest() {
	if err := app.LoadLatestVersion(); err != nil {
		tmos.Exit(err.Error())
	}

	ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})


	if err := app.WasmKeeper.InitializePinnedCodes(ctx); err != nil {
		tmos.Exit(fmt.Sprintf("failed initialize pinned codes %s", err))
	}
}

func (app *App) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrd info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrd := range Upgrades {
		upgrd := upgrd
		if upgradeInfo.Name == upgrd.UpgradeName {
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &upgrd.StoreUpgrades))
		}
	}
}

func (app *App) setupUpgradeHandlers() {
	for _, upgrd := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrd.UpgradeName,
			upgrd.CreateUpgradeHandler(
				app.mm,
				app.configurator,
				&upgrades.UpgradeKeepers{
					AccountKeeper:       app.AccountKeeper,
					FeeBurnerKeeper:     app.FeeBurnerKeeper,
					CronKeeper:          app.CronKeeper,
					IcqKeeper:           app.InterchainQueriesKeeper,

					SlashingKeeper:      app.SlashingKeeper,
					ParamsKeeper:        app.ParamsKeeper,
					CapabilityKeeper:    app.CapabilityKeeper,

					ContractManager:     app.ContractManagerKeeper,
					AdminModule:         app.AdminmoduleKeeper,
					ConsensusKeeper:     &app.ConsensusParamsKeeper,
					ConsumerKeeper:      &app.ConsumerKeeper,




					IbcRateLimitKeeper:  app.RateLimitingICS4Wrapper.IbcratelimitKeeper,
					ChannelKeeper:       &app.IBCKeeper.ChannelKeeper,
					TransferKeeper:      app.TransferKeeper.Keeper,
					GlobalFeeSubspace:   app.GetSubspace(globalfee.ModuleName),
					CcvConsumerSubspace: app.GetSubspace(ccvconsumertypes.ModuleName),
				},
				app,
				app.AppCodec(),
			),
		)
	}
}

func (app *App) AutoCliOpts() autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range app.mm.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(app.mm.Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	}
}





func (app *App) CheckTx(req *abci.RequestCheckTx) (*abci.ResponseCheckTx, error) {
	return app.checkTxHandler(req)
}


func (app *App) SetCheckTx(handler checktx.CheckTx) {
	app.checkTxHandler = handler
}


func (app *App) Name() string { return app.BaseApp.Name() }


func (app *App) GetBaseApp() *baseapp.BaseApp { return app.BaseApp }


func (app *App) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}


func (app *App) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}


func (app *App) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		return nil, err
	}
	if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap()); err != nil {
		return nil, err
	}
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}



func (app *App) TestInitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}




	err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	if err != nil {
		return nil, fmt.Errorf("failed to set module version map: %w", err)
	}
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}


func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}


func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}



func (app *App) BlockedAddrs() map[string]bool {



	bankBlockedAddrs := app.ModuleAccountAddrs()
	delete(bankBlockedAddrs, authtypes.NewModuleAddress(
		ccvconsumertypes.ConsumerToSendToProviderName).String())

	return bankBlockedAddrs
}





func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}





func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}




func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}




func (app *App) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}




func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}




func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}



func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)


	if apiConfig.Swagger {
		app.RegisterSwaggerUI(apiSvr)
	}
}


func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}


func (app *App) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}


func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName).WithKeyTable(authtypes.ParamKeyTable())
	paramsKeeper.Subspace(banktypes.ModuleName).WithKeyTable(banktypes.ParamKeyTable())
	paramsKeeper.Subspace(slashingtypes.ModuleName).WithKeyTable(slashingtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName).WithKeyTable(crisistypes.ParamKeyTable())
	paramsKeeper.Subspace(ibctransfertypes.ModuleName).WithKeyTable(ibctransfertypes.ParamKeyTable())

	keyTable := ibcclienttypes.ParamKeyTable()
	keyTable.RegisterParamSet(&ibcconnectiontypes.Params{})
	paramsKeeper.Subspace(ibchost.ModuleName).WithKeyTable(keyTable)

	paramsKeeper.Subspace(icacontrollertypes.SubModuleName).WithKeyTable(icacontrollertypes.ParamKeyTable())
	paramsKeeper.Subspace(icahosttypes.SubModuleName).WithKeyTable(icahosttypes.ParamKeyTable())

	paramsKeeper.Subspace(pfmtypes.ModuleName).WithKeyTable(pfmtypes.ParamKeyTable())

	paramsKeeper.Subspace(globalfee.ModuleName).WithKeyTable(globalfeetypes.ParamKeyTable())

	paramsKeeper.Subspace(ccvconsumertypes.ModuleName).WithKeyTable(ccv.ParamKeyTable())


	paramsKeeper.Subspace(crontypes.StoreKey).WithKeyTable(crontypes.ParamKeyTable())
	paramsKeeper.Subspace(feeburnertypes.StoreKey).WithKeyTable(feeburnertypes.ParamKeyTable())
	paramsKeeper.Subspace(feetypes.StoreKey).WithKeyTable(feetypes.ParamKeyTable())
	paramsKeeper.Subspace(interchainqueriesmoduletypes.StoreKey).WithKeyTable(interchainqueriesmoduletypes.ParamKeyTable())
	paramsKeeper.Subspace(interchaintxstypes.StoreKey).WithKeyTable(interchaintxstypes.ParamKeyTable())
	paramsKeeper.Subspace(lockuptypes.ModuleName).WithKeyTable(lockuptypes.ParamKeyTable())
	paramsKeeper.Subspace(incentivestypes.ModuleName).WithKeyTable(incentivestypes.ParamKeyTable())

    // autolp currently has no param subspace



	return paramsKeeper
}


func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

func (app *App) RegisterSwaggerUI(apiSvr *api.Server) {
	staticSubDir, err := fs.Sub(docs.Docs, "static")
	if err != nil {
		app.Logger().Error(fmt.Sprintf("failed to register swagger-ui route: %s", err))
		return
	}

	staticServer := http.FileServer(http.FS(staticSubDir))
	apiSvr.Router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}




func (app *App) GetTxConfig() client.TxConfig {
	return app.encodingConfig.TxConfig
}


func (app *App) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}


func (app *App) GetStakingKeeper() ibctestingtypes.StakingKeeper {
	return app.ConsumerKeeper
}


func (app *App) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}


func (app *App) GetConsumerKeeper() ccvconsumerkeeper.Keeper {
	return app.ConsumerKeeper
}

func (app *App) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}



func overrideWasmVariables() {

	wasmtypes.MaxWasmSize = 1_677_722
	wasmtypes.MaxProposalWasmSize = wasmtypes.MaxWasmSize
}













func (app *App) WireICS20PreWasmKeeper(
	appCodec codec.Codec,
) {

	app.PFMKeeper = pfmkeeper.NewKeeper(
		appCodec,
		app.keys[pfmtypes.StoreKey],
		app.TransferKeeper.Keeper,
		app.IBCKeeper.ChannelKeeper,
		app.FeeBurnerKeeper,
		&app.BankKeeper,
		app.IBCKeeper.ChannelKeeper,
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)

	wasmHooks := ibchooks.NewWasmHooks(nil, sdk.GetConfig().GetBech32AccountAddrPrefix())
	app.Ics20WasmHooks = &wasmHooks
	app.HooksICS4Wrapper = ibchooks.NewICS4Middleware(
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&wasmHooks,
	)

	ibcratelimitKeeper := ibcratelimitkeeper.NewKeeper(appCodec, app.keys[ibcratelimittypes.ModuleName], authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String())

	rateLimitingICS4Wrapper := ibcratelimit.NewICS4Middleware(
		app.HooksICS4Wrapper,
		&app.AccountKeeper,

		nil,
		&app.BankKeeper,
		&ibcratelimitKeeper,
	)
	app.RateLimitingICS4Wrapper = &rateLimitingICS4Wrapper


	app.TransferKeeper = wrapkeeper.NewKeeper(
		appCodec,
		app.keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.RateLimitingICS4Wrapper,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		&app.BankKeeper,
		app.ScopedTransferKeeper,
		app.FeeKeeper,
		contractmanager.NewSudoLimitWrapper(app.ContractManagerKeeper, &app.WasmKeeper),
		authtypes.NewModuleAddress(adminmoduletypes.ModuleName).String(),
	)

	app.PFMKeeper.SetTransferKeeper(app.TransferKeeper.Keeper)



	var ibcStack ibcporttypes.IBCModule = packetforward.NewIBCMiddleware(
		transferSudo.NewIBCModule(
			app.TransferKeeper,
			contractmanager.NewSudoLimitWrapper(app.ContractManagerKeeper, &app.WasmKeeper),

		),
		app.PFMKeeper,
		0,
		pfmkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		pfmkeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)

	ibcStack = gmpmiddleware.NewIBCMiddleware(ibcStack)

	rateLimitingTransferModule := ibcratelimit.NewIBCModule(ibcStack, app.RateLimitingICS4Wrapper)


	hooksTransferModule := ibchooks.NewIBCMiddleware(&rateLimitingTransferModule, &app.HooksICS4Wrapper)
	app.TransferStack = &hooksTransferModule
}
