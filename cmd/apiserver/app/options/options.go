package options

import (
	"flag"
	"strings"

	genericoptions "linkedcare.io/linkedcare/pkg/server/options"

	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog"
	"linkedcare.io/linkedcare/pkg/simple/client/k8s"
	"linkedcare.io/linkedcare/pkg/simple/client/servicemesh"
)

type ServerRunOptions struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions

	KubernetesOptions  *k8s.KubernetesOptions
	ServiceMeshOptions *servicemesh.ServiceMeshOptions
}

func NewServerRunOptions() *ServerRunOptions {

	s := ServerRunOptions{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),

		KubernetesOptions:  k8s.NewKubernetesOptions(),
		ServiceMeshOptions: servicemesh.NewServiceMeshOptions(),
	}

	return &s
}

func (s *ServerRunOptions) Flags() (fss cliflag.NamedFlagSets) {

	s.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))

	s.KubernetesOptions.AddFlags(fss.FlagSet("kubernetes"))
	s.ServiceMeshOptions.AddFlags(fss.FlagSet("servicemesh"))

	fs := fss.FlagSet("klog")
	local := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(local)
	local.VisitAll(func(fl *flag.Flag) {
		fl.Name = strings.Replace(fl.Name, "_", "-", -1)
		fs.AddGoFlag(fl)
	})

	return fss
}
