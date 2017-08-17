package main

import (
    "os"
    "github.com/simonassank/go_ws2811"
    "strconv"
)

func main() {
    color, _ := strconv.ParseUint(os.Args[1], 0, 64)

    ws2811.Init(18, 144, 255)
    
    for i := 0; i < 144; i++ {
        ws2811.SetLed(i, uint32(color))
    }

    ws2811.Render()
}