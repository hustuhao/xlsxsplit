package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/metafates/xlsxsplit/cmd"
	"github.com/metafates/xlsxsplit/config"
	"github.com/metafates/xlsxsplit/logger"
	"github.com/samber/lo"
)

func handlePanic() {
	if err := recover(); err != nil {
		log.Error("crashed", "err", err)
		os.Exit(1)
	}
}

func main() {
	defer handlePanic()

	// prepare config and logs
	lo.Must0(config.Init())
	lo.Must0(logger.Init())

	// run the app
	cmd.Execute()
}
