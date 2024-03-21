.PHONY: prereqs
prereqs:
	brew install buildkit

.PHONY: buildkit
buildkit:
	docker run \
		--name buildkitd \
		-d \
		-p 3001:3001 \
		--security-opt seccomp=unconfined \
		--security-opt apparmor=unconfined \
		--device /dev/fuse \
		moby/buildkit:rootless --oci-worker-no-process-sandbox --addr tcp://0.0.0.0:3001

.PHONY: build
build:
	buildctl --addr tcp://localhost:3000 build \
		--frontend dockerfile.v0 \
		--local context=. \
		--local dockerfile=. \
		--opt build-arg:BUILDKIT_INLINE_CACHE=1 \
		--output type=image,name=test-image:latest,push=false

.PHONY: start
start:
	go run ./cmd/buildkit-proxy