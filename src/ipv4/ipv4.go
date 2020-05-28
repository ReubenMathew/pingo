package ipv4

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"pingo/src/ping"
	"syscall"
	"time"
)

// PacketSend : Send packets to an ipv4 address
func PacketSend(targetAddr string, timeout int) {

	timeout = 100

	// fmt.Println(timeout)

	packetSent, packetRecv := 0, 0

	// Handle CTRL-C interrupt, else infinite loop
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() { // start a go routine that will run until interrupt
		<-c
		packetLoss := (packetSent - packetRecv)
		fmt.Printf("\nPackets: Sent = %d, Recieved = %d, Lost = %d ", packetSent, packetRecv, packetLoss)
		os.Exit(0)
	}()

	for {
		ping := func(addr string) {
			sent, rtt, err := ping.Ping(addr, timeout)
			packetSent++
			if err != nil {
				log.Printf("[FAILED] Ping %s: %s\nPacket was not sent\n", addr, err)
				return
			}
			resolvedIP, _ := net.ResolveIPAddr("ip4", addr)
			if sent {
				fmt.Printf("Reply from %s: bytes=32 time=%s\n", resolvedIP, rtt)
				packetRecv++
			} else {
				log.Printf("[FAILED] Ping %s: %s\nPacket was not sent\n", addr, err)
			}
		}
		ping(targetAddr)
		time.Sleep(3 * time.Second)
	}
}
