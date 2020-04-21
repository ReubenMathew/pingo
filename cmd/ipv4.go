/*
Copyright © 2020 NAME HERE reubenninan@outlook.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// ipv4Cmd represents the ipv4 command
var ipv4Cmd = &cobra.Command{
	Use:   "ipv4",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		targetAddr := args[0]

		ping := func(addr string) {
			sent, rtt, err := Ping(addr, 50)
			if err != nil {
				log.Printf("[FAILED] Ping %s: %s\nPacket was not sent\n", addr, err)
				return
			}
			resolvedIP, _ := net.ResolveIPAddr("ip4", addr)
			if sent {
				fmt.Printf("Reply from %s: bytes=32 time=%s\n", resolvedIP, rtt)
			} else {
				log.Printf("[FAILED] Ping %s: %s\nPacket was not sent\n", addr, err)
			}
		}

		ping(targetAddr)
	},
}

// Ping : main function to ping ipv4 addresses
// positional arguments :
// 	addr : target ipv4 address
// 	timeout : duration in milliseconds for timeout
// return fields :
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

	// read message ....
	// fmt.Println(m)

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

func init() {
	rootCmd.AddCommand(ipv4Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ipv4Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	ipv4Cmd.Flags().BoolP("timeout", "t", false, "set timeout in milliseconds")
}
