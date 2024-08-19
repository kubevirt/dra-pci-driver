#!/usr/bin/env bash

# Copyright 2024 The KubeVirt Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Apply CRDs
kubectl apply -f ../deployments/native/dra-pci-driver/crds/nas.pci.resource.kubevirt.io_nodeallocationstates.yaml
kubectl apply -f ../deployments/native/dra-pci-driver/crds/pci.resource.kubevirt.io_pciclaimparameters.yaml
kubectl apply -f ../deployments/native/dra-pci-driver/crds/pci.resource.kubevirt.io_deviceclassparameters.yaml


# Apply other Kubernetes objects
kubectl apply -f ../deployments/native/dra-pci-driver/templates/namespace.yaml
kubectl apply -f ../deployments/native/dra-pci-driver/templates/serviceaccount.yaml
kubectl apply -f ../deployments/native/dra-pci-driver/templates/clusterrole.yaml
kubectl apply -f ../deployments/native/dra-pci-driver/templates/clusterrolebinding.yaml
kubectl apply -f ../deployments/native/dra-pci-driver/templates/resourceclass.yaml
kubectl apply -f ../deployments/native/dra-pci-driver/templates/controller.yaml
kubectl apply -f ../deployments/native/dra-pci-driver/templates/kubeletplugin.yaml