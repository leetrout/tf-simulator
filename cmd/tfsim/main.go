// Command tfsim runs the statesim server: the Nimbus fake-cloud API, the
// simulator/observer API, the websocket event stream, and the embedded UI — all
// from a single binary.
package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/leetrout/terraform-sim/internal/cloud"
	"github.com/leetrout/terraform-sim/internal/sim"
	"github.com/leetrout/terraform-sim/internal/webserver"
	"github.com/leetrout/terraform-sim/internal/ws"
	"github.com/pkg/browser"
)

func main() {
	addrFlag := flag.String("addr", ":9321", "address for the web server to listen on")
	workDirFlag := flag.String("work-dir", ".", "working directory watched for terraform.tfstate")
	seedFlag := flag.String("seed", "demo", "demo|empty: seed the cloud on first run")
	noOpenFlag := flag.Bool("no-open", false, "do not open a browser on startup")
	flag.Parse()

	store := cloud.NewStore(cloud.DefaultStateFile)
	if err := store.Load(); err != nil {
		log.Fatalf("load cloud state: %v", err)
	}

	simulator := sim.New(store, *workDirFlag)

	// Seed the demo topology only on a fresh cloud (so restarts preserve state)
	// and only when the work dir has no Terraform config. A dir containing .tf
	// files signals a real terraform workflow: seeding (which writes a demo
	// terraform.tfstate) would collide with the user's own apply, so we skip it.
	if *seedFlag == "demo" && len(store.All()) == 0 {
		if hasTerraformConfig(*workDirFlag) {
			fmt.Printf("+++ %s has .tf config — skipping demo seed (run terraform yourself; use --seed empty to silence)\n", *workDirFlag)
		} else {
			if err := simulator.SeedDemo(true); err != nil {
				log.Fatalf("seed demo: %v", err)
			}
			fmt.Println("+++ seeded demo topology (cloud + terraform.tfstate)")
		}
	}

	ws.Init()
	simulator.Start()
	defer simulator.Stop()

	mux := webserver.NewMux(store, simulator)

	onListen := func(bound string) {
		fmt.Printf("+++ statesim server listening at http://%s\n", browserHost(bound))
		if !*noOpenFlag {
			_ = browser.OpenURL("http://" + browserHost(bound))
		}
	}

	if err := webserver.Serve(*addrFlag, mux, onListen); err != nil {
		log.Fatalf("server: %v", err)
	}
}

// hasTerraformConfig reports whether dir contains any .tf files (i.e. the user
// intends to drive terraform themselves).
func hasTerraformConfig(dir string) bool {
	matches, err := filepath.Glob(filepath.Join(dir, "*.tf"))
	return err == nil && len(matches) > 0
}

// browserHost turns a bound listener address (e.g. "[::]:9321" or ":9321") into a
// host:port a browser can open.
func browserHost(bound string) string {
	host := bound
	if i := strings.LastIndex(bound, ":"); i >= 0 {
		port := bound[i+1:]
		h := bound[:i]
		if h == "" || h == "[::]" || h == "0.0.0.0" || h == "::" {
			h = "localhost"
		}
		host = h + ":" + port
	}
	return host
}
