package main

import (
	gws "github.com/gowebsecure/goWebSecure-go"
	"github.com/gowebsecure/goWebSecure-go/oauth"
)

func main() {
	const (
		TestNode = iota
	)

	gws.AddAPINode("testnode", TestNode, true)

	gws.AddHandler("testAction", TestNode, func(info *gws.HandlerInfo) []byte {
		return []byte("Testresponse")
	})

	// serverinit is blocking
	gws.ServerInit(func(id string) (oauth.CustomClientInfo, error) {
		// todo code for receiving secret for specific if from db
		return oauth.CustomClientInfo{
			ClientInfo: nil,
			ID:         "",
			Secret:     "",
			Domain:     "",
			UserID:     "",
		}, nil
	}, 8080)
}
