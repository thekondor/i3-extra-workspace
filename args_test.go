package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type argsSut struct {
	args
}

func TestArgumentsParsing(t *testing.T) {
	var fatalMsg string
	fatalf := func(fmtArg string, args ...interface{}) {
		fatalMsg = fmt.Sprintf(fmtArg, args...)
	}

	t.Run(`empty arguments`, func(t *testing.T) {
		fatalMsg = ""
		sut := args{fatalfFn: fatalf}

		sut.mustParse([]string{"app.bin"})
		require.NotEmpty(t, fatalMsg)
		assert.Equal(t, "No argument(s) provided", fatalMsg)
	})

	t.Run(`empty action`, func(t *testing.T) {
		fatalMsg = ""
		sut := args{fatalfFn: fatalf}

		sut.mustParse([]string{"app.bin", "-ws-prefix", "+"})
		require.NotEmpty(t, fatalMsg)
		assert.Equal(t, "Failed to parsed arguments: Action must be provided as a single latest argument", fatalMsg)
	})

	t.Run(`no prefix`, func(t *testing.T) {
		sut := args{}
		sut.mustParse([]string{"app.bin", "flip"})

		assert.Equal(t, "+", sut.ExtensionPrefix)
		assert.Equal(t, "flip", sut.WsAction)
	})

	t.Run(`prefix & action`, func(t *testing.T) {
		sut := args{}
		sut.mustParse([]string{"app.bin", "-ws-prefix", "#", "toggle"})

		assert.Equal(t, "#", sut.ExtensionPrefix)
		assert.Equal(t, "toggle", sut.WsAction)
	})
}
