#!/usr/bin/env bash

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

source "${script_dir}/env"

podman_id=$(podman ps --filter "name=wiremock" --filter "status=running" --quiet)

if [[ -z "$podman_id" ]]; then
    echo "${CYAN}Starting WireMock...${NC}"
    if ! podman run \
        --detach \
        --name wiremock \
        --publish 8443:8443 \
        --rm \
        --volume "$PWD/testdata:/home/wiremock:ro" \
        wiremock/wiremock:2.35.0 \
        --disable-banner \
        --disable-http \
        --global-response-templating \
        --https-port 8443 \
        >/dev/null; then
        echo "[${RED}ERROR${NC}] Failed to start Wiremock"
        exit 1
    fi
fi

echo "[${GREEN}OK${NC}] WireMock is running"
