package models

// Workload actually a pod.
type Workload struct {
	// Name of the pod.
	Name           string
	// Controller of the pod. Like ReplicaSet, statefulSet, daemonSet, etc.
	Controller     string
	// Service of the pod.
	Service        []string
	// Secrets used by the pod.
	Secrets        []*Secret
	// ConfigMaps used by the pod.
	ConfigMaps     []*ConfigMap
	// Pvc Bound to the pod.
	Pvc            []string
	// ServiceAccount used by the pod.
	ServiceAccount *ServiceAccount
	//
	Additions      *Additions
}

const (
	sourceTypeEnvFrom = "envFrom"
	sourceTypeVolume = "Volume"
)

type SourceType string

type Secret struct {
	Name       string
	SourceType SourceType
}

type ConfigMap struct {
	Name       string
	SourceType SourceType
}

type ServiceAccount struct {
	Name       string
	SecretName string
	Bindings   []*RoleBinding
}

type RoleBinding struct {
	Type string
	Name string
	Role string
}

type Additions struct {

}