// Copyright (c) 2015, Jacques Supcik, HEIA-FR
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in the
//       documentation and/or other materials provided with the distribution.
//     * Neither the name of the <organization> nor the
//       names of its contributors may be used to endorse or promote products
//       derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

/*
Interface to ws2811 chip (neopixel driver). Make sure that you have
ws2811.h and pwm.h in a GCC include path (e.g. /usr/local/include) and
libws2811.a in a GCC library path (e.g. /usr/local/lib).
See https://github.com/jgarff/rpi_ws281x for instructions
*/

package ws2811

/*
#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lws2811
#include "ws2811.go.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

var (
    gammaTable = [257]uint32 {
        0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,
        0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  0,  1,  1,  1,  1,
        1,  1,  1,  1,  1,  1,  1,  1,  1,  2,  2,  2,  2,  2,  2,  2,
        2,  3,  3,  3,  3,  3,  3,  3,  4,  4,  4,  4,  4,  5,  5,  5,
        5,  6,  6,  6,  6,  7,  7,  7,  7,  8,  8,  8,  9,  9,  9, 10,
       10, 10, 11, 11, 11, 12, 12, 13, 13, 13, 14, 14, 15, 15, 16, 16,
       17, 17, 18, 18, 19, 19, 20, 20, 21, 21, 22, 22, 23, 24, 24, 25,
       25, 26, 27, 27, 28, 29, 29, 30, 31, 32, 32, 33, 34, 35, 35, 36,
       37, 38, 39, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 50,
       51, 52, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 66, 67, 68,
       69, 70, 72, 73, 74, 75, 77, 78, 79, 81, 82, 83, 85, 86, 87, 89,
       90, 92, 93, 95, 96, 98, 99,101,102,104,105,107,109,110,112,114,
      115,117,119,120,122,124,126,127,129,131,133,135,137,138,140,142,
      144,146,148,150,152,154,156,158,160,162,164,167,169,171,173,175,
      177,180,182,184,186,189,191,193,196,198,200,203,205,208,210,213,
      215,218,220,223,225,228,231,233,236,239,241,244,247,249,252,255 }
)

func Init(gpioPin int, ledCount int, brightness int) error {
	C.ledstring.channel[0].gpionum = C.int(gpioPin)
	C.ledstring.channel[0].count = C.int(ledCount)
	C.ledstring.channel[0].brightness = C.uint8_t(brightness)
	res := int(C.ws2811_init(&C.ledstring))
	if res == 0 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Error ws2811.init.%d", res))
	}
}

func Fini() {
	C.ws2811_fini(&C.ledstring)
}

func Render() error {
	res := int(C.ws2811_render(&C.ledstring))
	if res == 0 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Error ws2811.render.%d", res))
	}
}

func Wait() error {
	res := int(C.ws2811_wait(&C.ledstring))
	if res == 0 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Error ws2811.wait.%d", res))
	}
}

func SetLed(index int, value uint32, correctGamma bool) {
    fixedColor := recalculateColor(value, correctGamma)
	C.ws2811_set_led(&C.ledstring, C.int(index), C.uint32_t(fixedColor))
}

func Clear() {
	C.ws2811_clear(&C.ledstring)
}

func SetBitmap(a []uint32) {
	C.ws2811_set_bitmap(&C.ledstring, unsafe.Pointer(&a[0]), C.int(len(a)*4))
}

func recalculateColor(value uint32, correctGamma bool) uint32 {
    red := value >> 16 & 0xFF
    green := value >> 8 & 0xFF
    blue := value >> 0 & 0xFF
    
    // Gamma is off by default.
    if (correctGamma) {
        red = gammaTable[int(red)]
        green = gammaTable[int(green)]
        blue = gammaTable[int(blue)]
    }

    // Red and Green is mixed up, when Rendering on LED strip. Idk why.
    // This kind of works.
    regenR := red << 8 & 0xFFFFFF
    regenG := green << 16 & 0xFFFFFF
    regenB := blue << 0 & 0xFFFFFF

    return regenR + regenG + regenB
}