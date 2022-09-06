package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var bankFilename string = ""
var binFilenames []string

func parseCommandLineParameters() bool {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"\nUsage of ROMBuilder:\n"+
				" $ %s -out <bankname.bin> [prg1.bin prg2.bin ... prg8.bin]\n",
			os.Args[0])

		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
	}

	flag.StringVar(&bankFilename, "out", bankFilename, "Target bank file (bin)")
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Printf("  Type \"./%s -help\" for more info.\n", filepath.Base(os.Args[0]))
		return false
	}

	binFilenames = flag.Args()

	if len(binFilenames) > 8 {
		fmt.Println("  No more than 8 binary files can be concatenated into a EEPROM file.")
		return false
	}

	if bankFilename == "" {
		fmt.Println("  No target filename specified. Use the '-out' parameter.")
		return false
	}

	return true
}

func main() {
	fmt.Printf("* ROMBuilder utility v0.01\n")
	fmt.Printf("  Concatenate up to 8 binary programs into an uploadable EPROM binary bank\n")

	if !parseCommandLineParameters() {
		return
	}

	bytes := make([]uint8, 8*512)
	idx := 0

	for nr, fn := range binFilenames {
		fmt.Printf(" - #%d) Loading '%s'\n", nr, fn)
		b, err := ioutil.ReadFile(fn) // b has type []byte
		if err != nil {
			fmt.Printf("   ERROR reading file: %s\n", err)
			return
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
	err := ioutil.WriteFile(bankFilename, bytes, 777)
	if err != nil {
		fmt.Printf("   ERROR writing file: %s\n", err)
		return
	}
}
