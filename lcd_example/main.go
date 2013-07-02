package main

import (
	"fmt"
	"github.com/ziutek/lcd"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printlnEdge(l int) {
	os.Stdout.Write([]byte{'+'})
	for x := 0; x < l; x++ {
		os.Stdout.Write([]byte{'-'})
	}
	os.Stdout.Write([]byte{'+', '\n'})

}

func printWindow(w lcd.TextWindow) {
	r := w.Runes()
	printlnEdge(w.Width())
	for y := 0; y < w.Height(); y++ {
		os.Stdout.Write([]byte{'|'})
		for x := 0; x < w.Width(); x++ {
			fmt.Printf("%c", r[y*w.Width()+x])
		}
		os.Stdout.Write([]byte{'|', '\n'})
	}
	printlnEdge(w.Width())
}

func main() {
	w := lcd.NewTextWindow(20, 4)
	w.SetCPos(3, 2)
	fmt.Fprintln(w, "Halo to ja!")
	fmt.Fprintln(w, "Ja też!")

	printWindow(w)
}
