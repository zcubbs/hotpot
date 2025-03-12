package recipe

import (
	"errors"
	"testing"
)

func TestCheckPrerequisites(t *testing.T) {
	tests := []struct {
		name    string
		recipe  *Recipe
		sysInfo *mockSystemInfo
		wantErr bool
	}{
		{
			name: "all checks pass",
			recipe: &Recipe{
				Node: Node{
					SupportedOs:   []string{"linux"},
					SupportedArch: []string{"amd64"},
					MinMemory:     "4Gi",
					MinCpu:        2,
					MinDiskSize:   []Disk{{Path: "/", Size: "10Gi"}},
					Curl:          []string{"https://example.com"},
				},
			},
			sysInfo: &mockSystemInfo{},
			wantErr: false,
		},
		{
			name: "os check fails",
			recipe: &Recipe{
				Node: Node{
					SupportedOs: []string{"linux"},
				},
			},
			sysInfo: &mockSystemInfo{
				osErr: errors.New("unsupported OS"),
			},
			wantErr: true,
		},
		{
			name: "arch check fails",
			recipe: &Recipe{
				Node: Node{
					SupportedOs:   []string{"linux"},
					SupportedArch: []string{"amd64"},
				},
			},
			sysInfo: &mockSystemInfo{
				archErr: errors.New("unsupported architecture"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkPrerequisites(tt.recipe, tt.sysInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkPrerequisites() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstallK3s(t *testing.T) {
	tests := []struct {
		name    string
		recipe  *Recipe
		k3sMgr  *mockK3sManager
		helmMgr *mockHelmManager
		fs      *mockFileSystem
		wantErr bool
	}{
		{
			name: "successful installation",
			recipe: &Recipe{
				K3s: K3sConfig{
					Enabled: true,
				},
			},
			k3sMgr:  &mockK3sManager{},
			helmMgr: &mockHelmManager{},
			fs:      &mockFileSystem{},
			wantErr: false,
		},
		{
			name: "k3s install fails",
			recipe: &Recipe{
				K3s: K3sConfig{
					Enabled: true,
				},
			},
			k3sMgr:  &mockK3sManager{installErr: errors.New("install failed")},
			helmMgr: &mockHelmManager{},
			fs:      &mockFileSystem{},
			wantErr: true,
		},
		{
			name: "helm install fails",
			recipe: &Recipe{
				K3s: K3sConfig{
					Enabled: true,
				},
			},
			k3sMgr:  &mockK3sManager{},
			helmMgr: &mockHelmManager{isInstalledResult: false, installErr: errors.New("helm install failed")},
			fs:      &mockFileSystem{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := installK3s(tt.recipe, tt.k3sMgr, tt.helmMgr, tt.fs)
			if (err != nil) != tt.wantErr {
				t.Errorf("installK3s() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstallCertManager(t *testing.T) {
	tests := []struct {
		name    string
		recipe  *Recipe
		certMgr *mockCertManager
		wantErr bool
	}{
		{
			name: "successful installation",
			recipe: &Recipe{
				CertManager: CertManagerConfig{
					Enabled: true,
				},
			},
			certMgr: &mockCertManager{},
			wantErr: false,
		},
		{
			name: "installation fails",
			recipe: &Recipe{
				CertManager: CertManagerConfig{
					Enabled: true,
				},
			},
			certMgr: &mockCertManager{installErr: errors.New("install failed")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := installCertManager(tt.recipe, tt.certMgr)
			if (err != nil) != tt.wantErr {
				t.Errorf("installCertManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstallTraefik(t *testing.T) {
	tests := []struct {
		name       string
		recipe     *Recipe
		traefikMgr *mockTraefikManager
		wantErr    bool
	}{
		{
			name: "successful installation",
			recipe: &Recipe{
				Traefik: TraefikConfig{
					Enabled: true,
				},
			},
			traefikMgr: &mockTraefikManager{},
			wantErr:    false,
		},
		{
			name: "installation fails",
			recipe: &Recipe{
				Traefik: TraefikConfig{
					Enabled: true,
				},
			},
			traefikMgr: &mockTraefikManager{installErr: errors.New("install failed")},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := installTraefik(tt.recipe, tt.traefikMgr)
			if (err != nil) != tt.wantErr {
				t.Errorf("installTraefik() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstallArgoCD(t *testing.T) {
	tests := []struct {
		name      string
		recipe    *Recipe
		argocdMgr *mockArgoCDManager
		wantErr   bool
	}{
		{
			name: "successful installation",
			recipe: &Recipe{
				ArgoCD: ArgoCDConfig{
					Enabled: true,
				},
			},
			argocdMgr: &mockArgoCDManager{},
			wantErr:   false,
		},
		{
			name: "installation fails",
			recipe: &Recipe{
				ArgoCD: ArgoCDConfig{
					Enabled: true,
				},
			},
			argocdMgr: &mockArgoCDManager{installErr: errors.New("install failed")},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := installArgocd(tt.recipe, tt.argocdMgr)
			if (err != nil) != tt.wantErr {
				t.Errorf("installArgoCD() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstallRancher(t *testing.T) {
	tests := []struct {
		name       string
		recipe     *Recipe
		rancherMgr *mockRancherManager
		wantErr    bool
	}{
		{
			name: "successful installation",
			recipe: &Recipe{
				Rancher: RancherConfig{
					Enabled: true,
				},
			},
			rancherMgr: &mockRancherManager{},
			wantErr:    false,
		},
		{
			name: "installation fails",
			recipe: &Recipe{
				Rancher: RancherConfig{
					Enabled: true,
				},
			},
			rancherMgr: &mockRancherManager{installErr: errors.New("install failed")},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := installRancher(tt.recipe, tt.rancherMgr)
			if (err != nil) != tt.wantErr {
				t.Errorf("installRancher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstallK9s(t *testing.T) {
	tests := []struct {
		name    string
		recipe  *Recipe
		k9sMgr  *mockK9sManager
		wantErr bool
	}{
		{
			name: "successful installation",
			recipe: &Recipe{
				K9s: K9sConfig{
					Enabled: true,
				},
			},
			k9sMgr:  &mockK9sManager{},
			wantErr: false,
		},
		{
			name: "installation fails",
			recipe: &Recipe{
				K9s: K9sConfig{
					Enabled: true,
				},
			},
			k9sMgr:  &mockK9sManager{installErr: errors.New("install failed")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := installK9s(tt.recipe, tt.k9sMgr)
			if (err != nil) != tt.wantErr {
				t.Errorf("installK9s() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
