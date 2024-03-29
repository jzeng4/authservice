package main

import (
	"crypto/tls"
	//"crypto/x509"
	//"io/ioutil"
	"log"
        "fmt"
	"net/http"
)

var mp = map[string]string {
    "client1" : "client1",
    "client2" : "client2",
    "master1" : "master1",
}

func main() {
	//caCert, err := ioutil.ReadFile("client.crt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//caCertPool := x509.NewCertPool()
	//caCertPool.AppendCertsFromPEM(caCert)
	cfg := &tls.Config{
		//ClientAuth: tls.RequireAndVerifyClientCert,
		//ClientCAs:  caCertPool,
	}
	srv := &http.Server{
		Addr:      ":443",
		Handler:   &handler{},
		TLSConfig: cfg,
	}
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
        fmt.Printf("id=%s", req.URL.Query()["id"][0])
	w.Write([]byte("PONG\n"))
}

