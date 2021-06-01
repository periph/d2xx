// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !cgo no_d2xx
// +build !windows no_d2xx

package d2xx

// Available is true if the library is available on this system.
const Available = false

// Library functions.

func version() (uint8, uint8, uint8) {
	return 0, 0, 0
}

func createDeviceInfoList() (int, Err) {
	return 0, NoCGO
}

func open(i int) (Handle, Err) {
	return handle(0), NoCGO
}

func (h handle) Close() Err {
	return NoCGO
}

func (h handle) ResetDevice() Err {
	return NoCGO
}

func (h handle) GetDeviceInfo() (uint32, uint16, uint16, Err) {
	return unknown, 0, 0, NoCGO
}

func (h handle) EEPROMRead(devType uint32, ee *EEPROM) Err {
	return NoCGO
}

func (h handle) EEPROMProgram(e *EEPROM) Err {
	return NoCGO
}

func (h handle) EraseEE() Err {
	return NoCGO
}

func (h handle) WriteEE(offset uint8, value uint16) Err {
	return NoCGO
}

func (h handle) EEUASize() (int, Err) {
	return 0, NoCGO
}

func (h handle) EEUARead(ua []byte) Err {
	return NoCGO
}

func (h handle) EEUAWrite(ua []byte) Err {
	return NoCGO
}

func (h handle) SetChars(eventChar byte, eventEn bool, errorChar byte, errorEn bool) Err {
	return NoCGO
}

func (h handle) SetUSBParameters(in, out int) Err {
	return NoCGO
}

func (h handle) SetFlowControl() Err {
	return NoCGO
}

func (h handle) SetTimeouts(readMS, writeMS int) Err {
	return NoCGO
}

func (h handle) SetLatencyTimer(delayMS uint8) Err {
	return NoCGO
}

func (h handle) SetBaudRate(hz uint32) Err {
	return NoCGO
}

func (h handle) GetQueueStatus() (uint32, Err) {
	return 0, NoCGO
}

func (h handle) Read(b []byte) (int, Err) {
	return 0, NoCGO
}

func (h handle) Write(b []byte) (int, Err) {
	return 0, NoCGO
}

func (h handle) GetBitMode() (byte, Err) {
	return 0, NoCGO
}

func (h handle) SetBitMode(mask, mode byte) Err {
	return NoCGO
}
