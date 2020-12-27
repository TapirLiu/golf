package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var (
	port    = flag.String("port", "9999", "listening port")
	browser = flag.Bool("b", false, "open first page in a browser automatically")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	pwd, err := os.Getwd()
	if err != nil {
		log.Println("Can't get current path.")
		os.Exit(1)
	}

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
			for _, v := range etagHeaders {
				if r.Header.Get(v) != "" {
					r.Header.Del(v)
				}
			}
			for k, v := range noCacheHeaders {
				w.Header().Set(k, v)
			}
			h.ServeHTTP(w, r)
		}
	}

	go func() {
		time.Sleep(time.Second)

		log.Println("Serving folder:")
		log.Println("   " + pwd)
		log.Println("Running at:")
		log.Println("   http://localhost:" + *port)

		if addrs, err := net.InterfaceAddrs(); err == nil {
			for _, a := range addrs {
				if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					log.Println("   http://" + ipnet.IP.String() + ":" + *port)
				}
			}
		}

		if !*browser {
			return
		}

		// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
		var cmd string
		var args []string
		switch runtime.GOOS {
		case "windows":
			cmd = "cmd"
			args = []string{"/c", "start"}
		case "darwin":
			cmd = "open"
		default: // "linux", "freebsd", "openbsd", "netbsd"
			cmd = "xdg-open"
		}
		_ = exec.Command(cmd, append(args, "http://localhost:"+*port)...).Start()
	}()

	handler := NoCacheHandler(http.FileServer(http.Dir(pwd)))
  addr := ":"+*port
  if runtime.GOOS == "darwin" {
    addr = "localhost:" + *port
  }
	if err = http.ListenAndServe(addr, handler); err != nil {
		log.Printf("Failed to start server: %v\n", err)
	}
}
