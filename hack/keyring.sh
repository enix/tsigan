#!/bin/bash

[ -d hack ] || {
  echo "Run this script from the project root with: ./hack/$(basename $0)" >&2
  exit 1
}

[ -d .keyring ] && {
  echo "Development keyring already exists." >&2
}

set -xe

mkdir .keyring

NAME=$(dnssec-keygen -K .keyring/ -a ED25519 -q tsigan)
ln -sr .keyring/$NAME.key     .keyring/tsigan-ed25519.key
ln -sr .keyring/$NAME.private .keyring/tsigan-ed25519.private
