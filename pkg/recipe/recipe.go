package recipe

type Recipe struct {
	Name       string `mapstructure:"name" json:"name" yaml:"name"`
	Kubeconfig string `mapstructure:"kubeconfig" json:"kubeconfig" yaml:"kubeconfig"`
	Debug      bool   `mapstructure:"debug" json:"debug" yaml:"debug"`

	Ingredients Ingredients `mapstructure:"ingredients" json:"ingredients" yaml:"ingredients"`
}

type Ingredients struct {
	Node        Node              `mapstructure:"node" json:"node" yaml:"node"`
	CertManager CertManagerConfig `mapstructure:"certManager" json:"certManager" yaml:"certManager"`
	Traefik     TraefikConfig     `mapstructure:"traefik" json:"traefik" yaml:"traefik"`
	K3s         K3sConfig         `mapstructure:"k3s" json:"k3s" yaml:"k3s"`
	ArgoCD      ArgoCDConfig      `mapstructure:"argocd" json:"argocd" yaml:"argocd"`
}

type CertManagerConfig struct {
	Enabled bool `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
}

type Node struct {
	Ip          string   `mapstructure:"ip" json:"ip" yaml:"ip"`
	MinDiskSize string   `mapstructure:"minDiskSize" json:"minDiskSize" yaml:"minDiskSize"`
	MinCpu      string   `mapstructure:"minCpu" json:"minCpuCores" yaml:"minCpuCores"`
	MinMemory   string   `mapstructure:"minMemory" json:"minMemorySize" yaml:"minMemorySize"`
	CheckUrls   []string `mapstructure:"checkUrls" json:"checkUrls" yaml:"checkUrls"`
}

type K3sConfig struct {
	Enabled                 bool     `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	Disable                 []string `mapstructure:"disable" json:"disable" yaml:"disable"`
	Version                 string   `mapstructure:"version" json:"version" yaml:"version"`
	TlsSan                  []string `mapstructure:"tlsSan" json:"tlsSan" yaml:"tlsSan"`
	DataDir                 string   `mapstructure:"dataDir" json:"dataDir" yaml:"dataDir"`
	DefaultLocalStoragePath string   `mapstructure:"defaultLocalStoragePath" json:"defaultLocalStoragePath" yaml:"defaultLocalStoragePath"`
	WriteKubeconfigMode     string   `mapstructure:"writeKubeconfigMode" json:"writeKubeconfigMode" yaml:"writeKubeconfigMode"`
	IsHA                    bool     `mapstructure:"isHA" json:"isHA" yaml:"isHA"`
	IsServer                bool     `mapstructure:"isServer" json:"isServer" yaml:"isServer"`
	ClusterToken            string   `mapstructure:"clusterToken" json:"clusterToken" yaml:"clusterToken"`
	ExtraArgs               []string `mapstructure:"extraArgs" json:"extraArgs" yaml:"extraArgs"`
	PurgeExisting           bool     `mapstructure:"purgeExisting" json:"purgeExisting" yaml:"purgeExisting"`
}

type TraefikConfig struct {
	Enabled                   bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	EndpointsWeb              string `mapstructure:"endpointsWeb" json:"endpointsWeb" yaml:"endpointsWeb"`
	EndpointsWebsecure        string `mapstructure:"endpointsWebsecure" json:"endpointsWebsecure" yaml:"endpointsWebsecure"`
	EnableAccessLog           bool   `mapstructure:"enableAccessLog" json:"enableAccessLog" yaml:"enableAccessLog"`
	EnableDashboard           bool   `mapstructure:"enableDashboard" json:"enableDashboard" yaml:"enableDashboard"`
	ForwardHeaders            bool   `mapstructure:"forwardHeaders" json:"forwardHeaders" yaml:"forwardHeaders"`
	ForwardHeadersInsecure    bool   `mapstructure:"forwardHeadersInsecure" json:"forwardHeadersInsecure" yaml:"forwardHeadersInsecure"`
	ForwardHeadersTrustedIPs  string `mapstructure:"forwardHeadersTrustedIPs" json:"forwardHeadersTrustedIPs" yaml:"forwardHeadersTrustedIPs"`
	ProxyProtocol             bool   `mapstructure:"proxyProtocol" json:"proxyProtocol" yaml:"proxyProtocol"`
	ProxyProtocolEdgeIp       string `mapstructure:"proxyProtocolEdgeIp" json:"proxyProtocolEdgeIp" yaml:"proxyProtocolEdgeIp"`
	ProxyProtocolInsecure     bool   `mapstructure:"proxyProtocolInsecure" json:"proxyProtocolInsecure" yaml:"proxyProtocolInsecure"`
	ProxyProtocolTrustedIPs   string `mapstructure:"proxyProtocolTrustedIPs" json:"proxyProtocolTrustedIPs" yaml:"proxyProtocolTrustedIPs"`
	DnsChallenge              bool   `mapstructure:"dnsChallenge" json:"dnsChallenge" yaml:"dnsChallenge"`
	DnsChallengeProvider      string `mapstructure:"dnsChallengeProvider" json:"dnsChallengeProvider" yaml:"dnsChallengeProvider"`
	DnsChallengeDelay         int    `mapstructure:"dnsChallengeDelay" json:"dnsChallengeDelay" yaml:"dnsChallengeDelay"`
	DnsChallengeResolverEmail string `mapstructure:"dnsChallengeResolverEmail" json:"dnsChallengeResolverEmail" yaml:"dnsChallengeResolverEmail"`
	DnsChallengeTZ            string `mapstructure:"dnsChallengeTZ" json:"dnsChallengeTZ" yaml:"dnsChallengeTZ"`
	TransportInsecure         bool   `mapstructure:"transportInsecure" json:"transportInsecure" yaml:"transportInsecure"`
	Debug                     bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
	PurgeExisting             bool   `mapstructure:"purgeExisting" json:"purgeExisting" yaml:"purgeExisting"`
}

type ArgoCDConfig struct {
	Enabled       bool              `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	Projects      []Project         `mapstructure:"projects" json:"projects" yaml:"projects"`
	Credentials   ArgoCDCredentials `mapstructure:"credentials" json:"credentials" yaml:"credentials"`
	PurgeExisting bool              `mapstructure:"purgeExisting" json:"purgeExisting" yaml:"purgeExisting"`
}

type ArgoCDCredentials struct {
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	UseVault bool   `mapstructure:"useVault" json:"useVault" yaml:"useVault"`
	UseEnv   bool   `mapstructure:"useEnv" json:"useEnv" yaml:"useEnv"`
}

type Project struct {
	Name         string       `mapstructure:"name" json:"name" yaml:"name"`
	Repositories []GitopsRepo `mapstructure:"repositories" json:"repositories" yaml:"repositories"`
	Apps         []App        `mapstructure:"apps" json:"apps" yaml:"apps"`
}

type App struct {
	Name      string `mapstructure:"name" json:"name" yaml:"name"`
	Repo      string `mapstructure:"repo" json:"repo" yaml:"repo"`
	Revision  string `mapstructure:"revision" json:"revision" yaml:"revision"`
	Path      string `mapstructure:"path" json:"path" yaml:"path"`
	Namespace string `mapstructure:"namespace" json:"namespace" yaml:"namespace"`

	Charts    []AppChart    `mapstructure:"chart" json:"chart" yaml:"chart"`
	Manifests []AppManifest `mapstructure:"manifest" json:"manifest" yaml:"manifest"`
}

type AppChart struct {
	Repo        string   `mapstructure:"repo" json:"repo" yaml:"repo"`
	Revision    string   `mapstructure:"revision" json:"revision" yaml:"revision"`
	Path        string   `mapstructure:"path" json:"path" yaml:"path"`
	ValuesFiles []string `mapstructure:"valuesFiles" json:"valuesFiles" yaml:"valuesFiles"`
}

type AppManifest struct {
	Repo     string `mapstructure:"repo" json:"repo" yaml:"repo"`
	Revision string `mapstructure:"revision" json:"revision" yaml:"revision"`
	Path     string `mapstructure:"path" json:"path" yaml:"path"`
}

type GitopsRepo struct {
	Name        string                `mapstructure:"name" json:"name" yaml:"name"`
	Url         string                `mapstructure:"url" json:"url" yaml:"url"`
	Type        GitopsRepoType        `mapstructure:"type" json:"type" yaml:"type"`
	Credentials GitopsRepoCredentials `mapstructure:"credentials" json:"credentials" yaml:"credentials"`
}

type GitopsRepoCredentials struct {
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	UseVault bool   `mapstructure:"useVault" json:"useVault" yaml:"useVault"`
	UseEnv   bool   `mapstructure:"useEnv" json:"useEnv" yaml:"useEnv"`
}

type ContainerRegistryCredentials struct {
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Url      string `mapstructure:"url" json:"url" yaml:"url"`
	UseVault bool   `mapstructure:"useVault" json:"useVault" yaml:"useVault"`
	UseEnv   bool   `mapstructure:"useEnv" json:"useEnv" yaml:"useEnv"`
}

type GitopsRepoType string

const (
	GitopsRepoTypeHelm GitopsRepoType = "helm"
	GitopsRepoTypeGit  GitopsRepoType = "git"
)
