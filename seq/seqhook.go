package seq

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/fatih/structs"
)


type SeqHook struct {
	BaseUrl string
	ApiKey  string
}

type event struct {
	Timestamp 		time.Time
	Level 	  		string
	MessageTemplate string
	Properties 		map[string]interface{}
	Exception		string
}


type seqEvent []event

func (seqhook *SeqHook) Info(msg string, s interface{}) {
	m := mapProps(s)

	event := event{
		Timestamp: 		time.Now().UTC(),
		Level: 			"Information",
		MessageTemplate: msg,
		Properties:		 m,
	}
	seqhook.log(event)
}

func (seqhook *SeqHook) Warning(msg string, s interface{}) {
	m := mapProps(s)

	event := event{
		Timestamp:       time.Now().UTC(),
		Level:           "Warning",
		MessageTemplate: msg,
		Properties:      m,
	}
	seqhook.log(event)	
}

func (seqhook *SeqHook) Fatal(err error, s interface{}) {
	seqhook.Error(err, s)
	panic(err)
}


func (seqhook *SeqHook) Error(err error, s interface{}) {
	m := mapProps(s)

	event := event{
		Timestamp:       time.Now().UTC(),
		Level:           "Error",
		MessageTemplate: err.Error(),
		Properties:      m,
		Exception:       fmt.Sprintf("%+v", err),
	}
	seqhook.log(event)
}


func (seqhook *SeqHook) log(ev event) {
	var httpClient = &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout: 30 * time.Second,
		},
	}

	sc := seqClient{
		baseUrl: seqhook.BaseUrl,
		apiKey:  seqhook.ApiKey,
	}

	var se seqEvent
	se = append(se, ev)

	err := sc.send(&se, httpClient)
	if err != nil {
		println(err.Error())
	}
}


func mapProps(s interface{}) map[string]interface{} {
	var m map[string]interface{}
	if s != nil {
		switch reflect.ValueOf(s).Kind() {
		case reflect.Struct:
			m = structs.Map(s)
		default:
			m = make(map[string]interface{})
			m["key"] = fmt.Sprintf("%+v", s)
		}
	}
	return m
}