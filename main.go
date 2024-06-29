package main

import (
	"flag"
)

func main() {
	importedFile := flag.String("import", "", "import tdesktop exported json file")
	configFile := flag.String("config", "config.json", "config file")
	flag.Parse()

	if *importedFile != "" {
		importData(*importedFile)
		return
	}

	StartBot(*configFile)
}
