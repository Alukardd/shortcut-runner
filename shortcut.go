package main

import (
	"fmt"
	"github.com/gvalkov/golang-evdev"
)

const (
	device_glob = "/dev/input/event*"
)

func waitShortcut(inputdev string, shortcut []string) (err error) {
	var dev *evdev.InputDevice
	var events []evdev.InputEvent

	dev, err = evdev.Open(inputdev)
	if err != nil {
		l.Err(fmt.Sprintf("%s", err))
		return
	}

	l.Info(fmt.Sprintf("%s", dev))
	l.Info(fmt.Sprintf("Listening for events ...\n"))

	// Numerical equivalent of the string value key
	var ishortcut []int
	for key := range shortcut {
		ishortcut = append(ishortcut, keymaps[shortcut[key]])
	}

	key_sum := 0
	waiting_key_sum := 0
	for key := range ishortcut {
		waiting_key_sum += ishortcut[key]
	}

	for {
		events, err = dev.Read()
		for i := range events {
			ecode := int(events[i].Code)
			if events[i].Type == evdev.EV_KEY {
				for key := range ishortcut {
					if ecode == ishortcut[key] {
						switch events[i].Value {
						case 0:
							key_sum -= ecode
						case 1:
							key_sum += ecode
						}
					}
				}
			}
		}
		if key_sum == waiting_key_sum {
			return
		}
	}
}
