package main

import (
	"log"

	"github.com/leetrout/terraform-sim/internal/salmonstack"
)

func main() {
	log.Fatal(salmonstack.StartServer())
}
