package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	content, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(fmt.Errorf("Error opening file: %s", err))
	}

	program := string(content)
	program = strings.TrimSpace(program)
	startingOpcodes, err := ParseProgram(program)
	if err != nil {
		panic(fmt.Errorf("Error parsing program: %s", err))
	}

	router := NewRouter()
	for i := int64(0); i < 50; i++ {
		router.AddOutput(i, networkDevice(startingOpcodes, i, router.Input))
	}

	router.Run()
}

type Packet struct {
	Address int64
	X       int64
	Y       int64
}

type PacketQueue []int64

func NewPacketQueue() *PacketQueue {
	r := []int64{}

	return (*PacketQueue)(&r)
}

func (pq *PacketQueue) Add(p Packet) {
	*pq = append(*pq, p.X, p.Y)
}

func (pq *PacketQueue) Get() int64 {
	if len(*pq) == 0 {
		return int64(-1)
	} else {
		r := (*pq)[0]
		*pq = (*pq)[1:]
		return r
	}
}

type Router struct {
	Input  chan Packet
	Output map[int64](chan Packet)
}

func NewRouter() *Router {
	i := make(chan Packet)
	o := make(map[int64](chan Packet))

	r := Router{i, o}
	return &r
}

func (ro *Router) AddOutput(address int64, c chan Packet) {
	ro.Output[address] = c
}

func (ro *Router) Run() {
	for {
		p := <-ro.Input
		fmt.Println(p)
		if p.Address == 255 {
			fmt.Printf("Y = %d\n", p.Y)
			os.Exit(0)
		}
		ro.Output[p.Address] <- p
		fmt.Fprintf(os.Stderr, "Sent %v\n", p)
	}
}

func networkDevice(program map[int64]int64, address int64, routerInput chan Packet) chan Packet {
	nicInput := make(chan Packet)
	input := make(chan int64)
	output := make(chan int64)

	programCopy := CopyProgram(program)

	go func() {
		if err := ExecuteProgram(programCopy, input, output); err != nil {
			panic(fmt.Errorf("Error executing program: %s", err))
		}
	}()

	// Send network address
	input <- address

	// Send output of program to router input
	go func() {
		buffer := make([]int64, 0)
		for v := range output {
			buffer = append(buffer, v)
			if len(buffer) == 3 {
				p := Packet{buffer[0], buffer[1], buffer[2]}
				routerInput <- p
				buffer = buffer[:0]
			}
		}
	}()

	// Buffer NIC input and send input to program
	go func() {
		pq := NewPacketQueue()
		toSend := int64(-1)
		for {
			select {
			case p := <-nicInput:
				pq.Add(p)
				if toSend == -1 {
					toSend = pq.Get()
				}
			case input <- toSend:
				toSend = pq.Get()
			}
		}
	}()

	return nicInput
}
