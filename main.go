package main

import (
	"bytes"
	"log"
	"image/png"
	"net/http"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func main() {
	http.HandleFunc("/qr/", func(w http.ResponseWriter, r *http.Request) {
		// Get address
        url := r.URL.Path[len("/qr/"):]
		// Generate QR
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