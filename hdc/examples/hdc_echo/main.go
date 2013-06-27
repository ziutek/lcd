// Pin connection between FT232R and LCD (HD44780 compatible):
// TxD  (DBUS0) <--> D4
// RxD  (DBUS1) <--> D5
// RTS# (DBUS2) <--> D6
// CTS# (DBUS3) <--> D7
// DTR# (DBUS4) <--> E
// DSR# (DBUS5) <--> R/W#
// DCD# (DBUS6) <--> RS
package main

import (
	"flag"
	"fmt"
	"github.com/ziutek/ftdi"
	"github.com/ziutek/lcd/hdc"
	"os"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const (
	baudrate  = 1 << 17
	waitTicks = 6
)

var (
	x = flag.Int("x", 0, "column number")
	y = flag.Int("y", 0, "row number")
	c = flag.Bool("c", false, "clear display")
)

func main() {
	flag.Parse()
	text := strings.Join(flag.Args(), " ")

	d, err := ftdi.OpenFirst(0x0403, 0x6001, ftdi.ChannelAny)
	checkErr(err)
	defer d.Close()
	checkErr(d.SetBitmode(0xff, ftdi.ModeBitbang))
	checkErr(d.SetBaudrate(baudrate / 16))

	lcd := hdc.NewDevice(hdc.NewBitbang(d, waitTicks), 4, 20)
	checkErr(lcd.Init())
	checkErr(lcd.SetDisplay(hdc.DisplayOn))

	if *c {
		checkErr(lcd.ClearDisplay())
	}

	checkErr(lcd.MoveCursor(*x, *y))
	_, err = fmt.Fprintf(lcd, text)
	checkErr(err)
}
