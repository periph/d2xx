// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build cgo
// +build !darwin,!amd64
// +build !linux,!amd64
// +build !linux,!arm
// +build !windows
// +build !no_d2xx

package d2xx

// This assumes the library is installed and available for linking.

/*
#cgo LDFLAGS: -lftd2xx
*/
import "C"
