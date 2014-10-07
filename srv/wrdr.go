package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
//	"runtime"
	"strconv"
	"time"
)

func uploadFileChunkSrvr(w http.ResponseWriter, req *http.Request) {

	rmsg := "{"
	vstart := time.Now()

	// Read headers
	start := time.Now()
	fname := req.Header.Get("X-File-Name")
	if fname == "" {
		fmt.Println("X-File-Name not found in headers")
		http.Error(w, "X-File-Name not found in headers", 502)
		return
	}
	offset := req.Header.Get("X-File-Offset")
	if offset == "" {
		fmt.Println("X-File-Offset not found in headers")
		http.Error(w, "X-File-Offset not found in headers", 502)
		return
	}
	noffset, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "X-File-Offset invalid:"+offset, 502)
		return
	}
	dur := time.Since(start)
	rmsg = rmsg + "'Header-Parsing-Time' : '" + strconv.FormatFloat(dur.Seconds(), 'f',7,64) + "', "

    //Read the body
	start = time.Now()
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not read file data", 502)
		return
	}
	dur = time.Since(start)
	rmsg = rmsg + "'Body-Size' : '" + strconv.FormatInt(int64(len(data)),10) + "', "
	rmsg = rmsg + "'Body-Reading-Time' : '" + strconv.FormatFloat(dur.Seconds(), 'f',7,64) + "', "

    //Open file
	start = time.Now()
	fpath := path.Join("/tmp", "chunked-"+fname)
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not open file", 502)
		return
	}
	defer f.Close()
	dur = time.Since(start)
	rmsg = rmsg + "'File-Opening-Time' : '" + strconv.FormatFloat(dur.Seconds(), 'f',7,64) + "', "

	//Write body to file
	start = time.Now()
	_, errx := f.WriteAt(data, noffset)
	if errx != nil {
		fmt.Println(err)
		http.Error(w, "Could not write to file", 502)
		return
	}
	dur = time.Since(start)
	rmsg = rmsg + "'File-Writing-Time' : '" + strconv.FormatFloat(dur.Seconds(), 'f',7,64) + "', "

	//Wrap up...
	dur = time.Since(vstart)
	rmsg = rmsg + "'Total' : '" + strconv.FormatFloat(dur.Seconds(), 'f',7,64) + "', "
	rmsg = rmsg + "}"
	w.WriteHeader(200)
	w.Write([]byte(rmsg))

}
/*
func uploadFileChunkSrvr(w http.ResponseWriter, req *http.Request) {
	r := make(chan int)
	go uploadFileChunkSrvrDo(w, req, r)
	<-r
}
*/

func uploadSrvr(w http.ResponseWriter, req *http.Request) {
	file, handler, err := req.FormFile("file")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not find file element", 502)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not read data", 502)
	}
	err = ioutil.WriteFile("/tmp/"+handler.Filename, data, 0777)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Could not open target file", 502)
	} else {
	    // w.WriteHeader(200)
	    http.Redirect(w, req, "/wrdr/freader.html", http.StatusFound)
	}
}

func main() {
	//runtime.GOMAXPROCS(8)
	http.HandleFunc("/wrdr/upload", uploadSrvr)
	http.HandleFunc("/wrdr/uploadChunk", uploadFileChunkSrvr)
	http.HandleFunc("/wrdr/uploadBlock", uploadFileChunkSrvr)
	http.Handle("/wrdr/", http.StripPrefix("/wrdr/", http.FileServer(http.Dir("/home/ubuntu/src/wrdr"))))
	err := http.ListenAndServe(":80", nil)
	log.Printf("%v", err)
}

