package hub

import (
	"encoding/json"
)

type (
	method string

	CommonRequest struct {
		Method method          `json:"method"`
		Params json.RawMessage `json:"params"`
		Id     int             `json:"id"`
	}

	AuthRequest struct {
		Method method `json:"method"`
		Params struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"params"`
		Id int `json:"id"`
	}

	SubRequest struct {
		Method method   `json:"method"`
		Params []string `json:"params"`
		Id     int      `json:"id"`
	}

	SubRespone struct {
		Result []string `json:"result"`
		Id     int      `json:"id"`
	}

	ErrorMessage struct {
		Code int
		Msg  string `json:"msg"`
	}
)

var (
	methodSubscribe         method = "SUBSCRIBE"
	methodUnsubscribe       method = "UNSUBSCRIBE"
	methodListSubscriptions method = "LIST_SUBSCRIPTIONS"
	methodAuth              method = "AUTH"

	codeInvalidMethod = 1
	codeInvalidAuth   = 2
)
