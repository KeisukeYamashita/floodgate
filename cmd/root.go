package cmd

import (
	"io"

	"github.com/codilime/floodgate/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootOptions store root command options
type RootOptions struct {
	configFile string
	quiet      bool
}

// Execute execute command
func Execute(out io.Writer) error {
	rootCmd := NewRootCmd(out)
	return rootCmd.Execute()
}

// NewRootCmd create new root command
func NewRootCmd(out io.Writer) *cobra.Command {
	options := RootOptions{}

	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.String(),
	}
	cmd.PersistentFlags().StringVar(&options.configFile, "config", "", "path to config file (default $HOME/.config/floodgate/config.yaml)")
	cmd.PersistentFlags().BoolVarP(&options.quiet, "quiet", "q", false, "squelch non-essential output")

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		level := "debug"
		if options.quiet {
			level = "error"
		}
		if err := setUpLogs(out, level); err != nil {
			return err
		}
		return nil
	}

	cmd.AddCommand(NewSyncCmd(out))
	cmd.AddCommand(NewCompareCmd(out))
	cmd.AddCommand(NewHydrateCmd(out))
	cmd.AddCommand(NewInspectCmd(out))
	cmd.AddCommand(NewRenderCmd(out))

	return cmd
}

// setUpLogs set the log output and the log level
func setUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}
