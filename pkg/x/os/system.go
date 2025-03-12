package os

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"os/exec"
)

func IsOs(expected string) error {
	info, err := host.Info()
	if err != nil {
		return err
	}

	if info.OS != expected {
		return fmt.Errorf("only %s is supported. found %s", expected, info.OS)
	}

	return nil
}

func GetOS() (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", err
	}
	return info.OS, nil
}

func IsArch(expected string) error {
	info, err := host.Info()
	if err != nil {
		return err
	}

	if info.KernelArch != expected {
		return fmt.Errorf("only %s arch is supported. found %s", expected, info.KernelArch)
	}

	return nil
}

func IsArchIn(expected []string) error {
	info, err := host.Info()
	if err != nil {
		return err
	}

	for _, v := range expected {
		if info.KernelArch == v {
			return nil
		}
	}

	return fmt.Errorf("only %s arch is supported. found %s", expected, info.KernelArch)
}

func IsArchAMD64() (bool, error) {
	info, err := host.Info()
	if err != nil {
		return false, err
	}
	return info.KernelArch == "amd64", nil
}

func GetArch() (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", err
	}
	return info.KernelArch, nil
}

func IsSupportedDistro(platform, version string) error {
	info, err := host.Info()
	if err != nil {
		return err
	}

	if info.Platform != platform || info.PlatformVersion != version {
		return fmt.Errorf("only %s %s is supported. found %s %s", platform, version, info.Platform, info.PlatformVersion)
	}

	return nil
}

func GetDistro() (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", err
	}
	return info.Platform + " " + info.PlatformVersion, nil
}

func IsRAMEnough(minRAM uint64) error { // minRAM in bytes
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	if v.Total < minRAM {
		return fmt.Errorf("minimum memory required is %s but found %s", BytesToString(minRAM), BytesToString(v.Total))
	}

	return nil
}

func GetRAM() (uint64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return v.Total, nil
}

func IsCPUEnough(minCPUs int) error {
	counts, err := cpu.Counts(true)
	if err != nil {
		return err
	}

	if counts < minCPUs {
		return fmt.Errorf("minimum cpu required is %d but found %d", minCPUs, counts)
	}

	return nil
}

func GetCPU() (int, error) {
	counts, err := cpu.Counts(true)
	if err != nil {
		return 0, err
	}
	return counts, nil
}

func IsDiskSpaceEnough(minSpace uint64) error { // minSpace in bytes
	usage, err := disk.Usage("/")
	if err != nil {
		return err
	}

	if usage.Free < minSpace {
		return fmt.Errorf("minimum disk space required is %s but found %s", BytesToString(minSpace), BytesToString(usage.Free))
	}

	return nil
}

func IsDiskSpaceEnoughForPath(path string, minSpace uint64) error { // minSpace in bytes
	usage, err := disk.Usage(path)
	if err != nil {
		return err
	}

	if usage.Free < minSpace {
		return fmt.Errorf("minimum disk space required is %s but found %s", BytesToString(minSpace), BytesToString(usage.Free))
	}

	return nil
}

func GetDiskSpaceForPath(path string) (uint64, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return 0, err
	}
	return usage.Free, nil
}

func GetDiskSpace() (uint64, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return 0, err
	}

	var total uint64
	for _, part := range parts {
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			continue
		}
		total += usage.Free
	}

	return total, nil
}

func IsTelnetOK(address string) bool {
	cmd := exec.Command("telnet", address)
	err := cmd.Run()
	return err == nil
}

func IsCurlOK(url string) bool {
	cmd := exec.Command("curl", "-s", url)
	err := cmd.Run()
	return err == nil
}

func IsSSHOK(ip string) bool {
	cmd := exec.Command("ssh", "-o", "BatchMode=yes", "-o", "ConnectTimeout=5", ip, "echo ok > /dev/null")
	err := cmd.Run()
	return err == nil
}
