# HotPot: Cooking Your Cluster to Perfection 🍲

`HotPot` is your go-to CLI utility that marries the simplicity of cooking with the robustness of Kubernetes deployments. Drawing inspiration from crafting and culinary arts, HotPot serves up k3s clusters based on your specific recipe (configuration). 

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
- [x] Setup Argocd and configure applications, projects, and repositories

...And much more!


## Installation
```bash
curl -sfL https://raw.githubusercontent.com/zcubbs/hotpot/main/scripts/install/install.sh | bash
```

## Configuration

### Setup Let's Encrypt DNS Provider

### Setup Let's Encrypt DNS Provider

| Provider  | Environment Variables                                                                                        | Example Setting                                                       | Recipe Preparation                                       |
|-----------|--------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------|----------------------------------------------------------|
| **OVH**   | `OVH_ENDPOINT`, `OVH_APPLICATION_KEY`, `OVH_APPLICATION_SECRET`, `OVH_CONSUMER_KEY`                          | ```bash export OVH_ENDPOINT=ovh-eu export OVH_APPLICATION_KEY=... ``` | ```yaml ... traefik: ... dnsChallengeProvider: ovh ```   |
| **Azure** | `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`, `AZURE_SUBSCRIPTION_ID`, `AZURE_TENANT_ID`, `AZURE_RESOURCE_GROUP` | ```bash export AZURE_CLIENT_ID=... ```                                | ```yaml ... traefik: ... dnsChallengeProvider: azure ``` |


