package cmd

import (
	"github.com/hobbyfarm/hfcli/pkg/info"
	command "github.com/rancher/wrangler-cli"
	"github.com/spf13/cobra"
)

type Info struct{}

func NewInfo() *cobra.Command {
	infoCmd := command.Command(&Info{}, cobra.Command{
		Use:   "info",
		Short: "perform info operations, valid options are accesscode and email",
	})
	infoCmd.AddCommand(
		NewInfoEmail(),
		NewInfoAccessCode(),
	)
	return infoCmd
}

func (i *Info) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	return nil
}

type InfoEmail struct{}

func NewInfoEmail() *cobra.Command {
	infoEmailCmd := command.Command(&InfoEmail{}, cobra.Command{
		Use:   "email",
		Short: "get info about session and infra associated with email address",
	})

	return infoEmailCmd
}

func (ie *InfoEmail) Run(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return cmd.Help()
	}

	return info.GetEmail(args[0], Namespace, HfClient)
}

type InfoAccessCode struct {
	Stats bool `usage:"stats" default:"false" short:"s"`
}

func NewInfoAccessCode() *cobra.Command {
	infoAccessCodeCmd := command.Command(&InfoAccessCode{}, cobra.Command{
		Use:   "accesscode",
		Short: "hfcli info accesscode ACCESS_CODE",
	})

	return infoAccessCodeCmd
}

func (iac *InfoAccessCode) Run(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return cmd.Help()
	}

	return info.GetAccessCode(args[0], Namespace, HfClient, iac.Stats)
}
