package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jessewiles/jwixtac/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "jwixtac",
	Long: `Boost shortcut sprint helper CLI`,
}

func main() {
	var (
		openBrowser bool

		cmdUI = &cobra.Command{
			Use:   "ui",
			Short: "HTTP server for interacting with jwixtac",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				if openBrowser {
					go func() {
						exec.Command("open", "http://localhost:8088").Start()
					}()
				}
				server.SPA()
			},
		}
	)

	cmdUI.Flags().BoolVarP(&openBrowser, "open-browser", "o", false, "Open a browser window to the UI server.")

	rootCmd.AddCommand(cmdUI)

	log.SetLevel(log.DebugLevel)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
