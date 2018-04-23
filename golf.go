package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"net"
	"time"
)

var (
	port = flag.String("port", "9999", "listening port")
)

func main() {
	flag.Parse()
	
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Can't get current directory.")
		os.Exit(1)
	}

	go func() {
		time.Sleep(time.Second)
		
		fmt.Println("Serving folder:")
		fmt.Println("   " + pwd)
		fmt.Println("Running at:")
		fmt.Println("   http://localhost:" + *port)
	
		if addrs, err := net.InterfaceAddrs(); err == nil {
			for _, a := range addrs {
				if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					fmt.Println("   http://" + ipnet.IP.String() + ":" + *port)
				}
			}
		}
	}()
	
	// http://stackoverflow.com/questions/33880343/go-webserver-dont-cache-files-using-timestamp
	var epoch = time.Unix(0, 0).Format(time.RFC1123)
	var noCacheHeaders = map[string]string{
		"Expires":         epoch,
		"Cache-Control":   "no-cache, private, max-age=0",
		"Pragma":          "no-cache",
		"X-Accel-Expires": "0",
	}
	var etagHeaders = []string{
		"ETag",
		"If-Modified-Since",
		"If-Match",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
	}
	
	NoCacheHandler := func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Delete any ETag headers that may have been set
			for _, v := range etagHeaders {
				if r.Header.Get(v) != "" {
					r.Header.Del(v)
				}
			}

			// Set our NoCache headers
			for k, v := range noCacheHeaders {
				w.Header().Set(k, v)
			}
			
			// Serve with the actual handler.
			h.ServeHTTP(w, r)
		}
	} 

	handler := NoCacheHandler(http.FileServer(http.Dir(pwd)))
	if err = http.ListenAndServe(":"+*port, handler); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
