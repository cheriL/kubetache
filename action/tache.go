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

	// TODO check whether secrets and configmaps exist
	workLoad.Secrets, workLoad.ConfigMaps = makeSecretsAndConfigMaps(pod)

	if pod.Spec.ServiceAccountName != "" {

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

func makeSecretsAndConfigMaps(pod *v1.Pod) (secrets []*models.Secret, configMaps []*models.ConfigMap) {
	secretMappings := make(map[string]models.SourceType)
	configMappings := make(map[string]models.SourceType)

	var envFromList []v1.EnvFromSource
	for _, container := range pod.Spec.Containers {
		envFromList = append(envFromList, container.EnvFrom...)
	}

	for _, secret := range pod.Spec.ImagePullSecrets {
		name := secret.Name
		secretMappings[name] = models.SourceTypeImagePullSecret
	}

	for _, envFrom := range envFromList{
		if envFrom.SecretRef != nil {
			name := envFrom.SecretRef.Name
			secretMappings[name] = models.SourceTypeEnvFrom
		} else if envFrom.ConfigMapRef != nil {
			name := envFrom.ConfigMapRef.Name
			configMappings[name] = models.SourceTypeEnvFrom
		}
	}

	for _, volume := range pod.Spec.Volumes {
		if volume.Secret != nil {
			name := volume.Secret.SecretName
			secretMappings[name] = models.SourceTypeVolume
		} else if volume.ConfigMap != nil {
			name := volume.ConfigMap.Name
			configMappings[name] = models.SourceTypeVolume
		}
	}

	for secretName, sourceType := range secretMappings {

		secret := &models.Secret{
			Name:       secretName,
			SourceType: sourceType,
		}

		secrets = append(secrets, secret)
	}

	for configMapName, sourceType := range configMappings {

		configMap := &models.ConfigMap{
			Name:       configMapName,
			SourceType: sourceType,
		}

		configMaps = append(configMaps, configMap)
	}

	return
}