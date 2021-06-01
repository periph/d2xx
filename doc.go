// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package d2xx is a thin Go wrapper for the Future Technology "D2XX" driver.
//
// This package is not Go idiomatic. You want to use
// https://periph.io/x/host/v3/ftdi (or later) instead.
//
// A static library of the d2xx driver is included for linux and macOS. They are
// from https://ftdichip.com/drivers/d2xx-drivers/. See third_party/README.md
// for more details.
//
// Configuration
//
// See https://periph.io/device/ftdi/ for more details, and how to configure
// the host to be able to use this driver.
package d2xx
