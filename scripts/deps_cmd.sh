#!/usr/bin/env bash

set -e

source ./scripts/logger.sh

deps() {
  if [ $# -eq 0 ]; then
    info "Download and optimize dependencies..."
    go mod tidy && go mod download
    return
  fi

  local cmd="$1"
  shift

  case "$cmd" in
    help|--help|-h)
      deps_help
      ;;
    *)
      error "Invalid arguments: $cmd $@"
      ;;
  esac
}

deps_help() {
  printf "${BLACK_B}download and optimize dependencies\n\n"
  printf "$1${BLUE}ARGUMENTS${RESET}\n"
  printf "$1\thelp|-h|--help$BLACK_B - show help$RESET\n"
}

deps_signature() {
  printf "deps|dependencies"
}
