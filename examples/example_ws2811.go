package main

import (
    "os"
    "github.com/simonassank/go_ws2811"
    "strconv"
    _ "fmt"
)

func main() {
    color, _ := strconv.ParseUint(os.Args[1], 0, 64)

    config := ws2811.DefaultConfig
    config.Brightness = 255
    config.Pin = 12

    c, _ := ws2811.New(144*4, &config)

    defer c.Close()
    _ = c.Init()

    for i := 0; i < 144*4; i++ {
        c.SetLed(i, uint32(color))
    }

    c.Render()
}
