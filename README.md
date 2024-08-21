# KubeVirt DRA PCI Driver

**This repository holds currently a Proof of Concept and DRA isn't currently
supported by KubeVirt**

<p align="center">
<img src="https://github.com/kubevirt/community/raw/main/logo/KubeVirt_icon.png" width="100">
</p>

This repository contains a Dynamic Resource Allocation (DRA) device driver for PCI passthrough devices to be used with KubeVirt.

## About the Resource Driver

The KubeVirt DRA Driver is designed as an alternative to the device-plugin framework-based host device management in KubeVirt's `virt-handler`. It provides better control over devices on KubeVirt VMs by leveraging the [Dynamic Resource Allocation (DRA)](https://kubernetes.io/docs/concepts/scheduling-eviction/dynamic-resource-allocation/) framework in Kubernetes.

### Prerequisites

* [GNU Make 3.81+](https://www.gnu.org/software/make/)
* [GNU Tar 1.34+](https://www.gnu.org/software/tar/)
* [Docker v20.10+ (including buildx)](https://docs.docker.com/engine/install/)

### Documentation

- [Install and Use KubeVirt with the DRA PCI Driver](doc/KV_SETUP.md)
- [How to Change and Rebuild the Driver](doc/BUILD.md)

### Prior Art

- [kubernetes-sigs/dra-example-driver](https://github.com/kubernetes-sigs/dra-example-driver) provides a practical implementation of how a resource entity should be handled to ensure proper scheduling of Pods according to resources present in a Node.
- [Device Manager Implementation](https://github.com/kubevirt/kubevirt/tree/main/pkg/virt-handler/device-manager) from KubeVirt provides an implementation of handling several devices, including PCI devices for VMIs, using the device-plugin framework as part of `virt-handler`.
