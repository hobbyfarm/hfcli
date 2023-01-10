package cmd

import (
	"github.com/hobbyfarm/hfcli/pkg/scenario"
	command "github.com/rancher/wrangler-cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Delete struct{}

func NewDelete() *cobra.Command {
	delete := command.Command(&Delete{}, cobra.Command{
		Use:   "delete",
		Short: "delete objects, valid options are scenario",
	})
	delete.AddCommand(
		NewDeleteScenario(),
	)
	return delete
}

func (a *Delete) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	return nil
}

type DeleteScenario struct{}

func NewDeleteScenario() *cobra.Command {
	deleteScenario := command.Command(&DeleteScenario{}, cobra.Command{
		Use:   "scenario",
		Short: "delete scenario NAME",
	})
	return deleteScenario
}

func (sc *DeleteScenario) Run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		logrus.Error("not enough arguments supplied")
		return cmd.Help()
	}

	if len(args) > 1 {
		logrus.Error("too many arguments supplied")
		return cmd.Help()
	}

	logrus.Info(args[0])
	return scenario.Delete(args[0], HfClient)
}
