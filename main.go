package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	prt "github.com/alexbrainman/printer"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	var re = regexp.MustCompile(os.Getenv("PRINTER"))
	printers, err := prt.ReadNames()
	if err != nil {
		panic(err)
	}

	fmt.Println("PRINTERS:")
	fmt.Println(printers) //Imprime en consola los nombres de las impresoras disponibles

	printer := ""
	for _, p := range printers {
		if re.MatchString(p) {
			printer = p
		}
	}
	if printer != "" {
		log.Println("Printer found ", printer)
		content, err := ioutil.ReadFile(os.Getenv("FILE"))
		if err != nil {
			panic(err)
		}
		printContent("raw", content)
	} else {
		log.Println("Printer not found")
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func printContent(datatype string, content []byte) {
	name, err := prt.Default() // Regresa el nombre de la impresora por default
	checkErr(err)
	fmt.Println(name)

	p, err := prt.Open(name) // Abre la impresora y regresa su struct
	checkErr(err)

	err = p.StartDocument("test", datatype)
	checkErr(err)

	err = p.StartPage()
	checkErr(err)

	n, err := p.Write(content) //Enviar documento para impresión
	checkErr(err)
	fmt.Println("Número de bytes:", n)

	err = p.EndPage()
	checkErr(err)

	err = p.EndDocument()
	checkErr(err)

	err = p.Close()
	checkErr(err)

}
