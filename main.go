package main

import (
	"flag"
)

func main() {
	importedFile := flag.String("import", "", "import tdesktop exported json file")
	configFile := flag.String("config", "config.yaml", "config file")
	databaseFile := flag.String("database", "data.db", "database file")
	flag.Parse()

	if *importedFile != "" {
		importData(*databaseFile, *importedFile)
		return
	}

	StartBot(*databaseFile, *configFile)
}
