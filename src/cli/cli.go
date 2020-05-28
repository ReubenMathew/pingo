package cli

import (
	"fmt"
	"pingo/src/ipv4"
	"pingo/src/validate"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

// the questions to ask
var qs = []*survey.Question{
	{
		Name:      "addr",
		Prompt:    &survey.Input{Message: "Address to ping: "},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "protocol",
		Prompt: &survey.Select{
			Message: "Choose an IP protocol: ",
			Options: []string{"ipv4", "ipv6", "exit"},
			Default: "ipv4",
		},
	},
	{
		Name:      "timeout",
		Prompt:    &survey.Input{Message: "Enter a timeout in ms: "},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
}

// ServiceSelector : CLI selector for SodaDB services
func ServiceSelector() (string, string, int) {
	answers := struct {
		Address  string `survey:"addr"`
		Protocol string `survey:"protocol"`
		Timeout  int    `survey:"timeout"`
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	return answers.Address, answers.Protocol, answers.Timeout
}

func exitMsg() {
	fmt.Printf("\nexiting pingo...")
}

// Launch : Launches a SodaDB service from an interactive CLI selector
func Launch() {

	defer exitMsg()

	address, protocol, timeout := ServiceSelector()

	if timeout == 0 {
		fmt.Println("Timeout set to 100ms")
		timeout = 100
	}

	switch protocol {
	case "ipv4":

		valid := validate.Ipv4(address)
		if valid {
			fmt.Printf("\nYou chose to ping %s:%s\n", protocol, address)
			time.Sleep(3 * time.Second)
			fmt.Printf("Press CTRL+C to stop pingo\n\n")
			time.Sleep(3 * time.Second)
			ipv4.PacketSend(address, timeout)
		} else {
			fmt.Printf("Invalid hostname, please use a valid IPV4 address")
		}

	case "ipv6":
		fmt.Printf("\nSorry, ipv6 echo requests have not been implemented yet.\n")
	case "exit":
		fmt.Printf("\nYou chose to exit pingo\n")
	default:
		fmt.Printf("\nYou chose to exit pingo\n")
	}
	fmt.Println()
}
