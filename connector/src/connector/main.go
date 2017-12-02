package main

import (
	"flag"

	"github.com/matrix-org/gomatrix"
)

type aliasesContent struct {
	Aliases []string `json:"aliases"`
}

const roomAlias = "#informo:matrix.org"

var (
	homeserver = flag.String("homeserver", "127.0.0.1", "The URL of the Matrix homeserver")
	port       = flag.String("port", "443", "The port at which the homeserver can be reached")
	tls        = flag.Bool("tls", true, "If set to false, traffic will be sent with no TLS (plain HTTP)")
)

func main() {
	flag.Parse()

	homeserverURL := "http"
	if *tls {
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

	respJoin, err := client.JoinRoom(roomAlias, "", nil)
	if err != nil {
		panic(err)
	}

	println("Joined")

	content := aliasesContent{
		Aliases: []string{"#" + randSeq(10, false) + ":" + *homeserver},
	}

	_, err = client.SendStateEvent(respJoin.RoomID, "m.room.aliases", *homeserver, content)
	if err != nil {
		panic(err)
	}

	println("Set alias")
}
