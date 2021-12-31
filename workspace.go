package main

import (
	"fmt"
	"strings"

	i3 "go.i3wm.org/i3/v4"
)

const (
	extensionPrefix = `_`
)

type workspaces struct {
	i3wss []i3.Workspace
}

func (wss workspaces) GetFocused() (workspace, error) {
	for _, ws := range wss.i3wss {
		if ws.Focused {
			return workspace{ws}, nil
		}
	}

	return workspace{}, fmt.Errorf("No focused workspace among %d available workspaces", len(wss.i3wss))
}

func listWorkspaces() (workspaces, error) {
	i3wss, err := i3.GetWorkspaces()
	if err != nil {
		return workspaces{}, fmt.Errorf("failed to list workspaces: %w", err)
	}

	return workspaces{i3wss: i3wss}, nil
}

type workspace struct {
	i3.Workspace
}

func (ws workspace) id() string {
	return fmt.Sprintf("num=%d:name=%s", ws.Num, ws.Name)
}

func (ws workspace) String() string {
	return ws.id()
}

var (
	i3msg func(string) ([]i3.CommandResult, error) = i3.RunCommand
)

func (ws workspace) Focus() error {
	if cr, err := i3msg(fmt.Sprintf("workspace %s", ws.Name)); err != nil {
		return fmt.Errorf("failed to focus workspace '%s': %+v", ws.id(), cr)
	}
	return nil
}

func (ws workspace) BorrowFocusedContainer() error {
	if cr, err := i3msg(fmt.Sprintf("move container to workspace %s", ws.Name)); err != nil {
		return fmt.Errorf("failed to focus workspace '%s': %+v", ws.id(), cr)
	}
	return nil
}

func (ws workspace) IsExtension() bool {
	return ws.hasExtensionPrefix()
}

func (ws workspace) IsMain() bool {
	return !ws.hasExtensionPrefix()
}

func (ws workspace) hasExtensionPrefix() bool {
	return strings.HasPrefix(ws.Name, extensionPrefix)
}

func (ws workspace) Main() (workspace, error) {
	if !ws.hasExtensionPrefix() {
		return workspace{}, fmt.Errorf("workspace %s is already main", ws.id())
	}

	name := strings.TrimPrefix(ws.Name, extensionPrefix)
	return workspace{Workspace: i3.Workspace{Name: name}}, nil
}

func (ws workspace) Extension() (workspace, error) {
	if ws.hasExtensionPrefix() {
		return workspace{}, fmt.Errorf("workspace %s is already an extension", ws.id())
	}

	name := fmt.Sprintf("%s%s", extensionPrefix, ws.Name)
	return workspace{Workspace: i3.Workspace{Name: name}}, nil
}
