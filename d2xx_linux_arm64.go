// Copyright 2021 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

//go:build cgo && !no_d2xx
// +build cgo,!no_d2xx

package d2xx

/*
#cgo LDFLAGS: ${SRCDIR}/third_party/libftd2xx_linux_arm8hf_v1.4.27.a
*/
import "C"
