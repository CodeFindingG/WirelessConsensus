package main

import (
	"strconv"
	"time"
)

// 全局更新的时隙
var timeSlot = 1

var slotDuration = time.Duration(40) * time.Millisecond

func timeGenerate() {
	for {
		slotChange <- "Slot" + strconv.Itoa(timeSlot)
		ifSlotChange = true // break condition
		time.Sleep(slotDuration)
		timeSlot++
	}
	select {}
}
