CONF ?=$(shell pwd)/configs
PERMIT_PORT ?=40001
.PHONY: generate
generate:
	go generate ./...

.PHONY: run-form
run-form:generate
	dapr run -d ${CONF}/samples --app-id form -- go run cmd/form/main.go --config ${CONF}/config.yml

.PHONY: run-permit
run-permit:generate
	dapr run -d ${CONF}/samples --app-id permit -p ${PERMIT_PORT} -- go run cmd/permit/main.go --config ${CONF}/permit.yml