// Copyright 2021 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package d2xxtest defines logging wrapper and fake for unit testing.
package d2xxtest

import (
	"time"

	"github.com/UoS-PPEMD4ALL/d2xx-wrapper"
)

// Fake implements a fake d2xx.Handle.
type Fake struct {
	DevType uint32
	Vid     uint16
	Pid     uint16
	Data    [][]byte
	UA      []byte
	E       d2xx.EEPROM
}

// Close implements d2xx.Handle.
func (f *Fake) Close() d2xx.Err {
	return 0
}

// ResetDevice implements d2xx.Handle.
func (f *Fake) ResetDevice() d2xx.Err {
	return 0
}

// GetDeviceInfo implements d2xx.Handle.
func (f *Fake) GetDeviceInfo() (uint32, uint16, uint16, d2xx.Err) {
	return f.DevType, f.Vid, f.Pid, 0
}

// EEPROMRead implements d2xx.Handle.
func (f *Fake) EEPROMRead(devType uint32, e *d2xx.EEPROM) d2xx.Err {
	*e = f.E
	return 0
}

// EEPROMProgram implements d2xx.Handle.
func (f *Fake) EEPROMProgram(e *d2xx.EEPROM) d2xx.Err {
	f.E = *e
	return 0
}

// EraseEE implements d2xx.Handle.
func (f *Fake) EraseEE() d2xx.Err {
	return 0
}

// WriteEE implements d2xx.Handle.
func (f *Fake) WriteEE(offset uint8, value uint16) d2xx.Err {
	return 1
}

// EEUASize implements d2xx.Handle.
func (f *Fake) EEUASize() (int, d2xx.Err) {
	return len(f.UA), 0
}

// EEUARead implements d2xx.Handle.
func (f *Fake) EEUARead(UA []byte) d2xx.Err {
	copy(UA, f.UA)
	return 0
}

// EEUAWrite implements d2xx.Handle.
func (f *Fake) EEUAWrite(ua []byte) d2xx.Err {
	f.UA = make([]byte, len(ua))
	copy(f.UA, ua)
	return 0
}

// SetChars implements d2xx.Handle.
func (f *Fake) SetChars(eventChar byte, eventEn bool, errorChar byte, errorEn bool) d2xx.Err {
	return 0
}

// SetUSBParameters implements d2xx.Handle.
func (f *Fake) SetUSBParameters(in, out int) d2xx.Err {
	return 0
}

// SetFlowControl implements d2xx.Handle.
func (f *Fake) SetFlowControl() d2xx.Err {
	return 0
}

// SetTimeouts implements d2xx.Handle.
func (f *Fake) SetTimeouts(readMS, writeMS int) d2xx.Err {
	return 0
}

// SetLatencyTimer implements d2xx.Handle.
func (f *Fake) SetLatencyTimer(delayMS uint8) d2xx.Err {
	return 0
}

// SetBaudRate implements d2xx.Handle.
func (f *Fake) SetBaudRate(hz uint32) d2xx.Err {
	return 0
}

// GetQueueStatus implements d2xx.Handle.
func (f *Fake) GetQueueStatus() (uint32, d2xx.Err) {
	if len(f.Data) == 0 {
		return 0, 0
	}
	// This is to work around flushPending().
	l := len(f.Data[0])
	if l == 0 {
		f.Data = f.Data[1:]
	}
	return uint32(l), 0
}

// Read implements d2xx.Handle.
func (f *Fake) Read(b []byte) (int, d2xx.Err) {
	if len(f.Data) == 0 {
		return 0, 0
	}
	l := len(b)
	if j := len(f.Data[0]); j < l {
		l = j
	}
	if l == 0 {
		f.Data = f.Data[1:]
		return 0, 0
	}
	copy(b, f.Data[0])
	f.Data[0] = f.Data[0][l:]
	if len(f.Data[0]) == 0 {
		f.Data = f.Data[1:]
	}
	return l, 0
}

// Write implements d2xx.Handle.
func (f *Fake) Write(b []byte) (int, d2xx.Err) {
	return 0, 0
}

// GetBitMode implements d2xx.Handle.
func (f *Fake) GetBitMode() (byte, d2xx.Err) {
	return 0, 0
}

// SetBitMode implements d2xx.Handle.
func (f *Fake) SetBitMode(mask, mode byte) d2xx.Err {
	return 0
}

// Log adds logging to a d2xx.Handle to help diagnose issues with the d2xx
// driver.
type Log struct {
	H      d2xx.Handle
	Printf func(format string, v ...interface{})
}

// Close implements d2xx.Handle.
func (l *Log) Close() d2xx.Err {
	defer l.logDefer("Close()")()
	return l.H.Close()
}

// ResetDevice implements d2xx.Handle.
func (l *Log) ResetDevice() d2xx.Err {
	defer l.logDefer("ResetDevice()")()
	return l.H.ResetDevice()
}

// GetDeviceInfo implements d2xx.Handle.
func (l *Log) GetDeviceInfo() (uint32, uint16, uint16, d2xx.Err) {
	defer l.logDefer("GetDeviceInfo()")()
	return l.H.GetDeviceInfo()
}

// EEPROMRead implements d2xx.Handle.
func (l *Log) EEPROMRead(devType uint32, e *d2xx.EEPROM) d2xx.Err {
	defer l.logDefer("EEPROMRead(%d, %d bytes)")(devType, len(e.Raw))
	return l.H.EEPROMRead(devType, e)
}

// EEPROMProgram implements d2xx.Handle.
func (l *Log) EEPROMProgram(e *d2xx.EEPROM) d2xx.Err {
	defer l.logDefer("EEPROMProgram(%#x)")(e)
	return l.H.EEPROMProgram(e)
}

// EraseEE implements d2xx.Handle.
func (l *Log) EraseEE() d2xx.Err {
	defer l.logDefer("EraseEE()")()
	return l.H.EraseEE()
}

// WriteEE implements d2xx.Handle.
func (l *Log) WriteEE(offset uint8, value uint16) d2xx.Err {
	defer l.logDefer("WriteEE(%d, %d)")(offset, value)
	return l.H.WriteEE(offset, value)
}

// EEUASize implements d2xx.Handle.
func (l *Log) EEUASize() (int, d2xx.Err) {
	defer l.logDefer("EEUASize()")()
	return l.H.EEUASize()
}

// EEUARead implements d2xx.Handle.
func (l *Log) EEUARead(ua []byte) d2xx.Err {
	defer l.logDefer("EEUARead(%d bytes)")(len(ua))
	return l.H.EEUARead(ua)
}

// EEUAWrite implements d2xx.Handle.
func (l *Log) EEUAWrite(ua []byte) d2xx.Err {
	defer l.logDefer("EEUAWrite(%#x)")(ua)
	return l.H.EEUAWrite(ua)
}

// SetChars implements d2xx.Handle.
func (l *Log) SetChars(eventChar byte, eventEn bool, errorChar byte, errorEn bool) d2xx.Err {
	defer l.logDefer("SetChars(%d, %t, %d, %t)")(eventChar, eventEn, errorChar, errorEn)
	return l.H.SetChars(eventChar, eventEn, errorChar, errorEn)
}

// SetUSBParameters implements d2xx.Handle.
func (l *Log) SetUSBParameters(in, out int) d2xx.Err {
	defer l.logDefer("SetUSBParameters(%d, %d)")(in, out)
	return l.H.SetUSBParameters(in, out)
}

// SetFlowControl implements d2xx.Handle.
func (l *Log) SetFlowControl() d2xx.Err {
	defer l.logDefer("SetFlowControl()")()
	return l.H.SetFlowControl()
}

// SetTimeouts implements d2xx.Handle.
func (l *Log) SetTimeouts(readMS, writeMS int) d2xx.Err {
	defer l.logDefer("SetTimeouts(%d, %d)")(readMS, writeMS)
	return l.H.SetTimeouts(readMS, writeMS)
}

// SetLatencyTimer implements d2xx.Handle.
func (l *Log) SetLatencyTimer(delayMS uint8) d2xx.Err {
	defer l.logDefer("SetLatencyTimer(%d)")(delayMS)
	return l.H.SetLatencyTimer(delayMS)
}

// SetBaudRate implements d2xx.Handle.
func (l *Log) SetBaudRate(hz uint32) d2xx.Err {
	defer l.logDefer("SetBaudRate(%d)")(hz)
	return l.H.SetBaudRate(hz)
}

// GetQueueStatus implements d2xx.Handle.
func (l *Log) GetQueueStatus() (uint32, d2xx.Err) {
	f := l.logDefer("GetQueueStatus() = %d, %d")
	p, e := l.H.GetQueueStatus()
	f(p, e)
	return p, e
}

// Read implements d2xx.Handle.
func (l *Log) Read(b []byte) (int, d2xx.Err) {
	f := l.logDefer("Read(%d bytes) = %#x")
	n, e := l.H.Read(b)
	f(len(b), b[:n])
	return n, e
}

// Write implements d2xx.Handle.
func (l *Log) Write(b []byte) (int, d2xx.Err) {
	defer l.logDefer("Write(%#x)")(b)
	return l.H.Write(b)
}

// GetBitMode implements d2xx.Handle.
func (l *Log) GetBitMode() (byte, d2xx.Err) {
	f := l.logDefer("GetBitMode() = %02X")
	b, e := l.H.GetBitMode()
	f(b)
	return b, e
}

// SetBitMode implements d2xx.Handle.
func (l *Log) SetBitMode(mask, mode byte) d2xx.Err {
	f := l.logDefer("SetBitMode(0x%02X, 0x%02X) = %d")
	e := l.H.SetBitMode(mask, mode)
	f(mask, mode, e)
	return e
}

func (l *Log) logDefer(fmt string) func(args ...interface{}) {
	var start time.Time
	f := func(args ...interface{}) {
		l.Printf("%7s "+fmt, append([]interface{}{roundDuration(time.Since(start))}, args...)...)
	}
	start = time.Now()
	return f
}

//

// log10 is a cheap way to find the most significant digit
func log10(i int64) uint {
	switch {
	case i < 10:
		return 0
	case i < 100:
		return 1
	case i < 1000:
		return 2
	case i < 10000:
		return 3
	case i < 100000:
		return 4
	case i < 1000000:
		return 5
	case i < 10000000:
		return 6
	case i < 100000000:
		return 7
	case i < 1000000000:
		return 8
	case i < 10000000000:
		return 9
	case i < 100000000000:
		return 10
	case i < 1000000000000:
		return 11
	case i < 10000000000000:
		return 12
	case i < 100000000000000:
		return 13
	case i < 1000000000000000:
		return 14
	case i < 10000000000000000:
		return 15
	default:
		return 16
	}
}

func roundDuration(d time.Duration) time.Duration {
	if l := log10(int64(d)); l > 3 {
		m := time.Duration(1)
		for i := uint(3); i < l; i++ {
			m *= 10
		}
		d = (d + (m / 2)) / m * m
	}
	return d
}

var _ d2xx.Handle = &Fake{}
var _ d2xx.Handle = &Log{}
