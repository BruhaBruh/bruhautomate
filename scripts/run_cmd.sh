#!/usr/bin/env bash

set -e

source ./scripts/deps_cmd.sh
source ./scripts/logger.sh

run() {
  deps
  info "Run application..."
  go run . "$@"
}

run_help() {
  printf "${BLACK_B}run application\n\n"
}

run_signature() {
  printf "run"
}
