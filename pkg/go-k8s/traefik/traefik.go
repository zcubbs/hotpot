package traefik

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/zcubbs/hotpot/pkg/go-k8s/helm"
	"github.com/zcubbs/hotpot/pkg/go-k8s/kubernetes"
	"github.com/zcubbs/hotpot/pkg/secret"
	"github.com/zcubbs/hotpot/pkg/x/yaml"
	"os"
	"time"
)

const (
	traefikHelmRepoName = "traefik"
	traefikHelmRepoUrl  = "https://helm.traefik.io/traefik"
	traefikChartName    = "traefik"
	traefikChartVersion = "" // latest
	traefikNamespace    = "traefik"

	traefikDefaultResolver = "letsencrypt"

	traefikEndpointWeb       = "80"
	traefikEndpointWebsecure = "443"

	traefikDnsTZ = "Europe/Paris"
)

func Install(values Values, kubeconfig string, debug bool) error {
	if err := validateValues(&values); err != nil {
		return err
	}

	if values.DnsChallengeEnabled {
		if err := configureDNSChallengeVars(values, kubeconfig, debug); err != nil {
			return err
		}
	}

	valuesPath := getTmpFilePath("values")

	// create traefik values.yaml from template
	configFileContent, err := yaml.ApplyTmpl(traefikValuesTmpl, values, debug)
	if err != nil {
		return fmt.Errorf("failed to apply template \n %w", err)
	}

	// write tmp manifest
	err = os.WriteFile(valuesPath, configFileContent, 0600)
	if err != nil {
		return fmt.Errorf("failed to write traefik values.yaml \n %w", err)
	}

	// helm install traefik
	helmClient := helm.NewClient()
	helmClient.Settings.KubeConfig = kubeconfig
	helmClient.Settings.SetNamespace(traefikNamespace)
	helmClient.Settings.Debug = debug

	// add traefik helm repo
	err = helmClient.RepoAddAndUpdate(traefikHelmRepoName, traefikHelmRepoUrl)
	if err != nil {
		return fmt.Errorf("failed to add helm repo: %w", err)
	}

	// install traefik
	err = helmClient.InstallChart(helm.Chart{
		ChartName:       traefikChartName,
		ReleaseName:     traefikChartName,
		RepoName:        traefikHelmRepoName,
		Values:          nil,
		ValuesFiles:     []string{valuesPath},
		Debug:           debug,
		CreateNamespace: true,
		Upgrade:         true,
	})
	if err != nil {
		return fmt.Errorf("failed to install traefik \n %w", err)
	}

	// wait for traefik deployment to be ready
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	err = kubernetes.IsDeploymentReady(
		ctxWithTimeout,
		kubeconfig,
		traefikNamespace,
		[]string{"traefik"},
		debug,
	)
	if err != nil {
		return fmt.Errorf("failed to wait for traefik deployment to be ready \n %w", err)
	}

	// prepare default certificate secret
	if values.DefaultCertificateEnabled {
		err := createDefaultCertificateSecret(&values, kubeconfig, debug)
		if err != nil {
			return fmt.Errorf("failed to create default certificate secret \n %w", err)
		}

		// restart traefik
		err = kubernetes.RestartPods(kubeconfig, traefikNamespace, []string{"traefik"}, debug)
		if err != nil {
			return fmt.Errorf("failed to restart traefik \n %w", err)
		}
	}
	return nil
}

func Uninstall(kubeconfig string, debug bool) error {
	helmClient := helm.NewClient()
	helmClient.Settings.KubeConfig = kubeconfig
	helmClient.Settings.SetNamespace(traefikNamespace)
	helmClient.Settings.Debug = debug

	// delete traefik
	return helmClient.UninstallChart(traefikChartName)
}

func createDefaultCertificateSecret(values *Values, kubeconfig string, debug bool) error {
	// create namespace
	err := kubernetes.CreateNamespace(
		kubeconfig,
		[]string{traefikNamespace},
	)
	if err != nil {
		return fmt.Errorf("failed to create namespace %s \n %w", traefikNamespace, err)
	}

	cert, err := secret.Provide(values.DefaultCertificateCert)
	if err != nil {
		return fmt.Errorf("failed to provide default certificate crt \n %w", err)
	}

	key, err := secret.Provide(values.DefaultCertificateKey)
	if err != nil {
		return fmt.Errorf("failed to provide default certificate key \n %w", err)
	}

	values.DefaultCertificateCert = cert
	values.DefaultCertificateKey = key

	// apply template
	if err := applyDefaultCertificateSecret(*values, kubeconfig, debug); err != nil {
		return err
	}

	return nil
}

func applyDefaultCertificateSecret(values Values, _ string, debug bool) error {
	// apply default TLS store
	err := kubernetes.ApplyManifest(
		defaultTlsStoreTmpl,
		struct {
			Namespace string
		}{
			Namespace: traefikNamespace,
		},
		debug,
	)
	if err != nil {
		return fmt.Errorf("failed to apply default tls store \n %w", err)
	}

	// apply default TLS option
	if values.DefaultCertificateTlsOptionEnabled {
		err = kubernetes.ApplyManifest(
			defaultTlsOptionTmpl,
			struct {
				Namespace    string
				TlsStrictSNI bool
			}{
				Namespace:    traefikNamespace,
				TlsStrictSNI: values.TlsStrictSNI,
			},
			debug,
		)
		if err != nil {
			return fmt.Errorf("failed to apply default tls option \n %w", err)
		}
	}

	// read cert and key from paths
	cert, err := readFileContent(values.DefaultCertificateCert)
	if err != nil {
		return fmt.Errorf("failed to read default certificate crt \n %w", err)
	}
	key, err := readFileContent(values.DefaultCertificateKey)
	if err != nil {
		return fmt.Errorf("failed to read default certificate key \n %w", err)
	}

	cert = base64.StdEncoding.EncodeToString([]byte(cert))
	key = base64.StdEncoding.EncodeToString([]byte(key))

	// Add Default Certificate Secret
	if values.DefaultCertificateEnabled {
		err = kubernetes.ApplyManifest(
			defaultCertificateSecretTmpl,
			DefaultCertificateValues{
				Enabled: values.DefaultCertificateEnabled,
				Base64EncodedCertificate: struct {
					Crt string
					Key string
				}{
					Crt: cert,
					Key: key,
				},
				Namespace: traefikNamespace,
			}, debug)
		if err != nil {
			return fmt.Errorf("failed to apply default certificate secret \n %w", err)
		}
	}

	return nil
}

func readFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getTmpFilePath(name string) string {
	return os.TempDir() + "/" + name + "-" + time.Now().Format("20060102150405") + ".yaml"
}

func validateValues(values *Values) error {

	if values.IngressProvider != "" && values.DnsChallengeEnabled {
		return fmt.Errorf("can't set both ingressProvider and dnsProvider")
	} else if values.IngressProvider != "" && values.TlsChallengeEnabled {
		return fmt.Errorf("can't set both ingressProvider and tlsChallenge")
	} else if values.DnsChallengeEnabled && values.TlsChallengeEnabled {
		return fmt.Errorf("can't set both dnsChallenge and tlsChallenge")
	}

	if values.DnsChallengeEnabled {
		if values.DnsResolver == "" {
			values.DnsResolver = traefikDefaultResolver
		}

		if values.DnsResolverEmail == "" {
			return fmt.Errorf("dnsResolverEmail is required")
		}
	}

	if values.TlsChallengeEnabled {
		if values.TlsResolver == "" {
			values.TlsResolver = traefikDefaultResolver
		}

		if values.TlsResolverEmail == "" {
			return fmt.Errorf("tlsResolverEmail is required")
		}
	}

	if values.EndpointsWeb == "" {
		values.EndpointsWeb = traefikEndpointWeb
	}

	if values.EndpointsWebsecure == "" {
		values.EndpointsWebsecure = traefikEndpointWebsecure
	}

	if values.DnsTZ == "" {
		values.DnsTZ = traefikDnsTZ
	}

	if values.DefaultCertificateEnabled {
		if values.DefaultCertificateCert == "" {
			return fmt.Errorf("defaultCertificate.base64EncodedCertificate.crt is required")
		}
		if values.DefaultCertificateKey == "" {
			return fmt.Errorf("defaultCertificate.base64EncodedCertificate.key is required")
		}
	}

	return nil
}

type Values struct {
	AdditionalArguments                []string
	IngressProvider                    string
	TlsStrictSNI                       bool
	TlsChallengeEnabled                bool
	TlsResolver                        string
	TlsResolverEmail                   string
	DnsChallengeEnabled                bool
	DnsProvider                        string
	DnsResolver                        string
	DnsResolverIPs                     string
	DnsResolverEmail                   string
	EnableDashboard                    bool
	EnableAccessLog                    bool
	DebugLog                           bool
	EndpointsWeb                       string
	EndpointsWebsecure                 string
	ServersTransportInsecureSkipVerify bool
	ForwardedHeaders                   bool
	ForwardedHeadersInsecure           bool
	ForwardedHeadersTrustedIPs         string
	ProxyProtocol                      bool
	ProxyProtocolInsecure              bool
	ProxyProtocolTrustedIPs            string
	DnsTZ                              string

	DefaultCertificateEnabled          bool
	DefaultCertificateTlsOptionEnabled bool
	DefaultCertificateCert             string
	DefaultCertificateKey              string
}

var traefikValuesTmpl = `
additionalArguments:
  - "--global.checknewversion=false"
  - "--global.sendanonymoususage=false"
  - "--entrypoints.web.http.redirections.entrypoint.to=websecure"
  - "--entrypoints.web.http.redirections.entrypoint.scheme=https"
  - "--entrypoints.web.http.redirections.entrypoint.permanent=false"
  - "--entrypoints.web.http.redirections.entrypoint.priority=1"
  {{- if .DebugLog }}
  - "--log.level=DEBUG"
  {{- else }}
  - "--log.level=INFO"
  {{- end }}
  {{- if .EnableAccessLog }}
  - "--accesslog=true"
  {{- end }}
  {{- range $i, $arg := .AdditionalArguments }}
  - "{{ printf "%s" . }}"
  {{- end }}
  {{- if .ServersTransportInsecureSkipVerify }}
  - "--serversTransport.insecureSkipVerify"
  {{- end }}
  {{- if .ForwardedHeaders }}
  {{- if .ForwardedHeadersInsecure }}
  - "--entrypoints.websecure.forwardedHeaders.insecure"
  {{- end }}
  {{- if .ForwardedHeadersTrustedIPs }}
  - "--entrypoints.websecure.forwardedHeaders.trustedIPs=127.0.0.1/32,{{ .ForwardedHeadersTrustedIPs }}"
  - "--entrypoints.web.forwardedHeaders.trustedIPs=127.0.0.1/32,{{ .ForwardedHeadersTrustedIPs }}"
  {{- end }}
  {{- end }}
  {{- if .ProxyProtocol }}
  {{- if .ProxyProtocolInsecure }}
  - "--entrypoints.websecure.proxyProtocol.insecure"
  {{- end }}
  {{- if .ProxyProtocolTrustedIPs }}
  - "--entrypoints.websecure.proxyProtocol.trustedIPs=127.0.0.1/32,{{ .ProxyProtocolTrustedIPs }}"
  {{- end }}
  {{- end }}
  {{- if .TlsChallengeEnabled }}
  - "--certificatesresolvers.{{ .TlsResolver }}.acme.tlschallenge=true"
  - "--certificatesresolvers.{{ .TlsResolver }}.acme.storage=/data/acme.json"
  - "--certificatesresolvers.{{ .TlsResolver }}.acme.email={{ .TlsResolverEmail }}"
  - "--certificatesresolvers.{{ .TlsResolver }}.acme.caserver=https://acme-v02.api.letsencrypt.org/directory"
  - "--certificatesresolvers.{{ .TlsResolver }}-staging.acme.tlschallenge=true"
  - "--certificatesresolvers.{{ .TlsResolver }}-staging.acme.storage=/data/acme.json"
  - "--certificatesresolvers.{{ .TlsResolver }}-staging.acme.email={{ .TlsResolverEmail }}"
  - "--certificatesresolvers.{{ .TlsResolver }}-staging.acme.caserver=https://acme-staging-v02.api.letsencrypt.org/directory"
  {{- end }}
  {{- if .DnsChallengeEnabled }}
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.dnschallenge=true"
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.dnschallenge.provider={{ .DnsProvider }}"
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.dnschallenge.delayBeforeCheck=10"
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.email={{ .DnsResolverEmail }}"
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.storage=/data/acme.json"
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.caserver=https://acme-staging-v02.api.letsencrypt.org/directory"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.dnschallenge=true"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.dnschallenge.provider={{ .DnsProvider }}"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.dnschallenge.delayBeforeCheck=10"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.email={{ .DnsResolverEmail }}"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.storage=/data/acme.json"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.caserver=https://acme-v02.api.letsencrypt.org/directory"
  {{- if .DnsResolverIPs }}
  - "--certificatesresolvers.{{ .DnsResolver }}-staging.acme.dnschallenge.resolvers={{ .DnsResolverIPs }}"
  - "--certificatesresolvers.{{ .DnsResolver }}.acme.dnschallenge.resolvers={{ .DnsResolverIPs }}"
  {{- end }}
  {{- end }}

service:
  enabled: true
  type: LoadBalancer
rbac:
  enabled: true

ports:
  web:
    exposedPort: {{ .EndpointsWeb }}
  websecure:
    exposedPort: {{ .EndpointsWebsecure }}
    tls:
      enabled: true
      {{- if .DnsChallengeEnabled }}
      certResolver: {{ .DnsResolver }}
      {{- end }}

persistence:
  enabled: true
  accessMode: ReadWriteOnce
  size: 128Mi
  path: /data
  annotations: { }

ingressRoute:
  dashboard:
    enabled: {{ .EnableDashboard }}

pilot:
  enabled: false

{{- if .TlsChallengeEnabled }}
providers:
  kubernetesCRD: {}
{{- end }}

{{- if .IngressProvider }}
providers:
  kubernetesIngress:
    enabled: true
    ingressClass: {{ .IngressProvider }}
  kubernetesCRD:
    enabled: true
    ingressClass: {{ .IngressProvider }}
{{- end }}

deployment:
  initContainers:
    - name: volume-permissions
      image: busybox:1.31.1
      command: ["sh", "-c", "touch /data/acme.json; chmod -Rv 0600 /data/acme.json; cat /data/acme.json"]
      volumeMounts:
        - name: data
          mountPath: /data

{{- if .DnsChallengeEnabled }}
envFrom:
  - secretRef:
      name: traefik-dns-provider-credentials
{{- end }}
`

type DefaultCertificateValues struct {
	Enabled                  bool
	Base64EncodedCertificate struct {
		Crt string
		Key string
	}

	Namespace string
}

var defaultTlsStoreTmpl = `
apiVersion: traefik.io/v1alpha1
kind: TLSStore
metadata:
  name: default
  namespace: {{ .Namespace }}

spec:
  defaultCertificate:
    secretName: default-certificate
`

var defaultTlsOptionTmpl = `
apiVersion: traefik.io/v1alpha1
kind: TLSOption
metadata:
  name: default
  namespace: {{ .Namespace }}

spec:
  minVersion: VersionTLS12
  sniStrict: {{ .TlsStrictSNI }}
  cipherSuites:
    - TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	- TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	- TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	- TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	- TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	- TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305
`

var defaultCertificateSecretTmpl = `

---
apiVersion: v1
kind: Secret
metadata:
  name: default-certificate
  namespace: {{ .Namespace }}
type: Opaque
data:
  tls.crt: {{ .Base64EncodedCertificate.Crt }}
  tls.key: {{ .Base64EncodedCertificate.Key }}
`
