package commands

import (
	"log"
	"os"
)

func pullTemplates() error {
	var err error
	exists, err := os.Stat("./template")
	if err != nil || exists == nil {
		log.Println("No templates found in current directory.")

		err = fetchTemplates()
		if err != nil {
			log.Println("Unable to download templates from Github.")
			return err
		}
	}
	return err
}
