package main

import (
	"bytes"
	"log"
	"image/png"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/aztec"
	"github.com/boombuler/barcode/qr"
)

func handleData(str string) string {
    if m, _ := regexp.MatchString(`:\/`, str); m == true {
    	// Fix strings with protocols (would only be :/ instead of ://)
		re := regexp.MustCompile(`:\/`)
		slice := re.Split(str, 2)
		slice[0] = slice[0] + string("://")
		return strings.Join(slice, "")
    } else {
    	// Remove trailing slash
    	if strings.HasSuffix(str, "/") {
    		str = str[:len(str)-1]
    	}
    	return str
    }
}
func stdout(w http.ResponseWriter, data barcode.Barcode) {
	// Create a buffer to store barcode data
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, data); err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to generate QR."))
	}

	// Set headers and write to body of return request
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf.Bytes())))
	if _, err := w.Write(buf.Bytes()); err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Failed to display QR."))
	}
}

func main() {
	http.HandleFunc("/aztec/", func(w http.ResponseWriter, r *http.Request) {
		data := handleData(r.URL.Path[len("/aztec/"):])
		code, _ := aztec.Encode([]byte(data), 33, 0)
		code, _ = barcode.Scale(code, 500, 500)
		stdout(w, code)
	})
	http.HandleFunc("/qr/", func(w http.ResponseWriter, r *http.Request) {
		data := handleData(r.URL.Path[len("/qr/"):])
		code, _ := qr.Encode(strings.ToUpper(data), qr.M, qr.Auto)
		code, _ = barcode.Scale(code, 500, 500)
		stdout(w, code)
	})

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	log.Fatal(http.ListenAndServe(":9001", nil))
}

// localhost:9001/qr/https://apps.apple.com/gb/app/for-good/id1045549833