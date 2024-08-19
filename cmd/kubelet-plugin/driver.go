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
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	drapbv1 "k8s.io/kubelet/pkg/apis/dra/v1alpha3"
	nascrd "kubevirt.io/dra-pci-driver/api/kubevirt.io/resource/pci/nas/v1alpha1"
	nasclient "kubevirt.io/dra-pci-driver/api/kubevirt.io/resource/pci/nas/v1alpha1/client"
	clientset "kubevirt.io/dra-pci-driver/pkg/kubevirt.io/resource/clientset/versioned"
)

const (
	//TODO , now we need to keep the name of the class parameter as this
	pciClassParametersRefName = "pci-params"
)

var _ drapbv1.NodeServer = &driver{}

type driver struct {
	nascrd    *nascrd.NodeAllocationState
	nasclient *nasclient.Client
	state     *DeviceState
	clientset clientset.Interface
}

func NewDriver(ctx context.Context, config *Config) (*driver, error) {
	var d *driver
	client := nasclient.New(config.nascr, config.exampleclient.NasV1alpha1())

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err := client.GetOrCreate(ctx)
		if err != nil {
			return err
		}

		err = client.UpdateStatus(ctx, nascrd.NodeAllocationStateStatusNotReady)
		if err != nil {
			return err
		}

		supportedAllocatedDevices, err := getAllocatableDevicesFromResourceClass(ctx, config.clientSets.Example)
		if err != nil {
			return fmt.Errorf("error getting allocatable devices from resource class: %v", err)
		}

		possibleDevices, err := enumerateAllPossibleDevices(supportedAllocatedDevices)
		if err != nil {
			return fmt.Errorf("error enumerating all possible devices: %v", err)
		}

		state, err := NewDeviceState(config, possibleDevices)
		if err != nil {
			return err
		}

		updatedSpec, err := state.GetUpdatedSpec(&config.nascr.Spec)
		if err != nil {
			return fmt.Errorf("error getting updated CR spec: %v", err)
		}

		err = client.Update(ctx, updatedSpec)
		if err != nil {
			return err
		}

		err = client.UpdateStatus(ctx, nascrd.NodeAllocationStateStatusReady)
		if err != nil {
			return err
		}

		d = &driver{
			nascrd:    config.nascr,
			nasclient: client,
			state:     state,
			clientset: config.clientSets.Example,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *driver) Shutdown(ctx context.Context) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err := d.nasclient.Get(ctx)
		if err != nil {
			return err
		}
		return d.nasclient.UpdateStatus(ctx, nascrd.NodeAllocationStateStatusNotReady)
	})
}

func (d *driver) NodePrepareResources(ctx context.Context, req *drapbv1.NodePrepareResourcesRequest) (*drapbv1.NodePrepareResourcesResponse, error) {
	logger := klog.FromContext(ctx)
	logger.Info("NodePrepareResource", "numClaims", len(req.Claims))
	preparedResources := &drapbv1.NodePrepareResourcesResponse{Claims: map[string]*drapbv1.NodePrepareResourceResponse{}}

	// In production version some common operations of d.nodeUnprepareResources
	// should be done outside of the loop, for instance updating the CR could
	// be done once after all HW was prepared.
	for _, claim := range req.Claims {
		preparedResources.Claims[claim.Uid] = d.nodePrepareResource(ctx, claim)
	}

	return preparedResources, nil
}

func (d *driver) nodePrepareResource(ctx context.Context, claim *drapbv1.Claim) *drapbv1.NodePrepareResourceResponse {
	logger := klog.FromContext(ctx)
	var err error
	var prepared []string
	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		prepared, err = d.prepare(ctx, claim.Uid)
		if err != nil {
			return fmt.Errorf("error allocating devices for claim '%v': %v", claim.Uid, err)
		}

		updatedSpec, err := d.state.GetUpdatedSpec(&d.nascrd.Spec)
		if err != nil {
			return fmt.Errorf("error getting updated CR spec: %v", err)
		}

		err = d.nasclient.Update(ctx, updatedSpec)
		if err != nil {
			if err := d.state.Unprepare(claim.Uid); err != nil {
				logger.Error(err, "Failed to unprepare after Update", "claim", claim.Uid)
			}
			return err
		}

		return nil
	})

	if err != nil {
		return &drapbv1.NodePrepareResourceResponse{
			Error: fmt.Sprintf("error preparing resource: %v", err),
		}
	}

	klog.FromContext(ctx).Info("Prepared devices", "claim", claim.Uid)
	return &drapbv1.NodePrepareResourceResponse{CDIDevices: prepared}
}

func (d *driver) NodeUnprepareResources(ctx context.Context, req *drapbv1.NodeUnprepareResourcesRequest) (*drapbv1.NodeUnprepareResourcesResponse, error) {
	logger := klog.FromContext(ctx)
	logger.Info("NodeUnprepareResource", "numClaims", len(req.Claims))
	unpreparedResources := &drapbv1.NodeUnprepareResourcesResponse{
		Claims: map[string]*drapbv1.NodeUnprepareResourceResponse{},
	}

	// In production version some common operations of d.nodeUnprepareResources
	// should be done outside of the loop, for instance updating the CR could
	// be done once after all HW was unprepared.
	for _, claim := range req.Claims {
		unpreparedResources.Claims[claim.Uid] = d.nodeUnprepareResource(ctx, claim)
	}

	return unpreparedResources, nil
}

func (d *driver) nodeUnprepareResource(ctx context.Context, claim *drapbv1.Claim) *drapbv1.NodeUnprepareResourceResponse {
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err := d.unprepare(ctx, claim.Uid)
		if err != nil {
			return fmt.Errorf("error unpreparing devices for claim '%v': %v", claim.Uid, err)
		}

		updatedSpec, err := d.state.GetUpdatedSpec(&d.nascrd.Spec)
		if err != nil {
			return fmt.Errorf("error getting updated CR spec: %v", err)
		}

		err = d.nasclient.Update(ctx, updatedSpec)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return &drapbv1.NodeUnprepareResourceResponse{
			Error: fmt.Sprintf("error unpreparing resource: %v", err),
		}
	}

	klog.FromContext(ctx).Info("Unprepared devices", "claim", claim.Uid)
	return &drapbv1.NodeUnprepareResourceResponse{}
}

func (d *driver) prepare(ctx context.Context, claimUID string) ([]string, error) {
	err := d.nasclient.Get(ctx)
	if err != nil {
		return nil, err
	}
	prepared, err := d.state.Prepare(claimUID, d.nascrd.Spec.AllocatedClaims[claimUID])
	if err != nil {
		return nil, err
	}
	return prepared, nil
}

func (d *driver) unprepare(ctx context.Context, claimUID string) error {
	err := d.nasclient.Get(ctx)
	if err != nil {
		return err
	}
	err = d.state.Unprepare(claimUID)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Accessing a cluster level object like ClassParameters in node-local plugin should be avoided
// for KubeVirt usecase ,a possible solution is using a grpc server for passing  ClassParameters from controller to plugin
func getAllocatableDevicesFromResourceClass(ctx context.Context, clientset clientset.Interface) (AllocatableDevices, error) {
	dc, err := clientset.PciV1alpha1().DeviceClassParameters().Get(ctx, pciClassParametersRefName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting DeviecClassParameters called '%v': %v", pciClassParametersRefName, err)
	}

	devices := make(AllocatableDevices)
	for _, device := range dc.Spec.DeviceSelector {
		devices[device.PCIVendorSelector] = &AllocatableDeviceInfo{
			PCIDevice: &PCIDevice{
				vendorSelector: device.PCIVendorSelector,
				resourceName:   device.ResourceName,
			},
		}
	}
	return devices, nil
}
