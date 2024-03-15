CONFIG_PATH=${HOME}/.proglog-example

.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}

.PHONY: gencert
gencert:
	cfssl gencert \
		-initca test/ca-csr.json | cfssljson -bare ca

	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=server \
		test/server-csr.json | cfssljson -bare server

	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=client \
		-cn="root" \
		test/client-csr.json | cfssljson -bare root-client

	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=client \
		-cn="nobody" \
		test/client-csr.json | cfssljson -bare nobody-client

	mv *.pem *.csr ${CONFIG_PATH}

# Generate certs for other CA to test that the client can't connect to the server when using a different CA
.PHONY: gencert-other
gencert-other:
	cfssl gencert \
		-initca test/other-ca-csr.json | cfssljson -bare other-ca

	cfssl gencert \
		-ca=other-ca.pem \
		-ca-key=other-ca-key.pem \
		-config=test/ca-config.json \
		-profile=server \
		test/server-csr.json | cfssljson -bare other-server

	cfssl gencert \
		-ca=other-ca.pem \
		-ca-key=other-ca-key.pem \
		-config=test/ca-config.json \
		-profile=client \
		test/client-csr.json | cfssljson -bare other-client

	mv *.pem *.csr ${CONFIG_PATH}

.PHONY: lscerts
lscerts:
	ls -l ${CONFIG_PATH}

.PHONY: compile
compile:
	protoc api/v1/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.

$(CONFIG_PATH)/model.conf:
	cp test/model.conf ${CONFIG_PATH}/model.conf

$(CONFIG_PATH)/policy.csv:
	cp test/policy.csv ${CONFIG_PATH}/policy.csv

.PHONY: test
test: $(CONFIG_PATH)/model.conf $(CONFIG_PATH)/policy.csv
	go test -race ./...
