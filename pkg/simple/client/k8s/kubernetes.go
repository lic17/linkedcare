package k8s

import (
	"strings"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	linkedcare "linkedcare.io/linkedcare/pkg/client/clientset/versioned"
	applicationclientset "sigs.k8s.io/application/pkg/client/clientset/versioned"
)

type KubernetesClient struct {
	// kubernetes client interface
	k8s *kubernetes.Clientset

	// discovery client
	discoveryClient *discovery.DiscoveryClient

	// generated clientset
	lc *linkedcare.Clientset

	application *applicationclientset.Clientset

	master string

	config *rest.Config
}

// NewKubernetesClientOrDie creates KubernetesClient and panic if there is an error
func NewKubernetesClientOrDie(options *KubernetesOptions) *KubernetesClient {
	config, err := clientcmd.BuildConfigFromFlags("", options.KubeConfig)
	if err != nil {
		panic(err)
	}

	config.QPS = options.QPS
	config.Burst = options.Burst

	k := &KubernetesClient{
		k8s:             kubernetes.NewForConfigOrDie(config),
		discoveryClient: discovery.NewDiscoveryClientForConfigOrDie(config),
		lc:              linkedcare.NewForConfigOrDie(config),
		application:     applicationclientset.NewForConfigOrDie(config),
		master:          config.Host,
		config:          config,
	}

	if options.Master != "" {
		k.master = options.Master
	}
	// The https prefix is automatically added when using sa.
	// But it will not be set automatically when reading from kubeconfig
	// which may cause some problems in the client of other languages.
	if !strings.HasPrefix(k.master, "http://") && !strings.HasPrefix(k.master, "https://") {
		k.master = "https://" + k.master
	}
	return k
}

// NewKubernetesClient creates a KubernetesClient
func NewKubernetesClient(options *KubernetesOptions) (*KubernetesClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", options.KubeConfig)
	if err != nil {
		return nil, err
	}

	config.QPS = options.QPS
	config.Burst = options.Burst

	var k KubernetesClient
	k.k8s, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	k.lc, err = linkedcare.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	k.application, err = applicationclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	k.master = options.Master
	k.config = config

	return &k, nil
}

func (k *KubernetesClient) Kubernetes() kubernetes.Interface {
	return k.k8s
}

func (k *KubernetesClient) Discovery() discovery.DiscoveryInterface {
	return k.discoveryClient
}

func (k *KubernetesClient) Linkedcare() linkedcare.Interface {
	return k.lc
}

func (k *KubernetesClient) Application() applicationclientset.Interface {
	return k.application
}

// master address used to generate kubeconfig for downloading
func (k *KubernetesClient) Master() string {
	return k.master
}

func (k *KubernetesClient) Config() *rest.Config {
	return k.config
}
