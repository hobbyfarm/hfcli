package main

import (
	hf "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	"github.com/hobbyfarm/hfcli/pkg/cmd"
	command "github.com/rancher/wrangler-cli"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = hf.AddToScheme(scheme)
}
func main() {

	command.Main(cmd.App())
}
