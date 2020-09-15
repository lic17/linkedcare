package client

import (
	"fmt"
	"sync"

	"linkedcare.io/linkedcare/pkg/simple/client/k8s"
	"linkedcare.io/linkedcare/pkg/simple/client/linkedcare"
)

type ClientSetNotEnabledError struct {
	err error
}

func (e ClientSetNotEnabledError) Error() string {
	return fmt.Sprintf("client set not enabled: %v", e.err)
}

type ClientSetOptions struct {
	kubernetesOptions *k8s.KubernetesOptions
	linkedcareOptions *linkedcare.LinkedcareOptions
}

func NewClientSetOptions() *ClientSetOptions {
	return &ClientSetOptions{
		kubernetesOptions: k8s.NewKubernetesOptions(),
		linkedcareOptions: linkedcare.NewLinkedcareOptions(),
	}
}

func (c *ClientSetOptions) SetKubernetesOptions(options *k8s.KubernetesOptions) *ClientSetOptions {
	c.kubernetesOptions = options
	return c
}

func (c *ClientSetOptions) SetLinkedcareOptions(options *linkedcare.LinkedcareOptions) *ClientSetOptions {
	c.linkedcareOptions = options
	return c
}

// ClientSet provide best of effort service to initialize clients,
// but there is no guarantee to return a valid client instance,
// so do validity check before use
type ClientSet struct {
	csoptions *ClientSetOptions
	stopCh    <-chan struct{}

	k8sClient        *k8s.KubernetesClient
	linkedcareClient *linkedcare.LinkedcareClient
}

var mutex sync.Mutex

// global clientsets instance
var sharedClientSet *ClientSet

func ClientSets() *ClientSet {
	return sharedClientSet
}

func NewClientSetFactory(c *ClientSetOptions, stopCh <-chan struct{}) *ClientSet {
	sharedClientSet = &ClientSet{csoptions: c, stopCh: stopCh}

	if c.kubernetesOptions != nil {
		sharedClientSet.k8sClient = k8s.NewKubernetesClientOrDie(c.kubernetesOptions)
	}

	if c.linkedcareOptions != nil {
		sharedClientSet.linkedcareClient = linkedcare.NewLinkedcareClient(c.linkedcareOptions)
	}

	return sharedClientSet
}

// since kubernetes client is required, we will
// create it on setup
func (cs *ClientSet) K8s() *k8s.KubernetesClient {
	return cs.k8sClient
}

func (cs *ClientSet) Linkedcare() *linkedcare.LinkedcareClient {
	return cs.linkedcareClient
}
