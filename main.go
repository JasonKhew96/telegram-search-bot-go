package main

import (
	"flag"
)

func main() {
	importedFile := flag.String("import", "", "import tdesktop exported json file")
	flag.Parse()

	if *importedFile != "" {
		importData(*importedFile)
		return
	}

	StartBot()
}
