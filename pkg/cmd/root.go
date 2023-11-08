package cmd

import (
	// hfClientSet "github.com/hobbyfarm/gargantua/pkg/client/clientset/versioned/typed/hobbyfarm.io/v1"
	hfClientSet "github.com/hobbyfarm/gargantua/pkg/client/clientset/versioned"
	command "github.com/rancher/wrangler-cli"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// Hfcli is the root command struct
type Hfcli struct {
	Kubeconfig string `usage:"kubeconfig for authentication" short:"k" env:"KUBECONFIG"`
	Namespace  string `usage:"namespace" env:"NAMESPACE"  default:"gargantua" short:"n"`
}

var (
	Namespace string
	HfClient  *hfClientSet.Clientset
)

func App() *cobra.Command {
	root := command.Command(&Hfcli{}, cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
	})
	root.AddCommand(
		NewApply(),
		NewDelete(),
		NewGet(),
		NewInfo(),
	)
	return root
}

func (h *Hfcli) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	return nil
}

func (h *Hfcli) PersistentPre(cmd *cobra.Command, args []string) error {
	var err error
	RestConfig, err := clientcmd.BuildConfigFromFlags("", h.Kubeconfig)
	if err != nil {
		return err
	}

	HfClient, err = hfClientSet.NewForConfig(RestConfig)
	if err != nil {
		return err
	}
	Namespace = h.Namespace
	return nil
}
