package assets

import "embed"

//go:embed css/* img/* js/* fonts/*
var Assets embed.FS
