package web

import (
	"embed"
)

//go:embed  index.html js/* css/*
var Content embed.FS
