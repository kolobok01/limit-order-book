package fix

// Config holds the FIX client configuration
type Config struct {
	AccountID     string `json:"account_id" yaml:"account_id"`
	TargetCompID  string `json:"target_comp_id" yaml:"target_comp_id"`
	APIKey        string `json:"api_key" yaml:"api_key"`
	APISecret     string `json:"api_secret" yaml:"api_secret"`
	APIPassphrase string `json:"api_passphrase" yaml:"api_passphrase"`
	PortfolioID   string `json:"portfolio_id" yaml:"portfolio_id"`
}
