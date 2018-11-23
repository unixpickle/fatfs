// Command codegen generates a file called spec.go with
// specifications for each available environment.
package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/unixpickle/essentials"
)

func main() {
	file, err := os.Open("codegen/boot_sector.csv")
	essentials.Must(err)
	defer file.Close()
	fields, err := csv.NewReader(file).ReadAll()
	essentials.Must(err)

	outFile, err := os.Create("boot_sector.go")
	essentials.Must(err)
	defer outFile.Close()
	outFile.Write([]byte("package fatfs\n\ntype BootSector [512]byte\n"))
	for _, fieldInfo := range fields {
		name := fieldInfo[0]
		start := fieldInfo[1]
		size := fieldInfo[2]
		funcPrefix := ""
		if size == "1" || size == "2" || size == "4" {
			funcPrefix = "Raw"
		}
		outFile.Write([]byte(fmt.Sprintf("\nfunc (b *BootSector) %s%s() []byte {\n"+
			"\treturn b[%s : %s+%s]\n}\n", funcPrefix, name, start, start, size)))
		if size == "1" {
			outFile.Write([]byte(fmt.Sprintf("\nfunc (b *BootSector) %s() uint8 {\n"+
				"\treturn b[%s]\n}\n", name, start)))
			outFile.Write([]byte(fmt.Sprintf("\nfunc (b *BootSector) Set%s(x uint8) {\n"+
				"\tb[%s] = x\n}\n", name, start)))
		} else if size == "2" {
			outFile.Write([]byte(fmt.Sprintf("\nfunc (b *BootSector) %s() uint16 {\n"+
				"\treturn uint16(b[%s]) | (uint16(b[%s+1]) << 8)\n}\n", name, start, start)))
			outFile.Write([]byte(fmt.Sprintf("\nfunc (b *BootSector) Set%s(x uint16) {\n"+
				"\tb[%s] = uint8(x)\n\tb[%s+1] = uint8(x >> 8)\n}\n", name, start, start)))
		} else if size == "4" {
			outFile.Write([]byte(fmt.Sprintf("\nfunc (b *BootSector) %s() uint32 {\n"+
				"\treturn uint32(b[%s]) | (uint32(b[%s+1]) << 8) |\n"+
				"\t\t(uint32(b[%s+2]) << 16) | (uint32(b[%s+3]) << 24)\n}\n",
				name, start, start, start, start)))
			outFile.Write([]byte(fmt.Sprintf("\nfunc (b *BootSector) Set%s(x uint32) {\n"+
				"\tb[%s] = uint8(x)\n\tb[%s+1] = uint8(x >> 8)\n"+
				"\tb[%s+2] = uint8(x >> 16)\n\tb[%s+3] = uint8(x >> 24)\n}\n",
				name, start, start, start, start)))
		}
	}
}
