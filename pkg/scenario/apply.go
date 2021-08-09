package scenario

import (
	"context"
	hf "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	hfClientSet "github.com/hobbyfarm/gargantua/pkg/client/clientset/versioned/typed/hobbyfarm.io/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Apply(s *hf.Scenario, hfc *hfClientSet.HobbyfarmV1Client) (err error) {
	logrus.Infof("creating scenario %s", s.GetName())
	_, err = hfc.Scenarios().Create(context.TODO(), s, v1.CreateOptions{})
	return err
}
