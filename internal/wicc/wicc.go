package wicc

import (
	"net/http"
	"io/ioutil"
	"bytes"
	"fmt"
	"encoding/base64"
	"encoding/hex"
)


func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

//var url string = "http://211.159.185.115:6967";
var url string = "https://testnode.wiccdev.org"
var pubKey string ="wMT8G9hYHxdSsE2x31JcNDpGGLDfCKW2Cj";
// https://baas-test.wiccdev.org/v2/api/swagger-ui.html#!/contract-controller/getContradtDataUsingPOST
// https://www.wiccdev.org/tool/smartContractTool.html
//url := "https://baas-test.wiccdev.org/v2/api"

func toHex(input string) string {
	src := []byte(input)

	dst := hex.EncodeToString(src)
	fmt.Printf("%s\n", dst)

	return dst
}

func Send(key, val string) {
	//jsonByte := []byte(`{"jsonrpc":"2.0","id":"curltext","method":"callcontracttx","params":["webJp3474phQ9db4bYMuvaSerYWqZE98m5","206653-2",0,"f01700006e616d653a626f63656e",1000000]}`)
	//jsonByte := []byte(`{"jsonrpc":"2.0","id":"curltext","method":"callcontracttx","params":["wMT8G9hYHxdSsE2x31JcNDpGGLDfCKW2Cj","206653-2",0,"f01700006e616d653a626f63656e",1000000]}`)

	encoded := "f0170000" + toHex(key) + "3a" + toHex(val)
	jsonByte := []byte(`{"jsonrpc":"2.0","id":"curltext","method":"callcontracttx","params":["` + pubKey + `","206653-2",0,"` + encoded + `",1000000]}`)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
    //req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Basic " + basicAuth("wiccuser", "123456"))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}

func Fetch(key string) {
	//url := "http://211.159.185.115:6967"
	//jsonByte := []byte(`{"jsonrpc":"2.0","id":"curltext","method":"callcontracttx","params":["webJp3474phQ9db4bYMuvaSerYWqZE98m5","206653-2",0,"f01700006e616d653a626f63656e",1000000]}`)
	//jsonByte := []byte(`{"jsonrpc":"2.0","id":"curltext","method":"getcontractdata","params":["206653-2","wMT8G9hYHxdSsE2x31JcNDpGGLDfCKW2Cj.name"]}`)

	jsonByte := []byte(`{"jsonrpc":"2.0","id":"curltext","method":"getcontractdata","params":["206653-2","` + pubKey + `.` + key + `"]}`)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
    //req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Basic " + basicAuth("wiccuser", "123456"))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}
