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
	"strings"
)

func enumerateAllPossibleDevices(supportedAllocatedDevices AllocatableDevices) (AllocatableDevices, error) {
	supportedPCIDeviceMap := make(map[string]string)
	for _, device := range supportedAllocatedDevices {
		supportedPCIDeviceMap[strings.ToLower(device.PCIDevice.vendorSelector)] = device.PCIDevice.resourceName
	}

	allDevices := make(AllocatableDevices)
	discoveredDevices, err := DiscoverPermittedHostPCIDevices(supportedPCIDeviceMap)
	if err != nil {
		return allDevices, err
	}

	for _, device := range supportedAllocatedDevices {
		vendorSelector := device.PCIDevice.vendorSelector
		if devices, supported := discoveredDevices[vendorSelector]; supported {
			for _, discoveredDevice := range devices {
				newDevice := &AllocatableDeviceInfo{
					PCIDevice: discoveredDevice,
				}
				allDevices[newDevice.PCIDevice.uuid] = newDevice
			}
		}
	}

	return allDevices, nil
}
