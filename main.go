package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	timestampFormat = "02/01/2006:15:04:05 -0700"
)

var (
	listenAddr    string
	useProxy      bool
	useCloudflare bool
	logRequests   bool
	redirectURL   *url.URL
	redirectCode  int
)

func init() {
	e := os.Getenv

	var err error

	listenAddr = e("REDIRECT_LISTEN")

	redirectURL, err = url.Parse(e("REDIRECT_URL"))
	if err != nil {
		log.Fatalf("ERROR: Invalid URL supplied: %s", err)
	}

	redirectCode, err = strconv.Atoi(e("REDIRECT_CODE"))
	if err != nil {
		log.Fatalf("ERROR: invalid value: %s", err)
	}

	logRequests, err = strconv.ParseBool(e("REDIRECT_LOG"))
	if err != nil {
		log.Fatalf("ERROR: Invalid value: %s", err)
	}

	useCloudflare, err = strconv.ParseBool(e("REDIRECT_BEHIND_CLOUDFLARE"))
	if err != nil {
		log.Fatalf("ERROR: Invalid value: %s", err)
	}

	useProxy, err = strconv.ParseBool(e("REDIRECT_BEHIND_PROXY"))
	if err != nil {
		log.Fatalf("ERROR: Invalid value: %s", err)
	}
}

func main() {
	var h http.Handler

	h = http.RedirectHandler(redirectURL.String(), redirectCode)

	if logRequests {
		h = logMiddleware(h)
	}

	log.Printf("Listening on %s ...", listenAddr)
	err := http.ListenAndServe(listenAddr, h)
	log.Fatalf("Error running server: %s", err)
}

func logMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// handle early
		next.ServeHTTP(w, r)
		// log request in common log format
		fmt.Printf("%s - - [%s] \"%s %s\" %d -\n",
			getRemoteAddr(r),
			time.Now().Format(timestampFormat),
			r.Method, r.URL.Path, redirectCode,
		)
	}

	if logRequests {
		return http.HandlerFunc(fn)
	}

	return next
}

func getRemoteAddr(r *http.Request) string {
	if useCloudflare {
		if ip := r.Header.Get("Cf-Connecting-Ip"); ip != "" {
			return ip
		}
	}

	if useProxy {
		if ip := r.Header.Get("X-Real-Ip"); ip != "" {
			return ip
		}
	}

	return r.RemoteAddr
}
