package cluster

type Cluster struct {
	Name     string    `mapstructure:"name" json:"name" yaml:"name"`
	Config   Config    `mapstructure:"cluster" json:"config" yaml:"config"`
	Projects []Project `mapstructure:"projects" json:"projects" yaml:"projects"`
}

type Config struct {
	NodeIp  string  `mapstructure:"nodeIp" json:"nodeIp" yaml:"nodeIp"`
	Traefik Traefik `mapstructure:"traefik" json:"traefik" yaml:"traefik"`
}

type Traefik struct {
	Enabled                   bool     `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	EnableAccessLog           bool     `mapstructure:"enableAccessLog" json:"enableAccessLog" yaml:"enableAccessLog"`
	EnableDashboard           bool     `mapstructure:"enableDashboard" json:"enableDashboard" yaml:"enableDashboard"`
	ForwardHeaders            bool     `mapstructure:"forwardHeaders" json:"forwardHeaders" yaml:"forwardHeaders"`
	ForwardHeadersInsecure    bool     `mapstructure:"forwardHeadersInsecure" json:"forwardHeadersInsecure" yaml:"forwardHeadersInsecure"`
	ForwardHeadersTrustedIPs  []string `mapstructure:"forwardHeadersTrustedIPs" json:"forwardHeadersTrustedIPs" yaml:"forwardHeadersTrustedIPs"`
	ProxyProtocol             bool     `mapstructure:"proxyProtocol" json:"proxyProtocol" yaml:"proxyProtocol"`
	ProxyProtocolEdgeIp       string   `mapstructure:"proxyProtocolEdgeIp" json:"proxyProtocolEdgeIp" yaml:"proxyProtocolEdgeIp"`
	ProxyProtocolInsecure     bool     `mapstructure:"proxyProtocolInsecure" json:"proxyProtocolInsecure" yaml:"proxyProtocolInsecure"`
	ProxyProtocolTrustedIPs   []string `mapstructure:"proxyProtocolTrustedIPs" json:"proxyProtocolTrustedIPs" yaml:"proxyProtocolTrustedIPs"`
	DnsChallenge              bool     `mapstructure:"dnsChallenge" json:"dnsChallenge" yaml:"dnsChallenge"`
	DnsChallengeProvider      string   `mapstructure:"dnsChallengeProvider" json:"dnsChallengeProvider" yaml:"dnsChallengeProvider"`
	DnsChallengeDelay         int      `mapstructure:"dnsChallengeDelay" json:"dnsChallengeDelay" yaml:"dnsChallengeDelay"`
	DnsChallengeResolverEmail string   `mapstructure:"dnsChallengeResolverEmail" json:"dnsChallengeResolverEmail" yaml:"dnsChallengeResolverEmail"`
	TransportInsecure         bool     `mapstructure:"transportInsecure" json:"transportInsecure" yaml:"transportInsecure"`
	Debug                     bool     `mapstructure:"debug" json:"debug" yaml:"debug"`
	PurgePreviousInstall      bool     `mapstructure:"purgePreviousInstall" json:"purgePreviousInstall" yaml:"purgePreviousInstall"`
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

type GitopsRepoType string

const (
	GitopsRepoTypeHelm GitopsRepoType = "helm"
	GitopsRepoTypeGit  GitopsRepoType = "git"
)
