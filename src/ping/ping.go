package ping

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// Ping : main function to ping ipv4 addresses
// positional arguments:
// 	addr : target ipv4 address
// 	timeout : duration in milliseconds for timeout
// return fields:
// 	packetSent, bool : did the packet send correctly?
// 	resolvedAddress, IPAddr : potentially resolved (modified) IP Address
// 	latency, Duration : Round Trip Time
// 	err, error : potential errors
func Ping(addr string, timeout int) (bool, time.Duration, error) {

	// local listening address
	listeningAddr := "0.0.0.0"

	// privileged raw ICMP endpoint requires ip4 followed by colon and ICMP protocol
	c, err := icmp.ListenPacket("ip4:icmp", listeningAddr)
	if err != nil {
		log.Fatal("Error with local port, exiting ...")
		return false, 0, err
	}

	// returns an ip endpoint if addr is not a literal IP address, otherwise it is parsed as a regular IP address
	targetAddrResolved, ipResolveErr := net.ResolveIPAddr("ip4", addr)

	if ipResolveErr != nil {
		log.Fatal(err, "Unresolvable IP Address")
		return false, 0, err
	}

	// ICMP packet of 32 bits
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho, // 8 bit
		Code: 0,                 // 0 bits
		Body: &icmp.Echo{
			ID:   0x0,                                //0 bit ID http://networksorcery.com/enp/protocol/icmp/msg8.htm#Identifier
			Seq:  0x0,                                //0 bit sequence number http://networksorcery.com/enp/protocol/icmp/msg8.htm#Sequence%20number
			Data: []byte("This-is-a-test-message__"), // 24 bit message
		},
		Checksum: 227, //pi symbol checksum (c-code 16 bit ascii), message is sent with checksum appended to the end, the bytes added with ones complement should add up to all ones
	}

	bytesSent, err := m.Marshal(nil)
	if err != nil {
		return false, 0, err
	}
	// fmt.Println(m)

	// RTT timer start
	start := time.Now()

	// Sending the message
	n, err := c.WriteTo(bytesSent, targetAddrResolved)
	// fmt.Printf("Bytes Sent %v; Write Back %v\n", n, len(bytesSent))
	if err != nil {
		return false, 0, err
	} else if n != len(bytesSent) { // If all bytes that were sent are not the same in the write back then throw error
		return false, 0, fmt.Errorf("Bytes Sent %v; Write Back %v", n, len(bytesSent))
	}

	// reply buffer of 1024 bytes
	replyMessage := make([]byte, 1024)

	// set timeout
	err = c.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))

	if err != nil {
		return false, 0, err
	}

	recv, err := icmp.ParseMessage(1, replyMessage[:n])
	if err != nil {
		return false, 0, err
	}

	duration := time.Since(start)

	c.Close()

	if recv.Type == ipv4.ICMPTypeEchoReply {
		return true, duration, nil
	}

	return false, 0, fmt.Errorf("Something went wrong")

}
