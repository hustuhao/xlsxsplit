package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/metafates/xlsxsplit/style"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/metafates/xlsxsplit/app"
	"github.com/metafates/xlsxsplit/filesystem"
	"github.com/metafates/xlsxsplit/icon"
	"github.com/metafates/xlsxsplit/where"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().BoolP("version", "v", false, app.Name+" version")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   strings.ToLower(app.Name),
	Short: app.DescriptionShort,
	Long:  app.DescriptionLong,
	Run: func(cmd *cobra.Command, args []string) {
		// 参数为空则提示帮助信息
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		splitCmd.Run(cmd, args)
	},
}

func Execute() {
	// Setup colored cobra
	cc.Init(&cc.Config{
		RootCmd:         rootCmd,
		Headings:        cc.HiCyan + cc.Bold + cc.Underline,
		Commands:        cc.HiYellow + cc.Bold,
		Example:         cc.Italic,
		ExecName:        cc.Bold,
		Flags:           cc.Bold,
		FlagsDataType:   cc.Italic + cc.HiBlue,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	// Clears temp files on each run.
	// It should not affect startup time since it's being run as goroutine.
	go func() {
		_ = filesystem.Api().RemoveAll(where.Temp())
	}()

	_ = rootCmd.Execute()
}

// handleErr will stop program execution and logger error to the stderr
// if err is not nil
func handleErr(err error) {
	if err == nil {
		return
	}

	log.Error(err)
	_, _ = fmt.Fprintf(
		os.Stderr,
		"%s %s\n",
		style.Failure(icon.Cross),
		strings.Trim(err.Error(), " \n"),
	)
	os.Exit(1)
}
