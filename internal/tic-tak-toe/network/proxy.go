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

func NewPlayerProxyConfig(id string) *PlayerProxyConfig {
	return &PlayerProxyConfig{
		UserChanName:    fmt.Sprintf("user:%s", id),
		ControlChanName: fmt.Sprintf("user:control:%s", id),
	}
}

type GameProxyConfig struct {
	ChanName string
}

func NewGameProxyConfig(id string) *GameProxyConfig {
	return &GameProxyConfig{
		ChanName: fmt.Sprintf("game:%s", id),
	}
}
