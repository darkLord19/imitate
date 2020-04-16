package chip8

import (
	"fmt"
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
	v          [16]uint8 //CPU registers
	i          uint16    //index register
	pc         uint16    //program counter
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

func (c *Chip8) clearScreen() {
	for i := 0; i < 64*32; i++ {
		c.GFX[i] = 0
	}
}

// EmulateCycle emulates cpu cycles
func (c *Chip8) EmulateCycle() {
	opcode := c.fetch()
	switch opcode & 0xF000 {
	case 0x0:
		if opcode == 0x00E0 { // 0x00E0: Clears the screen
			c.clearScreen()
			c.pc += 2
			break
		} else if opcode == 0x00EE { // 0x00EE: Returns from subroutine
			c.sp--
			c.pc = c.stack[c.sp]
			c.pc += 2
			break
		}
		fmt.Printf("Unknown opcode 0x%X", opcode)
		break
	case 0x1: // 0x1NNN: Jumps to address NNN
		c.pc = opcode & 0x0FFF
		break
	case 0x2: // 0x2NNN: Calls subroutine at NNN
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = opcode & 0x0FFF
		break
	case 0x3: // 0x3XNN: Skips the next instruction if VX equals NN
		VX := c.v[opcode&0x0F00>>8]
		NN := uint8(opcode & 0x00FF)
		if VX == NN {
			c.pc += 4
		} else {
			c.pc += 2
		}
		break
	case 0x4: // 0x4XNN: Skips the next instruction if VX doesn't equal NN
		VX := c.v[opcode&0x0F00>>8]
		NN := uint8(opcode & 0x00FF)
		if VX != NN {
			c.pc += 4
		} else {
			c.pc += 2
		}
		break
	case 0x5: // 0x5XY0: Skips the next instruction if VX equals VY
		VX := c.v[opcode&0x0F00>>8]
		VY := c.v[opcode&0x00F0>>8]
		if VX == VY {
			c.pc += 4
		} else {
			c.pc += 2
		}
		break
	case 0x6: // 0x6XNN: Sets VX to NN
		X := opcode & 0x0F00 >> 8
		NN := uint8(opcode & 0x00FF)
		c.v[X] = NN
		c.pc += 2
		break
	case 0x7: // 0x7XNN: Adds NN to VX
		X := opcode & 0x0F00 >> 8
		NN := uint8(opcode & 0x00FF)
		c.v[X] += NN
		c.pc += 2
		break
	case 0x8:
		switch opcode & 0x000F {
		case 0x0: // 0x8XY1: Sets VX to the value of VY
			X := opcode & 0x0F00 >> 8
			Y := opcode & 0x00F0 >> 4
			c.v[X] = c.v[Y]
			break
		case 0x1: // 0x8XY1: Sets VX to VX or VY
			X := opcode & 0x0F00 >> 8
			Y := opcode & 0x00F0 >> 4
			c.v[X] |= c.v[Y]
			break
		case 0x2: // 0x8XY2: Sets VX to VX and VY
			X := opcode & 0x0F00 >> 8
			Y := opcode & 0x00F0 >> 4
			c.v[X] &= c.v[Y]
			break
		case 0x3: // 0x8XY3: Sets VX to VX xor VY
			X := opcode & 0x0F00 >> 8
			Y := opcode & 0x00F0 >> 4
			c.v[X] ^= c.v[Y]
			break
		case 0x4: // 0x8XY4: Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't
			X := opcode & 0x0F00 >> 8
			Y := opcode & 0x00F0 >> 4
			if c.v[Y] > (0xFF - c.v[X]) {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[X] += c.v[Y]
			break
		case 0x5: // 0x8XY5: VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			X := opcode & 0x0F00 >> 8
			Y := opcode & 0x00F0 >> 4
			if c.v[X] < c.v[Y] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[X] -= c.v[Y]
			break
		case 0x6: // 0x8XY6: Stores the least significant bit of VX in VF and then shifts VX to the right by 1
			X := opcode & 0x0F00 >> 8
			c.v[0xF] = c.v[X] & 0x1
			c.v[X] >>= c.v[X]
			break
		case 0x7: // 0x8XY7: Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't
			X := opcode & 0x0F00 >> 8
			Y := opcode & 0x00F0 >> 4
			if c.v[X] > c.v[Y] {
				c.v[0xF] = 0
			} else {
				c.v[0xF] = 1
			}
			c.v[X] = c.v[Y] - c.v[X]
			break
		case 0xE: // 0x8XYE: Stores the most significant bit of VX in VF and then shifts VX to the left by 1
			X := opcode & 0x0F00 >> 8
			c.v[0xF] = c.v[X] >> 16
			c.v[X] <<= 1
			break
		}
		c.pc += 2
		break
	case 0x9: // 0x9XY0: Skips the next instruction if VX doesn't equal VY
		X := opcode & 0x0F00 >> 8
		Y := opcode & 0x00F0 >> 4
		if c.v[X] != c.v[Y] {
			c.pc += 4
		} else {
			c.pc += 2
		}
		break
	case 0xA: // 0xANNN: Sets I to the address NNN
		c.i = opcode & 0x0FFF
		c.pc += 2
		break
	case 0xB: // 0xBNNN: Jumps to the address NNN plus V0
		c.pc = opcode&0x0FFF + uint16(c.v[0])
		c.pc += 2
		break
	case 0xC: // 0xCXNN: Sets VX to the result of a bitwise and operation on a random number and NN
		X := opcode & 0x0F00 >> 8
		NN := uint8(opcode & 0x00FF)
		random := uint8(rand.Intn(256))
		c.v[X] = random & NN
		c.pc += 2
		break
	case 0xD: // 0xDXYN: Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and
		//a height of N pixels. Each row of 8 pixels is read as bit-coded starting from memory location I;
		//I value doesn’t change after the execution of this instruction. VF is set to 1 if any
		//screen pixels are flipped from set to unset when the sprite is drawn,
		//and to 0 if that doesn’t happen
		c.pc += 2
		break
	case 0xE:
		switch opcode & 0x000F {
		case 0xE: // 0xEX9E: Skips the next instruction if the key stored in VX is pressed
			X := opcode & 0x0F00 >> 8
			if c.Key[c.v[X]] == 1 {
				c.pc += 4
			} else {
				c.pc += 2
			}
			break
		case 0x1: // 0xEXA1: Skips the next instruction if the key stored in VX isn't pressed
			X := opcode & 0x0F00 >> 8
			if c.Key[c.v[X]] != 1 {
				c.pc += 4
			} else {
				c.pc += 2
			}
			break
		}
		c.pc += 2
		break
	case 0xF:
		switch opcode & 0x00FF {
		case 0x07: // 0xFX07: Sets VX to the value of the delay timer
			X := opcode & 0x0F00 >> 8
			c.v[X] = c.delayTimer
			break
		case 0x0A: // 0xFX0A: A key press is awaited, and then stored in VX
			X := opcode & 0x0F00 >> 8
			keypress := false
			for i := 0; i < 16; i++ {
				if c.Key[i] != 0 {
					c.v[X] = uint8(i)
					keypress = true
					break
				}
			}
			if !keypress {
				return
			}
			c.pc += 2
			break
		case 0x15: // 0xFX15: Sets the delay timer to VX
			X := opcode & 0x0F00 >> 8
			c.delayTimer = c.v[X]
			break
		case 0x18: // 0xFX18: Sets the sound timer to VX
			X := opcode & 0x0F00 >> 8
			c.soundTimer = c.v[X]
			break
		case 0x1E: // 0xFX1E: Adds VX to I. VF is set to 1 when there is a range overflow (I+VX>0xFFF),
			// and to 0 when there isn't
			X := opcode & 0x0F00 >> 8
			if c.i+uint16(c.v[X]) > 0xFFF {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.i += uint16(c.v[X])
			break
		case 0x29: // 0xFX29: Sets I to the location of the sprite for the character in VX.
			// Characters 0-F (in hexadecimal) are represented by a 4x5 font
			break
		case 0x33: // 0xFX33: Stores the binary-coded decimal representation of VX,
			// with the most significant of three digits at the address in I, the middle digit at I plus 1,
			// and the least significant digit at I plus 2
			break
		case 0x55: // 0xFX55: Stores V0 to VX (including VX) in memory starting at address I.
			// The offset from I is increased by 1 for each value written, but I itself is left unmodified
			X := opcode & 0x0F00 >> 8
			cnt := c.i
			for i := uint16(0); i <= X; i++ {
				c.memory[cnt] = c.v[i]
				cnt++
			}
			break
		case 0x65: // 0xFX65: Fills V0 to VX (including VX) with values from memory starting at address I.
			// The offset from I is increased by 1 for each value written, but I itself is left unmodified
			X := opcode & 0x0F00 >> 8
			cnt := c.i
			for i := uint16(0); i <= X; i++ {
				c.v[i] = c.memory[cnt]
				cnt++
			}
			break
		}
		c.pc += 2
		break
	}
}
