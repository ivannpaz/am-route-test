package cli

import (
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/template"
)

// alertCmd represents the alert command
var (
	checkConfigCmd = app.Command("check-config", "Validate alertmanager config files")
	checkFiles     = checkConfigCmd.Arg("check-files", "Files to be validated").ExistingFiles()
)

func init() {
	checkConfigCmd.Action(checkConfig)
	longHelpText["check-config"] = `Validate alertmanager config files

Will validate the syntax and schema for alertmanager config file
and associated templates. Non existing templates will not trigger
errors`
}

func checkConfig(element *kingpin.ParseElement, ctx *kingpin.ParseContext) error {
	return CheckConfig(*checkFiles)
}

func CheckConfig(args []string) error {
	failed := 0

	for _, arg := range args {
		fmt.Printf("Checking '%s'", arg)
		config, _, err := config.LoadFile(arg)
		if err != nil {
			fmt.Printf("  FAILED: %s\n", err)
			failed++
		} else {
			fmt.Printf("  SUCCESS\n")
		}

		if config != nil {
			fmt.Printf("Found %d templates: ", len(config.Templates))
			if len(config.Templates) > 0 {
				_, err = template.FromGlobs(config.Templates...)
				if err != nil {
					fmt.Printf("  FAILED: %s\n", err)
					failed++
				} else {
					fmt.Printf("  SUCCESS\n")
				}
			}
		}
		fmt.Printf("\n")
	}
	if failed > 0 {
		return fmt.Errorf("failed to validate %d file(s)", failed)
	}
	return nil
}
