/*
 * Copyright 2024 The KubeVirt Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Common functions used by the plugin.
// Original code referenced from:
// https://github.com/kubevirt/kubevirt/blob/acafce385505da5862e1e54c4293370d97a8e845/pkg/virt-handler/device-manager/common.go

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	cdispec "github.com/container-orchestrated-devices/container-device-interface/specs-go"
)

type DeviceHandler interface {
	GetDeviceIOMMUGroup(basepath string, pciAddress string) (string, error)
	GetDeviceDriver(basepath string, pciAddress string) (string, error)
	GetDeviceNumaNode(basepath string, pciAddress string) (numaNode int)
	GetDevicePCIID(basepath string, pciAddress string) (string, error)
}

type deviceUtilsHandler struct{}

func NewDeviceHandler() DeviceHandler {
	return &deviceUtilsHandler{}
}

var Handler DeviceHandler

// getDeviceIOMMUGroup gets devices iommu_group
// e.g. /sys/bus/pci/devices/0000\:65\:00.0/iommu_group -> ../../../../../kernel/iommu_groups/45
func (h *deviceUtilsHandler) GetDeviceIOMMUGroup(basepath string, pciAddress string) (string, error) {
	iommuLink := filepath.Join(basepath, pciAddress, "iommu_group")
	iommuPath, err := os.Readlink(iommuLink)
	if err != nil {
		log.Printf("failed to read iommu_group link %s for device %s", iommuLink, pciAddress)
		return "", err
	}
	_, iommuGroup := filepath.Split(iommuPath)
	return iommuGroup, nil
}

// gets device driver
func (h *deviceUtilsHandler) GetDeviceDriver(basepath string, pciAddress string) (string, error) {
	driverLink := filepath.Join(basepath, pciAddress, "driver")
	driverPath, err := os.Readlink(driverLink)
	if err != nil {
		log.Printf("failed to read driver link %s for device %s", driverLink, pciAddress)
		return "", err
	}
	_, driver := filepath.Split(driverPath)
	return driver, nil
}

func (h *deviceUtilsHandler) GetDeviceNumaNode(basepath string, pciAddress string) (numaNode int) {
	numaNode = -1
	numaNodePath := filepath.Join(basepath, pciAddress, "numa_node")
	// #nosec No risk for path injection. Reading static path of NUMA node info
	numaNodeStr, err := os.ReadFile(numaNodePath)
	if err != nil {
		log.Printf("failed to read numa_node %s for device %s", numaNodePath, pciAddress)
		return
	}
	numaNodeStr = bytes.TrimSpace(numaNodeStr)
	numaNode, err = strconv.Atoi(string(numaNodeStr))
	if err != nil {
		log.Printf("failed to convert numa node value %v of device %s", numaNodeStr, pciAddress)
		return
	}
	return
}

func (h *deviceUtilsHandler) GetDevicePCIID(basepath string, pciAddress string) (string, error) {
	log.Printf("GetDevicePCIID: basepath: %s, pciAddress: %s", basepath, pciAddress)
	// #nosec No risk for path injection. Reading static path of PCI data
	file, err := os.Open(filepath.Join(basepath, pciAddress, "uevent"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PCI_ID") {
			equal := strings.Index(line, "=")
			value := strings.TrimSpace(line[equal+1:])
			return strings.ToLower(value), nil
		}
	}
	return "", fmt.Errorf("no pci_id is found")
}

func formatVFIODeviceSpecs(devID string) []*cdispec.DeviceNode {
	// always add /dev/vfio/vfio device as well
	devSpecs := make([]*cdispec.DeviceNode, 0)
	uid := uint32(107) //QEMU user UID
	gid := uint32(107) //QEMU user GID
	devSpecs = append(devSpecs, &cdispec.DeviceNode{
		HostPath:    vfioMount,
		Path:        vfioMount,
		Permissions: "mrw",
		UID:         &uid,
		GID:         &gid,
	})

	vfioDevice := filepath.Join(vfioDevicePath, devID)
	devSpecs = append(devSpecs, &cdispec.DeviceNode{
		HostPath:    vfioDevice,
		Path:        vfioDevice,
		Permissions: "mrw",
		UID:         &uid,
		GID:         &gid,
	})
	return devSpecs
}

func initHandler() {
	if Handler == nil {
		Handler = NewDeviceHandler()
	}
}
