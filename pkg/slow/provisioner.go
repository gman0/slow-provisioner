package slow

import (
	"fmt"

	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
)

type Provisioner struct {
	clientset    clientset.Interface
	nodeplugin   string
	delayOptions *DelayOptions
}

func NewProvisioner(c clientset.Interface, nodeplugin string, delayOptions *DelayOptions) *Provisioner {
	return &Provisioner{
		clientset:    c,
		nodeplugin:   nodeplugin,
		delayOptions: delayOptions,
	}
}

func (p *Provisioner) Provision(volOptions controller.VolumeOptions) (*v1.PersistentVolume, error) {
	if volOptions.PVC.Spec.Selector != nil {
		return nil, fmt.Errorf("claim Selector is not supported")
	}

	runDelay("Provision", p.delayOptions.Provision)

	return buildPersistentVolume(p.nodeplugin, &volOptions), nil
}

func (p *Provisioner) Delete(pv *v1.PersistentVolume) error {
	runDelay("Delete", p.delayOptions.Delete)
	return nil
}
