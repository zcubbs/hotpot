package host

import (
	"fmt"
	xos "github.com/zcubbs/x/os"
)

func IsOS(expected string) error {
	err := xos.IsOs(expected)
	if err != nil {
		return fmt.Errorf("failed to check os \n %w", err)
	}

	return nil
}

func IsArchIn(expected []string) error {
	err := xos.IsArchIn(expected)
	if err != nil {
		return fmt.Errorf("failed to check arch \n %w", err)
	}

	return nil
}

func IsRAMEnough(minRam string) error {
	minRamBytes, err := xos.StringToBytes(minRam)
	if err != nil {
		return fmt.Errorf("failed to parse ram \n %w", err)
	}

	err = xos.IsRAMEnough(minRamBytes)
	if err != nil {
		return fmt.Errorf("failed to check ram \n %w", err)
	}

	return nil
}

func IsCPUEnough(minCpu int) error {
	err := xos.IsCPUEnough(minCpu)
	if err != nil {
		return fmt.Errorf("failed to check cpu \n %w", err)
	}

	return nil
}

func IsDiskSpaceEnough(path, minDiskSize string) error {
	minDiskBytes, err := xos.StringToBytes(minDiskSize)
	if err != nil {
		return fmt.Errorf("failed to parse disk size \n %w", err)
	}

	err = xos.IsDiskSpaceEnoughForPath(path, minDiskBytes)
	if err != nil {
		return fmt.Errorf("failed to check disk space \n %w", err)
	}

	return nil
}

func IsCurlOK(urls []string) error {
	for _, v := range urls {
		ok := xos.IsCurlOK(v)
		if !ok {
			return fmt.Errorf("curl failed for %s", v)
		}
	}
	return nil
}
