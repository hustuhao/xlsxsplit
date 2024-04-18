package cmd

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/log"

	"github.com/charmbracelet/lipgloss"

	"github.com/metafates/xlsxsplit/app"
	"github.com/metafates/xlsxsplit/color"
	"github.com/samber/lo"

	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

var (
	// Used for flags.
	outputDir string
	file      string
)

func init() {
	rootCmd.AddCommand(splitCmd)
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolP("short", "s", false, "print the version number only")
	splitCmd.Flags().StringVarP(&file, "file", "f", "", "target xlsx file")
	splitCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")

}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the " + app.Name,
	Run: func(cmd *cobra.Command, args []string) {
		if lo.Must(cmd.Flags().GetBool("short")) {
			_, err := cmd.OutOrStdout().Write([]byte(app.Version + "\n"))
			handleErr(err)
			return
		}

		versionInfo := struct {
			Version  string
			OS       string
			Arch     string
			App      string
			Compiler string
		}{
			Version:  app.Version,
			App:      app.Name,
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			Compiler: runtime.Compiler,
		}

		t, err := template.New("version").Funcs(map[string]any{
			"faint":   lipgloss.NewStyle().Faint(true).Render,
			"bold":    lipgloss.NewStyle().Bold(true).Render,
			"magenta": lipgloss.NewStyle().Foreground(color.Purple).Render,
		}).Parse(`{{ magenta "▇▇▇" }} {{ magenta .App }} 

  {{ faint "Version" }}  {{ bold .Version }}
  {{ faint "Platform" }} {{ bold .OS }}/{{ bold .Arch }}
  {{ faint "Compiler" }} {{ bold .Compiler }}
`)
		handleErr(err)
		handleErr(t.Execute(cmd.OutOrStdout(), versionInfo))
	},
}

var splitCmd = &cobra.Command{
	Use:     "split",
	Short:   "Split xlsx file into multiple files",
	Example: "splitxlsx -f examples.xlsx",
	Run: func(cmd *cobra.Command, args []string) {
		xlsxFile := file
		// Open the xlsx file
		f, err := excelize.OpenFile(xlsxFile)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}

		// Iterate through each sheet
		for _, sheetName := range f.GetSheetList() {
			log.Debugf("Splitting sheet: %s", sheetName)
			// Create a new xlsx file
			newFileName := sheetName + ".xlsx"
			newFile := excelize.NewFile()
			// 创建新sheet
			if err := newFile.SetSheetName("Sheet1", sheetName); err != nil {
				log.Fatalf("Error creating sheet: %v", err)
			}

			// Copy the content of the current sheet into the new file
			rows, err := f.GetRows(sheetName)
			if err != nil {
				log.Fatalf("Error getting rows from sheet: %v", err)
			}
			newFile.Styles = f.Styles
			// Retrieve all comments
			comments, err := f.GetComments(sheetName)
			if err != nil {
				log.Fatalf("Error getting comments from sheet: %v", err)
			}
			commentMap := make(map[string]excelize.Comment, len(comments))
			for i := range comments {
				commentMap[comments[i].Cell] = comments[i]
			}

			for rowIndex, row := range rows {
				log.Debugf("Row %d", rowIndex)
				h, err := f.GetRowHeight(sheetName, rowIndex+1)
				if err != nil {
					log.Fatalf("Error getting row height: %v", err)
				}
				err = newFile.SetRowHeight(sheetName, rowIndex+1, h)
				if err != nil {
					log.Fatalf("Error setting row height: %v", err)
				}

				for colIndex, cellValue := range row {
					cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
					if err != nil {
						log.Fatalf("Error converting coordinates: %v", err)
					}
					colName, err := excelize.ColumnNumberToName(colIndex + 1)
					if err != nil {
						log.Fatalf("Error converting column number to name: %v", err)
					}

					// Set the dimensions of each cell
					width, err := f.GetColWidth(sheetName, colName)
					if err != nil {
						log.Fatalf("Error getting column width: %v", err)
					}
					err = newFile.SetColWidth(sheetName, colName, colName, width)
					if err != nil {
						log.Fatalf("Error setting column width: %v", err)
					}

					// Set the value of a cell
					if err := newFile.SetCellValue(sheetName, cellName, cellValue); err != nil {
						log.Fatalf("Error setting cell value: %v", err)
					}

					// Copy formulas
					formula, err := f.GetCellFormula(sheetName, cellName)
					if err != nil {
						log.Fatalf("Error getting cell formula: %v", err)
					}
					if err := newFile.SetCellFormula(sheetName, cellName, formula); err != nil {
						log.Fatalf("Error setting cell formula: %v", err)
					}

					// Copy comments
					if comment, ok := commentMap[cellName]; ok {
						if err := newFile.AddComment(sheetName, comment); err != nil {
							log.Fatalf("Error setting cell comment: %v", err)
						}
					}

					// Copy cell formatting
					styleID, err := f.GetCellStyle(sheetName, cellName)
					if err != nil {
						log.Fatalf("Error getting cell style: %v", err)
					}

					// Apply background color to cells in the new file
					fillColor, err := f.GetStyle(styleID)
					if err != nil || fillColor == nil {
						log.Fatalf("Error getting cell background color: %v", err)
					}
					if err := newFile.SetCellStyle(sheetName, cellName, cellName, styleID); err != nil {
						log.Fatalf("Error setting cell style: %v", err)
					}

				}
			}

			if err := saveNewFile(newFile, newFileName); err != nil {
				log.Fatalf("Error saving new file: %v", err)
			}
			log.Infof("Sheet '%s' has been saved as '%s'\n", sheetName, newFileName)
		}
		log.Info("Finished!'\n")
	},
}

func saveNewFile(file *excelize.File, fileName string) error {
	outputFolder := outputDir
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		if err := os.Mkdir(outputFolder, 0755); err != nil {
			return fmt.Errorf("error creating output folder: %v", err)
		}
	}

	// 保存新文件
	outputPath := filepath.Join(outputFolder, fileName)
	if err := file.SaveAs(outputPath); err != nil {
		return fmt.Errorf("error saving new file: %v", err)
	}
	return nil
}
