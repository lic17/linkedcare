package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
	"linkedcare.io/linkedcare/pkg/simple/client/k8s"
	"linkedcare.io/linkedcare/pkg/simple/client/linkedcare"
	"linkedcare.io/linkedcare/pkg/simple/client/servicemesh"
	"linkedcare.io/linkedcare/pkg/utils/reflectutils"
)

func newTestConfig() *Config {
	conf := &Config{
		KubernetesOptions: &k8s.KubernetesOptions{
			KubeConfig: "/Users/zry/.kube/config",
			Master:     "https://127.0.0.1:6443",
			QPS:        1e6,
			Burst:      1e6,
		},
		ServiceMeshOptions: &servicemesh.ServiceMeshOptions{
			IstioPilotHost:  "http://istio-pilot.istio-system.svc:9090",
			JaegerQueryHost: "http://jaeger-query.istio-system.svc:80",
		},
		LinkedcareOptions: &linkedcare.LinkedcareOptions{
			APIServer:     "http://ks-apiserver.linkedcare-system.svc",
			AccountServer: "http://ks-account.linkedcare-system.svc",
		},
	}
	return conf
}

func saveTestConfig(t *testing.T, conf *Config) {
	content, err := yaml.Marshal(conf)
	if err != nil {
		t.Fatalf("error marshal config. %v", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s.yaml", DefaultConfigurationName), content, 0640)
	if err != nil {
		t.Fatalf("error write configuration file, %v", err)
	}
}

func cleanTestConfig(t *testing.T) {
	file := fmt.Sprintf("%s.yaml", DefaultConfigurationName)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Log("file not exists, skipping")
		return
	}

	err := os.Remove(file)
	if err != nil {
		t.Fatalf("remove %s file failed", file)
	}

}

func TestGet(t *testing.T) {
	conf := newTestConfig()
	saveTestConfig(t, conf)
	defer cleanTestConfig(t)

	err := Load()
	if err != nil {
		t.Fatal(err)
	}
	conf2 := Get()

	if diff := reflectutils.Equal(conf, conf2); diff != nil {
		t.Fatal(diff)
	}
}

func TestLinkedcareOptions(t *testing.T) {
	conf := newTestConfig()

	t.Run("save nil linkedcare options", func(t *testing.T) {
		savedConf := *conf
		savedConf.LinkedcareOptions = nil
		saveTestConfig(t, &savedConf)
		defer cleanTestConfig(t)

		err := Load()
		if err != nil {
			t.Fatal(err)
		}
		loadedConf := Get()

		if diff := reflectutils.Equal(conf, loadedConf); diff != nil {
			t.Fatal(diff)
		}
	})

	t.Run("save partially linkedcare options", func(t *testing.T) {
		savedConf := *conf
		savedConf.LinkedcareOptions.APIServer = "http://example.com"
		savedConf.LinkedcareOptions.AccountServer = ""

		saveTestConfig(t, &savedConf)
		defer cleanTestConfig(t)

		err := Load()
		if err != nil {
			t.Fatal(err)
		}
		loadedConf := Get()

		savedConf.LinkedcareOptions.AccountServer = "http://ks-account.linkedcare-system.svc"

		if diff := reflectutils.Equal(&savedConf, loadedConf); diff != nil {
			t.Fatal(diff)
		}
	})
}
