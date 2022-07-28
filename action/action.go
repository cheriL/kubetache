package action

import (
	"flag"
	"path/filepath"

	"github.com/cheriL/kubetache/kube"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var client kube.Client

func Init(option int) {
	var c *kubernetes.Clientset
	var err error

	if option == 0 {
		// creates the in-cluster config
		//config, err := rest.InClusterConfig()
		//if err != nil {
		//	panic(err.Error())
		//}
	} else {
		config := outOfClusterConfig()
		// creates the clientset
		c, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}

	client = kube.NewClusterCache(c)
	client.Run()
}

func outOfClusterConfig() *rest.Config {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return config
}