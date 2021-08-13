package scenario

import (
	"context"
	"fmt"
	hf "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	hfClientSet "github.com/hobbyfarm/gargantua/pkg/client/clientset/versioned/typed/hobbyfarm.io/v1"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Apply(s *hf.Scenario, hfc *hfClientSet.HobbyfarmV1Client) (err error) {

	// check if scneario exists //
	sGet, err := hfc.Scenarios().Get(context.TODO(), s.GetName(), v1.GetOptions{})

	if err != nil {
		if apierrors.IsNotFound(err){
			logrus.Infof("creating scenario %s", s.GetName())
			_, err = hfc.Scenarios().Create(context.TODO(), s, v1.CreateOptions{})
			return err
		} else {
			return err
		}
	}

	if sGet != nil {
		key, ok := sGet.Annotations["managedBy"]
		if ok && key == "hfcli" {
			s.ObjectMeta.ResourceVersion = sGet.ObjectMeta.GetResourceVersion()
			logrus.Info("updating scenario %s", s.GetName())
			_, err = hfc.Scenarios().Update(context.TODO(), s, v1.UpdateOptions{})
		} else {
			err = fmt.Errorf("scenario %s already exists and is not managed by hfcli", sGet.GetName())
		}

	}

	return err
}
