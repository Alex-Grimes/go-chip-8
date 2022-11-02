
type VM struct {
	romlength              uint16
	pc, I, opcode, sp      uint16
	V                      [16]uint8
	memory                 [4096]uint8
	screen                 [32][8]uint8
	delayTimes, soundTimer uint8
	stack                  [16]uint16
	drawflog               bool
	wrapX                  string
	wrapY                  string
	clockSpeed             uint16
	timerSpeen             uint16
	screenBuffer           uint8
}