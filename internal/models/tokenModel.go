package models

import "time"

type TokenCheckArgs struct {
	ClientTime time.Time `json:"client_time"`
}

type TokenCheckResult struct {
	ClientTime time.Time `json:"client_time"`
	ServerTime time.Time `json:"server_time"`
	Result     Result    `json:"result"`
}