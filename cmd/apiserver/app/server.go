package app

import (
	"fmt"
	"net/http"

	kconfig "github.com/kiali/kiali/config"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog"
	"linkedcare.io/linkedcare/cmd/apiserver/app/options"
	"linkedcare.io/linkedcare/pkg/api"
	"linkedcare.io/linkedcare/pkg/apiserver/runtime"

	//"linkedcare.io/linkedcare/pkg/apiserver/servicemesh/tracing"
	"linkedcare.io/linkedcare/pkg/informers"
	"linkedcare.io/linkedcare/pkg/server"
	apiserverconfig "linkedcare.io/linkedcare/pkg/server/config"
	"linkedcare.io/linkedcare/pkg/server/filter"
	"linkedcare.io/linkedcare/pkg/simple/client"
	"linkedcare.io/linkedcare/pkg/utils/signals"
	"linkedcare.io/linkedcare/pkg/utils/term"
)

func NewAPIServerCommand() *cobra.Command {
	s := options.NewServerRunOptions()

	cmd := &cobra.Command{
		Use: "lc-apiserver",
		Long: `The Linkedcare API server validates and configures data for the api objects. 
The API Server services REST operations and provides the frontend to the
cluster's shared state through which all other components interact.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := apiserverconfig.Load()
			if err != nil {
				return err
			}

			err = Complete(s)
			if err != nil {
				return err
			}

			if errs := s.Validate(); len(errs) != 0 {
				return utilerrors.NewAggregate(errs)
			}

			return Run(s, signals.SetupSignalHandler())
		},
	}

	fs := cmd.Flags()
	namedFlagSets := s.Flags()

	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
	return cmd
}

func Run(s *options.ServerRunOptions, stopCh <-chan struct{}) error {

	err := CreateClientSet(apiserverconfig.Get(), stopCh)
	if err != nil {
		return err
	}

	err = WaitForResourceSync(stopCh)
	if err != nil {
		return err
	}

	initializeServicemeshConfig(s)

	err = CreateAPIServer(s)
	if err != nil {
		return err
	}

	return nil
}

func initializeServicemeshConfig(s *options.ServerRunOptions) {
	// Initialize kiali config
	config := kconfig.NewConfig()

	//tracing.JaegerQueryUrl = s.ServiceMeshOptions.JaegerQueryHost

	// Exclude system namespaces
	config.API.Namespaces.Exclude = []string{"istio-system", "linkedcare*", "kube*"}
	config.InCluster = true

	// Set istio pilot discovery service url
	config.ExternalServices.Istio.UrlServiceVersion = s.ServiceMeshOptions.IstioPilotHost

	kconfig.Set(config)
}

//
func CreateAPIServer(s *options.ServerRunOptions) error {
	var err error

	container := runtime.Container
	container.DoNotRecover(false)
	container.Filter(filter.Logging)
	container.RecoverHandler(server.LogStackOnRecover)

	api.InstallAPIs(container)

	// install config api
	apiserverconfig.InstallAPI(container)

	if s.GenericServerRunOptions.InsecurePort != 0 {
		klog.V(0).Infof("Server listening on insecure port %d.", s.GenericServerRunOptions.InsecurePort)
		err = http.ListenAndServe(fmt.Sprintf("%s:%d", s.GenericServerRunOptions.BindAddress, s.GenericServerRunOptions.InsecurePort), container)
		if err != nil {
			klog.Errorf("Server listening insecure err %s.", err.Error())
		}
	}

	if s.GenericServerRunOptions.SecurePort != 0 && len(s.GenericServerRunOptions.TlsCertFile) > 0 && len(s.GenericServerRunOptions.TlsPrivateKey) > 0 {
		klog.V(0).Infof("Server listening on secure port %d.", s.GenericServerRunOptions.SecurePort)
		err = http.ListenAndServeTLS(fmt.Sprintf("%s:%d", s.GenericServerRunOptions.BindAddress, s.GenericServerRunOptions.SecurePort), s.GenericServerRunOptions.TlsCertFile, s.GenericServerRunOptions.TlsPrivateKey, container)
		if err != nil {
			klog.Errorf("Server listening on secure err %s.", err.Error())
		}
	}

	return err
}

func CreateClientSet(conf *apiserverconfig.Config, stopCh <-chan struct{}) error {
	csop := &client.ClientSetOptions{}

	csop.SetKubernetesOptions(conf.KubernetesOptions).
		SetLinkedcareOptions(conf.LinkedcareOptions)

	client.NewClientSetFactory(csop, stopCh)

	return nil
}

func WaitForResourceSync(stopCh <-chan struct{}) error {
	klog.V(0).Info("Start cache objects")

	discoveryClient := client.ClientSets().K8s().Discovery()
	apiResourcesList, err := discoveryClient.ServerResources()
	if err != nil {
		return err
	}

	isResourceExists := func(resource schema.GroupVersionResource) bool {
		for _, apiResource := range apiResourcesList {
			if apiResource.GroupVersion == resource.GroupVersion().String() {
				for _, rsc := range apiResource.APIResources {
					if rsc.Name == resource.Resource {
						return true
					}
				}
			}
		}
		return false
	}

	informerFactory := informers.SharedInformerFactory()

	// resources we have to create informer first
	k8sGVRs := []schema.GroupVersionResource{
		{Group: "", Version: "v1", Resource: "namespaces"},
		{Group: "", Version: "v1", Resource: "nodes"},
		{Group: "", Version: "v1", Resource: "resourcequotas"},
		{Group: "", Version: "v1", Resource: "pods"},
		{Group: "", Version: "v1", Resource: "services"},
		{Group: "", Version: "v1", Resource: "persistentvolumeclaims"},
		{Group: "", Version: "v1", Resource: "secrets"},
		{Group: "", Version: "v1", Resource: "configmaps"},

		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"},

		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "apps", Version: "v1", Resource: "daemonsets"},
		{Group: "apps", Version: "v1", Resource: "replicasets"},
		{Group: "apps", Version: "v1", Resource: "statefulsets"},
		{Group: "apps", Version: "v1", Resource: "controllerrevisions"},

		{Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"},

		{Group: "batch", Version: "v1", Resource: "jobs"},
		{Group: "batch", Version: "v1beta1", Resource: "cronjobs"},

		{Group: "extensions", Version: "v1beta1", Resource: "ingresses"},

		{Group: "autoscaling", Version: "v2beta2", Resource: "horizontalpodautoscalers"},
	}

	for _, gvr := range k8sGVRs {
		if !isResourceExists(gvr) {
			klog.Warningf("resource %s not exists in the cluster", gvr)
		} else {
			_, err := informerFactory.ForResource(gvr)
			if err != nil {
				klog.Errorf("cannot create informer for %s", gvr)
				return err
			}
		}
	}

	informerFactory.Start(stopCh)
	informerFactory.WaitForCacheSync(stopCh)

	lcInformerFactory := informers.LcSharedInformerFactory()

	lcGVRs := []schema.GroupVersionResource{
		{Group: "servicemesh.linkedcare.io", Version: "v1alpha1", Resource: "strategies"},
		{Group: "servicemesh.linkedcare.io", Version: "v1alpha1", Resource: "servicepolicies"},
	}

	for _, gvr := range lcGVRs {
		if !isResourceExists(gvr) {
			klog.Warningf("resource %s not exists in the cluster", gvr)
		} else {
			_, err := lcInformerFactory.ForResource(gvr)
			if err != nil {
				return err
			}
		}
	}

	lcInformerFactory.Start(stopCh)
	lcInformerFactory.WaitForCacheSync(stopCh)

	appInformerFactory := informers.AppSharedInformerFactory()

	appGVRs := []schema.GroupVersionResource{
		{Group: "app.k8s.io", Version: "v1beta1", Resource: "applications"},
	}

	for _, gvr := range appGVRs {
		if !isResourceExists(gvr) {
			klog.Warningf("resource %s not exists in the cluster", gvr)
		} else {
			_, err := appInformerFactory.ForResource(gvr)
			if err != nil {
				return err
			}
		}
	}

	appInformerFactory.Start(stopCh)
	appInformerFactory.WaitForCacheSync(stopCh)

	klog.V(0).Info("Finished caching objects")

	return nil

}

// apply server run options to configuration
func Complete(s *options.ServerRunOptions) error {

	// loading configuration file
	conf := apiserverconfig.Get()

	conf.Apply(&apiserverconfig.Config{
		KubernetesOptions:  s.KubernetesOptions,
		ServiceMeshOptions: s.ServiceMeshOptions,
	})

	*s = options.ServerRunOptions{
		GenericServerRunOptions: s.GenericServerRunOptions,

		KubernetesOptions:  conf.KubernetesOptions,
		ServiceMeshOptions: conf.ServiceMeshOptions,
	}

	return nil
}
