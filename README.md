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
- [x] Check host prerequisites before creating a cluster, e.g. RAM, CPU, disk space, etc.
- [x] Setup and configure Helm
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
- [x] Recipe Sync Daemon
  - [x] Synchronize recipe files from Git repositories
  - [x] Support for private GitLab/GitHub repositories
  - [x] Token and SSH-based authentication
  - [x] Configurable sync frequency
  - [x] Systemd service integration

...And much more!


## Installation
> Supported platforms: `Linux`, `Mac`
```bash
curl -sfL https://raw.githubusercontent.com/zcubbs/hotpot/main/scripts/install.sh | bash
```

## Usage

### Cooking a Cluster

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

### Recipe Sync Daemon

The Recipe Sync Daemon allows you to keep your recipe files synchronized with a Git repository. It runs as a systemd service and can be configured using interactive prompts.

```bash
# Configure the sync daemon
> hotpot syncd config

🔧 Configuring hotpot-syncd...

Repository URL:
❯ https://github.com/user/repo

Branch:
❯ main

Auth Type (token/ssh):
❯ token

Token/SSH Key Path:
❯ ghp_xxxxxxxxxxxxxxxxxxxx

Local Path:
❯ /etc/hotpot/recipes/prod.yaml

Remote Path:
❯ recipes/prod.yaml

Sync Frequency (e.g., 5m, 1h):
❯ 5m

[ Submit ]

✅ Configuration saved successfully

# Enable and start the sync daemon
> hotpot syncd enable

🔌 Enabling hotpot-syncd service...
✅ Service enabled successfully

# Disable and stop the sync daemon
> hotpot syncd disable

🔌 Disabling hotpot-syncd service...
✅ Service disabled successfully
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

## Contributing

Contributions are welcome! If you find any issues, have suggestions, or would like to contribute code, please open an issue or a pull request on our GitHub page.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
