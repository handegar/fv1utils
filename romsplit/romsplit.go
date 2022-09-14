package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var bankFilename string = ""
var prgBinFilenames [8]string = [8]string{
	"", "", "", "",
	"", "", "", "",
} // Explicit program filenames 1..8
var filenamePattern string = "prg-%d.bin"
var numberOfPrograms int = 8
var programSize int = 128 * 4

func parseCommandLineParameters() bool {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"  Extract binary programs from a ROM/Bank file\n")
		fmt.Fprintf(os.Stderr, `
 Usage:
  $ %s -in <bankname.bin> [-pX <filename.bin> ...] [-pattern <string>]

 NOTE:
  Use either '-pattern' or '-pX' to name the output files.
  The '-pattern' parameter MUST contain a '%%d' part for the numbering.
`,
			os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")

		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
	}

	flag.StringVar(&filenamePattern, "pattern", filenamePattern, "Pattern for output filenames")

	flag.StringVar(&bankFilename, "in", bankFilename, "ROM/bank filename")

	flag.StringVar(&prgBinFilenames[0], "p1", prgBinFilenames[0], "Filename for Program #1")
	flag.StringVar(&prgBinFilenames[1], "p2", prgBinFilenames[1], "Filename for Program #2")
	flag.StringVar(&prgBinFilenames[2], "p3", prgBinFilenames[2], "Filename for Program #3")
	flag.StringVar(&prgBinFilenames[3], "p4", prgBinFilenames[3], "Filename for Program #4")
	flag.StringVar(&prgBinFilenames[4], "p5", prgBinFilenames[4], "Filename for Program #5")
	flag.StringVar(&prgBinFilenames[5], "p6", prgBinFilenames[5], "Filename for Program #6")
	flag.StringVar(&prgBinFilenames[6], "p7", prgBinFilenames[6], "Filename for Program #7")
	flag.StringVar(&prgBinFilenames[7], "p8", prgBinFilenames[7], "Filename for Program #8")

	flag.IntVar(&numberOfPrograms, "num", numberOfPrograms, "Number of programs in ROM/Bank file")
	flag.IntVar(&programSize, "program-size", programSize, "Size of each program (bytes)")

	flag.Parse()

	// Generate all output filename from "filenamePattern"
	for i := 0; i < numberOfPrograms; i += 1 {
		prgBinFilenames[i] = fmt.Sprintf(filenamePattern, i+1)
	}

	if flag.NFlag() == 0 {
		fmt.Printf("  Type \"./%s -help\" for more info.\n", filepath.Base(os.Args[0]))
		return false
	}

	if bankFilename == "" {
		fmt.Println("- No ROM/Bank filename specified. Use the '-in' parameter.")
		return false
	}

	if strings.Index(filenamePattern, "%d") == -1 {
		fmt.Printf("- The '-pattern' argument must contain a '%%d' for the numbering.\n")
		return false
	}

	return true
}

func extractProgram(romFile *os.File, idx int) ([]byte, error) {
	bytes := make([]uint8, programSize)
	_, err := romFile.ReadAt(bytes, int64(idx*programSize))

	allZero := true
	for _, b := range bytes {
		if b != 0 {
			allZero = false
		}
	}

	if allZero {
		fmt.Printf("- NOTE: Program #%d is all zero values.\n", idx)
	}

	return bytes, err
}

func main() {
	fmt.Printf("* ROMSplit utility v0.01\n")

	if !parseCommandLineParameters() {
		return
	}

	bankFile, err := os.Open(bankFilename)
	defer bankFile.Close()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	for i := 0; i < numberOfPrograms; i += 1 {
		bytes, err := extractProgram(bankFile, i)
		if err != nil {
			break
		}

		fmt.Printf("- Writing '%s'...\n", prgBinFilenames[i])
		err = ioutil.WriteFile(prgBinFilenames[i], bytes, 0666)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}
	}

}
