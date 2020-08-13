package main

func main() {
	window := Window{}
	config := Config{}
	// cpu := CPU{}
	config.fetchFlags()
	// load bootrom
	if config.bootrom {

	} else {
		// skip bootrom
		
	}
	window.init()
	window.loop()
}