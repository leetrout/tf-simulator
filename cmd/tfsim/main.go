package main

import (
	"log"

	"github.com/leetrout/terraform-sim/pkg/salmonstack"
)

func main() {
	log.Fatal(salmonstack.StartServer())
}
