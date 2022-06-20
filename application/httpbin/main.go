package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"k8s.io/klog/v2"
)

const (
	defaultAddr    string   = ":8080"
	defaultVersion string = "v1"

	tlsPath = "/etc/httpbin/tls/"
	caCertPath = "etc/httpbin/ca/ca.crt"
	tlsKey = tlsPath + "tls.key"
	tlsCert = tlsPath + "tls.crt"
)

var (
	args Args
	retryCount = 3
)

func printHostName(w http.ResponseWriter) {
	hostName := os.Getenv("HOSTNAME")
	_, _ = w.Write([]byte(hostName + "\n"))
}

func version(w http.ResponseWriter, req *http.Request) {
	printHostName(w)
	_, _ = w.Write([]byte(fmt.Sprintf("version: %s\n", args.Version)))
}

func header(w http.ResponseWriter, req *http.Request) {
	printHostName(w)

	header := req.Header
	header["Host"] = []string{req.Host}
	header["Path"] = []string{req.URL.String()}
	header["Protocol"] = []string{req.Proto}
	header["URL"] = []string{req.URL.String()}
	if req.TLS != nil {
		header["TLSHankShake"] = []string{strconv.FormatBool(req.TLS.HandshakeComplete)}
	}

	result, _ := json.MarshalIndent(header, "  ", "  ")
	_, _ = w.Write([]byte(fmt.Sprintf("headers: %s\n query param: %s\n", string(result), req.URL.RawQuery)))
}

func timeout(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	rawTime := req.Form["time"][0]
	t, err := strconv.Atoi(rawTime)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	time.Sleep(time.Duration(t) * time.Second)
	_, _ = w.Write([]byte("success"))
}

func retry(w http.ResponseWriter, _ *http.Request) {
	klog.Infof("Receive request, count: %d\n", retryCount)
	if retryCount != 0 {
		retryCount--
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(fmt.Sprintf("error, count: %d\n", retryCount)))
	} else {
		retryCount = 3
		_, _ = w.Write([]byte("success"))
	}
}

func main() {
	flag.StringVar(&args.ServerOptions.Addr, "addr", defaultAddr, "The addr this server listen on")
	flag.StringVar(&args.Version, "version", defaultVersion, "service version")
	flag.BoolVar(&args.ServerOptions.TLS.Enable, "tls", false, "open tls")
	flag.Parse()

	httpServerMux := http.NewServeMux()
	httpServerMux.HandleFunc("/version", version)
	httpServerMux.HandleFunc("/header", header)
	httpServerMux.HandleFunc("/timeout", timeout)
	httpServerMux.HandleFunc("/retry", retry)

	httpServer := &http.Server{
		Addr:    args.ServerOptions.Addr,
		Handler: httpServerMux,
	}

	listener, err := net.Listen("tcp", args.ServerOptions.Addr)
	if err != nil {
		panic(err)
	}

	if args.ServerOptions.TLS.Enable {
		if err := httpServer.ServeTLS(listener, tlsCert, tlsKey); err != nil {
			panic(err)
		}
	} else {
		if err := httpServer.Serve(listener); err != nil {
			panic(err)
		}
	}
}
