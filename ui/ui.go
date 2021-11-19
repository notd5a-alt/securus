package secui

import (
	"time"

	"github.com/pterm/pterm"
)

// prints the securus title
func Title() {
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("secur", pterm.NewStyle(pterm.FgYellow)),
		pterm.NewLettersFromStringWithStyle("us", pterm.NewStyle(pterm.FgCyan))).
		Render()
}

// Function to print out a debug string
func PrintDebug(debug string) {
	pterm.EnableDebugMessages()
	pterm.Debug.Println(debug)
}

// function to printout information for the user
func PrintInfo(info string) {
	pterm.Info.Println(info)
}

// function to print warning to console
func PrintWarn(warning string) {
	pterm.Warning.Println(warning)
}

// function to print an error to screen
func PrintError(err error) {
	pterm.Error.Println(err)
}

// not implementing a wrapper for FATAL errors since ill need them inline

// function to print success
func PrintSuccess(success string) {
	pterm.Success.Println(success)
}

// Sections printing
func PrintSection(text string) {
	pterm.DefaultSection.Println(text)
}

func PrintSectionWithLevel(text string, level int) {
	pterm.DefaultSection.WithLevel(level).Println(text)
}

func PrintInputPrefix(prefix string) {
	prefixStyle := pterm.NewStyle(pterm.FgCyan)
	prefixStyle.Print(prefix)
}

// Simualtes 2 seconds of loading time
func PrintSpinnerSuccess(message string) {
	spinner, _ := pterm.DefaultSpinner.Start(message)
	time.Sleep(time.Second * 2)
	spinner.Success()
}

func PrintSpinnerWarning(message string) {
	spinner, _ := pterm.DefaultSpinner.Start(message)
	time.Sleep(time.Second * 2)
	spinner.Warning()
}

func PrintSpinnerFail(message string) {
	spinner, _ := pterm.DefaultSpinner.Start(message)
	time.Sleep(time.Second * 2)
	spinner.Fail()
}

// initial wrappers for pterm created
