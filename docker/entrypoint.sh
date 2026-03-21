#!/bin/bash
set -e

if [ ! -z "${ASTRIA_UID}" ]; then
  if [ ! "$(id -u astria)" -eq "${ASTRIA_UID}" ]; then
    # Change the UID
    usermod -o -u "${ASTRIA_UID}" astria
  fi
fi

if [ ! -z "${ASTRIA_GID}" ]; then
  if [ ! "$(id -g astria)" -eq "${ASTRIA_GID}" ]; then
    groupmod -o -g "${ASTRIA_GID}" astria
  fi
fi

args=( "$@" )
# Login shell to properly set env vars
exec sudo -E -H -u astria "$@"

