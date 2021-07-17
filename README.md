# d2xx

Package d2xx is a thin Go wrapper for the Future Technology "D2XX" driver at
https://ftdichip.com/drivers/d2xx-drivers/.

This package is not Go idiomatic. You will want to use
https://periph.io/x/host/v3/ftdi (or later) instead.

But if you really want, here it goes:
[![PkgGoDev](https://pkg.go.dev/badge/periph.io/x/d2xx)](https://pkg.go.dev/periph.io/x/d2xx)

This Go package includes third party software. See
[third_party/README.md](third_party/README.md).

## Configuration

See https://periph.io/device/ftdi/ to configure the host to be able to use this
driver.

## Availability

On darwin_amd64, linux_amd64 linux_arm (v6, v7 compatible) and linux_arm64 (v8),
cgo is required. If cgo is disabled (via `CGO_ENABLED=0`), all functions in this
driver return error [NoCGO](https://periph.io/x/d2xx#NoCGO).

On Windows, cgo is not required. If the dynamic library is not found at runtime,
[Missing](https://periph.io/x/d2xx#Missing) is returned.

## bcm2385

On linux_arm (v6), hard-float is required. For cross compilation, this
means arm-linux-gnueabihf-gcc is preferred to arm-linux-gnueabi-gcc. Using
hardfloat causes a segfault on Raspberry Pi 1, Zero and Zero Wireless. It is
recommended to disable this driver if targeting these hosts, see below.

## Disabling

To disable this driver, build with tag `no_d2xx`, e.g.

```
go install -tags no_d2xx periph.io/x/cmd/gpio-list@latest
```

This will behave has if cgo was disabled, even on Windows.
