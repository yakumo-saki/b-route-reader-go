#!/bin/bash -eu

SCRIPT_DIR=$(cd $(dirname $0); pwd)

# Get version from global/version.go
VERSION=`cat ./src/global/version.go | grep -o "Version.*\".*\"" | sed -e 's/.*=\s//g' -e 's/"//g'`
COMMIT=`git rev-parse --short HEAD`
if ! git diff --quiet; then
  COMMIT="${COMMIT}-working"
fi

BIN_BASENAME=b-route-reader-go
ENTRYPOINT=main.go

BUILD_DIR=./build
BIN_DIR=${SCRIPT_DIR}/${BUILD_DIR}/bin
WORK_DIR=${SCRIPT_DIR}/${BUILD_DIR}/work
RELEASE_DIR=${SCRIPT_DIR}/${BUILD_DIR}/release

# MAIN
sed -i src/global/version.go -e 's/GitBuild = ".*"/GitBuild = "$COMMIT"/g'

# delete build dir
rm -rf ${BUILD_DIR}

mkdir -p ${BIN_DIR}
mkdir -p ${RELEASE_DIR}
mkdir -p ${WORK_DIR}

cp LICENSE ${WORK_DIR}
cp README.md ${WORK_DIR}

function build_unixlike () {
    # $1 OS $2 ARCH
    echo Building $1 $2 binary

    FINAL_PATH=${RELEASE_DIR}/${BIN_BASENAME}_${VERSION}_$1_$2.tar.gz
    GOOS=$1 GOARCH=$2 CGO_ENABLED=0 go build -o ${BIN_DIR}/${BIN_BASENAME} ${ENTRYPOINT}

    # copy bin to work
    cp ${BIN_DIR}/${BIN_BASENAME} ${WORK_DIR}/

    # copyback bin
    cp ${BIN_DIR}/${BIN_BASENAME} ${BIN_DIR}/${BIN_BASENAME}_$1_$2

    ORG_DIR=`pwd`
    cd ${WORK_DIR}
    tar -cvzf ${FINAL_PATH} --exclude *.tar.gz ./*
    cd ${ORG_DIR}
    
    echo "done => ${FINAL_PATH}"
}

# Unixlike
build_unixlike linux amd64
build_unixlike linux arm
build_unixlike linux arm64

git checkout src/global/version.go