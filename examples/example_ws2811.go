package main

import (
    "os"
    "github.com/simonassank/go_ws2811"
    "strconv"
    _ "fmt"
)

var (
led_cnt = 144*4
)

func main() {
    color, _ := strconv.ParseUint(os.Args[1], 0, 64)
    correction, _ := strconv.ParseUint(os.Args[2], 10, 64)
    isCorrectionEnabled := correction > 0
    ws2811.Init(18, led_cnt, 255)

    for i := 0; i < led_cnt; i++ {
        ws2811.SetLed(i, uint32(color), isCorrectionEnabled)
    }

    ws2811.Render()
}
