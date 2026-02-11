#!/usr/bin/env bash

set -e

source ./scripts/run_cmd.sh
source ./scripts/deps_cmd.sh
source ./scripts/test_cmd.sh
source ./scripts/logger.sh
source ./scripts/utilities.sh

log_level="info"

main() {
  if [ $# -eq 0 ]; then
      show_help
  fi

  local cmd="$1"
  shift

  case "$cmd" in
    run)
      run "$@"
      ;;
    deps|dependencies)
      deps "$@"
      ;;
    test)
      test "$@"
      ;;
    help|--help|-h)
      show_help
      ;;
    *)
      log "Invalid command"
      exit 1
      ;;
  esac
}

show_help() {
  printf "${GREEN}ctl$BLACK_B - cli for control project\n\n"

  printf "${BLUE}COMMANDS${RESET}\n"
  printf "\t$(run_signature)$BLACK_B - $(run_help "\t")\n\n"
  printf "\t$(deps_signature)$BLACK_B - $(deps_help "\t")\n\n"
  printf "\t$(test_signature)$BLACK_B - $(test_help "\t")\n\n"

  printf "${BLUE}ARGUMENTS${RESET}\n"
  printf "\thelp|-h|--help$BLACK_B - show help\n"
}

main "$@"
