package main

import (
	"net/http"
	"log"
	"os/exec"
	"encoding/json"
	"io/ioutil"
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
	http.Handle("/compile", http.HandlerFunc(Compile))
	http.Handle("/", http.FileServer(http.Dir("./html")))
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalf("omg-compiler listen error on http port 80: %v", err.Error())
	}
}

// handling a POST connection with a variable named doc
func Compile(w http.ResponseWriter, req *http.Request) {
	var code CodeS
	doc := req.FormValue("doc")
	err := json.Unmarshal([]byte(doc), &code)
	if err != nil {
		log.Printf("json.Unmarshal err = %s string = %s", err, doc)
		return
	}
	
	//writing to a tempfile
	fDoc, _ := ioutil.TempFile("tmp", "doc-")
	fDoc.WriteString(code.Code)
	fDoc.Close()
	
	//setting the code to filename
	code.Code = fDoc.Name()
	
	//compiling the posxml code
	if code.Type == "posxml" {
		err = CallWalkCompiler2(&code)
		if err != nil {
			log.Printf("CallWalkCompiler err = %s", err)
			return
		}
	}
	
	//echo
	if code.Type == "echo" {
		err = CallEcho(&code)
		if err != nil {
			log.Printf("CallEcho err = %s", err)
			return
		}
	}
	
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(code.Output))
}

func CallWalkCompiler2(c *CodeS) (error) {
	out, err := exec.Command("WALK_Compiler2.exe", "console", c.Code).Output()
	if err != nil {
		return err
	}
	c.Output = string(out)
	return nil
}

func CallEcho(c * CodeS) (error) {
	out, err := exec.Command("cat", c.Code).Output()
	if err != nil {
		return err
	}
	c.Output = string(out)
	return nil
}
