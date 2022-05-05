#!/usr/bin/env bash

# JWT_KEY required

BASE_DIR="$(realpath $(dirname "$0"))"
cd $BASE_DIR

function usage() {
    cat - <<EOF
Record script for chore

Usage: $0 <OPTIONS>
OPTIONS:
  --chore-url <URL>
    Show url of chore
    Ex: http://localhost:8080
        http://example.com/chore

  --mode <upload,download>
    Set mode 'upload' or 'download'
  --overwrite
    If already exist update it
    Works on upload mode

  --templates
    Templates operation
    Work in upload mode to upload all templates folder
  --template <NAME>
    In upload, download mode this is template name

  --controls
    Controls operation
    Work in upload mode
  --control <NAME>
    In upload, download mode this is control name

  --auths
    Auths operation
    Work in upload mode
  --auth
    In upload, download mode this is auth name

  -h, --help
    This help page
EOF
}

[[ $# -eq 0 ]] && {
    usage
    exit 0
}

set -o allexport
while [[ "$1" =~ ^- && ! "$1" == "--" ]]; do
  case "${1}" in
  --url)
    URL="${2}"
    shift
    ;;
  --mode)
    MODE="${2}"
    if [[ "${MODE}" != "download" && "${MODE}" != "upload" ]]; then
      echo "> mode should be 'download' or 'upload'"
      exit 1
    fi
    shift
    ;;
  --overwrite)
    OVERWRITE="Y"
    ;;
  --templates)
    TEMPLATES="Y"
    ;;
  --template)
    TEMPLATE="${2}"
    shift
    ;;
  --controls)
    CONTROLS="Y"
    ;;
  --control)
    CONTROL="${2}"
    shift
    ;;
  --auths)
    AUTHS="Y"
    ;;
  --auth)
    AUTH="${2}"
    shift
    ;;
  --test)
    TEST="Y"
    shift
    ;;
  -h | --help)
    usage
    exit 0
    ;;
  *)
    echo "> Not found $1"
    exit 1
    ;;
  esac
  shift 1
done
if [[ "$1" == '--' ]]; then shift; fi

if [[ -z ${TEST} ]]; then
  if [[ -z ${JWT_KEY} ]]; then
    echo "> JWT_KEY must be set"
    exit 1
  fi

  if [[ -z ${URL} ]]; then
    echo "> --url must be set"
    exit 1
  fi
fi

# $1 -> api
# $2 -> name
# $3 -> file_name
function requestUpload() {
  local CONVERTED_NAME=$(echo ${2} | sed s@/@%2F@g)
  curl -ksSL -X 'PUT' --data-binary @${3} \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${JWT_KEY}" \
    "${URL}/${1}?name=${CONVERTED_NAME}"
  echo "Uploaded ${3} to ${1}"
}

# $1 -> api
# $2 -> name
# $3 -> file_name
function requestDownload() {
  local CONVERTED_NAME=$(echo ${2} | sed s@/@%2F@g)
  curl -ksSL -X 'GET' --create-dirs -o "${3}" \
    -H "Authorization: Bearer ${JWT_KEY}" \
    "${URL}/${1}?name=${CONVERTED_NAME}&dump=true&pretty=true"
  echo "Dowloaded ${3} from ${1}"
}

# $1 -> namespace
# $2 -> mode
# $3 -> name
function process() {
  case "${1}" in
  auth)
    local API="api/v1/auth"
    local AUTH_FILE="auths/${3}.json"
    if [[ ${2} == "download" ]]; then
      requestDownload "${API}" "${3}" "${AUTH_FILE}"
    elif [[ ${2} == "upload" ]]; then
      requestUpload "${API}" "${3}" "${AUTH_FILE}"
    fi
    ;;
  control)
    local API="api/v1/control"
    local CONTROL_FILE="controls/${3}.json"
    if [[ ${2} == "download" ]]; then
      requestDownload "${API}" "${3}" "${CONTROL_FILE}"
    elif [[ ${2} == "upload" ]]; then
      requestUpload "${API}" "${3}" "${CONTROL_FILE}"
    fi
    ;;
  template)
    local API="api/v1/template"
    local TEMPLATE_FILE="templates/${3}.tmpl"
    if [[ ${2} == "download" ]]; then
      requestDownload "${API}" "${3}" "${TEMPLATE_FILE}"
    elif [[ ${2} == "upload" ]]; then
      requestUpload "${API}" "${3}" "${TEMPLATE_FILE}"
    fi
    ;;
  esac
}

# $1 -> extension
# $2 -> folder
function findIt() {
  local EXTENSION="${1}"
  local FOLDER_PATH="${2}"
  for FILE_ in $(find ${FOLDER_PATH}/ -name "*.${EXTENSION}" -type f -not -path '*/.*'); do
    local NAME=$(echo ${FILE_} | sed "s@${FOLDER_PATH}/\(.*\).${EXTENSION}@\1@g")
    echo "${NAME}"
  done
}

set +o allexport
# set error option
set -e

# single operations
if [[ -n "${AUTH}" ]]; then
  process auth "${MODE}" "${AUTH}"
fi
if [[ -n "${TEMPLATE}" ]]; then
  process template "${MODE}" "${TEMPLATE}"
fi
if [[ -n "${CONTROL}" ]]; then
  process control "${MODE}" "${CONTROL}"
fi

# folder operations
if [[ "${AUTHS}" == "Y" ]]; then
  findIt "json" "auths" | xargs -I {} bash -c 'process auth "${MODE}" "{}"'
fi
if [[ "${CONTROLS}" == "Y" ]]; then
  findIt "json" "controls" | xargs -I {} bash -c 'process control "${MODE}" "{}"'
fi
if [[ "${TEMPLATES}" == "Y" ]]; then
  findIt "tmpl" "templates" | xargs -I {} bash -c 'process template "${MODE}" "{}"'
fi
