package main

import (
	"fmt"
	"gh-user-activity/client"
	"gh-user-activity/tui"
	"gh-user-activity/utils"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		utils.PrintUsage()
		return
	}

	switch os.Args[1] {
	case "view":
		if len(os.Args) < 3 {
			fmt.Println("Usage: gh-activity view <username>")
			return
		}
		events, err := client.GetUserActivity(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		utils.PrintEvents(os.Args[2], events)
	case "help", "-h", "--help":
		utils.PrintUsage()

	case "tui":
		tui.RunTUI()

	default:
		utils.PrintUsage()
	}
}
