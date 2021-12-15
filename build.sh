#!/usr/bin/env bash

#######################
# Build and Test script
# Author: Deep-Core
#######################

BASE_DIR="$(realpath $(dirname "$0"))"
cd $BASE_DIR

APPNAME="chore"
VERSION="${IMAGE_TAG:-$(git describe --tags --first-parent --match "v*" 2> /dev/null || echo v0.0.0)}"

MAINGO="cmd/${APPNAME}/main.go"
PKG=$(head -n 1 go.mod | cut -d " " -f2)

FLAG=(
"${PKG}/internal/configs.appName=${APPNAME}"
"${PKG}/internal/configs.appVersion=${VERSION}"
)

FLAGS=$(echo ${FLAG[@]} | xargs -n 1 echo -n " -X")
OUTPUT_FOLDER="${BASE_DIR}/_out"

PLATFORMS="${PLATFORMS:-linux:amd64}"

# set docker
DOCKER_IMAGE_NAME=${DOCKER_IMAGE_NAME:-${APPNAME}}
export GO_IMAGE=${GO_IMAGE:-golang:1.16.11}
export BASE_IMAGE=${BASE_IMAGE:-alpine:3.15.0}
export IMAGE_TAG=${VERSION}

function usage() {
    cat - <<EOF
Build script for golang

Usage: $0 <OPTIONS>
OPTIONS:
  --docker-build
    Build docker image
    -v, --verbose
        use plain output
  --docker-name
    Return docker image name
  --run
    Run for dev
  --swag
    Build swagger docs
  --build
    Build application to various platforms
    --pack
      Pack output
    --install
      install required commands
  --clean
    Clean output folder
  --test
    Test code
    --cover
      Export coverage of test
    --html
      Show html

  -h, --help
    This help page
EOF
}

#######################
# Functions
function build() {
    echo "> Buiding ${APPNAME} for ${1}-${2}"
    OUTPUT_FOLDER_IN=${OUTPUT_FOLDER}/${1}
    mkdir -p ${OUTPUT_FOLDER_IN}
    CGO_ENABLED=0 GOOS=${1} GOARCH=${2} go build -trimpath -ldflags="-s -w ${FLAGS}" -o ${OUTPUT_FOLDER_IN}/${APPNAME} ${MAINGO}
    if [[ "${PACK}" == "Y" ]]; then
        (
            cd ${OUTPUT_FOLDER_IN}
            if [[ "${1}" == "windows" ]]; then
                zip ../${APPNAME}-${1}-${2}-${VERSION}.zip *
            else
                tar czf ../${APPNAME}-${1}-${2}-${VERSION}.tar.gz *
            fi
        )
    fi
}

function pre_build() {
    echo "> Checking swag command"
    if [[ ! $(command -v swag) ]]; then
        echo "> Command swag not found!"
        [[ ! ${AUTO_INSTALL} == "Y" ]] && return 1

        echo "> Installing swag"
        go install github.com/swaggo/swag/cmd/swag@latest
        return $?
    fi

    return 0
}
#######################

#######################
# Run
if [[ -z ${PLATFORMS} ]]; then
    # set default platforms
    PLATFORMS="linux:amd64"
fi

[[ $# -eq 0 ]] && {
    usage
    exit 0
}

while [[ "$1" =~ ^- && ! "$1" == "--" ]]; do
    case "${1}" in
    --docker-build)
        DOCKER_BUILD="Y"
        ;;
    -v | --verbose)
        export BUILDKIT_PROGRESS="plain"
        ;;
    --docker-name)
        echo ${DOCKER_IMAGE_NAME}:${IMAGE_TAG}
        exit 0
        ;;
    --run)
        RUN="Y"
        ;;
    --swag)
        SWAG="Y"
        ;;
    --build)
        BUILD="Y"
        ;;
    --install)
        AUTO_INSTALL="Y"
        ;;
    --pack)
        PACK="Y"
        ;;
    --clean)
        CLEAN="Y"
        ;;
    --test)
        TEST="Y"
        ;;
    --cover)
        COVER="Y"
        ;;
    --html)
        SHOW_HTML="Y"
        ;;
    -h | --help)
        usage
        exit 0
        ;;
    esac
    shift 1
done
if [[ "$1" == '--' ]]; then shift; fi

# docker build
if [[ "${DOCKER_BUILD}" == "Y" ]]; then
    set -e

    export SSH_KEY_CICD_64="$(cat ${SSH_KEY} | base64 | tr -d '\n')"
    # build command
    # sed 's@DEV-BUILD-EDIT@&\nWORKDIR /\nRUN --mount=type=ssh git config --global url."https://".insteadOf git:// \&\& \\\n git config --global http.sslVerify false \&\& \\\n mkdir -p -m 0600 ~/.ssh \&\& \\\n ssh-keyscan -H gitlab.test.igdcs.com >> ~/.ssh/known_hosts \&\& \\\n git clone git\@gitlab.test.igdcs.com:finops/devops/infra-certificates.git@g' Dockerfile | \
    # sed -e 's@COPY _infra-certificates/certs@COPY --from=builder /infra-certificates/certs@g' | \
    cat Dockerfile | \
    cat <(echo '# syntax=docker/dockerfile:experimental') - | \
    DOCKER_BUILDKIT=1 docker build \
        --add-host host.docker.internal:$(docker network inspect bridge | grep Gateway | tr -d '" ' | cut -d ":" -f2) \
        --build-arg GO_IMAGE \
        --build-arg BASE_IMAGE \
        --build-arg GOPROXY=$(go env GOPROXY | sed -e 's@localhost@host.docker.internal@g') \
        --build-arg GOPRIVATE=$(go env GOPRIVATE) \
        --build-arg IMAGE_TAG \
        --build-arg SSH_KEY_CICD_64 \
        --build-arg TRAEFIK_VERSION \
        -t ${DOCKER_IMAGE_NAME}:${IMAGE_TAG} -f - .
    echo "> image name => ${DOCKER_IMAGE_NAME}:${IMAGE_TAG}"
    set +e
    exit
fi

# Create output folder
mkdir -p ${OUTPUT_FOLDER}

# Clean output folder
if [[ "${CLEAN}" == "Y" ]]; then
    echo "> Cleaning builded files..."
	rm -rf ${OUTPUT_FOLDER}/* 2> /dev/null
fi

# Test
if [[ "${TEST}" == "Y" ]]; then
    echo "> Test started"
    [[ "${COVER}" == "Y" ]] && COVERAGE="-coverprofile=${OUTPUT_FOLDER}/cover.out"
	go test -v ./... ${COVERAGE}
    [[ "${SHOW_HTML}" == "Y" ]] && go tool cover -html=${OUTPUT_FOLDER}/cover.out
fi

# Build packages
if [[ "${BUILD}" == "Y" ]]; then
    set -e
    mkdir -p ${OUTPUT_FOLDER}
    pre_build
    [[ $? != 0 ]]; exit 5
    IFS=',' read -ra PLATFORMS_ARR <<< $(echo ${PLATFORMS} | tr -d ' ')
    for PLATFORM_A in "${PLATFORMS_ARR[@]}"; do
        PLATFORM=$(echo ${PLATFORM_A} | cut -d ':' -f 1)
        ARCHS=$(echo ${PLATFORM_A} | cut -d ':' -f 2)
        IFS='-' read -ra ARCHS_ARR <<< ${ARCHS}
        for ARCH in ${ARCHS_ARR[@]}; do
            build ${PLATFORM} ${ARCH}
        done
    done
    set +e
fi

if [[ "${SWAG}" == "Y" ]]; then
    swag init -g router.go --dir internal/server --output docs/
fi

# run
if [[ "${RUN}" == "Y" ]]; then
    set -x
    go run ${MAINGO} ${*}
    set +x
fi
###############
# END
