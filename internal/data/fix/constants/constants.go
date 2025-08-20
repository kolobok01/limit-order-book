package constants

import "github.com/quickfixgo/quickfix"

const (
	MsgTypeNew    = "D" // New Order
	MsgTypeStatus = "H" // Status
	MsgTypeCancel = "F" // Cancel
	MsgTypeLogon  = "A" // Logon

	FixTimeFormat = "20060102-15:04:05.000"

	DefaultTargetCompID = "COIN"

	OrdTypeLimit     = "LIMIT"
	OrdTypeMarket    = "MARKET"
	OrdTypeVwap      = "VWAP"
	SideBuy          = "BUY"
	SideSell         = "SELL"
	OrdTypeLimitFix  = "2" // Limit order type
	OrdTypeMarketFix = "1" // Market order type
	OrdTypeVwapFix   = "2" // VWAP order type (uses limit type)

	TimeInForceDay = "1" // Day
	TimeInForceIoc = "3" // Immediate or Cancel
	TimeInForceGtd = "6" // Good Till Date

	TargetStrategyLimit  = "L" // Limit strategy
	TargetStrategyMarket = "M" // Market strategy
	TargetStrategyVwap   = "V" // VWAP strategy

	SideBuyFix  = "1" // Buy side
	SideSellFix = "2" // Sell side
)

// FIX Tags
var (
	TagAccount           = quickfix.Tag(1)
	TagClOrdId           = quickfix.Tag(11)
	TagOrderId           = quickfix.Tag(37)
	TagOrderQty          = quickfix.Tag(38)
	TagOrdType           = quickfix.Tag(40)
	TagOrigClOrdId       = quickfix.Tag(41)
	TagTargetStrategy    = quickfix.Tag(847)
	TagPx                = quickfix.Tag(44)
	TagExecInst          = quickfix.Tag(18)
	TagSenderCompId      = quickfix.Tag(49)
	TagSendingTime       = quickfix.Tag(52)
	TagSide              = quickfix.Tag(54)
	TagSymbol            = quickfix.Tag(55)
	TagTargetCompId      = quickfix.Tag(56)
	TagTimeInForce       = quickfix.Tag(59)
	TagHmac              = quickfix.Tag(96)
	TagMsgType           = quickfix.Tag(35)
	TagExecType          = quickfix.Tag(150)
	TagCashOrderQty      = quickfix.Tag(152)
	TagPassword          = quickfix.Tag(554)
	TagDropCopyFlag      = quickfix.Tag(9406)
	TagAccessKey         = quickfix.Tag(9407)
	TagStartTime         = quickfix.Tag(168)
	TagExpireTime        = quickfix.Tag(126)
	TagParticipationRate = quickfix.Tag(849)
)
