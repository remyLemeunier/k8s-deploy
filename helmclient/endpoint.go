package helmclient

import (
	"fmt"
	"os"
	"path"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func GetTillerHosts(clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	tillers := []string{}

	endpoints, err := clientset.CoreV1().Endpoints("kube-system").Get("tiller", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	for _, sub := range endpoints.Subsets {
		for _, port := range sub.Ports {
			for _, addr := range sub.Addresses {
				tillers = append(tillers, fmt.Sprintf("%s:%d", addr.IP, port.Port))
			}
		}
	}
	if len(tillers) == 0 {
		return nil, fmt.Errorf("Endpoint 'tiller' has no ready pods")
	}
	return tillers, nil
}

func NewKubeClient(configPath string, cluster string) (*kubernetes.Clientset, error) {
	if configPath == "" {
		configPath = path.Join(os.Getenv("HOME"), ".kube", "config")
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.ExplicitPath = configPath

	configOverrides := &clientcmd.ConfigOverrides{
		Context: clientcmdapi.Context{
			Cluster: cluster,
		},
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	return kubernetes.NewForConfig(config)
}
