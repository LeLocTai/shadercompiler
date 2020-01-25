package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const fxcPath = "C:/Program Files (x86)/Windows Kits/10/bin/10.0.17763.0/x64/fxc.exe"

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/compile", compile)

	fmt.Println("Server listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type CompileRequest struct {
	Code string
}

type CompileResult struct {
	Err string
	Asm string
}

func compile(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var compileReq CompileRequest
	decodeErr := decoder.Decode(&compileReq)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusInternalServerError)
		return
	}

	_ = os.MkdirAll("./hlsl", os.ModeDir)
	tmpFile, writeErr := ioutil.TempFile("./hlsl", "code-*.hlsl")
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		return
	}

	defer os.Remove(tmpFile.Name())

	_, writeErr = tmpFile.WriteString(compileReq.Code)
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		return
	}

	closeErr := tmpFile.Close()
	if closeErr != nil {
		http.Error(w, closeErr.Error(), http.StatusInternalServerError)
		return
	}

	cmd := exec.Command(
		fxcPath,
		"/T", "ps_5_0",
		"/E", "PSMain",
		//"/Zi",
		tmpFile.Name())
	stdout, cmdErr := cmd.CombinedOutput()
	cmdErrString := ""
	if cmdErr != nil {
		cmdErrString = cmdErr.Error()
	}
	result := CompileResult{cmdErrString, string(stdout)}

	res, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
