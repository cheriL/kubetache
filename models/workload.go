package models

// Workload actually a pod.
type Workload struct {
	// Name of the pod.
	Name           string
	// Namespace of the pod.
	Namespace      string
	// Owners of the pod.
	Owners         []*Owner
	// Service of the pod.
	Services       []string
	// Secrets used by the pod.
	Secrets        []*Secret
	// ConfigMaps used by the pod.
	ConfigMaps     []*ConfigMap
	// Pvc Bound to the pod.
	Pvcs           []string
	// ServiceAccount used by the pod.
	ServiceAccount *ServiceAccount
	//
	Additions      *Additions
}

type SourceType string

const (
	sourceTypeEnvFrom = "envFrom"
	sourceTypeVolume = "Volume"
)

type Owner struct {
	Name string
	Kind string
}

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