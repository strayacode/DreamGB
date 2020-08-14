package main

import (
	"fmt"
	"flag"
	"strings"
	"os"
)

type Config struct {
	width int // comes from flag with multiplier. default is 1x
	height int // comes from flag with multiplier. default is 1x
	stepmode bool // allows the user to step through each instruction one at a time. default is false
	bootrom bool // allows the user to specify a bios to run. default is false
	bootrompath string // path to bios if specified by user
	rompath string // rom given by user
}

func (config *Config) fetchFlags() {
	var stepmode bool
	var bootrompath string
	flag.BoolVar(&stepmode, "step", false, "allows the user to go through each instruction step by step. helpful for debugging")
	flag.StringVar(&bootrompath, "bootrom", "", "allows user to specify a bootrom instead of skipping")
	flag.Parse()
	// only allow for gb rom as non-flag argument
	if !(len(flag.Args()) == 1 && strings.HasSuffix(flag.Args()[0], ".gb")) {
		fmt.Println("please provide a valid rom type such as .gb!")
		os.Exit(0)
	}
	// keep default for now
	config.width = 160
	config.height = 144
	config.stepmode = stepmode
	config.bootrompath = bootrompath
	if len(config.bootrompath) > 0 {
		config.bootrom = true
	}
	config.rompath = flag.Args()[0]
}