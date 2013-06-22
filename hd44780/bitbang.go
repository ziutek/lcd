package hd44780

// ClockedBitbang handles communication with HD44780 driven LCD over some
// controller that implements clocked output of received data onto the D4-D7, E,
// RS pins (eg. controlers based on FTDI USB converter chips like FT232R,
// FT245R). Additionaly it provides methods for user controll R/W and AUX bits
// (if you want to use R/W bit as second AUX, the real R/W pin should be
// connected to the VSS).
type ClockedBitbang struct {
	w io.Writer
	e, rw, aux byte
}
