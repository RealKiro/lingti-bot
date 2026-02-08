package cmd

import (
	"fmt"
	"os"

	"github.com/pltanton/lingti-bot/internal/logger"
	"github.com/spf13/cobra"
)

var (
	logLevel string
	debug    bool
)

var rootCmd = &cobra.Command{
	Use:   "lingti-bot",
	Short: "MCP server for system resources",
	Long: `lingti-bot is an MCP (Model Context Protocol) server that exposes
computer system resources to AI assistants.

It provides tools for:
  - File operations (read, write, list, search)
  - Shell command execution
  - System information (CPU, memory, disk)
  - Process management
  - Network information`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// If --debug is set, override log level to very-verbose
		if debug {
			logLevel = "very-verbose"
		}

		// Parse and set log level
		level, err := logger.ParseLevel(logLevel)
		if err != nil {
			return err
		}
		logger.SetLevel(level)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log", "info",
		"Log level: silent, info, verbose, very-verbose")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false,
		"Enable debug mode (sets log level to very-verbose and enables browser debug)")
}

// IsDebug returns true if debug mode is enabled globally
func IsDebug() bool {
	return debug
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
