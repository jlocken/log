package seq

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)


const (
	endpoint = "/api/events/raw"
)

type seqClient struct {
	baseUrl string
	apiKey  string
}

func (sc *seqClient) send(se *seqEvent, client *http.Client) error {
	fullUrl := sc.baseUrl + endpoint
	type seqrequestbody struct{
		Events  *seqEvent
	}
	reqbody := seqrequestbody{
		Events: se,
	}
	serializedevent, _ := json.Marshal(reqbody)
	str := string(serializedevent)
	println(str);

	request, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(serializedevent))
	if err != nil{return err}
	if len(sc.apiKey) > 1{
		request.Header.Set("X-Seq-ApiKey", sc.apiKey)
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := client.Do(request)
	defer request.Body.Close()
	if err != nil{
		return err
	}

	if response.StatusCode != 200 && response.StatusCode != 201{
		return errors.New(response.Status)
	}
	return nil
}