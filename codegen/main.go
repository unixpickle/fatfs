// Command codegen generates a file called spec.go with
// specifications for each available environment.
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/unixpickle/essentials"
)

func main() {
	generateDataType("codegen/boot_sector.csv", "boot_sector_gen.go", "BootSector", "Sector")
	generateDataType("codegen/dir_entry.csv", "dir_entry_gen.go", "DirEntry", "[32]byte")
}

func generateDataType(csvFile, outputFile, className, definition string) {
	file, err := os.Open(csvFile)
	essentials.Must(err)
	defer file.Close()
	fields, err := csv.NewReader(file).ReadAll()
	essentials.Must(err)

	outFile, err := os.Create(outputFile)
	essentials.Must(err)
	defer outFile.Close()
	outFile.Write([]byte("package fatfs\n\ntype " + className + " " + definition + "\n"))

	classLetter := strings.ToLower(className[:1])
	for _, fieldInfo := range fields {
		name := fieldInfo[0]
		start := fieldInfo[1]
		size := fieldInfo[2]
		funcPrefix := ""
		if size == "1" || size == "2" || size == "4" {
			funcPrefix = "Raw"
		}
		outFile.Write([]byte(fmt.Sprintf("\nfunc ("+classLetter+" *"+className+") %s%s() []byte {\n"+
			"\treturn "+classLetter+"[%s : %s+%s]\n}\n", funcPrefix, name, start, start, size)))
		if size == "1" {
			outFile.Write([]byte(fmt.Sprintf("\nfunc ("+classLetter+" *"+className+") %s() uint8 {\n"+
				"\treturn "+classLetter+"[%s]\n}\n", name, start)))
			outFile.Write([]byte(fmt.Sprintf("\nfunc ("+classLetter+" *"+className+") Set%s(x uint8) {\n"+
				"\t"+classLetter+"[%s] = x\n}\n", name, start)))
		} else if size == "2" {
			outFile.Write([]byte(fmt.Sprintf("\nfunc ("+classLetter+" *"+className+") %s() uint16 {\n"+
				"\treturn uint16("+classLetter+"[%s]) | (uint16("+classLetter+"[%s+1]) << 8)\n}\n", name, start, start)))
			outFile.Write([]byte(fmt.Sprintf("\nfunc ("+classLetter+" *"+className+") Set%s(x uint16) {\n"+
				"\t"+classLetter+"[%s] = uint8(x)\n\t"+classLetter+"[%s+1] = uint8(x >> 8)\n}\n", name, start, start)))
		} else if size == "4" {
			outFile.Write([]byte(fmt.Sprintf("\nfunc ("+classLetter+" *"+className+") %s() uint32 {\n"+
				"\treturn uint32("+classLetter+"[%s]) | (uint32("+classLetter+"[%s+1]) << 8) |\n"+
				"\t\t(uint32("+classLetter+"[%s+2]) << 16) | (uint32("+classLetter+"[%s+3]) << 24)\n}\n",
				name, start, start, start, start)))
			outFile.Write([]byte(fmt.Sprintf("\nfunc ("+classLetter+" *"+className+") Set%s(x uint32) {\n"+
				"\t"+classLetter+"[%s] = uint8(x)\n\t"+classLetter+"[%s+1] = uint8(x >> 8)\n"+
				"\t"+classLetter+"[%s+2] = uint8(x >> 16)\n\t"+classLetter+"[%s+3] = uint8(x >> 24)\n}\n",
				name, start, start, start, start)))
		}
	}
}
