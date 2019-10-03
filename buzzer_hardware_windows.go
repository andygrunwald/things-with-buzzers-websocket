// +build windows

package main

import (
	"fmt"
)

// This is only a dummy implementation for windows.
// On windows, gobot seems not to compile.
//
// 	⨯ release failed after 2.41s error=failed to build for windows_amd64: # gobot.io/x/gobot/sysfs
//	../Go/src/gobot.io/x/gobot/sysfs/i2c_device.go:70:3: undefined: syscall.SYS_IOCTL
//	../Go/src/gobot.io/x/gobot/sysfs/i2c_device.go:84:3: undefined: syscall.SYS_IOCTL
//	../Go/src/gobot.io/x/gobot/sysfs/i2c_device.go:201:3: undefined: syscall.SYS_IOCTL
//	../Go/src/gobot.io/x/gobot/sysfs/syscall.go:34:24: not enough arguments in call to syscall.Syscall
//		have (uintptr, uintptr, uintptr, uintptr)
//		want (uintptr, uintptr, uintptr, uintptr, uintptr)
//
// Therefor, we mock the hardware interface and only
// enable other features of twb-websocket on other platforms.

// NewHardwareBuzzer returns a new instance of
// the hardware buzzer stream
func NewHardwareBuzzer(buzzer chan buzzerHit) Buzzer {
	b := &HardwareBuzzer{
		buzzerStream: buzzer,
	}
	return b
}

// Initialize sets up everything that is needed to run the buzzer.
func (b *HardwareBuzzer) Initialize() error {
	err := fmt.Errorf("hardware buzzers are not supported under windows")
	return err
}

// Start boots up the buzzer and ensures that they are ready to hit.
func (b *HardwareBuzzer) Start() error {
	err := fmt.Errorf("hardware buzzers are not supported under windows")
	return err
}
