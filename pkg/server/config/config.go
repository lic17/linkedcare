package config

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog"
	"linkedcare.io/linkedcare/pkg/apiserver/runtime"
	"linkedcare.io/linkedcare/pkg/simple/client/k8s"
	"linkedcare.io/linkedcare/pkg/simple/client/linkedcare"
	"linkedcare.io/linkedcare/pkg/simple/client/servicemesh"
)

// Package config saves configuration for running Linkedcare components
//
// Config can be configured from command line flags and configuration file.
// Command line flags hold higher priority than configuration file. But if
// component Endpoint/Host/APIServer was left empty, all of that component
// command line flags will be ignored, use configuration file instead.
// For example, we have configuration file
//
// mysql:
//   host: mysql.linkedcare-system.svc
//   username: root
//   password: password
//
// At the same time, have command line flags like following:
//
// --mysql-host mysql.openpitrix-system.svc --mysql-username king --mysql-password 1234
//
// We will use `king:1234@mysql.openpitrix-system.svc` from command line flags rather
// than `root:password@mysql.linkedcare-system.svc` from configuration file,
// cause command line has higher priority. But if command line flags like following:
//
// --mysql-username root --mysql-password password
//
// we will `root:password@mysql.linkedcare-system.svc` as input, case
// mysql-host is missing in command line flags, all other mysql command line flags
// will be ignored.

// InstallAPI installs api for config
func InstallAPI(c *restful.Container) {
	ws := runtime.NewWebService(schema.GroupVersion{
		Group:   "",
		Version: "v1alpha1",
	})

	ws.Route(ws.GET("/configz").
		To(func(request *restful.Request, response *restful.Response) {
			var conf = *sharedConfig

			conf.stripEmptyOptions()

			response.WriteAsJson(convertToMap(&conf))
		}).
		Doc("Get system components configuration").
		Produces(restful.MIME_JSON).
		Writes(Config{}).
		Returns(http.StatusOK, "ok", Config{}))

	c.Add(ws)
}

// convertToMap simply converts config to map[string]bool
// to hide sensitive information
func convertToMap(conf *Config) map[string]bool {
	result := make(map[string]bool, 0)

	if conf == nil {
		return result
	}

	c := reflect.Indirect(reflect.ValueOf(conf))

	for i := 0; i < c.NumField(); i++ {
		name := strings.Split(c.Type().Field(i).Tag.Get("json"), ",")[0]
		if strings.HasPrefix(name, "-") {
			continue
		}

		if c.Field(i).IsNil() {
			result[name] = false
		} else {
			result[name] = true
		}
	}

	return result
}

// Load loads configuration after setup
func Load() error {
	sharedConfig = newConfig()

	viper.SetConfigName(DefaultConfigurationName)
	viper.AddConfigPath(DefaultConfigurationPath)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			klog.Warning("configuration file not found")
			return nil
		} else {
			panic(fmt.Errorf("error parsing configuration file %s", err))
		}
	}

	conf := newConfig()
	if err := viper.Unmarshal(conf); err != nil {
		klog.Error(fmt.Errorf("error unmarshal configuration %v", err))
		return err
	} else {
		// make sure linkedcare options always exists
		if conf.LinkedcareOptions == nil {
			conf.LinkedcareOptions = linkedcare.NewLinkedcareOptions()
		} else {
			ksOptions := linkedcare.NewLinkedcareOptions()
			conf.LinkedcareOptions.ApplyTo(ksOptions)
			conf.LinkedcareOptions = ksOptions
		}

		conf.Apply(shadowConfig)
		sharedConfig = conf
	}

	return nil
}

const (
	// DefaultConfigurationName is the default name of configuration
	DefaultConfigurationName = "linkedcare"

	// DefaultConfigurationPath the default location of the configuration file
	DefaultConfigurationPath = "/etc/linkedcare"
)

var (
	// sharedConfig holds configuration across linkedcare
	sharedConfig *Config

	// shadowConfig contains options from commandline options
	shadowConfig = &Config{}
)

type Config struct {
	KubernetesOptions  *k8s.KubernetesOptions          `json:"kubernetes,omitempty" yaml:"kubernetes,omitempty" mapstructure:"kubernetes"`
	ServiceMeshOptions *servicemesh.ServiceMeshOptions `json:"servicemesh,omitempty" yaml:"servicemesh,omitempty" mapstructure:"servicemesh"`

	// Options below are only loaded from configuration file, no command line flags for these options now.
	LinkedcareOptions *linkedcare.LinkedcareOptions `json:"-" yaml:"linkedcare,omitempty" mapstructure:"linkedcare"`
}

func newConfig() *Config {
	return &Config{
		KubernetesOptions:  k8s.NewKubernetesOptions(),
		ServiceMeshOptions: servicemesh.NewServiceMeshOptions(),
		LinkedcareOptions:  linkedcare.NewLinkedcareOptions(),
	}
}

func Get() *Config {
	return sharedConfig
}

func (c *Config) Apply(conf *Config) {
	shadowConfig = conf

	if conf.LinkedcareOptions != nil {
		conf.LinkedcareOptions.ApplyTo(c.LinkedcareOptions)
	}

	if conf.ServiceMeshOptions != nil {
		conf.ServiceMeshOptions.ApplyTo(c.ServiceMeshOptions)
	}

	if conf.KubernetesOptions != nil {
		conf.KubernetesOptions.ApplyTo(c.KubernetesOptions)
	}
}

func (c *Config) stripEmptyOptions() {
	if c.ServiceMeshOptions != nil && c.ServiceMeshOptions.IstioPilotHost == "" &&
		c.ServiceMeshOptions.ServicemeshPrometheusHost == "" &&
		c.ServiceMeshOptions.JaegerQueryHost == "" {
		c.ServiceMeshOptions = nil
	}
}
