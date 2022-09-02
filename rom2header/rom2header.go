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
var programSize int = 512
var numPrograms int = 8

func parseCommandLineParameters() bool {
	flag.StringVar(&romFilename, "in", romFilename, "Input file")
	flag.StringVar(&headerFilename, "out", headerFilename, "Header file")
	flag.StringVar(&structName, "prefix", structName, "C-Struct name")
	flag.IntVar(&programSize, "programsize", programSize, "Size of each program in bytes")
	flag.IntVar(&numPrograms, "numprograms", numPrograms, "Number of programs")

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
	fmt.Printf("* Rom2Header convert utility v0.01\n")

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

	//
	// Split bytes into separate chunks
	//

	finished := false
	bytesIdx := 0
	chunkIdx := 0
	var chunks [][]uint8
	for i := 0; i < numPrograms; i++ {
		chunks = append(chunks, make([]uint8, programSize))
		for j := 0; j < programSize; j++ {
			chunks[chunkIdx][j] = bytes[bytesIdx]
			bytesIdx++
			if bytesIdx >= len(bytes) {
				finished = true
				break
			}
		}

		if finished {
			break
		}
	}

	if finished {
		fmt.Printf(" - NOTE: Finished before %d*%d bytes had been processed.\n",
			programSize, numPrograms)
	}

	//
	// Convert chunks to a C-Header style text-buffer
	//

	var values [][]string
	for i := 0; i < len(chunks); i++ {
		values = append(values, make([]string, programSize))
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
*/
`, romFilename, programSize, len(programs))

	cHeaderText += fmt.Sprintf("const unsigned int NUM_PROGRAMS = %d;\n", len(programs))
	cHeaderText += fmt.Sprintf("const unsigned char %s[][] = {\n", structName)

	cHeaderText += strings.Join(programs, ",\n")
	cHeaderText += "\n}\n"

	//
	// Write text buffer to file
	//
	outfile, err := os.Create(headerFilename)
	if err != nil {
		fmt.Printf(" - Error creating c-header file: %s\n", err)
		return
	}

	outfile.WriteString(cHeaderText)
	fmt.Printf(" - Written to '%s'\n", headerFilename)

}
