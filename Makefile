CONF ?=$(shell pwd)/configs
.PHONY: run-form
generate:
	dapr run -d ${CONF}/samples --app-id form -- go run cmd/main.go --config ${CONF}/config.yml
