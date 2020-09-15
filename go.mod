module linkedcare.io/linkedcare

go 1.14

require (
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/emicklei/go-restful v2.9.6+incompatible
	github.com/emicklei/go-restful-openapi v1.3.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.4.1
	github.com/gorilla/websocket v1.4.0
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/jaegertracing/jaeger v1.17.1 // indirect
	github.com/json-iterator/go v1.1.8
	github.com/jstemmer/gotags v1.4.1 // indirect
	github.com/jtblin/go-ldap-client v0.0.0-20170223121919-b73f66626b33 // indirect
	github.com/kiali/k-charted v0.5.0 // indirect
	github.com/kiali/kiali v1.17.0
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/openshift/api v3.9.0+incompatible // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/ldap.v2 v2.5.1 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v2 v2.2.8
	istio.io/api v0.0.0-20200616090052-c19f5f1ec54d
	istio.io/client-go v0.0.0-20200615164228-d77b0b53b6a0
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/component-base v0.18.2
	k8s.io/klog v1.0.0
	sigs.k8s.io/application v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.6.0
)

replace sigs.k8s.io/application => github.com/lic17/application v0.8.2-0.20200429022105-d7f404d91eca

replace k8s.io/client-go => ../client-go

replace k8s.io/apimachinery => ../apimachinery

replace linkedcare.io/linkedcare => ./

replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
