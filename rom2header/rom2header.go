package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var romFilename string
var headerFilename string
var structName string = "PROGRAMS"
var programSize int = 128 * 4
var numPrograms int = 8
var useProgmem bool = true
var padEmpty bool = true

func parseCommandLineParameters() bool {
	flag.StringVar(&romFilename, "in", romFilename, "Input file")
	flag.StringVar(&headerFilename, "out", headerFilename, "Header file")
	flag.StringVar(&structName, "prefix", structName, "C-Struct name")
	flag.IntVar(&programSize, "programsize", programSize, "Size of each program in bytes")
	flag.IntVar(&numPrograms, "numprograms", numPrograms, "Number of programs")
	flag.BoolVar(&useProgmem, "progmem", useProgmem, "Use the PROGMEM statement")

	// FIXME: Parse structPrefix and chunkSize as well (20220902 handegar)

	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Printf("  Type \"./%s -help\" for more info.\n", filepath.Base(os.Args[0]))
		return false
	}

	if romFilename == "" {
		fmt.Println("  No ROM file specified. Use the '-in' parameter.")
		return false
	}

	if headerFilename == "" {
		fmt.Println("  No C-header filename specified. Use the '-out' parameter.")
		return false
	}

	return true
}

func main() {
	fmt.Printf("* ROM2Header convert utility v0.01\n")

	if !parseCommandLineParameters() {
		return
	}

	//
	// Read in binary file
	//

	fmt.Printf(" - Reading '%s'...\n", romFilename)
	file, err := os.Open(romFilename)

	if err != nil {
		fmt.Printf(" - Error opening ROM file: %s\n", err)
		return
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		fmt.Printf(" - Error stat'ing ROM file: %s\n", err)
		return
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	buf := bufio.NewReader(file)
	_, err = buf.Read(bytes)

	if err != nil {
		fmt.Printf(" - Error reading ROM file: %s\n", err)
		return
	}

	fmt.Printf("   - Read %d bytes\n", len(bytes))

	if len(bytes) < numPrograms*programSize && padEmpty {
		fmt.Printf("   - NOTE: Binary size is not %d*%d=%d bytes. Padding missing bytes with zeros.\n",
			numPrograms, programSize, numPrograms*programSize)
		numPad := numPrograms*programSize - len(bytes)
		for i := 0; i < numPad; i++ {
			bytes = append(bytes, 0)
		}
	}

	//
	// Split bytes into separate chunks
	//
	finished := false
	bytesIdx := 0
	var chunks [][]uint8
	for i := 0; i < numPrograms; i++ {
		chunks = append(chunks, make([]uint8, programSize))
		for j := 0; j < programSize; j++ {
			chunks[i][j] = bytes[bytesIdx]
			bytesIdx++
			if bytesIdx > len(bytes) {
				finished = true
				break
			}
		}

		if finished {
			break
		}
	}

	//
	// Convert chunks to a C-Header style text-buffer
	//
	var values [][]string
	for i := 0; i < len(chunks); i++ {
		values = append(values, make([]string, len(chunks[i])))
		for j := 0; j < len(chunks[i]); j++ {
			values[i][j] = fmt.Sprintf(" 0x%.2x", chunks[i][j])
		}
	}

	var programs []string
	for i := 0; i < len(values); i++ {
		programs = append(programs, " {\n  "+strings.Join(values[i], ",")+"\n }")
	}

	var cHeaderText string = fmt.Sprintf(`/*
*  Rom2Header.go
*   Origin file: '%s'
*   Program size: %d
*   Number of programs: %d
*   Total data size: %d bytes
*/
`, romFilename, programSize, len(programs), len(programs)*programSize)

	progmem := ""
	if useProgmem {
		progmem = "PROGMEM"
	}

	cHeaderText += fmt.Sprintf("const unsigned int NUM_PROGRAMS = %d;\n", len(programs))
	cHeaderText += fmt.Sprintf("const unsigned char %s[][%d] %s = {\n", structName, programSize, progmem)

	cHeaderText += strings.Join(programs, ",\n")
	cHeaderText += "\n};\n"

	//
	// Write text buffer to file
	//
	outfile, err := os.Create(headerFilename)
	if err != nil {
		fmt.Printf(" - Error creating c-header file: %s\n", err)
		return
	}

	outfile.WriteString(cHeaderText)
	fmt.Printf(" - Written %d programs to '%s'\n", len(chunks), headerFilename)

}
