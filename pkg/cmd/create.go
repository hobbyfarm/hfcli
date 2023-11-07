package cmd

import (
	"github.com/hobbyfarm/hfcli/pkg/scenario"
	command "github.com/rancher/wrangler-cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Apply struct{}

func NewApply() *cobra.Command {
	create := command.Command(&Apply{}, cobra.Command{
		Use:   "apply",
		Short: "apply objects, valid options are scenario",
	})
	create.AddCommand(
		NewApplyScenario(),
	)
	return create
}

func (a *Apply) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	return nil
}

type Scenario struct{}

func NewApplyScenario() *cobra.Command {
	createScenario := command.Command(&Scenario{}, cobra.Command{
		Use:   "scenario",
		Short: "create scenario NAME PATH_TO_SCENARIOS",
	})
	return createScenario
}

func (sc *Scenario) Run(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		logrus.Error("not enough arguments supplied")
		return cmd.Help()
	}

	if len(args) > 2 {
		logrus.Error("too many arguments supplied")
		return cmd.Help()
	}

	logrus.Info(args[0], args[1])
	s, err := scenario.ParseScenario(args[0], Namespace, args[1])

	if err != nil {
		return err
	}

	return scenario.Apply(s, Namespace, HfClient)
}
