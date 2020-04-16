package chip8

// Chip8 is representation of CPU
type Chip8 struct {
	GFX        [64 * 32]uint8
	Key        [16]uint8
	memory     [4096]uint8
	v          [16]uint16 //CPU registers
	i          uint16     //index register
	pc         uint16     //program counter
	delayTimer uint8
	soundTimer uint8
	stack      [16]uint16
	sp         uint8 //stack pointer
}
