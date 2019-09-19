package swap

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Server struct {
	IP   string
	Port string
}

type Result struct {
	Data json.RawMessage `json:"data"`
}

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Res  Result `json:"result,omitempty"`
}

func (s *Server) request(apiURL string, payload []byte, respData interface{}) error {
	url := "http://" + s.IP + ":" + s.Port + "/" + apiURL
	resp := new(response)
	if err := post(url, payload, resp); err != nil {
		return err
	}

	if resp.Code != 200 {
		return errors.New(resp.Msg)
	}

	return json.Unmarshal(resp.Res.Data, respData)
}

func post(url string, payload []byte, result interface{}) error {
	return PostWithHeader(url, nil, payload, result)
}

func PostWithHeader(url string, header map[string]string, payload []byte, result interface{}) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// set Content-Type in advance, and overwrite Content-Type if provided
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if result == nil {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}
