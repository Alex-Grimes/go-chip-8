package main

import (
	"fmt"
	"log"
)

type Display interface {
	clearDisplay()
	updateDisplay()
	drawPixel(x int32, y int32)
}

type Keyboard interface {
	waitForKeyPress() (uint8, bool)
	isKeyPressed(key uint8) bool
	specialKeyPressed(paused bool) (bool, bool)
}
type VM struct {
	romlength              uint16
	pc, I, opcode, sp      uint16
	V                      [16]uint8
	memory                 [4096]uint8
	screen                 [32][8]uint8 // this is my 64x32 bitmap
	delayTimer, soundTimer uint8
	stack                  [16]uint16
	drawflag               bool
	wrapX                  string
	wrapY                  string
	clockSpeed             uint16
	timerSpeed             uint16
	screenBuffer           uint8
}

func (vm *VM) printState() {
	fmt.Printf("PC: 0x%x\n, vm.pc")
	fmt.Printf("I: 0x%x\n", vm.I)
	fmt.Printf("Opcode: 0x%x\n", vm.opcode)
	fmt.Println("Memory:")
	fmt.Println(vm.memory)
	fmt.Println("Screen:")
	fmt.Println(vm.screen)
	fmt.Printf("DT: %d\n", vm.delayTimer)
	fmt.Printf("ST: %d\n", vm.soundTimer)
	fmt.Printf("SP: %d\n", vm.sp)
	fmt.Println("Stack:")
	fmt.Println(vm.stack)
}

func (vm *VM) initialiseFont() {
	vm.memory[0x050] = 0xF0
	vm.memory[0x051] = 0x90
	vm.memory[0x052] = 0x90
	vm.memory[0x053] = 0x90
	vm.memory[0x054] = 0xF0
	//1
	vm.memory[0x055] = 0x20
	vm.memory[0x056] = 0x60
	vm.memory[0x057] = 0x20
	vm.memory[0x058] = 0x20
	vm.memory[0x059] = 0x70
	//2
	vm.memory[0x05A] = 0xF0
	vm.memory[0x05B] = 0x10
	vm.memory[0x05C] = 0xF0
	vm.memory[0x05D] = 0x80
	vm.memory[0x05E] = 0xF0
	//3
	vm.memory[0x05F] = 0xF0
	vm.memory[0x060] = 0x10
	vm.memory[0x061] = 0xF0
	vm.memory[0x062] = 0x10
	vm.memory[0x063] = 0xF0
	//4
	vm.memory[0x064] = 0x90
	vm.memory[0x065] = 0x90
	vm.memory[0x066] = 0xF0
	vm.memory[0x067] = 0x10
	vm.memory[0x068] = 0x10
	//5
	vm.memory[0x069] = 0xF0
	vm.memory[0x06A] = 0x80
	vm.memory[0x06B] = 0xF0
	vm.memory[0x06C] = 0x10
	vm.memory[0x06D] = 0xF0
	//6
	vm.memory[0x06E] = 0xF0
	vm.memory[0x06F] = 0x80
	vm.memory[0x070] = 0xF0
	vm.memory[0x071] = 0x90
	vm.memory[0x072] = 0xF0
	//7
	vm.memory[0x073] = 0xF0
	vm.memory[0x074] = 0x10
	vm.memory[0x075] = 0x20
	vm.memory[0x076] = 0x40
	vm.memory[0x077] = 0x40
	//8
	vm.memory[0x078] = 0xF0
	vm.memory[0x079] = 0x90
	vm.memory[0x07A] = 0xF0
	vm.memory[0x07B] = 0x90
	vm.memory[0x07C] = 0xF0
	//9
	vm.memory[0x07D] = 0xF0
	vm.memory[0x07E] = 0x90
	vm.memory[0x07F] = 0xF0
	vm.memory[0x080] = 0x10
	vm.memory[0x081] = 0x10
	//A
	vm.memory[0x082] = 0xF0
	vm.memory[0x083] = 0x90
	vm.memory[0x084] = 0xF0
	vm.memory[0x085] = 0x90
	vm.memory[0x086] = 0x90
	//B
	vm.memory[0x087] = 0xE0
	vm.memory[0x088] = 0x90
	vm.memory[0x089] = 0xE0
	vm.memory[0x08A] = 0x90
	vm.memory[0x08B] = 0xE0
	//C
	vm.memory[0x08C] = 0xF0
	vm.memory[0x08D] = 0x80
	vm.memory[0x08E] = 0x80
	vm.memory[0x08F] = 0x80
	vm.memory[0x090] = 0xF0
	//D
	vm.memory[0x091] = 0xE0
	vm.memory[0x092] = 0x90
	vm.memory[0x093] = 0x90
	vm.memory[0x094] = 0x90
	vm.memory[0x095] = 0xE0
	//E
	vm.memory[0x096] = 0xF0
	vm.memory[0x097] = 0x80
	vm.memory[0x098] = 0xF0
	vm.memory[0x099] = 0x80
	vm.memory[0x09A] = 0xF0
	//F
	vm.memory[0x09B] = 0xF0
	vm.memory[0x09C] = 0x80
	vm.memory[0x09D] = 0xF0
	vm.memory[0x09E] = 0x80
	vm.memory[0x09F] = 0x80
}

func (vm *VM) loadROM(rombytes []byte) {
	vm.romlength = uint16(len(rombytes))
	for i, byt := range rombytes {
		vm.memory[0x200+i] = byt
	}
}

func (vm *VM) init(rombytes []byte, wrapX string, wrapY string, clockSpeed int, timerSpeed int, screenBuffer int) {
	vm.initialiseFont()
	vm.loadROM(rombytes)
	vm.pc = 0x200
	vm.drawflag = false
	vm.wrapX = wrapX
	vm.wrapY = wrapY
	vm.clockSpeed = uint16(clockSpeed)
	vm.timerSpeed = uint16(timerSpeed)
	vm.screenBuffer = uint8(screenBuffer)
}

func (vm *VM) parseOpcode(keyboard Keyboard) bool {
	var running bool
	vm.opcode = uint16(vm.memory[vm.pc])<<8 | uint16(vm.memory[vm.pc+1])
	vm.drawflag = false
	switch vm.opcode & 0xF000 {
	case 0x0000:
		switch vm.opcode & 0x00FF {
		case 0x00E0:
			//clear screen
			//fmt.Printf("Clear vm.screen % x, % d\n, vm.opcode, vm.pc")
			for yp := 0; yp < 32; yp++ {
				for xb := 0; xb < 8; xb++ {
					vm.screen[yp][xb] = 0
				}
			}
			vm.drawflag = true
			vm.pc += 2
		case 0x00EE:
			// 00EE - RET
			// Return from subroutine
			if vm.sp <= 0 {
				log.Fatal(fmt.Errorf("vm.stack pointer below 0"))
			}
			//fmt.Printf("RET pc: %x, new pc: %x\n", vm.pc, vm.stack[vm.sp]+2, vm.opcode)
			vm.sp--
			vm.pc = vm.stack[vm.sp] + 2
			// fmt.Printf("RET pc: %d/n", vm.pc)
		default:
			// fmt.Printf("SYS vm.opcode ignored: % x, % d\n," vm.opcode, vm.pc)
			vm.pc += 2
		}
	case 0x1000:
		// lnnn - JP addr
		// Jump to location nnn.
		vm.pc = 0x0FFF & vm.opcode
		if vm.pc < 0x200 || vm.pc > 0xFFF {
			log.Fatal(fmt.Errorf("illegal JMP instructions - PC: %x, opcode: %x", vm.pc, vm.opcode))
		}
		// fmt.Printf("JMP to : % x, % x\n", vm.pc, vm.opcode)
		// endless jumps used as halt

	case 0x2000:
		// 2nnn - CALL addr
		// Call subroutine at nnn
		vm.stack[vm.sp] = vm.pc
		vm.sp++
		// fmt.Printf("CALL pc: %x, new pc: %x, opcode: %x\n", vm.pc, (oxoFFF & vm.opsode), vm.opcode)
		vm.pc = 0x0fff & vm.opcode
		if vm.pc < 0x200 || vm.pc > 0xFFF {
			log.Fatal(fmt.Errorf("illegal JMP instructions - PC: %x, opcode: %x", vm.pc, vm.opcode))
		}
		// fmt.Printf("CALL stack: %s\n", fmt.Sprint(vm.stack))
	case 0x3000:
		//3xkk - SE vm.Vx, byte
		// skip next instruction if vm.Vx = kk
		if vm.V[0x0F00&vm.opcode>>8] == uint8(0x00FF&vm.opcode) {
			vm.pc += 4
		} else {
			vm.pc += 2
		}
	case 0x4000:
		// 4xkk - sne vm.Vx, byte
		// Skip next instruction if vm.Vx != kk
		if vm.V[0x0F00&vm.opcode>>8] != uint8(0x00FF&vm.opcode) {
			vm.pc += 4
		} else {
			vm.pc += 2
		}
	case 0x5000:
		// 5xy0 - SE vm.Vx, vm.Vy
		// Skip next instruction if vm.Vx = vm.Vy
		if vm.V[0x0F00&vm.opcode>>8] == vm.V[uint8(0x00F0&vm.opcode)>>4] {
			vm.pc += 4
		} else {
			vm.pc += 2
		}
	case 0x6000:
		// 6xkk - LD vm.Vx, byte
		// Set vm.Vx == kk
		vm.V[0x0F00&vm.opcode>>8] = uint8(0x00F & vm.opcode)
		vm.pc += 2

	case 0x7000:
		// 7xkk - ADD vm.Vx, byte
		// Set vm.Vx + kk
		vm.V[0x0F00&vm.opcode>>8] += uint8(0x0FF & vm.opcode)
		vm.pc += 2

	case 0x8000:
		switch vm.opcode & 0x00F {
		case 0x0000:
			//8xy0 - LD vm.Vx, vm.Vy
			//Set vm.Vx = vm.Vy
			vm.V[vm.opcode&0x0F00>>8] = vm.V[vm.opcode&0x00F0>>4]
			vm.pc = +2

		}
	}
}
