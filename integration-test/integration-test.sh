#!/bin/bash

set -ex

export BINARY=$1
export GOCOVERDIR=$2

main() {
  # Run the application in the background
  $BINARY ./.github/.env &
  PID=$!

  sleep 5 #TODO run the tests here

  # Shutdown the application
  kill $PID
}

main "$@"
