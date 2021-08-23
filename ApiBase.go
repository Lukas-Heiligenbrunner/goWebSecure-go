package goWebSecure_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gowebsecure/goWebSecure-go/oauth"
	"gopkg.in/oauth2.v3"
	"log"
	"net/http"
	"strconv"
)

const APIPREFIX = "/api"

type HandlerInfo struct {
	ID    uint32
	Token string
	Data  map[string]interface{}
}

type actionStruct struct {
	Action string
}

type Handler struct {
	action  string
	handler func(info *HandlerInfo) []byte
	apiNode int
}

var handlers = make(map[string]Handler)

func AddHandler(action string, apiNode int, h func(info *HandlerInfo) []byte) {
	// append new handler to the handlers
	handlers[fmt.Sprintf("%s/%d", action, apiNode)] = Handler{action, h, apiNode}
}

func ServerInit(userquery func(id string) (oauth.CustomClientInfo, error), port uint16) *error {
	// initialize oauth service and add corresponding auth routes
	oauth.InitOAuth(userquery)

	fmt.Printf("Server up and running on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return &err
}

// AddAPINode add create new API node (/api/<name>) on webserver
// define whether you want a secure node (with token check) or not with the Secure flag
func AddAPINode(name string, node int, secure bool) {
	if secure {
		http.Handle(APIPREFIX+"/"+name, oauth.ValidateToken(handlefunc, node))
	} else {
		http.HandleFunc(APIPREFIX+"/"+name, func(w http.ResponseWriter, r *http.Request) {
			handlefunc(w, r, node, nil)
		})
	}
}

func handleAPICall(action string, requestBody string, apiNode int, info *HandlerInfo) []byte {
	handler, ok := handlers[fmt.Sprintf("%s/%d", action, apiNode)]
	if !ok {
		// handler doesn't exist!
		fmt.Printf("no handler found for Action: %d/%s\n", apiNode, action)
		return nil
	}

	// check if info even exists
	if info == nil {
		info = &HandlerInfo{}
	}

	// parse the arguments
	var args map[string]interface{}
	err := json.Unmarshal([]byte(requestBody), &args)

	if err != nil {
		fmt.Printf("failed to decode arguments of action %s :: %s\n", action, requestBody)
	} else {
		// check if map has an action
		if _, ok := args["action"]; ok {
			delete(args, "action")
		}

		info.Data = args
	}

	// call the handler
	return handler.handler(info)
}

func handlefunc(rw http.ResponseWriter, req *http.Request, node int, tokenInfo *oauth2.TokenInfo) {
	// only allow post requests
	if req.Method != "POST" {
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	body := buf.String()

	var t actionStruct
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		fmt.Println("failed to read action from request! :: " + body)
	}

	if tokenInfo != nil {
		// load userid from received token object
		id, err := strconv.Atoi((*tokenInfo).GetClientID())
		if err != nil {
			log.Println("Error in UserId while loading Token from cache")
			return
		}

		userinfo := &HandlerInfo{
			ID:    uint32(id),
			Token: (*tokenInfo).GetCode(),
		}

		rw.Write(handleAPICall(t.Action, body, node, userinfo))
	} else {
		rw.Write(handleAPICall(t.Action, body, node, nil))
	}

}
