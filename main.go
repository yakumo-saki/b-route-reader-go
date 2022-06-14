package main

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/yakumo-saki/b-route-reader-go/src/mymod"
)

func main() {
	log.Info().Msg("Start")
	fmt.Println("test")
	mymod.Test()

}
