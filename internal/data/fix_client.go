package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kolobublik/limit-order-book/internal/conf"
	"github.com/kolobublik/limit-order-book/internal/data/fix"
	"github.com/quickfixgo/quickfix"
)

func NewFIXClient(c *conf.Data, logger log.Logger) (*fix.Client, error) {
	settings := quickfix.NewSettings()

	// Add default settings
	globalSettings := settings.GlobalSettings()
	globalSettings.Set("ConnectionType", "initiator")
	globalSettings.Set("ReconnectInterval", "60")
	globalSettings.Set("FileLogPath", "log")
	globalSettings.Set("FileStorePath", "store")
	globalSettings.Set("StartTime", "00:00:00")
	globalSettings.Set("EndTime", "00:00:00")
	globalSettings.Set("UseDataDictionary", "N")
	globalSettings.Set("ValidateUserDefinedFields", "N")
	globalSettings.Set("ValidateFieldsOutOfOrder", "N")
	globalSettings.Set("ValidateFieldsHaveValues", "N")
	globalSettings.Set("ValidateIncomingMessage", "N")

	// Add session settings
	sessionSettings := quickfix.NewSessionSettings()
	sessionSettings.Set("BeginString", "FIX.4.2")
	sessionSettings.Set("SenderCompID", c.Fix.AccountId)
	sessionSettings.Set("TargetCompID", c.Fix.TargetCompId)
	sessionSettings.Set("SocketConnectHost", "fix.prime.coinbase.com")
	sessionSettings.Set("SocketConnectPort", "4198")
	sessionSettings.Set("HeartBtInt", "30")
	sessionSettings.Set("ResetOnLogon", "Y")
	sessionSettings.Set("ResetOnLogout", "Y")
	sessionSettings.Set("ResetOnDisconnect", "Y")

	settings.AddSession(sessionSettings)

	config := &fix.Config{
		AccountID:     c.Fix.AccountId,
		TargetCompID:  c.Fix.TargetCompId,
		APIKey:        c.Fix.ApiKey,
		APISecret:     c.Fix.ApiSecret,
		APIPassphrase: c.Fix.ApiPassphrase,
		PortfolioID:   c.Fix.PortfolioId,
	}

	return fix.NewClient(settings, config)
}
