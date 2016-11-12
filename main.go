package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ginuerzh/gost"
	"github.com/golang/glog"
	"golang.org/x/net/http2"
	"io/ioutil"
	"os"
	"runtime"
)

var (
	options struct {
		ChainNodes, ServeNodes flagStringList
	}
)

func init() {
	var (
		configureFile string
		printVersion  bool
	)

	flag.StringVar(&configureFile, "C", "", "configure file")
	flag.Var(&options.ChainNodes, "F", "forward address, can make a forward chain")
	flag.Var(&options.ServeNodes, "L", "listen address, can listen on multiple ports")
	flag.BoolVar(&printVersion, "V", false, "print version")
	flag.Parse()
	fmt.Fprintf(os.Stderr, "%s\n", &options.ChainNodes)
	if err := loadConfigureFile(configureFile); err != nil {
		fmt.Fprintf(os.Stdout, "no configure file\n")
		glog.Fatal(err)
	}

	if glog.V(5) {
		http2.VerboseLogs = true
	}

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stdout, "no parameters\n")
		return
	}

	if printVersion {
		fmt.Fprintf(os.Stderr, "GOST %s (%s)\n", gost.Version, runtime.Version())
		return
	}
	fmt.Fprintf(os.Stdout, "init finish\n")
}

func main() {
	chain := gost.NewProxyChain()
	if err := chain.AddProxyNodeString(options.ChainNodes...); err != nil {
		glog.Fatal(err)
	}
	chain.Init()

	serverNode, err := gost.ParseProxyNode(options.ServeNodes[0])
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Fprintf(os.Stdout, "old bind address: ", serverNode)
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	bind := fmt.Sprintf("%s:%s", host, port)
	fmt.Fprintf(os.Stdout, "\nbind address: %s\n", bind)
	serverNode.Addr = bind

	server := gost.NewProxyServer(serverNode, chain, &tls.Config{})
	glog.Fatal(server.Serve())
}

func loadConfigureFile(configureFile string) error {
	if configureFile == "" {
		return nil
	}
	content, err := ioutil.ReadFile(configureFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, &options); err != nil {
		return err
	}
	return nil
}

type flagStringList []string

func (this *flagStringList) String() string {
	return fmt.Sprintf("%s", *this)
}
func (this *flagStringList) Set(value string) error {
	*this = append(*this, value)
	return nil
}
