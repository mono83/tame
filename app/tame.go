package main

import (
	"github.com/mono83/tame/app/cmd"
	"github.com/mono83/xray/std/xcobra"
)

func main() {
	xcobra.Start(cmd.TameCmd)
}
