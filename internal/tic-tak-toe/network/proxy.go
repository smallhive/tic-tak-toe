package network

import (
	"fmt"
)

type PlayerProxyConfig struct {
	// Chan name for game commands
	UserChanName string

	// Chan name for game control/system commands
	ControlChanName string
}

func NewPlayerProxyConfig(id int64) *PlayerProxyConfig {
	return &PlayerProxyConfig{
		UserChanName:    fmt.Sprintf("user:%d", id),
		ControlChanName: fmt.Sprintf("user:control:%d", id),
	}
}

type GameProxyConfig struct {
	ChanName string
}

func NewGameProxyConfig(id int64) *GameProxyConfig {
	return &GameProxyConfig{
		ChanName: fmt.Sprintf("game:%d", id),
	}
}
