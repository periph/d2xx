// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build cgo
// +build !no_d2xx

package d2xx

// TODO(maruel): https://github.com/golang/go/issues/7211 would help target the
// optimal ARM architecture.

/*
#cgo LDFLAGS: ${SRCDIR}/third_party/libftd2xx_linux_arm6hf_v1.4.24.a
*/
import "C"
