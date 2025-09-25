package autolp

import (
    "errors"

    sdk "github.com/cosmos/cosmos-sdk/types"
    capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
    channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
    porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
    ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

    "github.com/maany-xyz/maany-app/x/autolp/keeper"
)

// IBCModule implements ICS26 callbacks for autolp. It only observes events and emits them.
type IBCModule struct{ keeper keeper.Keeper }

func NewIBCModule(k keeper.Keeper) IBCModule { return IBCModule{keeper: k} }

// OnChanOpenInit not used.
func (im IBCModule) OnChanOpenInit(_ sdk.Context, _ channeltypes.Order, _ []string, _ string, _ string, _ *capabilitytypes.Capability, _ channeltypes.Counterparty, version string) (string, error) {
    return version, nil
}

// OnChanOpenTry not used.
func (im IBCModule) OnChanOpenTry(_ sdk.Context, _ channeltypes.Order, _ []string, _, _ string, _ *capabilitytypes.Capability, _ channeltypes.Counterparty, _ string) (string, error) {
    return "", nil
}

// OnChanOpenAck: record/emit open-ack
func (im IBCModule) OnChanOpenAck(ctx sdk.Context, portID, channelID, counterPartyChannelID, counterpartyVersion string) error {
    return im.keeper.HandleChanOpenAck(ctx, portID, channelID, counterPartyChannelID, counterpartyVersion)
}

// OnChanOpenConfirm not used.
func (im IBCModule) OnChanOpenConfirm(_ sdk.Context, _, _ string) error { return nil }

// OnChanCloseInit not used.
func (im IBCModule) OnChanCloseInit(_ sdk.Context, _, _ string) error { return nil }

// OnChanCloseConfirm not used.
func (im IBCModule) OnChanCloseConfirm(_ sdk.Context, _, _ string) error { return nil }

// OnRecvPacket is not supported.
func (im IBCModule) OnRecvPacket(_ sdk.Context, _ channeltypes.Packet, _ sdk.AccAddress) ibcexported.Acknowledgement {
    return channeltypes.NewErrorAcknowledgement(errors.New("autolp: recv not supported"))
}

// OnAcknowledgementPacket: emit ack event
func (im IBCModule) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
    return im.keeper.HandleAcknowledgement(ctx, packet, acknowledgement, relayer)
}

// OnTimeoutPacket: emit timeout event
func (im IBCModule) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
    return im.keeper.HandleTimeout(ctx, packet, relayer)
}

// Combine two IBCModules by calling both in order; stop if the first returns error.
type combinedIBC struct{ a, b porttypes.IBCModule }

func CombineIBCModules(a, b porttypes.IBCModule) porttypes.IBCModule { return combinedIBC{a: a, b: b} }

func (c combinedIBC) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, conns []string, portID, channelID string, cap *capabilitytypes.Capability, cp channeltypes.Counterparty, version string) (string, error) {
    v, err := c.a.OnChanOpenInit(ctx, order, conns, portID, channelID, cap, cp, version)
    if err != nil { return v, err }
    return c.b.OnChanOpenInit(ctx, order, conns, portID, channelID, cap, cp, v)
}
func (c combinedIBC) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, conns []string, portID, channelID string, cap *capabilitytypes.Capability, cp channeltypes.Counterparty, version string) (string, error) {
    v, err := c.a.OnChanOpenTry(ctx, order, conns, portID, channelID, cap, cp, version)
    if err != nil { return v, err }
    return c.b.OnChanOpenTry(ctx, order, conns, portID, channelID, cap, cp, v)
}
func (c combinedIBC) OnChanOpenAck(ctx sdk.Context, portID, channelID, counterPartyChannelID, counterpartyVersion string) error {
    if err := c.a.OnChanOpenAck(ctx, portID, channelID, counterPartyChannelID, counterpartyVersion); err != nil { return err }
    return c.b.OnChanOpenAck(ctx, portID, channelID, counterPartyChannelID, counterpartyVersion)
}
func (c combinedIBC) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
    if err := c.a.OnChanOpenConfirm(ctx, portID, channelID); err != nil { return err }
    return c.b.OnChanOpenConfirm(ctx, portID, channelID)
}
func (c combinedIBC) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
    if err := c.a.OnChanCloseInit(ctx, portID, channelID); err != nil { return err }
    return c.b.OnChanCloseInit(ctx, portID, channelID)
}
func (c combinedIBC) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
    if err := c.a.OnChanCloseConfirm(ctx, portID, channelID); err != nil { return err }
    return c.b.OnChanCloseConfirm(ctx, portID, channelID)
}
func (c combinedIBC) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
    // if first returns success ack, still forward to b? For safety, return first and ignore b.
    return c.a.OnRecvPacket(ctx, packet, relayer)
}
func (c combinedIBC) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
    if err := c.a.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer); err != nil { return err }
    return c.b.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}
func (c combinedIBC) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
    if err := c.a.OnTimeoutPacket(ctx, packet, relayer); err != nil { return err }
    return c.b.OnTimeoutPacket(ctx, packet, relayer)
}

