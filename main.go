package main

import (
	"fmt"
	"log"
	"os"
)

type wsExtensionAction string

const (
	toggleWsExtensionAction        wsExtensionAction = "toggle"
	flipContainerWsExtensionAction wsExtensionAction = "flip"
	defaultExtensionPrefix         string            = "+"
)

func getNextWorkspace(ws workspace) (nextWs workspace, err error) {
	if ws.IsExtension() {
		nextWs, err = ws.Main()
	} else if ws.IsMain() {
		nextWs, err = ws.Extension()
	} else {
		err = fmt.Errorf("the workspace '%d:%s' is neither main nor extension", ws.Num, ws.Name)
	}
	return
}

func flipContainer(extensionPrefix string) {
	nextWs := mustGetNextWs(extensionPrefix)
	if err := nextWs.FlipFocusedContainer(); err != nil {
		log.Fatalf("failed to flip focused container to workspace '%s': %v", nextWs, err)
	}
}

func toggle(extensionPrefix string) {
	nextWs := mustGetNextWs(extensionPrefix)
	if err := nextWs.Focus(); err != nil {
		log.Fatalf("failed to focus workspace '%s': %v", nextWs, err)
	}
}

func mustGetNextWs(extensionPrefix string) workspace {
	wss, err := listWorkspaces(extensionPrefix)
	if err != nil {
		log.Fatalf("failed to list workspaces: %v", err)
	}

	focusedWs, err := wss.GetFocused()
	if err != nil {
		log.Fatalf("failed to get focused workspace: %v", err)
	}

	nextWs, err := getNextWorkspace(focusedWs)
	if err != nil {
		log.Fatalf("failed to get calculate next workspace to toggle: %v", err)
	}

	return nextWs
}

func main() {
	args := args{}
	args.mustParse(os.Args)

	switch wsExtensionAction(args.WsAction) {
	case toggleWsExtensionAction:
		toggle(args.ExtensionPrefix)
	case flipContainerWsExtensionAction:
		flipContainer(args.ExtensionPrefix)

	default:
		log.Fatalf("Unknown action '%s' provided", args.WsAction)
	}
}
