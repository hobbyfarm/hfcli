package cmd

import (
	"github.com/hobbyfarm/hfcli/pkg/scenario"
	command "github.com/rancher/wrangler-cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Get struct{}

func NewGet() *cobra.Command {
	get := command.Command(&Get{}, cobra.Command{
		Use:   "get",
		Short: "get objects, valid options are scenario",
	})
	get.AddCommand(
		NewGetScenario(),
	)
	return get
}

func (a *Get) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	return nil
}

type GetScenario struct{}

func NewGetScenario() *cobra.Command {
	getScenario := command.Command(&GetScenario{}, cobra.Command{
		Use:   "scenario",
		Short: "get scenario NAME PATH_TO_SCENARIO",
		Args:  cobra.ExactArgs(2),
	})
	return getScenario
}

func (sc *GetScenario) Run(cmd *cobra.Command, args []string) error {
	logrus.Info(args[0], args[1])
	s, err := scenario.Get(args[0], Namespace, HfClient)

	if err != nil {
		return err
	}

	return scenario.DumpScenario(s, args[1])
}
