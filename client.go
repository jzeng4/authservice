package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
        "net/url"
)

func main() {
	caCert, err := ioutil.ReadFile("server.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	//cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	//if err != nil {
	//	log.Fatal(err)
	//}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				//Certificates: []tls.Certificate{cert},
                                //InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.PostForm("https://localhost:443?id=123", url.Values{"key": {"Value"}, "id": {"123"}})
	if err != nil {
		log.Println(err)
		return
	}

	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("%v\n", resp.Status)
	fmt.Printf(string(htmlData))
}
