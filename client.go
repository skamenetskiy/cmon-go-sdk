package cmon

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
)

func NewClient(cnf *ClientConfig) *Client {
	httpClient := &http.Client{}
	if cnf.Insecure {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	return &Client{
		cnf:    cnf,
		http:   httpClient,
		sesMu:  &sync.Mutex{},
		userMu: &sync.Mutex{},
		debug:  os.Getenv("DEBUG_MODE") != "",
	}
}

type Client struct {
	cnf    *ClientConfig
	http   *http.Client
	ses    *http.Cookie
	sesMu  *sync.Mutex
	user   *User
	userMu *sync.Mutex
	debug  bool
}

func (client *Client) Request(module string, req, res interface{}, retry bool) error {
	if client.ses == nil {
		if err := client.Authenticate(); err != nil {
			return err
		}
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(
		http.MethodPost,
		client.buildURI(module),
		bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if client.ses != nil {
		request.Header.Set("cookie", client.ses.String())
	}
	if client.debug {
		log.Println(httputil.DumpRequest(request, true))
	}
	response, err := client.http.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		if retry {
			return errors.New("retry failed after re-authentication")
		}
		if err := client.Authenticate(); err != nil {
			return err
		}
		return client.Request(module, req, res, true)
	}
	if client.debug {
		log.Println(httputil.DumpResponse(response, true))
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	return json.NewDecoder(response.Body).Decode(res)
}

func (client *Client) Authenticate() error {
	rd := &authenticateRequest{
		&WithOperation{"authenticateWithPassword"},
		client.cnf.Username,
		client.cnf.Password,
	}
	rb, err := json.Marshal(rd)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		client.buildURI(ModuleAuth),
		bytes.NewBuffer(rb))
	if err != nil {
		return err
	}
	res, err := client.http.Do(req)
	if err != nil {
		return err
	}
	if !client.saveSessionFromResponse(res) {
		return errors.New("failed to save session")
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println(err)
		}
	}()
	ar := &authenticateResponse{}
	if err := json.NewDecoder(res.Body).Decode(ar); err != nil {
		return err
	}
	if ar.RequestStatus != RequestStatusOk {
		return NewErrorFromResponseData(ar.WithResponseData)
	}
	client.userMu.Lock()
	client.user = ar.User
	client.userMu.Unlock()
	return nil
}

type authenticateRequest struct {
	*WithOperation `json:",inline"`

	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type authenticateResponse struct {
	*WithResponseData `json:",inline"`

	User *User `json:"user"`
}

func (client *Client) buildURI(module string) string {
	u := &url.URL{
		Host:   client.cnf.Host + ":" + client.cnf.Port,
		Scheme: "https",
		Path:   "/v2/" + module,
	}
	return u.String()
}

func (client *Client) saveSessionFromResponse(res *http.Response) bool {
	for _, c := range res.Cookies() {
		if c.Name == "cmon-sid" {
			client.sesMu.Lock()
			client.ses = c
			client.sesMu.Unlock()
			return true
		}
	}
	return false
}

type ClientConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Insecure bool
}
