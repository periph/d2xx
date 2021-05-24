// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package d2xx

import (
	"unsafe"
)

// EEPROM is the unprocessed EEPROM content.
//
// The EEPROM is in 3 parts: the defined struct, the 4 strings and the rest
// which is used as an 'user area'. The size of the user area depends on the
// length of the strings. The user area content is not included in this struct.
type EEPROM struct {
	// Raw is the raw EEPROM content. It excludes the strings.
	Raw []byte

	// The following condition must be true: len(Manufacturer) + len(Desc) <= 40.
	Manufacturer   string
	ManufacturerID string
	Desc           string
	Serial         string
}

func (e *EEPROM) asHeader() *eepromHeader {
	// sizeof(EEPROMHeader)
	if len(e.Raw) < 16 {
		return nil
	}
	return (*eepromHeader)(unsafe.Pointer(&e.Raw[0]))
}

// eepromHeader is the common 16 bytes header.
type eepromHeader struct {
	DeviceType uint32
	// The rest is not necessary here so it is skipped.
}
