BINARY_NAME=ludus-pinger
CONFIG_URL=https://public-configurations.s3.ap-southeast-1.amazonaws.com/ekko_sample_config.yaml

define create_pkg_folder
	mkdir -p ${PWD}/dist/ekko-$(1)
	cd ${PWD}/dist/ekko-$(1) && curl ${CONFIG_URL} --output config.yaml
	cd ${PWD}/dist/ekko-$(1) && mkdir -p logs
endef

define compress_pkg
	cd ${PWD}/dist && zip -r ekko-$(1).zip ekko-$(1)/*
endef

define build_pkg
	$(call create_pkg_folder,$(1))
	GOOS=$(1) GOARCH=$(2) go build github.com/soheltarir/ekko

	mv ${PWD}/ekko ${PWD}/dist/ekko-$(1)
	$(call compress_pkg,$(1))
endef

define build_windows_pkg
	$(call create_pkg_folder,windows)
	GOOS=windows GOARCH=$(1) go build github.com/soheltarir/ekko

	mv ${PWD}/ekko.exe ${PWD}/dist/ekko-windows
	$(call compress_pkg,windows)
endef

build:
	$(call build_pkg,darwin,amd64)
	$(call build_pkg,linux,amd64)
	$(call build_windows_pkg,amd64)

clean:
	go clean
	rm -rf dist
