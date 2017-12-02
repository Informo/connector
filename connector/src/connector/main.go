package main

import (
	"flag"
	"strings"

	"github.com/matrix-org/gomatrix"
)

type aliasesContent struct {
	Aliases []string `json:"aliases"`
}

var (
	homeserver = flag.String("homeserver", "127.0.0.1", "The URL of the Matrix homeserver")
	port       = flag.String("port", "443", "The port at which the homeserver can be reached")
	noTLS      = flag.Bool("no-tls", false, "If set to true, traffic will be sent with no TLS (plain HTTP)")
	entryPoint = flag.String("entrypoint", "#informo:matrix.org", "The entrypoint to the Informo network")
)

func main() {
	flag.Parse()

	if !strings.HasPrefix(*entryPoint, "#") {
		panic("Invalid entrypoint: " + *entryPoint)
	}

	homeserverURL := "http"
	if !*noTLS {
		homeserverURL = homeserverURL + "s"
	}
	homeserverURL = homeserverURL + "://" + *homeserver + ":" + *port

	username := randSeq(20, false)

	reqReq := gomatrix.ReqRegister{
		Username: username,
		Password: randSeq(20, true),
	}

	regClient, err := gomatrix.NewClient(homeserverURL, "", "")
	if err != nil {
		panic(err)
	}
	resp, err := regClient.RegisterDummy(&reqReq)
	if err != nil {
		panic(err)
	}

	println("Registered as " + username)

	client, err := gomatrix.NewClient(homeserverURL, (*resp).UserID, (*resp).AccessToken)
	if err != nil {
		panic(err)
	}

	println("Client reinit")

	respJoin, err := client.JoinRoom(*entryPoint, "", nil)
	if err != nil {
		panic(err)
	}

	println("Joined")

	var content aliasesContent
	if err = client.StateEvent(respJoin.RoomID, "m.room.aliases", *homeserver, &content); err != nil {
		panic(err)
	}

	println("Fetched previous entrypoints for this homeserver")

	content.Aliases = append(content.Aliases, "#"+randSeq(10, false)+":"+*homeserver)

	_, err = client.SendStateEvent(respJoin.RoomID, "m.room.aliases", *homeserver, content)
	if err != nil {
		panic(err)
	}

	println("Entrypoint added")
}
