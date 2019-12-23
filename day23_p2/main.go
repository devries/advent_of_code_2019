package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const idleRepeat = 10000

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
		router.AddOutput(i, networkDevice(startingOpcodes, i, router.Input, router.Idle))
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

type IdleAlert struct {
	Address int64
	Idle    bool
}

type Router struct {
	Input  chan Packet
	Output map[int64](chan Packet)
	Nat    Packet
	Idle   chan IdleAlert
}

func NewRouter() *Router {
	i := make(chan Packet)
	o := make(map[int64](chan Packet))
	m := make(chan IdleAlert)

	r := Router{i, o, Packet{0, 0, 0}, m}
	return &r
}

func (ro *Router) AddOutput(address int64, c chan Packet) {
	ro.Output[address] = c
}

func (ro *Router) Run() {
	sentY := make(map[int64]bool)

	for {
		select {
		case p := <-ro.Input:
			if p.Address == 255 {
				ro.Nat = p
				fmt.Fprintf(os.Stderr, "Received %v\n", ro.Nat)
			} else {
				ro.Output[p.Address] <- p
			}
		case <-time.After(1 * time.Second):
			ro.Output[0] <- ro.Nat
			fmt.Fprintf(os.Stderr, "Sent %v to 0\n", ro.Nat)
			if sentY[ro.Nat.Y] {
				fmt.Printf("First Y sent twice from NAT to addr 0: %d\n", ro.Nat.Y)
				os.Exit(0)
			} else {
				sentY[ro.Nat.Y] = true
			}

		}
	}
}

func networkDevice(program map[int64]int64, address int64, routerInput chan Packet, idleChan chan IdleAlert) chan Packet {
	nicInput := make(chan Packet)
	input := make(chan int64)
	output := make(chan int64)

	programCopy := CopyProgram(program)

	go func() {
		if err := ExecuteProgram(programCopy, input, output); err != nil {
			panic(fmt.Errorf("Error executing program: %s", err))
		}
	}()

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

	// Send network address
	input <- address

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
