// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !no_d2xx

package d2xx

import (
	"bytes"
	"syscall"
	"unsafe"
)

// Available is true if the library is available on this system.
var Available = false

func version() (uint8, uint8, uint8) {
	var v uint32
	if pGetLibraryVersion != nil {
		_, _, _ = pGetLibraryVersion.Call(uintptr(unsafe.Pointer(&v)))
	}
	return uint8(v >> 16), uint8(v >> 8), uint8(v)
}

func createDeviceInfoList() (int, Err) {
	if pCreateDeviceInfoList != nil {
		var num uint32
		r1, _, _ := pCreateDeviceInfoList.Call(uintptr(unsafe.Pointer(&num)))
		return int(num), Err(r1)
	}
	return 0, Missing
}

func open(i int) (Handle, Err) {
	var h handle
	if pOpen != nil {
		r1, _, _ := pOpen.Call(uintptr(i), uintptr(unsafe.Pointer(&h)))
		return h, Err(r1)
	}
	return h, Missing
}

func (h handle) Close() Err {
	r1, _, _ := pClose.Call(h.toH())
	return Err(r1)
}

func (h handle) ResetDevice() Err {
	r1, _, _ := pResetDevice.Call(h.toH())
	return Err(r1)
}

func (h handle) GetDeviceInfo() (uint32, uint16, uint16, Err) {
	var d uint32
	var id uint32
	if r1, _, _ := pGetDeviceInfo.Call(h.toH(), uintptr(unsafe.Pointer(&d)), uintptr(unsafe.Pointer(&id)), 0, 0, 0); r1 != 0 {
		return unknown, 0, 0, Err(r1)
	}
	return d, uint16(id >> 16), uint16(id), 0
}

func (h handle) EEPROMRead(devType uint32, ee *EEPROM) Err {
	var manufacturer [64]byte
	var manufacturerID [64]byte
	var desc [64]byte
	var serial [64]byte
	// Shortcuts.
	m := uintptr(unsafe.Pointer(&manufacturer[0]))
	mi := uintptr(unsafe.Pointer(&manufacturerID[0]))
	de := uintptr(unsafe.Pointer(&desc[0]))
	s := uintptr(unsafe.Pointer(&serial[0]))

	eepromVoid := unsafe.Pointer(&ee.Raw[0])
	hdr := ee.asHeader()
	// It MUST be set here. This is not always the case on posix.
	hdr.DeviceType = devType
	if r1, _, _ := pEEPROMRead.Call(h.toH(), uintptr(eepromVoid), uintptr(len(ee.Raw)), m, mi, de, s); r1 != 0 {
		return Err(r1)
	}

	ee.Manufacturer = toStr(manufacturer[:])
	ee.ManufacturerID = toStr(manufacturerID[:])
	ee.Desc = toStr(desc[:])
	ee.Serial = toStr(serial[:])
	return 0
}

func (h handle) EEPROMProgram(ee *EEPROM) Err {
	var cmanu [64]byte
	copy(cmanu[:], ee.Manufacturer)
	var cmanuID [64]byte
	copy(cmanuID[:], ee.ManufacturerID)
	var cdesc [64]byte
	copy(cdesc[:], ee.Desc)
	var cserial [64]byte
	copy(cserial[:], ee.Serial)
	r1, _, _ := pEEPROMProgram.Call(h.toH(), uintptr(unsafe.Pointer(&ee.Raw[0])), uintptr(len(ee.Raw)), uintptr(unsafe.Pointer(&cmanu[0])), uintptr(unsafe.Pointer(&cmanuID[0])), uintptr(unsafe.Pointer(&cdesc[0])), uintptr(unsafe.Pointer(&cserial[0])))
	return Err(r1)
}

func (h handle) EraseEE() Err {
	r1, _, _ := pEraseEE.Call(h.toH())
	return Err(r1)
}

func (h handle) WriteEE(offset uint8, value uint16) Err {
	r1, _, _ := pWriteEE.Call(h.toH(), uintptr(offset), uintptr(value))
	return Err(r1)
}

func (h handle) EEUASize() (int, Err) {
	var size uint32
	if r1, _, _ := pEEUASize.Call(h.toH(), uintptr(unsafe.Pointer(&size))); r1 != 0 {
		return 0, Err(r1)
	}
	return int(size), 0
}

func (h handle) EEUARead(ua []byte) Err {
	var size uint32
	if r1, _, _ := pEEUARead.Call(h.toH(), uintptr(unsafe.Pointer(&ua[0])), uintptr(len(ua)), uintptr(unsafe.Pointer(&size))); r1 != 0 {
		return Err(r1)
	}
	if int(size) != len(ua) {
		return 6 // FT_INVALID_PARAMETER
	}
	return 0
}

func (h handle) EEUAWrite(ua []byte) Err {
	r1, _, _ := pEEUAWrite.Call(h.toH(), uintptr(unsafe.Pointer(&ua[0])), uintptr(len(ua)))
	return Err(r1)
}

func (h handle) SetChars(eventChar byte, eventEn bool, errorChar byte, errorEn bool) Err {
	v := uintptr(0)
	if eventEn {
		v = 1
	}
	w := uintptr(0)
	if errorEn {
		w = 1
	}
	r1, _, _ := pSetChars.Call(h.toH(), uintptr(eventChar), v, uintptr(errorChar), w)
	return Err(r1)
}

func (h handle) SetUSBParameters(in, out int) Err {
	r1, _, _ := pSetUSBParameters.Call(h.toH(), uintptr(in), uintptr(out))
	return Err(r1)
}

func (h handle) SetFlowControl() Err {
	// FT_FLOW_RTS_CTS
	r1, _, _ := pSetFlowControl.Call(h.toH(), 0x0100, 0, 0)
	return Err(r1)
}

func (h handle) SetTimeouts(readMS, writeMS int) Err {
	r1, _, _ := pSetTimeouts.Call(h.toH(), uintptr(readMS), uintptr(writeMS))
	return Err(r1)
}

func (h handle) SetLatencyTimer(delayMS uint8) Err {
	r1, _, _ := pSetLatencyTimer.Call(h.toH(), uintptr(delayMS))
	return Err(r1)
}

func (h handle) SetBaudRate(hz uint32) Err {
	r1, _, _ := pSetBaudRate.Call(h.toH(), uintptr(hz))
	return Err(r1)
}

func (h handle) GetQueueStatus() (uint32, Err) {
	var v uint32
	r1, _, _ := pGetQueueStatus.Call(h.toH(), uintptr(unsafe.Pointer(&v)))
	return v, Err(r1)
}

func (h handle) Read(b []byte) (int, Err) {
	var bytesRead uint32
	r1, _, _ := pRead.Call(h.toH(), uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), uintptr(unsafe.Pointer(&bytesRead)))
	return int(bytesRead), Err(r1)
}

func (h handle) Write(b []byte) (int, Err) {
	var bytesSent uint32
	r1, _, _ := pWrite.Call(h.toH(), uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), uintptr(unsafe.Pointer(&bytesSent)))
	return int(bytesSent), Err(r1)
}

func (h handle) GetBitMode() (byte, Err) {
	var s uint8
	r1, _, _ := pGetBitMode.Call(h.toH(), uintptr(unsafe.Pointer(&s)))
	return s, Err(r1)
}

func (h handle) SetBitMode(mask, mode byte) Err {
	r1, _, _ := pSetBitMode.Call(h.toH(), uintptr(mask), uintptr(mode))
	return Err(r1)
}

func (h handle) toH() uintptr {
	return uintptr(h)
}

//

var (
	pClose                *syscall.Proc
	pCreateDeviceInfoList *syscall.Proc
	pEEPROMRead           *syscall.Proc
	pEEPROMProgram        *syscall.Proc
	pEraseEE              *syscall.Proc
	pWriteEE              *syscall.Proc
	pEEUASize             *syscall.Proc
	pEEUARead             *syscall.Proc
	pEEUAWrite            *syscall.Proc
	pGetBitMode           *syscall.Proc
	pGetDeviceInfo        *syscall.Proc
	pGetLibraryVersion    *syscall.Proc
	pGetQueueStatus       *syscall.Proc
	pOpen                 *syscall.Proc
	pRead                 *syscall.Proc
	pResetDevice          *syscall.Proc
	pSetBaudRate          *syscall.Proc
	pSetBitMode           *syscall.Proc
	pSetChars             *syscall.Proc
	pSetFlowControl       *syscall.Proc
	pSetLatencyTimer      *syscall.Proc
	pSetTimeouts          *syscall.Proc
	pSetUSBParameters     *syscall.Proc
	pWrite                *syscall.Proc
)

func init() {
	if dll, _ := syscall.LoadDLL("ftd2xx.dll"); dll != nil {
		// If any function is not found, disable the support.
		Available = true
		find := func(n string) *syscall.Proc {
			s, _ := dll.FindProc(n)
			if s == nil {
				Available = false
			}
			return s
		}
		pClose = find("FT_Close")
		pCreateDeviceInfoList = find("FT_CreateDeviceInfoList")
		pEEPROMRead = find("FT_EEPROM_Read")
		pEEPROMProgram = find("FT_EEPROM_Program")
		pEraseEE = find("FT_EraseEE")
		pWriteEE = find("FT_WriteEE")
		pEEUASize = find("FT_EE_UASize")
		pEEUARead = find("FT_EE_UARead")
		pEEUAWrite = find("FT_EE_UAWrite")
		pGetBitMode = find("FT_GetBitMode")
		pGetDeviceInfo = find("FT_GetDeviceInfo")
		pGetLibraryVersion = find("FT_GetLibraryVersion")
		pGetQueueStatus = find("FT_GetQueueStatus")
		pOpen = find("FT_Open")
		pRead = find("FT_Read")
		pResetDevice = find("FT_ResetDevice")
		pSetBaudRate = find("FT_SetBaudRate")
		pSetBitMode = find("FT_SetBitMode")
		pSetChars = find("FT_SetChars")
		pSetFlowControl = find("FT_SetFlowControl")
		pSetLatencyTimer = find("FT_SetLatencyTimer")
		pSetTimeouts = find("FT_SetTimeouts")
		pSetUSBParameters = find("FT_SetUSBParameters")
		pWrite = find("FT_Write")
	}
}

func toStr(c []byte) string {
	i := bytes.IndexByte(c, 0)
	if i != -1 {
		return string(c[:i])
	}
	return string(c)
}
