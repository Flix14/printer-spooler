package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	prt "github.com/alexbrainman/printer"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	version := r.Group("/v1")
	{
		version.GET("/print", printPdf)
	}

	// dbExample()

	err = r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
}

func printPdf(c *gin.Context) {
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
