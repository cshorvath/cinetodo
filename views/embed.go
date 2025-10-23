package views

import "embed"

// Files exposes the embedded templates.
//
//go:embed */*.html
var Files embed.FS
