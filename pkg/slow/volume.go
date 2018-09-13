package slow

import (
	"fmt"

	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildPersistentVolume(nodeplugin string, volOptions *controller.VolumeOptions) *v1.PersistentVolume {
	return &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: volOptions.PVName,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: volOptions.PersistentVolumeReclaimPolicy,
			AccessModes:                   volOptions.PVC.Spec.AccessModes,
			Capacity:                      v1.ResourceList{v1.ResourceStorage: resource.MustParse(fmt.Sprintf("%dG", volOptions.PVC.Spec.Size()))},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				CSI: &v1.CSIPersistentVolumeSource{
					Driver:       nodeplugin,
					ReadOnly:     false,
					VolumeHandle: "pvc-" + string(volOptions.PVC.GetUID()),
				},
			},
		},
	}
}
