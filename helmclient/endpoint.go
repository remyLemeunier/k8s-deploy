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
		//for a := range sub.Addresses {
		//	epAddrs = append(epAddrs, a)
		//}
	}
	return tillers, nil
}

//func RunTiller(env []string, arg ...string) {
//	cmd := exec.Command("helm", arg...)
//	cmd.Env = env
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//	cmd.Run()
//}

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

//func main() {
//	var clusterName string
//
//	flag.StringVar(&clusterName, "cluster", "", "The name of the kubeconfig cluster to use")
//	flag.Parse()
//
//	spew.Dump(flag.Args())
//
//	cl, err := NewKubeClient("/Users/olivierb/.kube/config", clusterName)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	tillers, err := GetTillerHosts(cl, "kube-system")
//	if err != nil {
//		panic(err.Error())
//	}
//
//	env := os.Environ()
//	env = append(env, fmt.Sprintf("HELM_HOST=%s", tillers[0]))
//	RunTiller(env, flag.Args()...)
//}
