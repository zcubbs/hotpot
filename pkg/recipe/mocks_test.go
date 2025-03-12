package recipe

import (
	"github.com/zcubbs/hotpot/pkg/go-k8s/argocd"
	"github.com/zcubbs/hotpot/pkg/go-k8s/certmanager"
	"github.com/zcubbs/hotpot/pkg/go-k8s/k3s"
	"github.com/zcubbs/hotpot/pkg/go-k8s/rancher"
	"github.com/zcubbs/hotpot/pkg/go-k8s/traefik"
)

type mockSystemInfo struct {
	osErr   error
	archErr error
	ramErr  error
	cpuErr  error
	diskErr error
	curlErr error
}

func (m *mockSystemInfo) IsOS(_ string) error                 { return m.osErr }
func (m *mockSystemInfo) IsArchIn(_ []string) error           { return m.archErr }
func (m *mockSystemInfo) IsRAMEnough(_ string) error          { return m.ramErr }
func (m *mockSystemInfo) IsCPUEnough(_ int) error             { return m.cpuErr }
func (m *mockSystemInfo) IsDiskSpaceEnough(_, _ string) error { return m.diskErr }
func (m *mockSystemInfo) IsCurlOK(_ []string) error           { return m.curlErr }

type mockK3sManager struct {
	installErr   error
	uninstallErr error
}

func (m *mockK3sManager) Install(_ k3s.Config, _ bool) error { return m.installErr }
func (m *mockK3sManager) Uninstall(_ bool) error             { return m.uninstallErr }

type mockHelmManager struct {
	isInstalledResult bool
	isInstalledErr    error
	installErr        error
}

func (m *mockHelmManager) IsHelmInstalled() (bool, error) {
	return m.isInstalledResult, m.isInstalledErr
}
func (m *mockHelmManager) InstallCli(_ bool) error { return m.installErr }

type mockCertManager struct {
	installErr   error
	uninstallErr error
}

func (m *mockCertManager) Install(_ certmanager.Values, _ string, _ bool) error { return m.installErr }
func (m *mockCertManager) Uninstall(_ string, _ bool) error                     { return m.uninstallErr }

type mockTraefikManager struct {
	installErr   error
	uninstallErr error
}

func (m *mockTraefikManager) Install(_ traefik.Values, _ string, _ bool) error { return m.installErr }
func (m *mockTraefikManager) Uninstall(_ string, _ bool) error                 { return m.uninstallErr }

type mockArgoCDManager struct {
	installErr    error
	uninstallErr  error
	createProjErr error
	createAppErr  error
	createRepoErr error
}

func (m *mockArgoCDManager) Install(_ argocd.Values, _ string, _ bool) error { return m.installErr }
func (m *mockArgoCDManager) Uninstall(_ string, _ bool) error                { return m.uninstallErr }
func (m *mockArgoCDManager) CreateProject(_ argocd.Project, _ string, _ bool) error {
	return m.createProjErr
}
func (m *mockArgoCDManager) CreateApplication(_ argocd.Application, _ string, _ bool) error {
	return m.createAppErr
}
func (m *mockArgoCDManager) CreateRepository(_ argocd.Repository, _ string, _ bool) error {
	return m.createRepoErr
}

type mockRancherManager struct {
	installErr   error
	uninstallErr error
}

func (m *mockRancherManager) Install(_ rancher.Values, _ string, _ bool) error { return m.installErr }
func (m *mockRancherManager) Uninstall(_ string, _ bool) error                 { return m.uninstallErr }

type mockK9sManager struct {
	installErr error
}

func (m *mockK9sManager) Install(_ bool) error { return m.installErr }

type mockFileSystem struct {
	removeAllErr error
}

func (m *mockFileSystem) RemoveAll(_ string) error { return m.removeAllErr }
