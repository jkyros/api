package v1

import (
	configv1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// MachineConfigRoleLabelKey is metadata key in the MachineConfig. Specifies the node role that config should be applied to.
// For example: `master` or `worker`
const MachineConfigRoleLabelKey = "machineconfiguration.openshift.io/role"

// KubeletConfigRoleLabelPrefix is the label that must be present in the KubeletConfig CR
const KubeletConfigRoleLabelPrefix = "pools.operator.machineconfiguration.openshift.io/"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControllerConfig describes configuration for MachineConfigController.
// This is currently only used to drive the MachineConfig objects generated by the TemplateController.
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +openshift:compatibility-gen:internal
type ControllerConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// TODO(jkyros): inconsistent historical generation resulted in the controllerconfig CRD being
	// generated with all fields required, while everything else was generated with optional

	// +kubebuilder:validation:Required
	Spec ControllerConfigSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status ControllerConfigStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

// ControllerConfigSpec is the spec for ControllerConfig resource.
type ControllerConfigSpec struct {
	// clusterDNSIP is the cluster DNS IP address
	// +kubebuilder:validation:Required
	ClusterDNSIP string `json:"clusterDNSIP" protobuf:"bytes,1,opt,name=clusterDNSIP"`

	// cloudProviderConfig is the configuration for the given cloud provider
	// +kubebuilder:validation:Required
	CloudProviderConfig string `json:"cloudProviderConfig" protobuf:"bytes,2,opt,name=cloudProviderConfig"`

	// platform is deprecated, use Infra.Status.PlatformStatus.Type instead
	// +optional
	Platform string `json:"platform,omitempty" protobuf:"bytes,3,opt,name=platform"`

	// etcdDiscoveryDomain is deprecated, use Infra.Status.EtcdDiscoveryDomain instead
	// +optional
	EtcdDiscoveryDomain string `json:"etcdDiscoveryDomain,omitempty" protobuf:"bytes,4,opt,name=etcdDiscoveryDomain"`

	// TODO: Use string for CA data

	// kubeAPIServerServingCAData managed Kubelet to API Server Cert... Rotated automatically
	// +kubebuilder:validation:Required
	KubeAPIServerServingCAData []byte `json:"kubeAPIServerServingCAData" protobuf:"bytes,5,opt,name=kubeAPIServerServingCAData"`

	// rootCAData specifies the root CA data
	// +kubebuilder:validation:Required
	RootCAData []byte `json:"rootCAData" protobuf:"bytes,6,opt,name=rootCAData"`

	// cloudProvider specifies the cloud provider CA data
	// +kubebuilder:validation:Required
	// +nullable
	CloudProviderCAData []byte `json:"cloudProviderCAData" protobuf:"bytes,7,opt,name=cloudProviderCAData"`

	// additionalTrustBundle is a certificate bundle that will be added to the nodes
	// trusted certificate store.
	// +kubebuilder:validation:Required
	// +nullable
	AdditionalTrustBundle []byte `json:"additionalTrustBundle" protobuf:"bytes,8,opt,name=additionalTrustBundle"`

	// imageRegistryBundleUserData is Image Registry Data provided by the user
	// +optional
	ImageRegistryBundleUserData []ImageRegistryBundle `json:"imageRegistryBundleUserData"`

	// imageRegistryBundleData is the ImageRegistryData
	// +optional
	ImageRegistryBundleData []ImageRegistryBundle `json:"imageRegistryBundleData"`

	// TODO: Investigate using a ConfigMapNameReference for the PullSecret and OSImageURL

	// pullSecret is the default pull secret that needs to be installed
	// on all machines.
	// +optional
	PullSecret *corev1.ObjectReference `json:"pullSecret,omitempty" protobuf:"bytes,9,opt,name=pullSecret"`

	// images is map of images that are used by the controller to render templates under ./templates/
	// +kubebuilder:validation:Required
	Images map[string]string `json:"images" protobuf:"bytes,10,rep,name=images"`

	// BaseOSContainerImage is the new-format container image for operating system updates.
	// +kubebuilder:validation:Required
	BaseOSContainerImage string `json:"baseOSContainerImage" protobuf:"bytes,11,opt,name=baseOSContainerImage"`

	// BaseOSExtensionsContainerImage is the matching extensions container for the new-format container
	// +optional
	BaseOSExtensionsContainerImage string `json:"baseOSExtensionsContainerImage" protobuf:"bytes,12,opt,name=baseOSExtensionsContainerImage"`

	// OSImageURL is the old-format container image that contains the OS update payload.
	// +optional
	OSImageURL string `json:"osImageURL" protobuf:"bytes,13,opt,name=osImageURL"`

	// releaseImage is the image used when installing the cluster
	// +kubebuilder:validation:Required
	ReleaseImage string `json:"releaseImage" protobuf:"bytes,14,opt,name=releaseImage"`

	// proxy holds the current proxy configuration for the nodes
	// +kubebuilder:validation:Required
	// +nullable
	Proxy *configv1.ProxyStatus `json:"proxy" protobuf:"bytes,15,opt,name=proxy"`

	// infra holds the infrastructure details
	// +kubebuilder:validation:EmbeddedResource
	// +kubebuilder:validation:Required
	// +nullable
	Infra *configv1.Infrastructure `json:"infra" protobuf:"bytes,16,opt,name=infra"`

	// dns holds the cluster dns details
	// +kubebuilder:validation:EmbeddedResource
	// +kubebuilder:validation:Required
	// +nullable
	DNS *configv1.DNS `json:"dns" protobuf:"bytes,17,opt,name=dns"`

	// ipFamilies indicates the IP families in use by the cluster network
	// +kubebuilder:validation:Required
	IPFamilies IPFamiliesType `json:"ipFamilies" protobuf:"bytes,18,opt,name=ipFamilies,casttype=IPFamiliesType"`

	// networkType holds the type of network the cluster is using
	// XXX: this is temporary and will be dropped as soon as possible in favor of a better support
	// to start network related services the proper way.
	// Nobody is also changing this once the cluster is up and running the first time, so, disallow
	// regeneration if this changes.
	// +optional
	NetworkType string `json:"networkType,omitempty" protobuf:"bytes,19,opt,name=networkType"`

	// Network contains additional network related information
	// +kubebuilder:validation:Required
	// +nullable
	Network *NetworkInfo `json:"network" protobuf:"bytes,20,opt,name=network"`
}

// ImageRegustryBundle contains information for writing image registry certificates
type ImageRegistryBundle struct {
	// file holds the name of the file where the bundle will be written to disk
	// +kubebuilder:validation:Required
	File string `json:"file"`
	// data holds the contents of the bundle that will be written to the file location
	// +kubebuilder:validation:Required
	Data []byte `json:"data"`
}

// IPFamiliesType indicates whether the cluster network is IPv4-only, IPv6-only, or dual-stack
type IPFamiliesType string

const (
	IPFamiliesIPv4                 IPFamiliesType = "IPv4"
	IPFamiliesIPv6                 IPFamiliesType = "IPv6"
	IPFamiliesDualStack            IPFamiliesType = "DualStack"
	IPFamiliesDualStackIPv6Primary IPFamiliesType = "DualStackIPv6Primary"
)

// Network contains network related configuration
type NetworkInfo struct {
	// MTUMigration contains the MTU migration configuration.
	// +kubebuilder:validation:Required
	// +nullable
	MTUMigration *configv1.MTUMigration `json:"mtuMigration" protobuf:"bytes,1,rep,name=mtuMigration"`
}

// ControllerConfigStatus is the status for ControllerConfig
type ControllerConfigStatus struct {
	// observedGeneration represents the generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// conditions represents the latest available observations of current state.
	// +optional
	Conditions []ControllerConfigStatusCondition `json:"conditions" protobuf:"bytes,2,rep,name=conditions"`
	// controllerCertificates represents the latest available observations of the automatically rotating certificates in the MCO.
	// +optional
	ControllerCertificates []ControllerCertificate `json:"controllerCertificates"`
}

// ControllerCertificate contains info about a specific cert.
type ControllerCertificate struct {
	// subject is the cert subject
	// +kubebuilder:validation:Required
	Subject string `json:"subject"`

	// signer is the  cert Issuer
	// +kubebuilder:validation:Required
	Signer string `json:"signer"`

	// notBefore is the lower boundary for validity
	// +kubebuilder:validation:Required
	NotBefore string `json:"notBefore"`

	// notAfter is the upper boundary for validity
	// +kubebuilder:validation:Required
	NotAfter string `json:"notAfter"`

	// bundleFile is the larger bundle a cert comes from
	// +kubebuilder:validation:Required
	BundleFile string `json:"bundleFile"`
}

// ControllerConfigStatusCondition contains condition information for ControllerConfigStatus
type ControllerConfigStatusCondition struct {
	// type specifies the state of the operator's reconciliation functionality.
	// +kubebuilder:validation:Required
	Type ControllerConfigStatusConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ControllerConfigStatusConditionType"`

	// status of the condition, one of True, False, Unknown.
	// +kubebuilder:validation:Required
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`

	// lastTransitionTime is the time of the last update to the current status object.
	// +kubebuilder:validation:Required
	// +nullable
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,3,opt,name=lastTransitionTime"`

	// reason is the reason for the condition's last transition.  Reasons are PascalCase
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`

	// message provides additional information about the current condition.
	// This is only to be consumed by humans.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// ControllerConfigStatusConditionType valid conditions of a ControllerConfigStatus
type ControllerConfigStatusConditionType string

const (
	// TemplateControllerRunning means the template controller is currently running.
	TemplateControllerRunning ControllerConfigStatusConditionType = "TemplateControllerRunning"

	// TemplateControllerCompleted means the template controller has completed reconciliation.
	TemplateControllerCompleted ControllerConfigStatusConditionType = "TemplateControllerCompleted"

	// TemplateControllerFailing means the template controller is failing.
	TemplateControllerFailing ControllerConfigStatusConditionType = "TemplateControllerFailing"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControllerConfigList is a list of ControllerConfig resources
//
// Compatibility level 4: No compatibility is provided, the API can change at any point for any reason. These capabilities should not be used by applications needing long term support.
// +openshift:compatibility-gen:level=4
// +openshift:compatibility-gen:internal
type ControllerConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Items []ControllerConfig `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfig defines the configuration for a machine
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type MachineConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// +optional
	Spec MachineConfigSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// MachineConfigSpec is the spec for MachineConfig
type MachineConfigSpec struct {
	// OSImageURL specifies the remote location that will be used to
	// fetch the OS.
	// +optional
	OSImageURL string `json:"osImageURL" protobuf:"bytes,1,opt,name=osImageURL"`

	// BaseOSExtensionsContainerImage specifies the remote location that will be used
	// to fetch the extensions container matching a new-format OS image
	// +optional
	BaseOSExtensionsContainerImage string `json:"baseOSExtensionsContainerImage" protobuf:"bytes,7,opt,name=baseOSExtensionsContainerImage"`

	// Config is a Ignition Config object.
	// +optional
	Config runtime.RawExtension `json:"config" protobuf:"bytes,2,opt,name=config"`

	// kernelArguments contains a list of kernel arguments to be added
	// +nullable
	// +optional
	KernelArguments []string `json:"kernelArguments" protobuf:"bytes,3,rep,name=kernelArguments"`

	// extensions contains a list of additional features that can be enabled on host
	// +optional
	Extensions []string `json:"extensions" protobuf:"bytes,4,rep,name=extensions"`

	// fips controls FIPS mode
	// +optional
	FIPS bool `json:"fips" protobuf:"varint,5,opt,name=fips"`

	// kernelType contains which kernel we want to be running like default
	// (traditional), realtime.
	// +optional
	KernelType string `json:"kernelType" protobuf:"bytes,6,opt,name=kernelType"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfigList is a list of MachineConfig resources
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type MachineConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Items []MachineConfig `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfigPool describes a pool of MachineConfigs.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type MachineConfigPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// +kubebuilder:validation:Required
	Spec MachineConfigPoolSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status MachineConfigPoolStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

// MachineConfigPoolSpec is the spec for MachineConfigPool resource.
type MachineConfigPoolSpec struct {
	// machineConfigSelector specifies a label selector for MachineConfigs.
	// Refer https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/ on how label and selectors work.
	// +optional
	MachineConfigSelector *metav1.LabelSelector `json:"machineConfigSelector,omitempty" protobuf:"bytes,1,opt,name=machineConfigSelector"`

	// nodeSelector specifies a label selector for Machines
	// +optional
	NodeSelector *metav1.LabelSelector `json:"nodeSelector,omitempty" protobuf:"bytes,2,opt,name=nodeSelector"`

	// paused specifies whether or not changes to this machine config pool should be stopped.
	// This includes generating new desiredMachineConfig and update of machines.
	// +optional
	Paused bool `json:"paused" protobuf:"varint,3,opt,name=paused"`

	// maxUnavailable defines either an integer number or percentage
	// of nodes in the pool that can go Unavailable during an update.
	// This includes nodes Unavailable for any reason, including user
	// initiated cordons, failing nodes, etc. The default value is 1.
	//
	// A value larger than 1 will mean multiple nodes going unavailable during
	// the update, which may affect your workload stress on the remaining nodes.
	// You cannot set this value to 0 to stop updates (it will default back to 1);
	// to stop updates, use the 'paused' property instead. Drain will respect
	// Pod Disruption Budgets (PDBs) such as etcd quorum guards, even if
	// maxUnavailable is greater than one.
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty" protobuf:"bytes,4,opt,name=maxUnavailable"`

	// The targeted MachineConfig object for the machine config pool.
	// +optional
	Configuration MachineConfigPoolStatusConfiguration `json:"configuration" protobuf:"bytes,5,opt,name=configuration"`
}

// MachineConfigPoolStatus is the status for MachineConfigPool resource.
type MachineConfigPoolStatus struct {
	// observedGeneration represents the generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// configuration represents the current MachineConfig object for the machine config pool.
	// +optional
	Configuration MachineConfigPoolStatusConfiguration `json:"configuration" protobuf:"bytes,2,opt,name=configuration"`

	// machineCount represents the total number of machines in the machine config pool.
	// +optional
	MachineCount int32 `json:"machineCount" protobuf:"varint,3,opt,name=machineCount"`

	// updatedMachineCount represents the total number of machines targeted by the pool that have the CurrentMachineConfig as their config.
	// +optional
	UpdatedMachineCount int32 `json:"updatedMachineCount" protobuf:"varint,4,opt,name=updatedMachineCount"`

	// readyMachineCount represents the total number of ready machines targeted by the pool.
	// +optional
	ReadyMachineCount int32 `json:"readyMachineCount" protobuf:"varint,5,opt,name=readyMachineCount"`

	// unavailableMachineCount represents the total number of unavailable (non-ready) machines targeted by the pool.
	// A node is marked unavailable if it is in updating state or NodeReady condition is false.
	// +optional
	UnavailableMachineCount int32 `json:"unavailableMachineCount" protobuf:"varint,6,opt,name=unavailableMachineCount"`

	// degradedMachineCount represents the total number of machines marked degraded (or unreconcilable).
	// A node is marked degraded if applying a configuration failed..
	// +optional
	DegradedMachineCount int32 `json:"degradedMachineCount" protobuf:"varint,7,opt,name=degradedMachineCount"`

	// conditions represents the latest available observations of current state.
	// +optional
	Conditions []MachineConfigPoolCondition `json:"conditions" protobuf:"bytes,8,rep,name=conditions"`

	// certExpirys keeps track of important certificate expiration data
	// +optional
	CertExpirys []CertExpiry `json:"certExpirys"`
}

// ceryExpiry contains the bundle name and the expiry date
type CertExpiry struct {
	// bundle is the name of the bundle in which the subject certificate resides
	// +kubebuilder:validation:Required
	Bundle string `json:"bundle"`
	// subject is the subject of the certificate
	// +kubebuilder:validation:Required
	Subject string `json:"subject"`
	// expiry is the date after which the certificate will no longer be valid
	// +kubebuilder:validation:Required
	Expiry string `json:"expiry"`
}

// MachineConfigPoolStatusConfiguration stores the current configuration for the pool, and
// optionally also stores the list of MachineConfig objects used to generate the configuration.
type MachineConfigPoolStatusConfiguration struct {
	corev1.ObjectReference `json:",inline"  protobuf:"bytes,1,rep,name=objectReference"`

	// source is the list of MachineConfig objects that were used to generate the single MachineConfig object specified in `content`.
	// +optional
	Source []corev1.ObjectReference `json:"source,omitempty" protobuf:"bytes,2,rep,name=source"`
}

// MachineConfigPoolCondition contains condition information for an MachineConfigPool.
type MachineConfigPoolCondition struct {
	// type of the condition, currently ('Done', 'Updating', 'Failed').
	// +optional
	Type MachineConfigPoolConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=MachineConfigPoolConditionType"`

	// status of the condition, one of ('True', 'False', 'Unknown').
	// +optional
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`

	// lastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	// +nullable
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,3,opt,name=lastTransitionTime"`

	// reason is a brief machine readable explanation for the condition's last
	// transition.
	// +optional
	Reason string `json:"reason" protobuf:"bytes,4,opt,name=reason"`

	// message is a human readable description of the details of the last
	// transition, complementing reason.
	// +optional
	Message string `json:"message" protobuf:"bytes,5,opt,name=message"`
}

// MachineConfigPoolConditionType valid conditions of a MachineConfigPool
type MachineConfigPoolConditionType string

const (
	// MachineConfigPoolUpdated means MachineConfigPool is updated completely.
	// When the all the machines in the pool are updated to the correct machine config.
	MachineConfigPoolUpdated MachineConfigPoolConditionType = "Updated"

	// MachineConfigPoolUpdating means MachineConfigPool is updating.
	// When at least one of machine is not either not updated or is in the process of updating
	// to the desired machine config.
	MachineConfigPoolUpdating MachineConfigPoolConditionType = "Updating"

	// MachineConfigPoolNodeDegraded means the update for one of the machine is not progressing
	MachineConfigPoolNodeDegraded MachineConfigPoolConditionType = "NodeDegraded"

	// MachineConfigPoolRenderDegraded means the rendered configuration for the pool cannot be generated because of an error
	MachineConfigPoolRenderDegraded MachineConfigPoolConditionType = "RenderDegraded"

	// MachineConfigPoolDegraded is the overall status of the pool based, today, on whether we fail with NodeDegraded or RenderDegraded
	MachineConfigPoolDegraded MachineConfigPoolConditionType = "Degraded"

	MachineConfigPoolBuildPending MachineConfigPoolConditionType = "BuildPending"

	MachineConfigPoolBuilding MachineConfigPoolConditionType = "Building"

	MachineConfigPoolBuildSuccess MachineConfigPoolConditionType = "BuildSuccess"

	MachineConfigPoolBuildFailed MachineConfigPoolConditionType = "BuildFailed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineConfigPoolList is a list of MachineConfigPool resources
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type MachineConfigPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Items []MachineConfigPool `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeletConfig describes a customized Kubelet configuration.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type KubeletConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// +kubebuilder:validation:Required
	Spec KubeletConfigSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status KubeletConfigStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

// KubeletConfigSpec defines the desired state of KubeletConfig
type KubeletConfigSpec struct {
	// +optional
	AutoSizingReserved *bool `json:"autoSizingReserved,omitempty" protobuf:"varint,1,opt,name=autoSizingReserved"`
	// +optional
	LogLevel *int32 `json:"logLevel,omitempty" protobuf:"varint,2,opt,name=logLevel"`

	// MachineConfigPoolSelector selects which pools the KubeletConfig shoud apply to.
	// A nil selector will result in no pools being selected.
	// +optional
	MachineConfigPoolSelector *metav1.LabelSelector `json:"machineConfigPoolSelector,omitempty" protobuf:"bytes,3,opt,name=machineConfigPoolSelector"`
	// kubeletConfig fields are defined in kubernetes upstream. Please refer to the types defined in the version/commit used by
	// OpenShift of the upstream kubernetes. It's important to note that, since the fields of the kubelet configuration are directly fetched from
	// upstream the validation of those values is handled directly by the kubelet. Please refer to the upstream version of the relevant kubernetes
	// for the valid values of these fields. Invalid values of the kubelet configuration fields may render cluster nodes unusable.
	// +optional
	KubeletConfig *runtime.RawExtension `json:"kubeletConfig,omitempty" protobuf:"bytes,4,opt,name=kubeletConfig"`

	// If unset, the default is based on the apiservers.config.openshift.io/cluster resource.
	// Note that only Old and Intermediate profiles are currently supported, and
	// the maximum available MinTLSVersions is VersionTLS12.
	// +optional
	TLSSecurityProfile *configv1.TLSSecurityProfile `json:"tlsSecurityProfile,omitempty" protobuf:"bytes,5,opt,name=tlsSecurityProfile"`
}

// KubeletConfigStatus defines the observed state of a KubeletConfig
type KubeletConfigStatus struct {
	// observedGeneration represents the generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// conditions represents the latest available observations of current state.
	// +optional
	Conditions []KubeletConfigCondition `json:"conditions" protobuf:"bytes,2,rep,name=conditions"`
}

// KubeletConfigCondition defines the state of the KubeletConfig
type KubeletConfigCondition struct {
	// type specifies the state of the operator's reconciliation functionality.
	// +optional
	Type KubeletConfigStatusConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=KubeletConfigStatusConditionType"`

	// status of the condition, one of True, False, Unknown.
	// +optional
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`

	// lastTransitionTime is the time of the last update to the current status object.
	// +optional
	// +nullable
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,3,opt,name=lastTransitionTime"`

	// reason is the reason for the condition's last transition.  Reasons are PascalCase
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`

	// message provides additional information about the current condition.
	// This is only to be consumed by humans.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// KubeletConfigStatusConditionType is the state of the operator's reconciliation functionality.
type KubeletConfigStatusConditionType string

const (
	// KubeletConfigSuccess designates a successful application of a KubeletConfig CR.
	KubeletConfigSuccess KubeletConfigStatusConditionType = "Success"

	// KubeletConfigFailure designates a failure applying a KubeletConfig CR.
	KubeletConfigFailure KubeletConfigStatusConditionType = "Failure"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeletConfigList is a list of KubeletConfig resources
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type KubeletConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Items []KubeletConfig `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContainerRuntimeConfig describes a customized Container Runtime configuration.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type ContainerRuntimeConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// +kubebuilder:validation:Required
	Spec ContainerRuntimeConfigSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status ContainerRuntimeConfigStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

// ContainerRuntimeConfigSpec defines the desired state of ContainerRuntimeConfig
type ContainerRuntimeConfigSpec struct {
	// MachineConfigPoolSelector selects which pools the ContainerRuntimeConfig shoud apply to.
	// A nil selector will result in no pools being selected.
	// +optional
	MachineConfigPoolSelector *metav1.LabelSelector `json:"machineConfigPoolSelector,omitempty" protobuf:"bytes,1,opt,name=machineConfigPoolSelector"`

	// +kubebuilder:validation:Required
	ContainerRuntimeConfig *ContainerRuntimeConfiguration `json:"containerRuntimeConfig,omitempty" protobuf:"bytes,2,opt,name=containerRuntimeConfig"`
}

// ContainerRuntimeConfiguration defines the tuneables of the container runtime
type ContainerRuntimeConfiguration struct {
	// pidsLimit specifies the maximum number of processes allowed in a container
	// +optional
	PidsLimit *int64 `json:"pidsLimit,omitempty" protobuf:"varint,1,opt,name=pidsLimit"`

	// logLevel specifies the verbosity of the logs based on the level it is set to.
	// Options are fatal, panic, error, warn, info, and debug.
	// +optional
	LogLevel string `json:"logLevel,omitempty" protobuf:"bytes,2,opt,name=logLevel"`

	// logSizeMax specifies the Maximum size allowed for the container log file.
	// Negative numbers indicate that no size limit is imposed.
	// If it is positive, it must be >= 8192 to match/exceed conmon's read buffer.
	// +optional
	LogSizeMax resource.Quantity `json:"logSizeMax,omitempty" protobuf:"bytes,3,opt,name=logSizeMax"`

	// overlaySize specifies the maximum size of a container image.
	// This flag can be used to set quota on the size of container images. (default: 10GB)
	// +optional
	OverlaySize resource.Quantity `json:"overlaySize,omitempty" protobuf:"bytes,4,opt,name=overlaySize"`

	// defaultRuntime is the name of the OCI runtime to be used as the default.
	// +optional
	DefaultRuntime ContainerRuntimeDefaultRuntime `json:"defaultRuntime,omitempty" protobuf:"bytes,5,opt,name=defaultRuntime"`
}

type ContainerRuntimeDefaultRuntime string

const (
	ContainerRuntimeDefaultRuntimeEmpty   = ""
	ContainerRuntimeDefaultRuntimeRunc    = "runc"
	ContainerRuntimeDefaultRuntimeCrun    = "crun"
	ContainerRuntimeDefaultRuntimeDefault = ContainerRuntimeDefaultRuntimeRunc
)

// ContainerRuntimeConfigStatus defines the observed state of a ContainerRuntimeConfig
type ContainerRuntimeConfigStatus struct {
	// observedGeneration represents the generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// conditions represents the latest available observations of current state.
	// +optional
	Conditions []ContainerRuntimeConfigCondition `json:"conditions" protobuf:"bytes,2,rep,name=conditions"`
}

// ContainerRuntimeConfigCondition defines the state of the ContainerRuntimeConfig
type ContainerRuntimeConfigCondition struct {
	// type specifies the state of the operator's reconciliation functionality.
	// +optional
	Type ContainerRuntimeConfigStatusConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ContainerRuntimeConfigStatusConditionType"`

	// status of the condition, one of True, False, Unknown.
	// +optional
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`

	// lastTransitionTime is the time of the last update to the current status object.
	// +nullable
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,3,opt,name=lastTransitionTime"`

	// reason is the reason for the condition's last transition.  Reasons are PascalCase
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`

	// message provides additional information about the current condition.
	// This is only to be consumed by humans.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// ContainerRuntimeConfigStatusConditionType is the state of the operator's reconciliation functionality.
type ContainerRuntimeConfigStatusConditionType string

const (
	// ContainerRuntimeConfigSuccess designates a successful application of a ContainerRuntimeConfig CR.
	ContainerRuntimeConfigSuccess ContainerRuntimeConfigStatusConditionType = "Success"

	// ContainerRuntimeConfigFailure designates a failure applying a ContainerRuntimeConfig CR.
	ContainerRuntimeConfigFailure ContainerRuntimeConfigStatusConditionType = "Failure"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContainerRuntimeConfigList is a list of ContainerRuntimeConfig resources
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type ContainerRuntimeConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Items []ContainerRuntimeConfig `json:"items" protobuf:"bytes,2,rep,name=items"`
}
