package ui

import "mime"

//go:generate pnpm install
//go:generate pnpm run build

func init() {
	mime.AddExtensionType(".js", "application/javascript")
}
