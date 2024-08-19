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

ARG GOLANG_VERSION=undefined
ARG BASE_IMAGE=undefined
FROM golang:${GOLANG_VERSION} as build

WORKDIR /build
COPY . .

RUN mkdir /artifacts
RUN make PREFIX=/artifacts cmds

FROM ${BASE_IMAGE}

LABEL io.k8s.display-name="KUBEVIRT DRA Driver"
LABEL name="KUBEVIRT DRA Driver"
LABEL vendor="kubevirt.io"
ARG VERSION="N/A"
LABEL version=${VERSION}
ARG RELEASE="N/A"
LABEL release="N/A"
LABEL summary="KUBEVIRT DRA Driver"
LABEL description="See summary"

COPY --from=build /artifacts/virt-dra-controller    /usr/bin/virt-dra-controller
COPY --from=build /artifacts/kubelet-plugin /usr/bin/kubelet-plugin
COPY --from=build /artifacts/set-nas-status            /usr/bin/set-nas-status
