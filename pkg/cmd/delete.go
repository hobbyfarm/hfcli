package cmd

import (
	"github.com/hobbyfarm/hfcli/pkg/scenario"
	command "github.com/rancher/wrangler-cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Delete struct{}

func NewDelete() *cobra.Command {
	deleteCommand := command.Command(&Delete{}, cobra.Command{
		Use:   "delete",
		Short: "delete objects, valid options are scenario",
	})
	deleteCommand.AddCommand(
		NewDeleteScenario(),
	)
	return deleteCommand
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
		Args:  cobra.ExactArgs(1),
	})
	return deleteScenario
}

func (sc *DeleteScenario) Run(cmd *cobra.Command, args []string) error {
	logrus.Info(args[0])
	return scenario.Delete(args[0], Namespace, HfClient)
}
