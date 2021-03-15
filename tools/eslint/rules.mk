eslint_cwd := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
eslint := $(eslint_cwd)/node_modules/.bin/eslint

$(eslint): $(eslint_cwd)/package.json
	$(info [eslint] installing package...)
	@cd $(eslint_cwd) && npm install --no-save --no-audit &> /dev/null
	@touch $@
