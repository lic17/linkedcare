package constants

const (
	APIVersion = "v1alpha1"

	KubeSystemNamespace           = "kube-system"
	OpenPitrixNamespace           = "openpitrix-system"
	LinkedcareDevOpsNamespace     = "linkedcare-devops-system"
	IstioNamespace                = "istio-system"
	LinkedcareMonitoringNamespace = "linkedcare-monitoring-system"
	LinkedcareLoggingNamespace    = "linkedcare-logging-system"
	LinkedcareNamespace           = "linkedcare-system"
	LinkedcareControlNamespace    = "linkedcare-controls-system"
	PorterNamespace               = "porter-system"
	IngressControllerNamespace    = LinkedcareControlNamespace
	AdminUserName                 = "admin"
	DataHome                      = "/etc/linkedcare"
	IngressControllerFolder       = DataHome + "/ingress-controller"
	IngressControllerPrefix       = "linkedcare-router-"

	WorkspaceLabelKey              = "linkedcare.io/workspace"
	DisplayNameAnnotationKey       = "linkedcare.io/alias-name"
	DescriptionAnnotationKey       = "linkedcare.io/description"
	CreatorAnnotationKey           = "linkedcare.io/creator"
	System                         = "system"
	OpenPitrixRuntimeAnnotationKey = "openpitrix_runtime"
	WorkspaceAdmin                 = "workspace-admin"
	ClusterAdmin                   = "cluster-admin"
	WorkspaceRegular               = "workspace-regular"
	WorkspaceViewer                = "workspace-viewer"
	WorkspacesManager              = "workspaces-manager"
	DevopsOwner                    = "owner"
	DevopsReporter                 = "reporter"

	UserNameHeader = "X-Token-Username"

	TenantResourcesTag         = "Tenant Resources"
	IdentityManagementTag      = "Identity Management"
	AccessManagementTag        = "Access Management"
	NamespaceResourcesTag      = "Namespace Resources"
	ClusterResourcesTag        = "Cluster Resources"
	ComponentStatusTag         = "Component Status"
	OpenpitrixTag              = "Openpitrix Resources"
	VerificationTag            = "Verification"
	RegistryTag                = "Docker Registry"
	UserResourcesTag           = "User Resources"
	DevOpsProjectTag           = "DevOps Project"
	DevOpsProjectCredentialTag = "DevOps Project Credential"
	DevOpsProjectMemberTag     = "DevOps Project Member"
	DevOpsPipelineTag          = "DevOps Pipeline"
	DevOpsWebhookTag           = "DevOps Webhook"
	DevOpsJenkinsfileTag       = "DevOps Jenkinsfile"
	DevOpsScmTag               = "DevOps Scm"
	ClusterMetricsTag          = "Cluster Metrics"
	NodeMetricsTag             = "Node Metrics"
	NamespaceMetricsTag        = "Namespace Metrics"
	PodMetricsTag              = "Pod Metrics"
	PVCMetricsTag              = "PVC Metrics"
	ContainerMetricsTag        = "Container Metrics"
	WorkloadMetricsTag         = "Workload Metrics"
	WorkspaceMetricsTag        = "Workspace Metrics"
	ComponentMetricsTag        = "Component Metrics"
	LogQueryTag                = "Log Query"
	FluentBitSetting           = "Fluent Bit Setting"
)

var (
	//	WorkSpaceRoles   = []string{WorkspaceAdmin, WorkspaceRegular, WorkspaceViewer}
	SystemNamespaces = []string{LinkedcareNamespace, LinkedcareLoggingNamespace, LinkedcareMonitoringNamespace, OpenPitrixNamespace, KubeSystemNamespace, IstioNamespace, LinkedcareDevOpsNamespace, PorterNamespace}
)
