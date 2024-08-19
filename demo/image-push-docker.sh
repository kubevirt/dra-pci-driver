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

# Check if an image name is provided as an argument
if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <image-name>"
  exit 1
fi

# Assign the provided image name to a variable
IMAGE=$1

K8S_VERSION=k8s-1.30
CONT=${K8S_VERSION}-dnsmasq

# Get the port number for the registry running in the specified container
PORT=$(sudo docker port ${CONT} 5000 | awk -F ":" '{print $2}')
IMAGE_PUSH=localhost:${PORT}/${IMAGE}

# Remove the image from the local registry if it exists, suppress errors
sudo docker rmi ${IMAGE_PUSH} || true

# Tag the local image for the local registry
sudo docker tag ${IMAGE} ${IMAGE_PUSH}

# Push the tagged image to the local registry
sudo docker push ${IMAGE_PUSH}