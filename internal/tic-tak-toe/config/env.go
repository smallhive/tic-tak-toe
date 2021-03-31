package config

import (
	"os"
	"sort"
	"strings"
)

const (
	EnvPrefix = "TTT"

	Addr      = "addr"
	RedisAddr = "redis_addr"
)

type Config struct {
	Addr      string
	RedisAddr string
}

func Load() *Config {
	envs := os.Environ()
	sort.Strings(envs)

	var c = &Config{
		Addr:      ":8080",
		RedisAddr: ":6379",
	}

	for _, key := range envs {
		if strings.HasPrefix(key, EnvPrefix) {
			k, v := parseKey(EnvPrefix, key)
			// lower key because ENV_KEY_NAME
			k = strings.ToLower(k)

			switch k {
			case Addr:
				c.Addr = v
			case RedisAddr:
				c.RedisAddr = v
			}
		}
	}

	return c
}

func parseKey(prefix, env string) (string, string) {
	skipSince := len(prefix) + 1
	pair := strings.SplitN(env, "=", 2)
	if len(pair) < 2 {
		return env[skipSince:], ""
	}

	return pair[0][skipSince:], pair[1]
}
