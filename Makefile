# make all builds both cloud and edge binaries

BINARIES=cloudcore \
	admission \
	edgecore \
	edgesite \
	keadm

COMPONENTS=cloud\
	edge

.EXPORT_ALL_VARIABLES:
OUT_DIR ?= _output

define ALL_HELP_INFO
# Build code.
#
# Args:
#   WHAT: binary names to build. support: $(BINARIES) 
#         the build will produce executable files under $(OUT_DIR)
#         If not specified, "everything" will be built.
#
# Example:
#   make
#   make all
#   make all HELP=y
#   make all WHAT=cloudcore
endef

.PHONY: all
ifeq ($(HELP),y)
all:
	@echo "$$ALL_HELP_INFO"
else
all: verify-golang
	hack/make-rules/build.sh $(WHAT)
endif


define VERIFY_HELP_INFO
# verify golang,vendor and codegen
#
# Example:
# make verify 
endef
.PHONY: verify
ifeq ($(HELP),y)
verify:
	@echo "$$VERIFY_HELP_INFO"
else
verify:verify-golang verify-vendor verify-codegen 
endif

.PHONY: verify-golang
verify-golang: 
	bash hack/verify-golang.sh

.PHONY: verify-vendor
verify-vendor: 
	bash hack/verify-vendor.sh
.PHONY: verify-codegen
verify-codegen: 
	bash cloud/hack/verify-codegen.sh



define TEST_HELP_INFO
# run golang test case.
#
# Args:
#   WHAT: Component names to be testd. support: $(COMPONENTS) 
#         If not specified, "everything" will be tested.
#
# Example:
#   make test 
#   make test HELP=y
#   make test WHAT=cloud
endef

.PHONY: test 
ifeq ($(HELP),y)
test:
	@echo "$$TEST_HELP_INFO"
else
test: 
	hack/make-rules/test.sh $(WHAT)
endif

LINTS=cloud \
	edge \
	keadm \
	bluetoothdevice


define LINT_HELP_INFO
# run golang lint check.
#
# Args:
#   WHAT: Component names to be lint check. support: $(LINTS) 
#         If not specified, "everything" will be lint check.
#
# Example:
#   make lint 
#   make lint HELP=y
#   make lint WHAT=cloud
endef

.PHONY: lint 
ifeq ($(HELP),y)
lint:
	@echo "$$LINT_HELP_INFO"
else
lint: 
	hack/make-rules/lint.sh $(WHAT)
endif


INTEGRATION_TEST_COMPONENTS=edge
define INTEGRATION_TEST_HELP_INFO
# run integration test.
#
# Args:
#   WHAT: Component names to be lint check. support: $(INTEGRATION_TEST_COMPONENTS) 
#         If not specified, "everything" will be lint check.
#
# Example:
#   make integrationtest 
#   make integrationtest HELP=y
endef

.PHONY: integrationtest 
ifeq ($(HELP),y)
integrationtest:
	@echo "$$INTEGRATION_TEST_HELP_INFO"
else
integrationtest: 
	hack/make-rules/build.sh edgecore 
	edge/test/integration/scripts/execute.sh
endif

CROSSBUILD_COMPONENTS=edgecore\
						edgesite
GOARM_VALUES=GOARM7 \
	GOARM8

define CROSSBUILD_HELP_INFO
# cross build components.
#
# Args:
#   WHAT: Component names to be lint check. support: $(CROSSBUILD_COMPONENTS) 
#         If not specified, "everything" will be lint check.
#
#	GOARM: go arm value, now support:$(GOARM_VALUES)
#			If not specified ,default use GOARM=GOARM8 
#
#
# Example:
#   make crossbuild 
#   make crossbuild HELP=y
#	make crossbuild WHAT=edgecore
#	make crossbuild WHAT=edgecore GOARM=GOARM7
#
endef

.PHONY: crossbuild 
ifeq ($(HELP),y)
crossbuild:
	@echo "$$CROSSBUILD_HELP_INFO"
else
crossbuild: 
	hack/make-rules/crossbuild.sh $(WHAT) $(GOARM)
endif


####################################


.PHONY: edge_small_build
edge_small_build:
	cd edge && $(MAKE) small_build

.PHONY: edgesite_small_build
edgesite_small_build:
	$(MAKE) -C edgesite small_build

.PHONY: e2e_test
e2e_test:
#	bash tests/e2e/scripts/execute.sh device_crd
#	This has been commented temporarily since there is an issue of CI using same master for all PRs, which is causing failures when run parallely
	bash tests/e2e/scripts/execute.sh

.PHONY: performance_test
performance_test:
	bash tests/performance/scripts/jenkins.sh

QEMU_ARCH ?= x86_64
ARCH ?= amd64
IMAGE_TAG ?= $(shell git describe --tags)

.PHONY: cloudimage
cloudimage:
	docker build -t kubeedge/cloudcore:${IMAGE_TAG} -f build/cloud/Dockerfile .

.PHONY: admissionimage
admissionimage:
	docker build -t kubeedge/admission:${IMAGE_TAG} -f build/admission/Dockerfile .

.PHONY: csidriverimage
csidriverimage:
	docker build -t kubeedge/csidriver:${IMAGE_TAG} -f build/csidriver/Dockerfile .

.PHONY: edgeimage
edgeimage:
	mkdir -p ./build/edge/tmp
	rm -rf ./build/edge/tmp/*
	curl -L -o ./build/edge/tmp/qemu-${QEMU_ARCH}-static.tar.gz https://github.com/multiarch/qemu-user-static/releases/download/v3.0.0/qemu-${QEMU_ARCH}-static.tar.gz
	tar -xzf ./build/edge/tmp/qemu-${QEMU_ARCH}-static.tar.gz -C ./build/edge/tmp
	docker build -t kubeedge/edgecore:${IMAGE_TAG} \
	--build-arg BUILD_FROM=${ARCH}/golang:1.12-alpine3.10 \
	--build-arg RUN_FROM=${ARCH}/docker:dind \
	-f build/edge/Dockerfile .

.PHONY: edgesiteimage
edgesiteimage:
	mkdir -p ./build/edgesite/tmp
	rm -rf ./build/edgesite/tmp/*
	curl -L -o ./build/edgesite/tmp/qemu-${QEMU_ARCH}-static.tar.gz https://github.com/multiarch/qemu-user-static/releases/download/v3.0.0/qemu-${QEMU_ARCH}-static.tar.gz
	tar -xzf ./build/edgesite/tmp/qemu-${QEMU_ARCH}-static.tar.gz -C ./build/edgesite/tmp
	docker build -t kubeedge/edgesite:${IMAGE_TAG} \
	--build-arg BUILD_FROM=${ARCH}/golang:1.12-alpine3.10 \
	--build-arg RUN_FROM=${ARCH}/docker:dind \
	-f build/edgesite/Dockerfile .


.PHONY: bluetoothdevice
bluetoothdevice:
	make -C mappers/bluetooth_mapper

.PHONY: bluetoothdevice_image
bluetoothdevice_image:
	make -C mappers/bluetooth_mapper bluetooth_mapper_image
