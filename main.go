package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CmdOpts struct {
	service  string
	address  string
	rpcType  string
	endpoint string
	port     int
	interval int
	batch    int
	multi    int
	list     bool
	help     bool
}

var (
	opts  CmdOpts
	Modes = map[string]map[string]int{
		"tsdb": map[string]int{
			"rpc":     8433,
			"grpc":    8434,
			"jsonrpc": 8435,
		},
	}
)

func listNodes() {
	fmt.Printf("\n%8s  %8s  %7s\n", "ServName", "RpcType", "RpcPort")
	fmt.Println()
	for serv, rpcMap := range Modes {
		for rpcType, port := range rpcMap {
			fmt.Printf("%8s  %8s  %7d\n", serv, rpcType, port)
		}
	}
	fmt.Println()
}

func getPort(service, rpcType string) (int, error) {
	if _, ok := Modes[service]; !ok {
		return 0, fmt.Errorf("Service %s not found", service)
	}

	port, ok := Modes[service][rpcType]
	if !ok {
		return 0, fmt.Errorf("Rpc type %s of service %s not found", rpcType, service)
	}
	return port, nil
}

func init() {
	flag.StringVar(&opts.service, "s", "tsdb", "service name: tsdb")
	flag.StringVar(&opts.address, "a", "127.0.0.1", "service address: 127.0.0.1")
	flag.StringVar(&opts.rpcType, "r", "jsonrpc", "rpc type: rpc/grpc/jsonrpc")
	flag.StringVar(&opts.endpoint, "e", "endpoint", "specify endpoint's prefix")
	flag.IntVar(&opts.port, "p", 0, "service port (default 8434)")
	flag.IntVar(&opts.interval, "i", 1000, "interval time(ms) create mock data every time")
	flag.IntVar(&opts.batch, "b", 5, "how many counters create every time")
	flag.IntVar(&opts.multi, "m", 1, "how many processes for mock")
	flag.BoolVar(&opts.list, "l", false, "list all service and type")
	flag.BoolVar(&opts.help, "h", false, "print usage")
}

func main() {
	flag.Parse()

	if opts.help {
		flag.Usage()
		os.Exit(0)
	}

	if opts.list {
		listNodes()
		os.Exit(0)
	}

	defaultPort, err := getPort(opts.service, opts.rpcType)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	if opts.port <= 0 {
		opts.port = defaultPort
	}

	mocker := NewMocker(&opts)
	mocker.Mock()
	mocker.Stats()
}
