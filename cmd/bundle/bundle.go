package main

import (
	"fmt"

	"github.com/Shopify/themekit/src/static"
)

func main() {
	if err := static.Bundle("theme-template", "cmd/static/theme.go"); err != nil {
		fmt.Println(err)
	}
}
