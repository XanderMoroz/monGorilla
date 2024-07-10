package models

type TokenCheckArgs struct {
	Token string `json:"token"`
}

type TokenCheckResult struct {
	Token   string `json:"token"`
	Expired string `json:"expired"`
	Result  Result `json:"result"`
}
