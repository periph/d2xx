#!/bin/bash
# Copyright 2017 The Periph Authors. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Builds the package on multiple OSes to confirm it builds fine.

set -eu

cd `dirname $0`

OPT=$*

# Do not set CGO_ENABLED, as we want the default behavior when cross compiling,
# which is different from when CGO_ENABLED=1.
export -n CGO_ENABLED

# Cleanup.
export -n GOOS
export -n GOARCH

function build {
  export GOOS=$1
  export GOARCH=$2
  echo "Building on $GOOS/$GOARCH"
  go build $OPT
  echo "Building on $GOOS/$GOARCH - no_d2xx"
  go build -tags no_d2xx $OPT
  echo "Building on $GOOS/$GOARCH - no cgo"
  CGO_ENABLED=0 go build $OPT
  echo "Building on $GOOS/$GOARCH - no cgo, no_d2xx"
  CGO_ENABLED=0 go build -tags no_d2xx $OPT
}

CGO_ENABLED=1 CC=x86_64-linux-gnu-gcc build linux amd64
# Requires: sudo apt install gcc-arm-linux-gnueabihf
CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc build linux arm
# Requires: sudo apt install gcc-aarch64-linux-gnu
CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc build linux arm64

build linux 386
build windows amd64
build darwin amd64
