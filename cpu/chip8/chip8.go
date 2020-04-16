package chip8

import (
	"math/rand"
	"time"
)

var fontset = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

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

// Init initialises chip8 cpu
func (c *Chip8) Init() {
	c.pc = 0x200
	for i := 0; i < 80; i++ {
		c.memory[i] = fontset[i]
	}
	rand.Seed(time.Now().UTC().UnixNano())
}

func (c *Chip8) fetch() uint16 {
	// In our Chip 8 emulator, data is stored in an array
	// in which each address contains one byte. As one opcode is
	// 2 bytes long, we will need to fetch two successive bytes and
	// merge them to get the actual opcode.
	// So here we shifted current byte at memory pc left 8 bits, which adds 8 zeros
	// Next we use the bitwise OR operation to merge them and get two bytes long opcode
	return uint16(c.memory[c.pc]<<8 | c.memory[c.pc+1])
}
