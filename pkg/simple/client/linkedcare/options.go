package linkedcare

import "github.com/spf13/pflag"

type LinkedcareOptions struct {
	APIServer     string `json:"apiServer" yaml:"apiServer"`
	AccountServer string `json:"accountServer" yaml:"accountServer"`
}

// NewLinkedcareOptions create a default options
func NewLinkedcareOptions() *LinkedcareOptions {
	return &LinkedcareOptions{
		APIServer:     "http://ks-apiserver.linkedcare-system.svc",
		AccountServer: "http://ks-account.linkedcare-system.svc",
	}
}

func (s *LinkedcareOptions) ApplyTo(options *LinkedcareOptions) {
	if s.AccountServer != "" {
		options.AccountServer = s.AccountServer
	}

	if s.APIServer != "" {
		options.APIServer = s.APIServer
	}
}

func (s *LinkedcareOptions) Validate() []error {
	errs := []error{}

	return errs
}

func (s *LinkedcareOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.APIServer, "linkedcare-apiserver-host", s.APIServer, ""+
		"Linkedcare apiserver host address.")

	fs.StringVar(&s.AccountServer, "linkedcare-account-host", s.AccountServer, ""+
		"Linkedcare account server host address.")
}
