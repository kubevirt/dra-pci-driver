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

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	pciBasePath       = "/sys/bus/pci/devices"
	PCIResourcePrefix = "PCI_RESOURCE"
)

type PCIDevice struct {
	uuid           string
	vendorSelector string
	resourceName   string
	pciAddress     string
	driver         string
	iommuGroup     string
	numaNode       int
	pciID          string
}

func DiscoverPermittedHostPCIDevices(supportedPCIDeviceMap map[string]string) (map[string][]*PCIDevice, error) {
	initHandler()

	iommuToPCIMap := make(map[string]string)

	pciDevicesMap := make(map[string][]*PCIDevice)
	err := filepath.Walk(pciBasePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		pciID, err := Handler.GetDevicePCIID(pciBasePath, info.Name())
		if err != nil {
			log.Printf("Failed to get vendor:device ID for device: %s, error: %v", info.Name(), err)
			return nil
		}

		if _, supported := supportedPCIDeviceMap[pciID]; supported {
			driver, err := Handler.GetDeviceDriver(pciBasePath, info.Name())
			if err != nil || driver != "vfio-pci" {
				log.Printf("Driver error: %v", err)
				return nil
			}

			pcidev := &PCIDevice{
				uuid:           uuid.New().String(),
				pciID:          pciID,
				pciAddress:     info.Name(),
				vendorSelector: pciID,
				resourceName:   supportedPCIDeviceMap[pciID],
			}

			iommuGroup, err := Handler.GetDeviceIOMMUGroup(pciBasePath, info.Name())
			if err != nil {
				log.Printf("IOMMU group error: %v", err)
				return nil
			}
			pcidev.iommuGroup = iommuGroup
			pcidev.driver = driver
			pcidev.numaNode = Handler.GetDeviceNumaNode(pciBasePath, info.Name())

			iommuToPCIMap[pcidev.iommuGroup] = pcidev.pciAddress

			pciDevicesMap[pciID] = append(pciDevicesMap[pciID], pcidev)
		}
		return nil
	})
	if err != nil {
		log.Printf("Failed to discover host devices, error: %v", err)
	}

	return pciDevicesMap, err
}

// MockDiscoverPermittedHostPCIDevices returns predefined data for testing
func MockDiscoverPermittedHostPCIDevices(supportedPCIDeviceMap map[string]string) (map[string][]*PCIDevice, error) {
	log.Printf("enter  MockDiscoverPermittedHostPCIDevices")
	pciDevicesMap := make(map[string][]*PCIDevice)
	iommuToPCIMap := make(map[string]string)

	for vendorSelector, resourceName := range supportedPCIDeviceMap {
		pcidev1 := &PCIDevice{
			uuid:           uuid.New().String(),
			vendorSelector: vendorSelector,
			resourceName:   resourceName,
			pciAddress:     "0000:00:1d." + vendorSelector,
			driver:         "vfio-pci",
			iommuGroup:     "20" + vendorSelector,
			numaNode:       0,
			pciID:          vendorSelector,
		}

		pcidev2 := &PCIDevice{
			uuid:           uuid.New().String(),
			vendorSelector: vendorSelector,
			resourceName:   resourceName,
			pciAddress:     "0000:00:1e." + vendorSelector,
			driver:         "vfio-pci",
			iommuGroup:     "30" + vendorSelector,
			numaNode:       0,
			pciID:          vendorSelector,
		}
		iommuToPCIMap["20"+vendorSelector] = "0000:00:1d." + vendorSelector
		iommuToPCIMap["30"+vendorSelector] = "0000:00:1e." + vendorSelector

		pciDevicesMap[vendorSelector] = append(pciDevicesMap[vendorSelector], pcidev1, pcidev2)

	}

	return pciDevicesMap, nil
}
