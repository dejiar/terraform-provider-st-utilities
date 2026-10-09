# The name of Terraform custom provider.
CUSTOM_PROVIDER_NAME ?= terraform-provider-st-utilities
# The url of Terraform provider.
CUSTOM_PROVIDER_URL ?= example.local/myklst/st-utilities
# The binary name produced by go install (module path's last segment)
GO_BINARY_NAME ?= st-utilities

.PHONY: install-local-custom-provider
install-local-custom-provider:
	export PROVIDER_LOCAL_PATH='$(CUSTOM_PROVIDER_URL)'
	go build -o $(GO_BINARY_NAME) .
	GO_INSTALL_PATH="$(pwd)"; \
	HOME_DIR="$(ls -d ~)"; \
	mkdir -p  $$HOME_DIR/.terraform.d/plugins/$(CUSTOM_PROVIDER_URL)/0.1.0/linux_amd64/; \
	cp $$GO_INSTALL_PATH/$(GO_BINARY_NAME) $$HOME_DIR/.terraform.d/plugins/$(CUSTOM_PROVIDER_URL)/0.1.0/linux_amd64/$(CUSTOM_PROVIDER_NAME)
	unset PROVIDER_LOCAL_PATH
