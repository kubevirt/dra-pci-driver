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
	"fmt"

	pcicrd "kubevirt.io/dra-pci-driver/api/kubevirt.io/resource/pci/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	resourcev1 "k8s.io/api/resource/v1alpha2"
	"k8s.io/dynamic-resource-allocation/controller"

	nascrd "kubevirt.io/dra-pci-driver/api/kubevirt.io/resource/pci/nas/v1alpha1"
)

type pcidriver struct {
	PendingAllocatedClaims *PerNodeAllocatedClaims
}

func NewPciDriver() *pcidriver {
	return &pcidriver{
		PendingAllocatedClaims: NewPerNodeAllocatedClaims(),
	}
}

func (p *pcidriver) ValidateClaimParameters(claimParams *pcicrd.PciClaimParametersSpec) error {
	if claimParams.DeviceName != "devices.kubevirt.io/nvme" {
		return fmt.Errorf("unsupported pci device type: %s", claimParams.DeviceName)
	}
	return nil
}

func (p *pcidriver) Allocate(crd *nascrd.NodeAllocationState, claim *resourcev1.ResourceClaim, claimParams *pcicrd.PciClaimParametersSpec, class *resourcev1.ResourceClass, classParams *pcicrd.DeviceClassParametersSpec, selectedNode string) (OnSuccessCallback, error) {
	claimUID := string(claim.UID)

	if !p.PendingAllocatedClaims.Exists(claimUID, selectedNode) {
		return nil, fmt.Errorf("no allocations generated for claim '%v' on node '%v' yet", claim.UID, selectedNode)
	}

	crd.Spec.AllocatedClaims[claimUID] = p.PendingAllocatedClaims.Get(claimUID, selectedNode)
	onSuccess := func() {
		p.PendingAllocatedClaims.Remove(claimUID)
	}

	return onSuccess, nil
}

func (p *pcidriver) Deallocate(crd *nascrd.NodeAllocationState, claim *resourcev1.ResourceClaim) error {
	claimUID := string(claim.UID)
	p.PendingAllocatedClaims.Remove(claimUID)
	return nil
}

func (p *pcidriver) UnsuitableNode(crd *nascrd.NodeAllocationState, pod *corev1.Pod, pcicas []*controller.ClaimAllocation, allcas []*controller.ClaimAllocation, potentialNode string) error {

	// Visit the node and update the allocated claims
	p.PendingAllocatedClaims.VisitNode(potentialNode, func(claimUID string, allocation nascrd.AllocatedDevices) {
		if _, exists := crd.Spec.AllocatedClaims[claimUID]; exists {
			p.PendingAllocatedClaims.Remove(claimUID)
		} else {
			crd.Spec.AllocatedClaims[claimUID] = allocation
		}
	})

	// Allocate resources
	allocated := p.allocate(crd, pod, pcicas, allcas, potentialNode)

	// Iterate over the PCI claim allocations
	for _, ca := range pcicas {
		claimUID := string(ca.Claim.UID)
		_, ok := ca.ClaimParameters.(*pcicrd.PciClaimParametersSpec)
		if !ok {
			return fmt.Errorf("invalid claim parameters for claim UID: %s", claimUID)
		}

		// Check if there is exactly one allocated device
		if len(allocated[claimUID]) != 1 {
			for _, ca := range allcas {
				ca.UnsuitableNodes = append(ca.UnsuitableNodes, potentialNode)
			}
			return nil
		}

		// Prepare the allocated device
		device := nascrd.AllocatedPci{
			UUID: allocated[claimUID][0],
		}

		// Create the allocated devices structure
		allocatedDevices := nascrd.AllocatedDevices{
			Pci: &nascrd.AllocatedPcis{
				Devices: []nascrd.AllocatedPci{device},
			},
		}

		// Set the pending allocated claims
		p.PendingAllocatedClaims.Set(claimUID, potentialNode, allocatedDevices)
	}

	return nil
}

func (p *pcidriver) allocate(crd *nascrd.NodeAllocationState, pod *corev1.Pod, pcicas []*controller.ClaimAllocation, allcas []*controller.ClaimAllocation, node string) map[string][]string {

	available := make(map[string]*nascrd.AllocatablePci)

	for _, device := range crd.Spec.AllocatableDevices {
		if device.Type() != nascrd.PciDeviceType {
			continue
		}
		if _, exist := crd.Spec.AllocatedClaims[device.Pci.UUID]; exist {
			continue
		}
		available[device.Pci.UUID] = device.Pci
	}

	allocated := make(map[string][]string)

	for _, ca := range pcicas {
		claimUID := string(ca.Claim.UID)

		if v, exists := crd.Spec.AllocatedClaims[claimUID]; exists {
			devices := v.Pci.Devices
			for _, device := range devices {
				allocated[claimUID] = append(allocated[claimUID], device.UUID)
			}
			continue
		}

		claimParams, _ := ca.ClaimParameters.(*pcicrd.PciClaimParametersSpec)

		for uuid, device := range available {
			// Check if the device type is the one requested
			if device.ResourceName == claimParams.DeviceName {
				allocated[claimUID] = []string{device.UUID}
				delete(available, uuid)
				break
			}
		}
	}

	return allocated
}
