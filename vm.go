
type VM struct {
	romlength              uint16
	pc, I, opcode, sp      uint16
	V                      [16]uint8
	memory                 [4096]uint8
	screen                 [32][8]uint8 // this is my 64x32 bitmap
	delayTimer, soundTimer uint8
	stack                  [16]uint16
	drawflog               bool
	wrapX                  string
	wrapY                  string
	clockSpeed             uint16
	timerSpeen             uint16
	screenBuffer           uint8
}

func (vm *VM) printState {
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

//4

//5
}