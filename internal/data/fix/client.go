package fix

import (
	"fmt"
	"log"
	"time"

	"github.com/kolobublik/limit-order-book/internal/data/fix/builder"
	"github.com/kolobublik/limit-order-book/internal/data/fix/constants"
	"github.com/kolobublik/limit-order-book/internal/data/fix/model"
	"github.com/quickfixgo/quickfix"
)

type Client struct {
	*quickfix.MessageRouter
	initiator *quickfix.Initiator
	msgChan   chan *quickfix.Message
	config    *Config
}

func NewClient(settings *quickfix.Settings, config *Config) (*Client, error) {
	client := &Client{
		MessageRouter: quickfix.NewMessageRouter(),
		msgChan:       make(chan *quickfix.Message, 100),
		config:        config,
	}

	// Update session settings with configuration values
	for sessionID := range settings.SessionSettings() {
		sessionSettings := settings.SessionSettings()[sessionID]
		sessionSettings.Set("SenderCompID", config.AccountID)
		sessionSettings.Set("TargetCompID", config.TargetCompID)
	}

	// Register message handlers
	client.MessageRouter.AddRoute("8", "FIX.4.2", client.onExecutionReport) // ExecutionReport
	client.MessageRouter.AddRoute("3", "FIX.4.2", client.onReject)          // Reject

	var err error
	client.initiator, err = quickfix.NewInitiator(
		client,
		quickfix.NewMemoryStoreFactory(),
		settings,
		quickfix.NewNullLogFactory(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create FIX initiator: %v", err)
	}

	return client, nil
}

func (c *Client) Start() error {
	return c.initiator.Start()
}

func (c *Client) Stop() {
	c.initiator.Stop()
}

func (c *Client) OnCreate(sessionID quickfix.SessionID) {
	log.Printf("Session created: %v\n", sessionID)
}

func (c *Client) OnLogon(sessionID quickfix.SessionID) {
	log.Printf("Logged on: %v\n", sessionID)
}

func (c *Client) OnLogout(sessionID quickfix.SessionID) {
	log.Printf("Logged out: %v\n", sessionID)
}

func (c *Client) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	msgType, err := msg.Header.GetString(constants.TagMsgType)
	if err != nil {
		log.Printf("Error getting message type: %v\n", err)
		return
	}

	if msgType == "A" { // Logon message
		ts := time.Now().UTC().Format(constants.FixTimeFormat)
		builder.BuildLogon(
			&msg.Body,
			ts,
			c.config.APIKey,
			c.config.APISecret,
			c.config.APIPassphrase,
			c.config.TargetCompID,
			c.config.PortfolioID,
		)
	}
}

func (c *Client) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	log.Printf("Sending message: %v\n", msg)
	return nil
}

func (c *Client) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	log.Printf("Admin message received: %v\n", msg)
	return nil
}

func (c *Client) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return c.Route(msg, sessionID)
}

func (c *Client) onExecutionReport(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	c.msgChan <- msg
	return nil
}

func (c *Client) onReject(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	c.msgChan <- msg
	return nil
}

func (c *Client) PlaceOrder(
	symbol, ordType, side, qtyType, qty, price, portfolio string, vwapParams ...string,
) error {
	msg, err := builder.BuildNew(symbol, ordType, side, qtyType, qty, price, portfolio, vwapParams...)
	if err != nil {
		return fmt.Errorf("failed to build order message: %v", err)
	}

	return quickfix.Send(msg)
}

func (c *Client) GetOrderStatus(info model.OrderInfo, portfolio string) error {
	msg := builder.BuildStatus(info.ClOrdId, info.OrderId, info.Side, info.Symbol)
	return quickfix.Send(msg)
}

func (c *Client) CancelOrder(info model.OrderInfo, portfolio string) error {
	msg := builder.BuildCancel(info, portfolio)
	return quickfix.Send(msg)
}

func (c *Client) Messages() chan *quickfix.Message {
	return c.msgChan
}
