#!/usr/bin/env bash

set -e

source ./scripts/deps_cmd.sh
source ./scripts/logger.sh

test() {
  local is_verbose=false
  local is_coverage=false
  local path=("./...")
  local args=()

  while [[ $# -gt 0 ]]; do
    case "$1" in
      verbose|--verbose|-v)
        is_verbose=true
        shift
        ;;
      coverage|--coverage|-c)
        is_coverage=true
        shift
        ;;
      path|--path|-p)
        shift
        path=()
        while [[ $# -gt 0 ]]; do
          case "$1" in
            verbose|--verbose|-v|coverage|--coverage|-c|path|--path|-p|help|--help|-h)
              break
              ;;
            *)
              path+=("$1")
              shift
              ;;
          esac
        done
        ;;
      help|--help|-h)
        test_help
        return
        ;;
      *)
        args+=("$1")
        shift
        ;;
    esac
  done

  cmd=(go test)
  $is_verbose && cmd+=(-v)
  $is_coverage && cmd+=(-cover)
  cmd+=("${args[@]}" "${path[@]}")

  info "Run tests..."
  info "${cmd[*]}"
  "${cmd[@]}"
}


test_help() {
  printf "${BLACK_B}test application\n\n"
  printf "$1${BLUE}ARGUMENTS${RESET}\n"
  printf "$1\tverbose|-v|--verbose$BLACK_B - verbose tests$RESET\n"
  printf "$1\tcoverage|-c|--coverage$BLACK_B - coverage tests$RESET\n"
  printf "$1\tpath|-p|--path [...paths]$BLACK_B - custom paths$RESET\n"
  printf "$1\thelp|-h|--help$BLACK_B - show help$RESET\n"
}

test_signature() {
  printf "test [-c] [-v] [-p [...paths]]  [...cmd args]"
}
