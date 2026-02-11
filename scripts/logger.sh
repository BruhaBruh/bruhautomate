#!/usr/bin/env bash

source ./scripts/utilities.sh

log_level="info"

declare -A LOG_PRIORITY=(
  [debug]=0
  [log]=1
  [info]=2
  [warning]=3
  [error]=4
)

log() {
  local level
  local first="$1"

  case "${first,,}" in
    debug|info|warning|warn|error|err)
      level="${first,,}"
      shift
      ;;
    *)
      level="log"
      ;;
  esac

  [[ "$level" == "warn" ]] && level="warning"
  [[ "$level" == "err" ]] && level="error"

  local current_priority="${LOG_PRIORITY[$log_level]}"
  local msg_priority="${LOG_PRIORITY[$level]}"
  (( msg_priority < current_priority )) && return 0  # не выводим

  case "$level" in
    debug)
      printf "${BLACK_B}[${GREEN}DEBUG${BLACK_B}]${RESET} "
      ;;
    info)
      printf "${BLACK_B}[${BLUE}INFO${BLACK_B}]${RESET} "
      ;;
    warning)
      printf "${BLACK_B}[${YELLOW}WARNING${BLACK_B}]${RESET} "
      ;;
    error)
      printf "${BLACK_B}[${RED}ERROR${BLACK_B}]${RESET} "
      ;;
    log)
      printf "${BLACK_B}[${MAGENTA}LOG${BLACK_B}]${RESET} "
      ;;
  esac

  printf "%s\n" "$*"
}

debug() {
  log debug "$@"
}

info() {
  log info "$@"
}

warn() {
  log warn "$@"
}

error() {
  log error "$@"
}
