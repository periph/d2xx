// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build cgo
// +build !no_d2xx

package d2xx

/*
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit ${SRCDIR}/third_party/libftd2xx_darwin_amd64_v1.4.4.a
*/
import "C"
