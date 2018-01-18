package ws2811

/*
#cgo CFLAGS: -std=c99 -Ivendor/rpi_ws281x
#cgo LDFLAGS: -lws2811 -Lvendor/rpi_ws281x
#include <stdint.h>
#include <string.h>
#include <ws2811.h>
void ws2811_set_led(ws2811_t *ws2811, int ch, int index, uint32_t value) {
  ws2811->channel[ch].leds[index] = value;
}

uint32_t ws2811_get_led(ws2811_t *ws2811, int ch, int index) {
    return ws2811->channel[ch].leds[index];
}
*/
import "C"

import (
  "fmt"
  "unsafe"
)

// StripType layout
type StripType int

const (
  // 4 color R, G, B and W ordering
  StripRGBW StripType = 0x18100800
  StripRBGW StripType = 0x18100008
  StripGRBW StripType = 0x18081000
  StripGBRW StripType = 0x18080010
  StripBRGW StripType = 0x18001008
  StripBGRW StripType = 0x18000810

  // 3 color R, G and B ordering
  StripRGB StripType = 0x00100800
  StripRBG StripType = 0x00100008
  StripGRB StripType = 0x00081000
  StripGBR StripType = 0x00080010
  StripBRG StripType = 0x00001008
  StripBGR StripType = 0x00000810
)

// DefaultConfig default WS281x configuration
var DefaultConfig = HardwareConfig{
  Pin:        18,
  Frequency:  800000, // 800khz
  DMA:        10,
  Brightness: 30,
  StripType:  StripGRB,
}

// HardwareConfig WS281x configuration
type HardwareConfig struct {
  Pin        int       // GPIO Pin with PWM alternate function, 0 if unused
  Frequency  int       // the frequency of the display signal in hertz, can go as low as 400000
  DMA        int       // the DMA channel to use
  Invert     bool      // specifying if the signal line should be inverted
  Channel    int       // PWM channel to us
  Brightness int       // brightness value between 0 and 255
  StripType  StripType // strip color layout
}

// WS281x matrix representation for ws281x
type WS281x struct {
  Config *HardwareConfig

  size   int
  leds   *C.ws2811_t
  closed bool
}

// NewWS281x returns a new matrix using the given size and config
func New(size int, config *HardwareConfig) (*WS281x, error) {
  c := &WS281x{
    Config: config,

    size: size,
    leds: (*C.ws2811_t)(C.malloc(C.sizeof_ws2811_t)),
  }

  if c.leds == nil {
    return nil, fmt.Errorf("unable to allocate memory")
  }

  C.memset(unsafe.Pointer(c.leds), 0, C.sizeof_ws2811_t)

  c.initializeChannels()
  c.initializeController()
  return c, nil
}

func (c *WS281x) initializeChannels() {
  for ch := 0; ch < 2; ch++ {
    c.setChannel(0, 0, 0, 0, false, StripRGB)
  }

  c.setChannel(
    c.Config.Channel,
    c.size,
    c.Config.Pin,
    c.Config.Brightness,
    c.Config.Invert,
    c.Config.StripType,
  )
}

func (c *WS281x) setChannel(ch, count, pin, brightness int, inverse bool, t StripType) {
  c.leds.channel[ch].count = C.int(count)
  c.leds.channel[ch].gpionum = C.int(pin)
  c.leds.channel[ch].brightness = C.uint8_t(brightness)
  c.leds.channel[ch].invert = C.int(btoi(inverse))
  c.leds.channel[ch].strip_type = C.int(int(t))
}

func (c *WS281x) initializeController() {
  c.leds.freq = C.uint32_t(c.Config.Frequency)
  c.leds.dmanum = C.int(c.Config.DMA)
}

// Initialize initialize library, must be called once before other functions are
// called.
func (c *WS281x) Init() error {
  if resp := int(C.ws2811_init(c.leds)); resp != 0 {
    return fmt.Errorf("ws2811_init failed with code: %d", resp)
  }

  return nil
}

// Render update the display with the data from the LED buffer
func (c *WS281x) Render() error {
  if resp := int(C.ws2811_render(c.leds)); resp != 0 {
    return fmt.Errorf("ws2811_render failed with code: %d", resp)
  }

  return nil
}

func (c *WS281x) GetLed(position int) uint32 {
  color := C.ws2811_get_led(c.leds, C.int(c.Config.Channel), C.int(position))
  return uint32(color)
}

func (c *WS281x) SetLed(position int, color uint32) {
  C.ws2811_set_led(c.leds, C.int(c.Config.Channel), C.int(position), C.uint32_t(color))
}

// Close finalizes the ws281x interface
func (c *WS281x) Close() error {
  if c.closed {
    return nil
  }

  c.closed = true
  C.ws2811_fini(c.leds)
  return nil
}

func btoi(b bool) int {
  if b {
    return 1
  }
  return 0
}

