package scenario

import (
	"context"
	"fmt"

	hf "github.com/hobbyfarm/gargantua/pkg/apis/hobbyfarm.io/v1"
	hfClientSet "github.com/hobbyfarm/gargantua/pkg/client/clientset/versioned"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Apply(s *hf.Scenario, Namespace string, hfc *hfClientSet.Clientset) (err error) {

	// check if scneario exists //
	sGet, err := hfc.HobbyfarmV1().Scenarios(Namespace).Get(context.TODO(), s.GetName(), v1.GetOptions{})

	if err != nil {
		if apierrors.IsNotFound(err) {
			logrus.Infof("creating scenario %s", s.GetName())
			_, err = hfc.HobbyfarmV1().Scenarios(Namespace).Create(context.TODO(), s, v1.CreateOptions{})
			return err
		} else {
			return err
		}
	}

	if sGet != nil {
		key, ok := sGet.Annotations["managedBy"]
		if ok && key == "hfcli" {
			s.ObjectMeta.ResourceVersion = sGet.ObjectMeta.GetResourceVersion()
			logrus.Infof("updating scenario %s", s.GetName())
			_, err = hfc.HobbyfarmV1().Scenarios(Namespace).Update(context.TODO(), s, v1.UpdateOptions{})
		} else {
			err = fmt.Errorf("scenario %s already exists and is not managed by hfcli", sGet.GetName())
		}

	}

	return err
}

func Get(name string, Namespace string, hfc *hfClientSet.Clientset) (s *hf.Scenario, err error) {
	logrus.Infof("downloading scenario %s", name)

	return hfc.HobbyfarmV1().Scenarios(Namespace).Get(context.TODO(), name, v1.GetOptions{})
}

func Delete(name string, Namespace string, hfc *hfClientSet.Clientset) (err error) {
	logrus.Infof("deleting scenario %s", name)

	return hfc.HobbyfarmV1().Scenarios(Namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
}
