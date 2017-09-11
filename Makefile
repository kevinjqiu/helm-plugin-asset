.PHONY: build dep clean install

dep:
	dep ensure

clean:
	rm -rf build/

build/helm-plugin-asset:
	go build -o build/helm-plugin-asset

install: clean build/helm-plugin-asset
	mkdir -p $(HOME)/.helm/plugins/asset
	unlink $(HOME)/.helm/plugins/asset/plugin.yaml
	ln -s $$(pwd)/plugin.yaml $(HOME)/.helm/plugins/asset/plugin.yaml
	unlink $(HOME)/.helm/plugins/asset/helm-plugin-asset
	ln -s $$(pwd)/build/helm-plugin-asset $(HOME)/.helm/plugins/asset/helm-plugin-asset
