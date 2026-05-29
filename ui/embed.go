package ui

import "embed"

// Assets contains the built frontend files
//
//go:embed all:dist
var Assets embed.FS
