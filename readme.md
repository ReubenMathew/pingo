# Pingo

[![Go Report Card](https://goreportcard.com/badge/github.com/ReubenMathew/pingo)](https://goreportcard.com/report/github.com/ReubenMathew/pingo)

Pingo is a lightweight ICMP echo request CLI.

## Installation

Check the *releases* page for platform specific download

### Build from source

Navigate to `/src` and run `go build -v -o pingo`

`go build -v -o pingo.exe` for Windows platforms


### Build for all platforms 
Navigate to `/tools` and run `./build.sh`

*Makefile coming soon*

## Usage

```bash
$ ./pingo
```

### Example
```bash
? Address to ping:  192.168.0.1
? Choose an IP protocol:  ipv4                                 
? Enter a timeout in ms:  200
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
