package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Buzz struct {
	Color string `json:"color"`
}

type HttpBuzzer struct {
	buzzerStream chan buzzerHit
}

func NewHttpBuzzer(buzzer chan buzzerHit) *HttpBuzzer {
	sb := &HttpBuzzer{buzzerStream: buzzer}
	return sb
}

func (hb *HttpBuzzer) buzz(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var buzz Buzz
	json.Unmarshal(reqBody, &buzz)
	switch buzz.Color {
	case "red":
		sendBuzzerHit(hb.buzzerStream, buzzerColorRed)
	case "green":
		sendBuzzerHit(hb.buzzerStream, buzzerColorGreen)
	case "blue":
		sendBuzzerHit(hb.buzzerStream, buzzerColorBlue)
	case "yellow":
		sendBuzzerHit(hb.buzzerStream, buzzerColorYellow)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
