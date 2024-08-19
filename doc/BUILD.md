# How to build Kubevirt DRA Driver container image

## Platforms supported

- Linux

## Prerequisites

- Docker

## Building

- Script to simply rebuild the image is already present in demo folder
 ```bash
cd demo
./build-driver.sh
```
# How to Build KubeVirt DRA Driver Container Image

## Platforms Supported

- Linux

## Prerequisites

- Docker

## Building

1. **Rebuild the image using the provided script:**

	The script to rebuild the image is located in the `demo` folder.

	```bash
	cd demo
	./build-driver.sh
	```

2. **Rebuild the CRD specifically:**

	If you want to rebuild the CRD to change the logic in the driver, you can use the `Makefile`. The `Makefile` automates this process. The only required tool is Docker, where the generation process takes place in a container and is copied to the configured path.

	```bash
	make docker-generate
	```

3. **Rebuild it entirely:**

	To rebuild the driver in a docker container, use the following command:

	```bash
	make docker-build
	```

