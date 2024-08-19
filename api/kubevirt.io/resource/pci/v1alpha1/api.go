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

package v1alpha1

import nascrd "kubevirt.io/dra-pci-driver/api/kubevirt.io/resource/pci/nas/v1alpha1"

const (
	GroupName = "pci.resource.kubevirt.io"
	Version   = "v1alpha1"

	PciClaimParametersKind = "PciClaimParameters"
)

func DefaultDeviceClassParametersSpec() *DeviceClassParametersSpec {
	return &DeviceClassParametersSpec{
		DeviceSelector: []DeviceSelector{
			{
				Type:              nascrd.PciDeviceType,
				ResourceName:      "*",
				PCIVendorSelector: "*",
			},
		},
	}
}

func DefaultPciClaimParametersSpec() *PciClaimParametersSpec {
	return &PciClaimParametersSpec{
		DeviceName: "*",
	}
}
