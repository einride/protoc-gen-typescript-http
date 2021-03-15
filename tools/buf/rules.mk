buf_cwd := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
buf := $(buf_cwd)/bin/buf

buf_version := 0.39.1

arch = $(shell uname -s)-$(shell uname -m)

# enforce x86 arch if mac m1 until tool has official support
ifeq ($(arch),Darwin-arm64)
arch = Darwin-x86_64
endif

buf_bin_url := https://github.com/bufbuild/buf/releases/download/v$(buf_version)/buf-$(arch)

$(buf): $(buf_cwd)/rules.mk
	$(info [buf] feching $(buf_version) binary...)
	@mkdir -p $(dir $@)
	@curl -sSL $(buf_bin_url) -o $@
	@chmod +x $@
	@touch $@
