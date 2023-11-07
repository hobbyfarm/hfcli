module github.com/hobbyfarm/hfcli

go 1.16

// github.com/hobbyfarm/gargantua => github.com/hobbyfarm/gargantua v0.2.2-0.20210823170529-e2466136c002
replace k8s.io/client-go => k8s.io/client-go v0.20.0

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/ghodss/yaml v1.0.0
	github.com/hobbyfarm/gargantua v1.0.0
	github.com/rancher/wrangler-cli v0.0.0-20210217230406-95cfa275f52f
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.2.1
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)
