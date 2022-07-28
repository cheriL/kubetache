package action

import (
	"github.com/cheriL/kubetache/models"
	v1 "k8s.io/api/core/v1"
)

func Tache(name, namespace string) {
	pod := client.GetPod(name, namespace)
	if pod == nil {
		// TODO
		return
	}

	workLoad := &models.Workload{
		Name:      pod.Name,
		Namespace: pod.Namespace,
	}

	workLoad.Owners = makeOwners(pod)

	services := client.GetServices(pod.Labels, pod.Namespace)
	for _, s := range services {
		workLoad.Services = append(workLoad.Services, s.Name)
	}
}

func makeOwners(pod *v1.Pod) (owners []*models.Owner) {
	ownerReferences := pod.OwnerReferences
	for _, v := range ownerReferences {
		ownerName := v.Name
		ownerKind := v.Kind

		managedByController := *v.Controller
		if managedByController == true {
			// find the `Deployment` if this pod is managed by a `ReplicaSet`
			if v.Kind == "ReplicaSet" {
				replicaSet := client.GetReplicaSet(ownerName, pod.Namespace)
				for _, v1 := range replicaSet.OwnerReferences {
					if v1.Kind == "Deployment" {
						ownerName = v1.Name
						ownerKind = v1.Kind
						break
					}
				}
			}
		}

		owner := &models.Owner{
			Name: ownerName,
			Kind: ownerKind,
		}

		owners = append(owners, owner)
	}

	return
}