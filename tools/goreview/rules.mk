goreview_cwd := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
goreview_version := 0.16.0
goreview := $(goreview_cwd)/$(goreview_version)/goreview

arch := $(shell uname -s)_$(shell uname -m)

# enforce x86 arch if mac m1 until tool has official support
ifeq ($(arch),Darwin_arm64)
arch = Darwin_x86_64
endif

goreview_archive_url := https://github.com/einride/goreview/releases/download/v$(goreview_version)/goreview_$(goreview_version)_$(arch).tar.gz

$(goreview): $(goreview_cwd)/rules.mk
	$(info [go-review] fetching $(goreview_version) binary...)
	@mkdir -p $(dir $@)
	@curl -sSL $(goreview_archive_url) -o - | tar -xz --directory $(dir $@)
	@chmod +x $@
	@touch $@

# go-review: review Go code for Einride-specific conventions
.PHONY: go-review
go-review: $(goreview)
	$(info [$@] reviewing Go code for Einride-specific conventions...)
	@$(goreview) -c 1 ./...
