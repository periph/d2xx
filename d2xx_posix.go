// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build cgo
// +build !windows
// +build !no_d2xx

package d2xx

/*
#include "third_party/ftd2xx.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// Available is true if the library is available on this system.
const Available = true

// Library functions.

func version() (uint8, uint8, uint8) {
	var v C.DWORD
	C.FT_GetLibraryVersion(&v)
	return uint8(v >> 16), uint8(v >> 8), uint8(v)
}

func createDeviceInfoList() (int, Err) {
	var num C.DWORD
	e := C.FT_CreateDeviceInfoList(&num)
	return int(num), Err(e)
}

func open(i int) (Handle, Err) {
	var h C.FT_HANDLE
	e := C.FT_Open(C.int(i), &h)
	if uintptr(h) == 0 && e == 0 {
		// 18 means FT_OTHER_ERROR. Kind of a hack but better than panic.
		e = 18
	}
	return handle(h), Err(e)
}

func (h handle) Close() Err {
	return Err(C.FT_Close(h.toH()))
}

func (h handle) ResetDevice() Err {
	return Err(C.FT_ResetDevice(h.toH()))
}

func (h handle) GetDeviceInfo() (uint32, uint16, uint16, Err) {
	var dev C.FT_DEVICE
	var id C.DWORD
	if e := C.FT_GetDeviceInfo(h.toH(), &dev, &id, nil, nil, nil); e != 0 {
		return unknown, 0, 0, Err(e)
	}
	return uint32(dev), uint16(id >> 16), uint16(id), 0
}

func (h handle) EEPROMRead(devType uint32, ee *EEPROM) Err {
	var manufacturer [64]C.char
	var manufacturerID [64]C.char
	var desc [64]C.char
	var serial [64]C.char
	eepromVoid := unsafe.Pointer(&ee.Raw[0])
	hdr := ee.asHeader()

	// There something odd going on here.
	//
	// On a ft232h, we observed that hdr.DeviceType MUST NOT be set, but on a
	// ft232r, it MUST be set. Since we can't know in advance what we must use,
	// just try both. ¯\_(ツ)_/¯
	hdr.DeviceType = devType
	if e := C.FT_EEPROM_Read(h.toH(), eepromVoid, C.DWORD(len(ee.Raw)), &manufacturer[0], &manufacturerID[0], &desc[0], &serial[0]); e != 0 {
		// FT_INVALID_PARAMETER
		if e == 6 {
			hdr.DeviceType = 0
			e = C.FT_EEPROM_Read(h.toH(), eepromVoid, C.DWORD(len(ee.Raw)), &manufacturer[0], &manufacturerID[0], &desc[0], &serial[0])
		}
		if e != 0 {
			return Err(e)
		}
	}

	ee.Manufacturer = C.GoString(&manufacturer[0])
	ee.ManufacturerID = C.GoString(&manufacturerID[0])
	ee.Desc = C.GoString(&desc[0])
	ee.Serial = C.GoString(&serial[0])
	return 0
}

func (h handle) EEPROMProgram(ee *EEPROM) Err {
	// len(manufacturer) + len(desc) <= 40.
	/*
		var cmanu [64]byte
		copy(cmanu[:], ee.manufacturer)
		var cmanuID [64]byte
		copy(cmanuID[:], ee.manufacturerID)
		var cdesc [64]byte
		copy(cdesc[:], ee.desc)
		var cserial [64]byte
		copy(cserial[:], ee.serial)
	*/
	cmanu := C.CString(ee.Manufacturer)
	defer C.free(unsafe.Pointer(cmanu))
	cmanuID := C.CString(ee.ManufacturerID)
	defer C.free(unsafe.Pointer(cmanuID))
	cdesc := C.CString(ee.Desc)
	defer C.free(unsafe.Pointer(cdesc))
	cserial := C.CString(ee.Serial)
	defer C.free(unsafe.Pointer(cserial))

	if len(ee.Raw) == 0 {
		return Err(C.FT_EEPROM_Program(h.toH(), unsafe.Pointer(uintptr(0)), 0, cmanu, cmanuID, cdesc, cserial))
	}
	return Err(C.FT_EEPROM_Program(h.toH(), unsafe.Pointer(&ee.Raw[0]), C.DWORD(len(ee.Raw)), cmanu, cmanuID, cdesc, cserial))
}

func (h handle) EraseEE() Err {
	return Err(C.FT_EraseEE(h.toH()))
}

func (h handle) WriteEE(offset uint8, value uint16) Err {
	return Err(C.FT_WriteEE(h.toH(), C.DWORD(offset), C.WORD(value)))
}

func (h handle) EEUASize() (int, Err) {
	var size C.DWORD
	if e := C.FT_EE_UASize(h.toH(), &size); e != 0 {
		return 0, Err(e)
	}
	return int(size), 0
}

func (h handle) EEUARead(ua []byte) Err {
	var size C.DWORD
	if e := C.FT_EE_UARead(h.toH(), (*C.UCHAR)(unsafe.Pointer(&ua[0])), C.DWORD(len(ua)), &size); e != 0 {
		return Err(e)
	}
	if int(size) != len(ua) {
		return 6 // FT_INVALID_PARAMETER
	}
	return 0
}

func (h handle) EEUAWrite(ua []byte) Err {
	if e := C.FT_EE_UAWrite(h.toH(), (*C.UCHAR)(&ua[0]), C.DWORD(len(ua))); e != 0 {
		return Err(e)
	}
	return 0
}

func (h handle) SetChars(eventChar byte, eventEn bool, errorChar byte, errorEn bool) Err {
	v := C.UCHAR(0)
	if eventEn {
		v = 1
	}
	w := C.UCHAR(0)
	if errorEn {
		w = 1
	}
	return Err(C.FT_SetChars(h.toH(), C.UCHAR(eventChar), v, C.UCHAR(errorChar), w))
}

func (h handle) SetUSBParameters(in, out int) Err {
	return Err(C.FT_SetUSBParameters(h.toH(), C.DWORD(in), C.DWORD(out)))
}

func (h handle) SetFlowControl() Err {
	return Err(C.FT_SetFlowControl(h.toH(), C.FT_FLOW_RTS_CTS, 0, 0))
}

func (h handle) SetTimeouts(readMS, writeMS int) Err {
	return Err(C.FT_SetTimeouts(h.toH(), C.DWORD(readMS), C.DWORD(writeMS)))
}

func (h handle) SetLatencyTimer(delayMS uint8) Err {
	return Err(C.FT_SetLatencyTimer(h.toH(), C.UCHAR(delayMS)))
}

func (h handle) SetBaudRate(hz uint32) Err {
	return Err(C.FT_SetBaudRate(h.toH(), C.DWORD(hz)))
}

func (h handle) GetQueueStatus() (uint32, Err) {
	var v C.DWORD
	e := C.FT_GetQueueStatus(h.toH(), &v)
	return uint32(v), Err(e)
}

func (h handle) Read(b []byte) (int, Err) {
	var bytesRead C.DWORD
	e := C.FT_Read(h.toH(), C.LPVOID(unsafe.Pointer(&b[0])), C.DWORD(len(b)), &bytesRead)
	return int(bytesRead), Err(e)
}

func (h handle) Write(b []byte) (int, Err) {
	var bytesSent C.DWORD
	e := C.FT_Write(h.toH(), C.LPVOID(unsafe.Pointer(&b[0])), C.DWORD(len(b)), &bytesSent)
	return int(bytesSent), Err(e)
}

func (h handle) GetBitMode() (byte, Err) {
	var s C.UCHAR
	e := C.FT_GetBitMode(h.toH(), &s)
	return uint8(s), Err(e)
}

func (h handle) SetBitMode(mask, mode byte) Err {
	return Err(C.FT_SetBitMode(h.toH(), C.UCHAR(mask), C.UCHAR(mode)))
}

func (h handle) toH() C.FT_HANDLE {
	return C.FT_HANDLE(h)
}
