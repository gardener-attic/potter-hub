# Required helm.sh/helm to be checked out in your gopath
src/shared/hapi/release.js:
	./node_modules/.bin/pbjs -t static-module -w commonjs -o src/shared/hapi/release.js -p $$GOPATH/src/helm.sh/helm/ $$GOPATH/src/helm.sh/helm/pkg/release/ && \
	./node_modules/.bin/pbts -o src/shared/hapi/release.d.ts src/shared/hapi/release.js && \
	echo "// tslint:disable\n$$(cat src/shared/hapi/release.js)" > src/shared/hapi/release.js && \
	echo "// tslint:disable\n$$(cat src/shared/hapi/release.d.ts)" > src/shared/hapi/release.d.ts
