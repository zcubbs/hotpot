# HotPot: Cooking Your Cluster to Perfection 🍲

`HotPot` is your go-to CLI utility that marries the simplicity of cooking with the robustness of Kubernetes deployments. Drawing inspiration from crafting and culinary arts, HotPot serves up k3s clusters based on your specific recipe (configuration). It aims to provide a reproducible, consistent, and reliable way to deploy your clusters and applications. It can also work with none k3s clusters by disabling the k3s feature. 

---
<p align="center">
</p>
<p align="center">
  <img width="850" src="docs/assets/splash.png">
</p>

---

## Features

- [x] Create a k3s cluster with yaml configuration
- [x] Delete a k3s cluster
- [x] Setup and configure Traefik
  - [x] Setup and configure Let's Encrypt
  - [x] Setup and configure CertManager
  - [x] Setup and configure IngressRoutes
  - [x] Configure support for DNS01 and HTTP01 challenges
  - [x] Configure Providers: Cloudflare, OVH, Azure
- [x] Setup and configure CertManager
- [x] Bootstrap Secrets: Container Registry Credentials, Generic Secrets
- [x] Setup Argocd and configure applications, projects, and repositories
- [x] Override any of the features above without recreating the cluster
- [x] Nuke a cluster

...And much more!


## Installation
```bash
curl -sfL https://raw.githubusercontent.com/zcubbs/hotpot/main/scripts/install/install.sh | bash
```

## Usage

```bash
> hotpot cook -r recipe.yaml

🍲 Cooking...
🍳 Checking prerequisites... 
    ├─ os: ok
    ├─ arch: ok
    ├─ ram: ok
    ├─ cpu: ok
    ├─ disk: ok
    ├─ curl: ok
    └─ prerequisites ok
🍕 Adding k3s... 
    └─ install ok
🍉 Adding helm cli... 
🌶️ Adding secrets... 
    ├─ container registry credentials: regcred 
    │  ├─ namespaces: [hub] ok
    │  └─ secret ok
    ├─ generic secret: my-secret 
    │  ├─ namespaces: hub ok
    │  └─ secret ok
    └─ secrets ok
🍙 Adding cert-manager... 
    └─ install ok
🍔 Adding traefik... 
    └─ install ok
🥪 Adding argocd... 
    ├─ argocd admin password: ok
    └─ install ok
🌭 Adding gitops... 
    ├─ project: hotpot ok
    │  ├─ repository: gitops-private-repo ok
    │  ├─ repository: helm-private-repo ok
    │  ├─ application: hub ok
    │  ├─ application: hub-manifests ok
    └─ gitops ok
 ok    completed
```

## Configuration

### ACME Providers (Let's Encrypt)

Refer to documentation: https://doc.traefik.io/traefik/https/acme/#providers

#### TLS Challenge using ALPN

> **Note**: TLS Challenge is not currently supported by CertManager. This is a Traefik only feature.

```yaml
traefik:
  tlsChallenge: true
```

#### DNS Challenge

To delegate ACME Challenges to CertManager, set `dnsChallenge` or `tlsChallenge` to `true` and configure the `certManager` section. And set `letsEncryptIngressClassResolver` to `traefik` in the `certManager` section. Also make sure Traefik is configured with `dnsChallenge` and `tlsChallenge` set to `false`.

Docs: https://cert-manager.io/docs/configuration/acme/

```yaml
traefik:
  tlsChallenge: false
  dnsChallenge: false
certManager:
  dnsChallengeEnabled: true
  dnsProvider: azure # ovh, azure, cloudflare or route53
  letsEncryptIngressClassResolver: traefik
```

#### Supported DNS Providers

| Provider  | Environment Variables                                                                                        | Recipe Config                    |
|-----------|--------------------------------------------------------------------------------------------------------------|----------------------------------|
| **OVH**   | `OVH_ENDPOINT`, `OVH_APPLICATION_KEY`, `OVH_APPLICATION_SECRET`, `OVH_CONSUMER_KEY`                          | `certManager.dnsProvider: ovh`   |
| **Azure** | `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`, `AZURE_SUBSCRIPTION_ID`, `AZURE_TENANT_ID`, `AZURE_RESOURCE_GROUP` | `certManager.dnsProvider: azure` |

> **Note**: future versions of HotPot will support AWS Route53, Cloudflare, and other DNS providers.

Example:
    
```yaml
certManager:
  dnsChallengeEnabled: true
  dnsProvider: azure
  dnsAzureClientID: env.HOTPOT_DNS_AZURE_CLIENT_ID
  dnsAzureClientSecret: env.HOTPOT_DNS_AZURE_CLIENT_SECRET
  dnsAzureHostedZoneName: example.com
  dnsAzureResourceGroupName: env.HOTPOT_DNS_AZURE_RESOURCE_GROUP_NAME
  dnsAzureSubscriptionID: env.HOTPOT_DNS_AZURE_SUBSCRIPTION_ID
  dnsAzureTenantID: env.HOTPOT_DNS_AZURE_TENANT_ID
```

**Note**: If you need to override CodeDNS Nameservers config (CoreDNS uses the default resolv.conf on the host), use this:

```yaml
certManager:
  dnsRecursiveNameservers:
    - 8.8.8.8:53
  dnsRecursiveNameserversOnly: true
```

