package apicupcfg

import (
	tls2 "crypto/tls"
	"encoding/base64"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"io/ioutil"
	"net/http"
	"strings"
)

func SomaReq(reqfile string, dpenv string, url string, tbox *rice.Box) (status string, statusCode int, reply string, err error) {

	// read datapower environment file
	username, password, err := loadDatapowerEnv(dpenv)
	if err != nil {
		return "", 0, "", err
	}

	// read soma request file
	somabytes := readFileBytes(reqfile)
	somasoap := strings.ReplaceAll(string(somabytes), "\n", "")

	return somaPost(url, somasoap, username, password)
}

func SomaUploadFile(file, dpdir, dpfile, dpdomain, dpenv, url string, tbox *rice.Box) (status string, statusCode int, reply string, err error) {

	// read datapower environment file
	username, password, err := loadDatapowerEnv(dpenv)
	if err != nil {
		return "", 0, "", err
	}

	// read upload file
	uploadbytes := readFileBytes(file)
	return somaSetFile(uploadbytes, dpdir, dpfile, dpdomain, url, username, password, tbox)
}

func loadDatapowerEnv(dpenv string) (username, password string, err error) {

	auth := string(readFileBytes(dpenv))
	up := strings.Split(auth, "\n")

	if len(up) < 2 {
		err := fmt.Errorf("invalid datapower environment file '%s'... username and password required", dpenv)
		return "", "", err
	}

	return up[0], up[1], nil
}

func somaSetFile(filebytes []byte, dpdir, dpfile string, dpdomain, url, user, password string, tbox *rice.Box) (status string, statusCode int, reply string, err error) {

	dp := DpFile{
		Domain:      dpdomain,
		Directory:   dpdir,
		FileName:    dpfile,
		FileContent: string(filebytes),
	}

	t := parseTemplate(tbox, tpdir(tbox) + "dp-set-file.tmpl")
	soma := executeTemplate2(t, dp)

	return somaPost(url, soma, user, password)
}

func somaPost(url string, somasoap string, user, password string) (status string, statusCode int, reply string, err error) {

	tls := &tls2.Config{
		Rand:                        nil,
		Time:                        nil,
		Certificates:                nil,
		NameToCertificate:           nil,
		GetCertificate:              nil,
		GetClientCertificate:        nil,
		GetConfigForClient:          nil,
		VerifyPeerCertificate:       nil,
		RootCAs:                     nil,
		NextProtos:                  nil,
		ServerName:                  "",
		ClientAuth:                  0,
		ClientCAs:                   nil,
		InsecureSkipVerify:          true, // do not verify server cert
		CipherSuites:                nil,
		PreferServerCipherSuites:    false,
		SessionTicketsDisabled:      false,
		SessionTicketKey:            [32]byte{},
		ClientSessionCache:          nil,
		MinVersion:                  0,
		MaxVersion:                  0,
		CurvePreferences:            nil,
		DynamicRecordSizingDisabled: false,
		Renegotiation:               0,
		KeyLogWriter:                nil,
	}

	tr := &http.Transport{
		Proxy:                  nil,
		DialContext:            nil,
		Dial:                   nil,
		DialTLS:                nil,
		TLSClientConfig:        tls,
		TLSHandshakeTimeout:    0,
		DisableKeepAlives:      false,
		DisableCompression:     false,
		MaxIdleConns:           0,
		MaxIdleConnsPerHost:    0,
		MaxConnsPerHost:        0,
		IdleConnTimeout:        0,
		ResponseHeaderTimeout:  0,
		ExpectContinueTimeout:  0,
		TLSNextProto:           nil,
		ProxyConnectHeader:     nil,
		MaxResponseHeaderBytes: 0,
		WriteBufferSize:        0,
		ReadBufferSize:         0,
		ForceAttemptHTTP2:      false,
	}

	client := &http.Client{
		Transport:     tr,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	// if body implements Closer it will be closed by client.Do()
	// *bytes.Buffer, *bytes.Reader, *strings.Reader => req.ContentLength is set to exact value
	bodyrdr := strings.NewReader(somasoap)

	req, err := http.NewRequest("POST", url, bodyrdr)
	if err != nil {
		return "", 0, "",nil
	}

	req.Header.Set("ContentType", "application/xml")
	req.Header.Set("ContentLength", fmt.Sprintf("%d", req.ContentLength))
	req.Header.Set("Authorization", basicAuthHeader(user, password))

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, "", err
	}
	defer func() {_ = resp.Body.Close()}()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, "", err
	}

	reply = string(bytes)

	return resp.Status, resp.StatusCode, reply, nil
}

func basicAuthHeader(user, password string) string {
	return fmt.Sprintf("Basic %s", basicAuthEncode(user, password))
}

func basicAuthEncode(user, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, password)))
}
