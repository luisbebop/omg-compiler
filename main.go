package main

import (
	"net/http"
	"log"
	"code.google.com/p/go.net/websocket"
	"os/exec"
)

// struct to share information with websocket clients
type CodeS struct {
	Type string
	Code string
	Output string
}

func main () {
	log.Printf("omg-compiler listening websocket and http on port 80")
	log.Printf("omg-compiler S2 U")
	http.Handle("/compile", websocket.Handler(Compile))
	http.Handle("/", http.FileServer(http.Dir("./html")))
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalf("omg-compiler listen error on http port 31415: %v", err.Error())
	}
}

// handle a websocket connection received from webserver
func Compile(ws *websocket.Conn) {
	log.Printf("compile %#v\n", ws.Config())
	for {
		var code CodeS
		
		// receives a text message serialized CodeS as JSON.
		err := websocket.JSON.Receive(ws, &code)
		if err != nil {
			log.Printf("%s", err)
			break
		}
		log.Printf("recv:%#v\n", code)
		
		// compiling the posxml code
		if code.Type == "posxml" {
			err = CallWalkCompiler2(&code)
			if err != nil {
				log.Printf("%s", err)
				break
			}
		}
		
		// send to the webclients
		err = websocket.JSON.Send(ws, code)
		if err != nil {
			log.Printf("%s", err)
			break
		}
		log.Printf("send:%#v\n", code)
	}
}

func CallWalkCompiler2(c *CodeS) (error){
	out, err := exec.Command("WALK_Compiler2.exe", "console", c.Code).Output()
	if err != nil {
		return err
	}
	c.Output = string(out)
	return nil
}