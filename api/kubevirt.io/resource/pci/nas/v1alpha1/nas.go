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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AllocatablePci represents an allocatable Pci on a node.
type AllocatablePci struct {
	UUID         string `json:"uuid"`
	ResourceName string `json:"resourceName"`
	PciAddress   string `json:"pciAddress"`
}

// AllocatableDevice represents an allocatable device on a node.
type AllocatableDevice struct {
	Pci *AllocatablePci `json:"pci,omitempty"`
}

// Type returns the type of AllocatableDevice this represents.
func (d AllocatableDevice) Type() string {
	if d.Pci != nil {
		return PciDeviceType
	}
	return UnknownDeviceType
}

// AllocatedPci represents an allocated PCI.
type AllocatedPci struct {
	UUID string `json:"uuid,omitempty"`
}

// AllocatedPcis represents a set of allocated PCIs.
type AllocatedPcis struct {
	Devices []AllocatedPci `json:"devices"`
}

// AllocatedDevices represents a set of allocated devices.
type AllocatedDevices struct {
	Pci *AllocatedPcis `json:"pci,omitempty"`
}

// Type returns the type of AllocatedDevices this represents.
func (r AllocatedDevices) Type() string {
	if r.Pci != nil {
		return PciDeviceType
	}
	return UnknownDeviceType
}

// PreparedPci represents a prepared PCI on a node.
type PreparedPci struct {
	UUID string `json:"uuid"`
}

// PreparedPcis represents a set of prepared PCIs on a node.
type PreparedPcis struct {
	Devices []PreparedPci `json:"devices"`
}

// PreparedDevices represents a set of prepared devices on a node.
type PreparedDevices struct {
	Pci *PreparedPcis `json:"pci,omitempty"`
}

// Type returns the type of PreparedDevices this represents.
func (d PreparedDevices) Type() string {
	if d.Pci != nil {
		return PciDeviceType
	}
	return UnknownDeviceType
}

// NodeAllocationStateSpec is the spec for the NodeAllocationState CRD.
type NodeAllocationStateSpec struct {
	AllocatableDevices []AllocatableDevice         `json:"allocatableDevices,omitempty"`
	AllocatedClaims    map[string]AllocatedDevices `json:"allocatedClaims,omitempty"`
	PreparedClaims     map[string]PreparedDevices  `json:"preparedClaims,omitempty"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:resource:singular=nas

// NodeAllocationState holds the state required for allocation on a node.
type NodeAllocationState struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeAllocationStateSpec `json:"spec,omitempty"`
	Status string                  `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeAllocationStateList represents the "plural" of a NodeAllocationState CRD object.
type NodeAllocationStateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []NodeAllocationState `json:"items"`
}
