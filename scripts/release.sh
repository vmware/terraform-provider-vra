#!/bin/bash

COMMANDS="gox github-release"
NAME=terraform-provider-vra
NOWDATE=$(TZ=UTC date +%FT%T%Z)
VERSION=$(git describe --tags --always)
VERSION_LONG=$(git describe --tags --long --always)

check_dependencies() {
    deps=$1
    for cmd in $deps; do
        which $cmd > /dev/null 2>&1
        if [ $? -ne 0 ]; then
            echo "Command $cmd not found, please install it."
            exit 1
        fi
    done
}

build_providers() {
    rm -rf build
    mkdir build
    GO111MODULE=on CGO_ENABLED=0 gox -ldflags "-X main.version=${VERSION} -X main.buildTime=${NOWDATE}" \
        -os "linux darwin windows" -arch "amd64" -output "build/{{.OS}}_{{.Arch}}/terraform-provider-vra_${VERSION}"
    if [ $? -ne 0 ]; then
        echo "Error building providers...exiting"
        exit 1
    fi
}

package(){
    rm -rf release
    mkdir release
    for f in build/*; do \
        g=`basename $f`; \
        tar -zcf release/${NAME}-${g}-${VERSION}.tgz -C build/${g} .; \
    done
}

release(){
	github-release release \
		--user vmware \
		--repo "${NAME}" \
		--target "$(git rev-parse --abbrev-ref HEAD)" \
		--tag "${VERSION}" \
		--name "${VERSION}"
	cd release && ls *.tgz | xargs -I FILE github-release upload \
		--user vmware \
		--repo "${NAME}" \
		--tag "${VERSION}" \
		--name FILE \
		--file FILE
}

# Verify GITHUB_TOKEN is available before doing any further processing
if [ -z "$GITHUB_TOKEN" ]
then
    echo "GITHUB_TOKEN not set in environment...exiting"
    exit 1
fi

check_dependencies "$COMMANDS"
build_providers
package
release
