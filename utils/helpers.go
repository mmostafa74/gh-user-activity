package utils

import (
	"fmt"
	"gh-user-activity/client"
	"gh-user-activity/formatter"
	"strings"
)

func PrintUsage() {
	fmt.Println(`gh-activity — View GitHub user activity from the CLI

Usage:
  gh-activity view <username>           Show recent GitHub activity for a user
  gh-activity view <username> --limit N Limit displayed events to N (default: 30)
  gh-activity tui                       Launch interactive TUI (coming soon)
  gh-activity help                      Show this help message

Examples:
  gh-activity view torvalds
  gh-activity view torvalds --limit 5`)
}

func PrintEvents(username string, events []client.Event) {
	fmt.Printf("\nActivity for %s (last %d events)\n", username, len(events))
	fmt.Println(strings.Repeat("─", 50))
	for _, e := range events {
		fmt.Printf("%s  • %s\n", e.CreatedAt.Format("2006-01-02 15:04"), formatter.FormatEvent(e))
	}
}
