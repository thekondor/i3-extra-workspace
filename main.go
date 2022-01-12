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

func flipContainer() {
	nextWs := mustGetNextWs()
	if err := nextWs.FlipFocusedContainer(); err != nil {
		log.Fatalf("failed to flip focused container to workspace '%s': %v", nextWs, err)
	}
}

func toggle() {
	nextWs := mustGetNextWs()
	if err := nextWs.Focus(); err != nil {
		log.Fatalf("failed to focus workspace '%s': %v", nextWs, err)
	}
}

func mustGetNextWs() workspace {
	wss, err := listWorkspaces()
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
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <%s|%s>\n", os.Args[0], toggleWsExtensionAction, flipContainerWsExtensionAction)
		log.Fatalf("No argument(s) provided")
	}

	wsAction := wsExtensionAction(os.Args[1])
	switch wsAction {
	case toggleWsExtensionAction:
		toggle()
	case flipContainerWsExtensionAction:
		flipContainer()

	default:
		log.Fatalf("Unknown action '%s' provided", wsAction)
	}
}
