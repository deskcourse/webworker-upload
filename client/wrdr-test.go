package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"bytes"
)

func rateMsg(action string, d time.Duration, bytes int64) (string){
	r := float64(999999999.0)
	ds := d.Seconds()
	if (d > 0) {
		r = (float64(bytes*8) / (1000*1000))/ ds;
	}
	msg := action + ":" + strconv.FormatFloat(r, 'f', 3, 64) + " Mbits/sec";
	msg = msg + " (" + strconv.FormatFloat(ds, 'f',7,64) + "s)"
	return msg;
}

func postChunk(file *os.File, b []byte, fname string, pos int64, rlen int64, target string, rch chan int) {
	fmt.Println(target, pos, rlen)
	req, _ := http.NewRequest("POST", target, bytes.NewReader(b) )
	req.Header.Set("X-File-Name", fname)
	req.Header.Set("X-File-Offset", strconv.FormatInt(pos,10))
	req.ContentLength = rlen
	start := time.Now()
	resp, errp := http.DefaultClient.Do(req)
	rm := rateMsg("XMIT", time.Since(start), rlen)
	if errp == nil {
		defer resp.Body.Close()
		bodyMsg, err := ioutil.ReadAll(resp.Body)
		if (err == nil) {
		  fmt.Println(resp.StatusCode, rm, string(bodyMsg[:]))
		}
	} else {
		fmt.Println(errp)
	}
	rch <- 0
}

func postFileInChunks(file string, target string) {

    numThreads := 1
	rsz := int64(10*1024*1024)
	messages := make(chan int)

	f, err := os.Open(file)
	if (err != nil) {
		fmt.Println(err)
		return
	}
	defer f.Close()
	fst, err := f.Stat()
	flen := fst.Size()
	fbname := fst.Name()
	cpos := int64(0)
	epos := int64(0)

    posts := 0
    done := false
	start := time.Now()
    for !done {
		posts = 0
	    //rstart := time.Now()
		for i:=0; i<numThreads; i++ {
			if cpos >= flen {
				break
			}
		    posts++
			epos = cpos + rsz
			if epos > flen {
				epos = flen
			}
			blen := epos - cpos;
            b := make([]byte, blen)
	        //rstart = time.Now()
            f.ReadAt(b, cpos)
	        //rrm := rateMsg("READ", time.Since(rstart), blen)
	        //fmt.Println(rrm)
			go postChunk(f, b, fbname, cpos, blen, target, messages)
			cpos = epos
		}
		if posts > 0 {
			for posts > 0 {
				<-messages
				posts--
			}
		}
	    fmt.Println("")
		if epos >= flen {
			done = true
	        fmt.Println("done")
		}
	}
	rm := rateMsg("Total Throughput", time.Since(start), flen)
	fmt.Println(rm)
}


func main() {
    postFileInChunks("/tmp/Test.txt", "http://a.b.c.d/wrdr/uploadBlock")
}

