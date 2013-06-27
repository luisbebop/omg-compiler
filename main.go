package main

import (
	"net/http"
	"log"
	"os"
	"os/exec"
	"encoding/json"
	"io/ioutil"
	"bitbucket.org/kardianos/service"
	"fmt"
)

// file path and temp directory to store documents
var Tempdir = "c:\\omg-compiler\\tmp"
var PathToWalkCompiler2 = "c:\\omg-compiler\\WALK_Compiler2.exe"

// struct to share information with websocket clients
type CodeS struct {
	Type string
	Code string
	Output string
}

func main() {
	var name = "omg-compiler"
	var displayName = "omg-compiler Service"
	var desc = "This is a omg-compiler service."

	var s, err = service.NewService(name, displayName, desc)
	
	if err != nil {
		fmt.Printf("%s unable to start: %s", displayName, err)
		return
	}

	if len(os.Args) > 1 {
		var err error
		verb := os.Args[1]
		switch verb {
		case "install":
			err = s.Install()
			if err != nil {
				fmt.Printf("Failed to install: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" installed.\n", displayName)
		case "remove":
			err = s.Remove()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", displayName)
		case "run":
			doWork()
		case "start":
			err = s.Start()
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" started.\n", displayName)
		case "stop":
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" stopped.\n", displayName)
		}
		return
	}
	err = s.Run(func() error {
		// start
		go doWork()
		return nil
	}, func() error {
		// stop
		stopWork()
		return nil
	})
	if err != nil {
		s.Error(err.Error())
	}
}

func doWork() {
	f, _ := os.OpenFile("c:\\omg-compiler\\log.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0755)
	log.SetOutput(f)
	log.Printf("omg-compiler running!")
	log.Printf("omg-compiler listening http on port 80")
	log.Printf("omg-compiler S2 U")
	http.Handle("/compile", http.HandlerFunc(Compile))
	http.Handle("/", http.FileServer(http.Dir("c:\\omg-compiler\\html")))
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalf("omg-compiler listen error on http port 80: %v", err.Error())
	}
}

func stopWork() {
	log.Printf("omg-compiler stopping!")
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
	fDoc, _ := ioutil.TempFile(Tempdir, "doc-")
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
	log.Printf("omg-compiler compiling: %s", c.Code)
	out, err := exec.Command(PathToWalkCompiler2, "console", c.Code).Output()
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
