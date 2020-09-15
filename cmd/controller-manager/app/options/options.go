package options

import (
	"flag"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/client-go/tools/leaderelection"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog"
)

type LinkedcareControllerManagerOptions struct {
	LeaderElection *leaderelection.LeaderElectionConfig
	Kubeconfig     string
}

func NewLinkedcareControllerManagerOptions() *LinkedcareControllerManagerOptions {
	s := &LinkedcareControllerManagerOptions{
		LeaderElection: &leaderelection.LeaderElectionConfig{
			LeaseDuration: 30 * time.Second,
			RenewDeadline: 15 * time.Second,
			RetryPeriod:   5 * time.Second,
		},
		Kubeconfig: "",
	}

	return s
}

func (s *LinkedcareControllerManagerOptions) Flags() cliflag.NamedFlagSets {
	fss := cliflag.NamedFlagSets{}

	fs := fss.FlagSet("leaderelection")
	s.bindLeaderElectionFlags(s.LeaderElection, fs)

	kcfs := fss.FlagSet("Kubeconfig")
	kcfs.StringVar(&s.Kubeconfig, "default-kubeconfig", s.Kubeconfig, "")

	kfs := fss.FlagSet("klog")
	local := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(local)
	local.VisitAll(func(fl *flag.Flag) {
		fl.Name = strings.Replace(fl.Name, "_", "-", -1)
		kfs.AddGoFlag(fl)
	})

	return fss
}

func (s *LinkedcareControllerManagerOptions) bindLeaderElectionFlags(l *leaderelection.LeaderElectionConfig, fs *pflag.FlagSet) {
	fs.DurationVar(&l.LeaseDuration, "leader-elect-lease-duration", l.LeaseDuration, ""+
		"The duration that non-leader candidates will wait after observing a leadership "+
		"renewal until attempting to acquire leadership of a led but unrenewed leader "+
		"slot. This is effectively the maximum duration that a leader can be stopped "+
		"before it is replaced by another candidate. This is only applicable if leader "+
		"election is enabled.")
	fs.DurationVar(&l.RenewDeadline, "leader-elect-renew-deadline", l.RenewDeadline, ""+
		"The interval between attempts by the acting master to renew a leadership slot "+
		"before it stops leading. This must be less than or equal to the lease duration. "+
		"This is only applicable if leader election is enabled.")
	fs.DurationVar(&l.RetryPeriod, "leader-elect-retry-period", l.RetryPeriod, ""+
		"The duration the clients should wait between attempting acquisition and renewal "+
		"of a leadership. This is only applicable if leader election is enabled.")
}
