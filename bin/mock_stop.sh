#!/usr/bin/env bash

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

source "${script_dir}/env.sh"

podman_id=$(podman ps --filter "name=wiremock" --filter "status=running" --quiet)

if [[ -n "$podman_id" ]]; then
    echo "${CYAN}Stopping WireMock...${NC}"
    if ! podman stop wiremock >/dev/null; then
        echo "[${RED}ERROR${NC}] Failed to stop Wiremock"
        exit 1
    fi
fi

echo "[${GREEN}OK${NC}] WireMock is stopped"
