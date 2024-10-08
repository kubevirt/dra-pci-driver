//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AllocatableDevice) DeepCopyInto(out *AllocatableDevice) {
	*out = *in
	if in.Pci != nil {
		in, out := &in.Pci, &out.Pci
		*out = new(AllocatablePci)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AllocatableDevice.
func (in *AllocatableDevice) DeepCopy() *AllocatableDevice {
	if in == nil {
		return nil
	}
	out := new(AllocatableDevice)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AllocatablePci) DeepCopyInto(out *AllocatablePci) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AllocatablePci.
func (in *AllocatablePci) DeepCopy() *AllocatablePci {
	if in == nil {
		return nil
	}
	out := new(AllocatablePci)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AllocatedDevices) DeepCopyInto(out *AllocatedDevices) {
	*out = *in
	if in.Pci != nil {
		in, out := &in.Pci, &out.Pci
		*out = new(AllocatedPcis)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AllocatedDevices.
func (in *AllocatedDevices) DeepCopy() *AllocatedDevices {
	if in == nil {
		return nil
	}
	out := new(AllocatedDevices)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AllocatedPci) DeepCopyInto(out *AllocatedPci) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AllocatedPci.
func (in *AllocatedPci) DeepCopy() *AllocatedPci {
	if in == nil {
		return nil
	}
	out := new(AllocatedPci)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AllocatedPcis) DeepCopyInto(out *AllocatedPcis) {
	*out = *in
	if in.Devices != nil {
		in, out := &in.Devices, &out.Devices
		*out = make([]AllocatedPci, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AllocatedPcis.
func (in *AllocatedPcis) DeepCopy() *AllocatedPcis {
	if in == nil {
		return nil
	}
	out := new(AllocatedPcis)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeAllocationState) DeepCopyInto(out *NodeAllocationState) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeAllocationState.
func (in *NodeAllocationState) DeepCopy() *NodeAllocationState {
	if in == nil {
		return nil
	}
	out := new(NodeAllocationState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeAllocationState) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeAllocationStateConfig) DeepCopyInto(out *NodeAllocationStateConfig) {
	*out = *in
	if in.Owner != nil {
		in, out := &in.Owner, &out.Owner
		*out = new(v1.OwnerReference)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeAllocationStateConfig.
func (in *NodeAllocationStateConfig) DeepCopy() *NodeAllocationStateConfig {
	if in == nil {
		return nil
	}
	out := new(NodeAllocationStateConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeAllocationStateList) DeepCopyInto(out *NodeAllocationStateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NodeAllocationState, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeAllocationStateList.
func (in *NodeAllocationStateList) DeepCopy() *NodeAllocationStateList {
	if in == nil {
		return nil
	}
	out := new(NodeAllocationStateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeAllocationStateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeAllocationStateSpec) DeepCopyInto(out *NodeAllocationStateSpec) {
	*out = *in
	if in.AllocatableDevices != nil {
		in, out := &in.AllocatableDevices, &out.AllocatableDevices
		*out = make([]AllocatableDevice, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AllocatedClaims != nil {
		in, out := &in.AllocatedClaims, &out.AllocatedClaims
		*out = make(map[string]AllocatedDevices, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.PreparedClaims != nil {
		in, out := &in.PreparedClaims, &out.PreparedClaims
		*out = make(map[string]PreparedDevices, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeAllocationStateSpec.
func (in *NodeAllocationStateSpec) DeepCopy() *NodeAllocationStateSpec {
	if in == nil {
		return nil
	}
	out := new(NodeAllocationStateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreparedDevices) DeepCopyInto(out *PreparedDevices) {
	*out = *in
	if in.Pci != nil {
		in, out := &in.Pci, &out.Pci
		*out = new(PreparedPcis)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreparedDevices.
func (in *PreparedDevices) DeepCopy() *PreparedDevices {
	if in == nil {
		return nil
	}
	out := new(PreparedDevices)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreparedPci) DeepCopyInto(out *PreparedPci) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreparedPci.
func (in *PreparedPci) DeepCopy() *PreparedPci {
	if in == nil {
		return nil
	}
	out := new(PreparedPci)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreparedPcis) DeepCopyInto(out *PreparedPcis) {
	*out = *in
	if in.Devices != nil {
		in, out := &in.Devices, &out.Devices
		*out = make([]PreparedPci, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreparedPcis.
func (in *PreparedPcis) DeepCopy() *PreparedPcis {
	if in == nil {
		return nil
	}
	out := new(PreparedPcis)
	in.DeepCopyInto(out)
	return out
}
