// Copyright 2021 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package d2xx_test

import (
	"fmt"

	"periph.io/x/d2xx"
)

func ExampleVersion() {
	// Print the d2xx driver version. It will be 0.0.0 if not found.
	major, minor, build := d2xx.Version()
	fmt.Printf("Using library %d.%d.%d\n", major, minor, build)
}
