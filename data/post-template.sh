#!/usr/bin/env bash

# JWT_KEY required

BASE_DIR="$(realpath $(dirname "$0"))"
cd $BASE_DIR

for template in templates/*.json; do
    name=$(basename ${template%.json})
    echo "> post ${name}"
    curl -X 'POST' \
        "http://localhost:8080/api/v1/template?name=internal%2F${name}" \
        -H "accept: application/json" \
        -H "Authorization: Bearer ${JWT_KEY}" \
        -d @${template}
    echo ""
done
