package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/leetrout/terraform-sim/internal/store"
	"github.com/leetrout/terraform-sim/internal/webserver"
	"github.com/leetrout/terraform-sim/internal/ws"
	"github.com/pkg/browser"
)

const defAddr = ":9321"

var addrFlag = flag.String("addr", defAddr, "address for webserver to listen")
var interactiveFlag = flag.Bool("i", false, "interactive mode (console)")

func main() {
	flag.Parse()

	addr := defAddr
	if addrFlag != nil {
		addr = *addrFlag
	}

	addrParts := strings.Split(addr, ":")
	if len(addrParts) != 2 {
		fmt.Fprint(os.Stderr, "addr in bad format, expected `<host>:<port>` or `:<port>`")
		os.Exit(2)
	}
	if addrParts[0] == "" {
		addrParts[0] = "localhost"
	}

	store.Initialize()
	if *interactiveFlag == true {
		go console()
	}

	timer := time.NewTimer(200 * time.Millisecond)
	go func() {
		<-timer.C
		browser.OpenURL(fmt.Sprintf("http://%s:%s", addrParts[0], addrParts[1])) // FIXME
	}()

	fmt.Printf("+++ API server starting at %s\n", addr)
	webserver.RunServer(addr)
}

func console() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("TerraSim")
	fmt.Println("exit: Close the program")
	fmt.Println("sock: Show websocket info")
	fmt.Println("sock send <string>: Broadcast the string to all WS clients")
	fmt.Println("---------------------")
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.ReplaceAll(text, "\n", "")

		if strings.Compare("sock", text) == 0 {
			fmt.Printf("%+v\n", ws.SOCKEX)
		}

		if strings.HasPrefix(text, "sock send") {
			ws.SOCKEX.Send(strings.ReplaceAll(text, "sock send ", ""))
		}

		if strings.Compare("exit", text) == 0 {
			os.Exit(0)
		}

	}
}
