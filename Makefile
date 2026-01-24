# Copyright (C) 2025 The Concord Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL := bash

PATH := $(GOBIN):$(PATH)
GOBIN := $(shell go env GOPATH)/bin
LDFLAGS=

DATE=$(shell date '+%Y-%m-%d')
HOSTNAME=$(shell hostname)
LOG_DIR=log

GIT_ROOT=github.com/cybergarage
PRODUCT_NAME=go-concord
MODULE_ROOT=${GIT_ROOT}/${PRODUCT_NAME}

PKG_NAME=concord
PKG_VER=$(shell git tag | tail -n 1)
PKG_COVER=${PKG_NAME}-cover
PKG_ROOT=${MODULE_ROOT}/${PKG_NAME}

PKG_SRC_ROOT=${PKG_NAME}
PKG=${MODULE_ROOT}/${PKG_SRC_ROOT}

TEST_SRC_ROOT=${PKG_NAME}test
TEST_PKG=${MODULE_ROOT}/${TEST_SRC_ROOT}

PHONY: test format vet lint clean doc
.IGNORE: lint

all: test

version:
	@pushd ${PKG_SRC_ROOT} && ./version.gen > version.go && popd
	-git commit ${PKG_SRC_ROOT}/version.go -m "Update version"

format: version
	gofmt -s -w ${PKG_SRC_ROOT} ${TEST_SRC_ROOT}

vet: format
	go vet ${PKG}/... ${TEST_PKG}/...

lint: format
	golangci-lint run ${PKG_SRC_ROOT}/... ${TEST_SRC_ROOT}/...

%.md : %.adoc
	asciidoctor -b docbook -a leveloffset=+1 -o - $< | pandoc -t markdown_strict --wrap=none -f docbook > $@
	-git add $@ $< 
	-git commit $@ $< -m "docs: update $@"

csvs := $(wildcard doc/*/*.csv doc/*/*/*.csv)
docs := $(patsubst %.adoc,%.md,$(wildcard *.adoc doc/*.adoc doc/*/*.adoc))
doc: $(docs) $(csvs)

test: lint
	rm -f ${PKG_COVER}.out ${PKG_COVER}.local.out ${PKG_COVER}.html
	go test -v -count=1 -p 1 -timeout 60m \
	-gcflags=${GCFLAGS} -ldflags=${LDFLAGS} \
	-cover -coverpkg=${PKG}/... -coverprofile=${PKG_COVER}.out \
	${PKG}/... ${TEST_PKG}/...
	@sed -e 's|^${MODULE_ROOT}/|${CURDIR}/|' ${PKG_COVER}.out > ${PKG_COVER}.local.out
	go tool cover -html=${PKG_COVER}.local.out -o ${PKG_COVER}.html

build:
	go build -v -gcflags=${GCFLAGS} -ldflags=${LDFLAGS} ${BINS}

log:
	git log ${PKG_VER}..HEAD --date=short --no-merges --pretty=format:"%s"

clean:
	go clean -i ${PKG}
	find . -name "*.log" -or -name "*.prof" | xargs -I{} rm -f {}
	rm -f ${PKG_COVER}.out ${PKG_COVER}.local.out ${PKG_COVER}.html
