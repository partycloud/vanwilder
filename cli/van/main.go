package main

import (
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gorilla/websocket"
)

var (
	sockFile = "/var/run/van.sock"
)

var AppHelpTemplate = `Usage:
  {{.Name}} <command> [arguments...]

Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Use "{{.Name}} <command> --help" for more information about a command.
`

var CommandHelpTemplate = `{{.Description}}

Options:
   {{range .Flags}}{{.}}
   {{end}}
`

func main() {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate

	app := cli.NewApp()
	app.Name = "van"
	app.Commands = []cli.Command{
		VanStartCommand,
	}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	app.Run(os.Args)

	// // websocket
	// conn, err := dialWebSocket("ws://api.partycloud.dev:3000/vw?token1234")
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println("connection established")
	//
	// for {
	// 	_, r, err := conn.NextReader()
	//
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	evtString, err := ioutil.ReadAll(r)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	fmt.Println(string(evtString))
	// }
}

func dialWebSocket(url string) (*websocket.Conn, error) {
	var dialer *websocket.Dialer
	var dialHeader http.Header

	conn, _, err := dialer.Dial(url, dialHeader)

	return conn, err
}
