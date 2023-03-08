Chip 8 Emulator in Go

A Chip 8 emulator built using the Go programming language. Chip 8 is an interpreted programming language that was first used on the COSMAC VIP and Telmac 1800 microcomputers in the late 1970s.
Features

    Emulates the Chip 8 instruction set
    Supports input from the original Chip 8 hexadecimal keyboard
    Ability to load and run Chip 8 programs
    Option to toggle the display of the original Chip 8 screen resolution

Installation

To install and run the Chip 8 emulator, you will need to have Go installed on your computer.
you can find instructions on doing this at https://go.dev/dl/

Clone the repository to your local machine:
``` bash
git clone https://github.com/alex-grimes/go-chip-8.git
```
Change into the project directory:

```go
cd chip8-emulator-go
```
Build the emulator:

```go
go build
```
Run the emulator:

```go
./go-chip-8 [path to Chip 8 program]
```
Usage

The Chip 8 emulator supports the original Chip 8 hexadecimal keyboard for input. The key mapping is as follows:

```css
1 2 3 4
q w e r
a s d f
z x c v
```
You can toggle the display of the original Chip 8 screen resolution by pressing the F1 key.
Contributing

Contributions are welcome! If you would like to contribute to this project, please open a pull request.
License

This project is licensed under the MIT License.
