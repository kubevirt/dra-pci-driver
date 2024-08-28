# KubeVirt Setup

This guide will walk you through the steps to install and use KubeVirt with the `dra-pci-driver`.

## Step 1: Install KubeVirt

1. **Download the KubeVirt repository:**

   ```bash
   git clone https://github.com/TheRealSibasishBehera/kubevirt/tree/host-rc
   cd kubevirt
   ```

2. **Start a Kubernetes cluster with 1 NICs:**

   ```bash
   export KUBEVIRT_PROVIDER_EXTRA_ARGS="--nvme 1G --nvme 500M" # Needed to emulate NVMe devices
   export KUBEVIRT_PROVIDER=k8s-1.30
   export FEATURE_GATES=HostDevices,DynamicResourceAllocation
   export KUBEVIRT_NUM_NODES=1
   export KUBEVIRT_NUM_SECONDARY_NICS=1
   make cluster-up
   make cluster-sync
   ```

3. **Verify that the nodes are up and running:**

   ```bash
   cluster-up/kubectl.sh get nodes
   ```

   It should show something like this:

   ```bash
   NAME     STATUS   ROLES                  AGE    VERSION
   node01   Ready    control-plane,worker   122m   v1.30.3
   ```

4. **SSH and verify that the nodes have the devices present:**

   ```bash
   cluster-up/ssh.sh node01
   ```

   ```bash
   $ lspci -Dnn | grep NVM
   ```

   It should show something like this:

   ```
   0000:00:07.0 Non-Volatile memory controller [0108]: Red Hat, Inc. QEMU NVM Express Controller [1b36:0010] (rev 02)
   0000:00:08.0 Non-Volatile memory controller [0108]: Red Hat, Inc. QEMU NVM Express Controller [1b36:0010] (rev 02)
   ```

## Step 2: Enable Dynamic Resource Allocation

1. **SSH into the control plane node:**

   ```bash
   cluster-up/ssh.sh node01
   ```

2. **Edit the manifest files to add the necessary parameters:**

   - **kube-apiserver:** Edit `/etc/kubernetes/manifests/kube-apiserver.yaml` and add the following to the `command` section under the `container` section:

     ```yaml
     - --feature-gates=DynamicResourceAllocation=true
     - --runtime-config=resource.k8s.io/v1alpha2=true
     ```

   - **kube-controller-manager:** Edit `/etc/kubernetes/manifests/kube-controller-manager.yaml` and add the following to the `command` section under the `container` section:

     ```yaml
     - --feature-gates=DynamicResourceAllocation=true
     ```

   - **kube-scheduler:** Edit `/etc/kubernetes/manifests/kube-scheduler.yaml` and add the following to the `command` section under the `container` section:

     ```yaml
     - --feature-gates=DynamicResourceAllocation=true
     ```

   - **kubelet config:** Edit `/var/lib/kubelet/config.yaml` and add the following:

     ```yaml
     featureGates:
       DynamicResourceAllocation: true
     ```

## Step 3: Bind Devices to `vfio-pci` Driver

1. **SSH into the control plane node:**

   ```bash
   cluster-up/ssh.sh node01
   ```

2. **Change GRUB file**

   ```bash
   sudo vi /etc/default/grub
   ```

   Add `intel_iommu=on` to `GRUB_CMDLINE_LINUX`

   ```bash
   GRUB_CMDLINE_LINUX="nofb splash=quiet console=tty0 ... intel_iommu=on
   ```

   Make the grub file

   ```bash
   sudo grub2-mkconfig -o /boot/grub2/grub.cfg
   ```

   Reboot the system

   ```bash
   sudo reboot
   ```

   Enable `vfio-pci`

   ```bash
   sudo modprobe vfio-pci
   ```

3. **Create a script to unbind the device from its default driver and bind it to the `vfio-pci` driver:**

   ```bash
   vi pci-nvme-bind.sh
   ```

   Paste the following script and save:

   ```bash
   #!/bin/bash
   # Check if PCI address is provided
   if [ -z "$1" ]; then
     echo "Usage: $0 <PCI_ADDRESS>"
     exit 1
   fi
   PCI_ADDRESS=$1

   # Unbind the device from its current driver
   echo "Unbinding the device $PCI_ADDRESS from its current driver..."
   echo $PCI_ADDRESS > /sys/bus/pci/drivers/nvme/unbind

   # Set the driver override to vfio-pci
   echo "Setting the driver override to vfio-pci for device $PCI_ADDRESS..."
   echo "vfio-pci" > /sys/bus/pci/devices/$PCI_ADDRESS/driver_override

   # Bind the device to the vfio-pci driver
   echo "Binding the device $PCI_ADDRESS to the vfio-pci driver..."
   echo $PCI_ADDRESS > /sys/bus/pci/drivers/vfio-pci/bind
   echo "Device $PCI_ADDRESS has been successfully bound to vfio-pci."
   ```

4. **Make the script executable and use it by passing the PCI address of the devices as an argument. The address can be found using `lspci -Dnn | grep NVM`:**

   ```bash
   chmod +x pci-nvme-bind.sh
   sudo ./pci-nvme-bind.sh 0000:00:07.0
   sudo ./pci-nvme-bind.sh 0000:00:08.0
   ```

5. **Verify if the devices are bound to the `vfio-pci` driver:**

   ```bash
   lspci -Dnnk -s 0000:00:07.0
   lspci -Dnnk -s 0000:00:08.0
   ```

   It should show `Kernel driver in use: vfio-pci` for the devices:

   ```bash
   0000:00:07.0 Non-Volatile memory controller [0108]: Red Hat, Inc. QEMU NVM Express Controller [1b36:0010] (rev 02)
       Subsystem: Red Hat, Inc. Device [1af4:1100]
       Kernel driver in use: vfio-pci
       Kernel modules: nvme
   0000:00:08.0 Non-Volatile memory controller [0108]: Red Hat, Inc. QEMU NVM Express Controller [1b36:0010] (rev 02)
       Subsystem: Red Hat, Inc. Device [1af4:1100]
       Kernel driver in use: vfio-pci
       Kernel modules: nvme
   ```

6. **Disable SELinux inside the node:**

   ```bash
   sudo setenforce 0
   ```

## Step 4: Verification for Resource API

1. **Check if your Kubernetes cluster supports dynamic resource allocation:**

   ```bash
   cluster-up/kubectl.sh get resourceclasses
   ```

   - If your cluster supports dynamic resource allocation, the response will either be a list of `ResourceClass` objects or:

     ```bash
     No resources found
     ```

   - If dynamic resource allocation is not supported, you will see the following error:

     ```bash
     error: the server does not have a resource type "resourceclasses"
     ```

## Step 5: Deploying the KubeVirt DRA PCI Driver

1. **Download the KubeVirt DRA PCI Driver repository:**

   ```bash
   git clone https://github.com/kubevirt/dra-pci-driver.git
   cd dra-pci-driver
   ```

2. **Build the driver:**

   ```bash
   cd demo
   ./build-driver.sh
   ```

   The driver will be saved as an image named `registry.example.com/dra-pci-driver:v0.1.0`.

3. **Push the driver image into the cluster where KubeVirt is running:**

   ```bash
   export K8S_VERSION=k8s-1.30 # Change this if you use a different version
   export CONT=${K8S_VERSION}-dnsmasq
   chmod +x image-push-docker.sh
   ./image-push-docker.sh registry.example.com/dra-pci-driver:v0.1.0
   ```

4. **Apply the DRA PCI Driver manifests:**

   ```bash
   export KUBECONFIG=$(path/to/kubevirt/cluster-up/kubeconfig.sh)
   ./deploy-native.sh
   ```

5. **Verify Node State:**

   ```bash
   kubectl describe nas node01 -n dra-pci-driver
   ```

   The Spec should contain only Allocatable Devices:

   ```yaml
   Spec:
     Allocatable Devices:
       Pci:
         Pci Address:    0000:00:07.0
         Resource Name:  devices.kubevirt.io/nvme
         Uuid:           bc628854-6471-463a-878d-b96b8c7022dd
       Pci:
         Pci Address:    0000:00:08.0
         Resource Name:  devices.kubevirt.io/nvme
         Uuid:           c98572f0-37a0-41bf-b4e0-70d8a12278a3
   ```

6. **Deploy an example VMI:**

   ```bash
   kubectl apply -f vmi-resource-claim.yaml
   ```

7. **Verify if the PCI device is allocated:**

   ```bash
   kubectl describe nas node01 -n dra-pci-driver
   ```

   The Spec should contain Allocated Devices:

   ```yaml
   Spec:
     Allocatable Devices:
       Pci:
         Pci Address:    0000:00:07.0
         Resource Name:  devices.kubevirt.io/nvme
         Uuid:           bc628854-6471-463a-878d-b96b8c7022dd
       Pci:
         Pci Address:    0000:00:08.0
         Resource Name:  devices.kubevirt.io/nvme
         Uuid:           c98572f0-37a0-41bf-b4e0-70d8a12278a3
     Allocated Claims:
       d313f2ab-a84f-449e-bf98-1a379f256ec3:
         Pci:
           Devices:
             Uuid:  bc628854-6471-463a-878d-b96b8c7022dd
     Prepared Claims:
       d313f2ab-a84f-449e-bf98-1a379f256ec3:
         Pci:
           Devices:
             Uuid:  bc628854-6471-463a-878d-b96b8c7022dd
   ```

8. **Verify if the `virt-launcher` pod is running:**

   ```bash
   kubectl get pods -A
   ```

   This should show a `virt-launcher` pod with the name `virt-launcher-vmi-nvme-xxx` in the `Running` state.
