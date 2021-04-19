#!/bin/bash
set -e

if [ "$1" = 'run' ]; then
  shift
  exec ./bootstrap "$@"
fi
exec "$@"