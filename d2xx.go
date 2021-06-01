// Copyright 2021 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package d2xx

import (
	"strconv"
)

// Err is the error type returned by d2xx functions.
type Err int

// These are additional synthetic error codes.
const (
	// NoCGO is returned when the package was compiled without cgo, thus the d2xx
	// library is unavailable or the library was disabled via the `no_d2xx` build
	// tag.
	NoCGO Err = -2
	// Missing is returned when the dynamic library is not available.
	Missing Err = -1
)

// String converts a error integer to a string representation of the error.
func (e Err) String() string {
	switch e {
	case Missing:
		// when the library d2xx couldn't be loaded at runtime.
		return "couldn't load driver; visit https://periph.io/device/ftdi/ for help"
	case NoCGO:
		return "can't be used without cgo"
	case 0: // FT_OK
		return ""
	case 1: // FT_INVALID_HANDLE
		return "invalid handle"
	case 2: // FT_DEVICE_NOT_FOUND
		return "device not found; see https://periph.io/device/ftdi/ for help"
	case 3: // FT_DEVICE_NOT_OPENED
		return "device busy; see https://periph.io/device/ftdi/ for help"
	case 4: // FT_IO_ERROR
		return "I/O error"
	case 5: // FT_INSUFFICIENT_RESOURCES
		return "insufficient resources"
	case 6: // FT_INVALID_PARAMETER
		return "invalid parameter"
	case 7: // FT_INVALID_BAUD_RATE
		return "invalid baud rate"
	case 8: // FT_DEVICE_NOT_OPENED_FOR_ERASE
		return "device not opened for erase"
	case 9: // FT_DEVICE_NOT_OPENED_FOR_WRITE
		return "device not opened for write"
	case 10: // FT_FAILED_TO_WRITE_DEVICE
		return "failed to write device"
	case 11: // FT_EEPROM_READ_FAILED
		return "eeprom read failed"
	case 12: // FT_EEPROM_WRITE_FAILED
		return "eeprom write failed"
	case 13: // FT_EEPROM_ERASE_FAILED
		return "eeprom erase failed"
	case 14: // FT_EEPROM_NOT_PRESENT
		return "eeprom not present"
	case 15: // FT_EEPROM_NOT_PROGRAMMED
		return "eeprom not programmed"
	case 16: // FT_INVALID_ARGS
		return "invalid argument"
	case 17: // FT_NOT_SUPPORTED
		return "not supported"
	case 18: // FT_OTHER_ERROR
		return "other error"
	case 19: // FT_DEVICE_LIST_NOT_READY
		return "device list not ready"
	default:
		return "unknown status " + strconv.Itoa(int(e))
	}
}

// unknown is a forward declaration of ftdi.DevType.
const unknown = 3

// handle is a d2xx handle.
//
// This is the base type which each OS specific implementation adds methods to.
type handle uintptr

// Handle is d2xx device handle.
type Handle interface {
	Close() Err
	// ResetDevice takes >1.2ms
	ResetDevice() Err
	GetDeviceInfo() (uint32, uint16, uint16, Err)
	EEPROMRead(devType uint32, e *EEPROM) Err
	EEPROMProgram(e *EEPROM) Err
	EraseEE() Err
	WriteEE(offset uint8, value uint16) Err
	EEUASize() (int, Err)
	EEUARead(ua []byte) Err
	EEUAWrite(ua []byte) Err
	SetChars(eventChar byte, eventEn bool, errorChar byte, errorEn bool) Err
	SetUSBParameters(in, out int) Err
	SetFlowControl() Err
	SetTimeouts(readMS, writeMS int) Err
	SetLatencyTimer(delayMS uint8) Err
	SetBaudRate(hz uint32) Err
	// GetQueueStatus takes >60µs
	GetQueueStatus() (uint32, Err)
	// Read takes <5µs if GetQueueStatus was called just before,
	// 300µs~800µs otherwise (!)
	Read(b []byte) (int, Err)
	// Write takes >0.1ms
	Write(b []byte) (int, Err)
	GetBitMode() (byte, Err)
	// SetBitMode takes >0.1ms
	SetBitMode(mask, mode byte) Err
}

var _ Handle = handle(0)

// Version returns the library's version.
//
// 0, 0, 0 is returned if the library is unavailable.
func Version() (uint8, uint8, uint8) {
	return version()
}

// CreateDeviceInfoList discovers the currently found devices.
//
// If the driver is disabled via build tag `no_d2xx`, or on posix
// `CGO_ENABLED=0` environment variable, NoCGO is returned.
//
// On Windows, Missing is returned if the dynamic library is not found at
// runtime.
func CreateDeviceInfoList() (int, Err) {
	return createDeviceInfoList()
}

// Open opens the ith device discovered.
//
// If the driver is disabled via build tag `no_d2xx`, or on posix
// `CGO_ENABLED=0` environment variable, NoCGO is returned.
//
// On Windows, Missing is returned if the dynamic library is not found at
// runtime.
func Open(i int) (Handle, Err) {
	return open(i)
}
