package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/maolonglong/brainfuck/internal"
)

const DefaultCellSize = 30000

func bfJumps(code []byte) (map[int]int, error) {
	stack := internal.NewStack()
	jumps := make(map[int]int)
	for i, n := 0, len(code); i < n; i++ {
		switch code[i] {
		case '[':
			stack.Push(i)
		case ']':
			j, err := stack.Pop()
			if err != nil {
				return nil, errors.New("unexpected closing bracket")
			}
			jumps[i] = j
			jumps[j] = i
		}
	}
	if !stack.Empty() {
		return nil, errors.New("excessive opening brackets")
	}
	return jumps, nil
}

func eval(file io.Reader, size int) error {
	code, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	in := bufio.NewReader(os.Stdin)

	jumps, err := bfJumps(code)
	if err != nil {
		return err
	}

	cell := make([]byte, size)
	pointer := 0

	for i, n := 0, len(code); i < n; i++ {
		switch code[i] {
		case '+':
			cell[pointer]++
		case '-':
			cell[pointer]--
		case '>':
			pointer++
			if pointer == size {
				return errors.New(fmt.Sprintf("pointer out of bounds: %v", pointer))
			}
		case '<':
			if pointer == 0 {
				return errors.New(fmt.Sprintf("pointer out of bounds: %v", pointer))
			}
			pointer--
		case '.':
			fmt.Printf("%c", cell[pointer])
		case ',':
			if cell[pointer], err = in.ReadByte(); err != nil {
				os.Exit(1)
			}
		case '[':
			if cell[pointer] == 0 {
				i = jumps[i]
			}
		case ']':
			if cell[pointer] > 0 {
				i = jumps[i]
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [file.bf]\n", os.Args[0])
		os.Exit(3)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	defer file.Close()

	err = eval(file, DefaultCellSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
