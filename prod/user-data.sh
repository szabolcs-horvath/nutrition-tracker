#!/bin/bash

# logging
exec >> /var/log/user-data.log
exec 2>&1

set -ex -o pipefail

export HOME=/root

install_go() {
  curl -LO https://go.dev/dl/go1.23.1.linux-amd64.tar.gz
  rm -rf /usr/local/go
  tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz
  echo "export PATH=$PATH:/usr/local/go/bin:/root/go/bin" >> /etc/profile
  echo "export GOPATH=/root/go" >> /etc/profile
  source /etc/profile
  go version
}

main() {
  apt-get update

  if ! which make || ! which gcc; then
    apt-get install -y build-essential
    which make
    which gcc
  fi

  if ! which go; then
    install_go
  fi

  if ! which sqlite3; then
    apt-get install -y sqlite3
    which sqlite3
  fi

  if [[ ! -d ~/nutrition-tracker ]]; then
    cd ~
    git clone https://github.com/szabolcs-horvath/nutrition-tracker.git
    cd nutrition-tracker
    make install
  else
    echo "Nutrition Tracker already installed."
  fi
}

main "$@"
