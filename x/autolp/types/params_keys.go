package types

var (
    KeyAllowedSubmitters = []byte("AllowedSubmitters")
    KeyEnsureIcaEnabled  = []byte("EnsureIcaEnabled")
    KeyConnectionId      = []byte("ConnectionId")
    KeyIcaId             = []byte("IcaId")
    KeyEnsureIcaInterval = []byte("EnsureIcaInterval")

    KeyBaseHostDenom           = []byte("BaseHostDenom")
    KeyQuoteHostDenom          = []byte("QuoteHostDenom")
    KeyCreatePoolRequired      = []byte("CreatePoolRequired")
    KeyAutoCreatePoolEnabled   = []byte("AutoCreatePoolEnabled")
    KeyICQUpdatePeriod         = []byte("ICQUpdatePeriod")
    KeyICATxTimeoutSeconds     = []byte("ICATxTimeoutSeconds")
    KeyTransferTimeoutSeconds  = []byte("TransferTimeoutSeconds")
    KeyHostAccountHRP          = []byte("HostAccountHRP")
    KeyBaseLocalDenom          = []byte("BaseLocalDenom")
    KeyBaseDexChannel          = []byte("BaseDexChannel")
    KeyQuoteLocalDenom         = []byte("QuoteLocalDenom")
    KeyQuoteDexChannel         = []byte("QuoteDexChannel")
)
