package main

import (
	"bytes"
	"fmt"
	"log"
	"image/png"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func main() {
	http.HandleFunc("/qr/", func(w http.ResponseWriter, r *http.Request) {
		// Get address
        url := r.URL.Path[len("/qr/"):]
        // Fix issues with URLs missing the second slash of a protocol
        if m, _ := regexp.MatchString(`:\/`, url); m == false {
        	re := regexp.MustCompile(`:\/`)
        	slice := re.Split(url, 2)
        	slice[0] = slice[0] + string("/")
        	fmt.Println(slice[0])
        	fmt.Println(slice)
        	url = strings.Join(slice, "")
        }
		// Generate QR
    	fmt.Println(url)
		qrCode, _ := qr.Encode(strings.ToUpper(url), qr.M, qr.Auto)
		qrCode, _ = barcode.Scale(qrCode, 500, 500)

		// Return to stdout
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, qrCode); err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Failed to generate QR."))
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buf.Bytes())))
		if _, err := w.Write(buf.Bytes()); err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Failed to display QR."))
		}
	})

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	log.Fatal(http.ListenAndServe(":9001", nil))
}

// localhost:9001/qr/https://apps.apple.com/gb/app/for-good/id1045549833