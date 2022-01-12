package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	i3 "go.i3wm.org/i3/v4"
)

func newWs(t *testing.T, name string) workspace {
	t.Helper()
	return workspace{i3.Workspace{Name: name}}
}

func TestWorkspace(t *testing.T) {
	t.Run(`i3msg`, func(t *testing.T) {
		original := i3msg
		defer func() { i3msg = original }()

		var cmd string
		i3msg = func(i3cmd string) ([]i3.CommandResult, error) {
			cmd = i3cmd
			return nil, nil
		}

		t.Run(`focus`, func(t *testing.T) {
			ws := newWs(t, `1:work`)

			require.NoError(t, ws.Focus())
			assert.Equal(t, `workspace 1:work`, cmd)

			nextWs, err := ws.Extension()
			require.NoError(t, err)
			require.NoError(t, nextWs.Focus())
			assert.Equal(t, `workspace _1:work`, cmd)
		})

		t.Run(`flip container`, func(t *testing.T) {
			ws := newWs(t, `1:work`)

			require.NoError(t, ws.FlipFocusedContainer())
			assert.Equal(t, `move container to workspace 1:work`, cmd)

			nextWs, err := ws.Extension()
			require.NoError(t, err)
			require.NoError(t, nextWs.FlipFocusedContainer())
			assert.Equal(t, `move container to workspace _1:work`, cmd)

		})
	})

	t.Run(`main`, func(t *testing.T) {
		ws := newWs(t, "1:work")

		t.Run(`predicates`, func(t *testing.T) {
			assert.True(t, ws.IsMain())
			assert.False(t, ws.IsExtension())
		})

		t.Run(`fails on another main`, func(t *testing.T) {
			_, err := ws.Main()
			assert.EqualError(t, err, "workspace num=0:name=1:work is already main")
		})

		t.Run(`returns extension`, func(t *testing.T) {
			wsExt, err := ws.Extension()

			require.NoError(t, err)
			assert.True(t, wsExt.IsExtension() && !wsExt.IsMain())
			assert.Equal(t, "_1:work", wsExt.Name)
		})
	})

	t.Run(`extension`, func(t *testing.T) {
		ws := newWs(t, "_1:work")

		t.Run(`predicates`, func(t *testing.T) {
			assert.True(t, ws.IsExtension())
			assert.False(t, ws.IsMain())
		})

		t.Run(`fails on another extension`, func(t *testing.T) {
			_, err := ws.Extension()
			assert.EqualError(t, err, "workspace num=0:name=_1:work is already an extension")
		})

		t.Run(`returns main`, func(t *testing.T) {
			wsMain, err := ws.Main()

			require.NoError(t, err)
			assert.True(t, wsMain.IsMain() && !wsMain.IsExtension())
			assert.Equal(t, "1:work", wsMain.Name)
		})
	})
}
