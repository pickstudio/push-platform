//go:build tools
// +build tools

package pushplatform

import (
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
	_ "github.com/golang/mock/mockgen"
)
