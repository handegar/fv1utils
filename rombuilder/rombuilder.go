package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const PROGRAM_SIZE int = 128 * 4

var bankFilename string = ""
var argsBinFilenames []string // Filenames souped up from the args
var prgBinFilenames [8]string = [8]string{
	"", "", "", "",
	"", "", "", "",
} // Explicit program filenames 1..8

func parseCommandLineParameters() bool {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"  Concatenate up to 8 binary programs into an uploadable EPROM binary bank\n")
		fmt.Fprintf(os.Stderr,
			"\n  Usage:\n"+
				"  $ %s -out <bankname.bin> [-pX <filename.bin> ...] [prg1.bin prg2.bin ... prg8.bin]\n",
			os.Args[0])
		fmt.Fprintf(os.Stderr,
			"    NOTE: Use bin-filename 'BLANK' to insert an empty program.\n\n")

		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
	}

	flag.StringVar(&bankFilename, "out", bankFilename, "Target bank file (bin)")

	flag.StringVar(&prgBinFilenames[0], "p1", prgBinFilenames[0], "Explicit bin file for Program #1")
	flag.StringVar(&prgBinFilenames[1], "p2", prgBinFilenames[1], "Explicit bin file for Program #2")
	flag.StringVar(&prgBinFilenames[2], "p3", prgBinFilenames[2], "Explicit bin file for Program #3")
	flag.StringVar(&prgBinFilenames[3], "p4", prgBinFilenames[3], "Explicit bin file for Program #4")
	flag.StringVar(&prgBinFilenames[4], "p5", prgBinFilenames[4], "Explicit bin file for Program #5")
	flag.StringVar(&prgBinFilenames[5], "p6", prgBinFilenames[5], "Explicit bin file for Program #6")
	flag.StringVar(&prgBinFilenames[6], "p7", prgBinFilenames[6], "Explicit bin file for Program #7")
	flag.StringVar(&prgBinFilenames[7], "p8", prgBinFilenames[7], "Explicit bin file for Program #8")

	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Printf("  Type \"./%s -help\" for more info.\n", filepath.Base(os.Args[0]))
		return false
	}

	argsBinFilenames = flag.Args()

	if len(argsBinFilenames) > 8 {
		fmt.Printf("- No more than 8 binary files can be concatenated into a EEPROM file (%d files was provided)\n",
			len(argsBinFilenames))
		return false
	}

	if len(argsBinFilenames) < 8 {
		missing := 8 - len(argsBinFilenames)
		appendix := make([]string, missing)
		argsBinFilenames = append(argsBinFilenames, appendix...)
	}

	if bankFilename == "" {
		fmt.Println("- No target filename specified. Use the '-out' parameter.")
		return false
	}

	// Override ARGS with explicit declared bin filenames?
	for i, fn := range prgBinFilenames {
		if fn != "" {
			argsBinFilenames[i] = fn
		}
	}

	return true
}

func main() {
	fmt.Printf("* ROMBuilder utility v0.01\n")

	if !parseCommandLineParameters() {
		return
	}

	bytes := make([]uint8, 8*PROGRAM_SIZE)
	idx := 0

	for nr, fn := range argsBinFilenames {
		var b []byte
		var err error
		if fn == "BLANK" || fn == "" {
			fmt.Printf(" - #%d) <Blank program>\n", nr+1)
			b = make([]byte, PROGRAM_SIZE)
			for i := 0; i < PROGRAM_SIZE; i += 4 {
				b[i] = 0x11 // NOP instruction
			}
		} else {
			fmt.Printf(" - #%d) Loading '%s'\n", nr+1, fn)
			b, err = ioutil.ReadFile(fn) // b has type []byte
			if err != nil {
				fmt.Printf("   ERROR reading file: %s\n", err)
				return
			}
		}

		if idx+len(b) > len(bytes) {
			fmt.Printf("   ERROR total size of bin-files exceed %d bytes\n", len(bytes))
			return
		}

		for i := 0; i < len(b); i++ {
			bytes[idx] = b[i]
			idx++
		}
	}

	fmt.Printf(" - Writing bank to '%s'\n", bankFilename)
	err := ioutil.WriteFile(bankFilename, bytes, 0666)
	if err != nil {
		fmt.Printf("   ERROR writing file: %s\n", err)
		return
	}
}
