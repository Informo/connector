package main

import (
	"flag"
	"regexp"
	"strconv"
	"strings"

	"github.com/matrix-org/gomatrix"
)

type aliasesContent struct {
	Aliases []string `json:"aliases"`
}

var (
	serverName = flag.String("server-name", "127.0.0.1", "The homeserver's name")
	fqdn       = flag.String("fqdn", "", "The FQDN at which the Matrix APIs for this homeserver are reachable")
	port       = flag.Int("port", 0, "The port at which the homeserver can be reached")
	noTLS      = flag.Bool("no-tls", false, "If set to true, traffic will be sent with no TLS (plain HTTP)")
	entryPoint = flag.String("entrypoint", "#informo:matrix.org", "The entrypoint to the Informo network")
)

func main() {
	flag.Parse()

	if *fqdn == "" && *port == 0 {
		var success bool
		success, *fqdn, *port = lookup(*serverName)
		if !success {
			*fqdn = *serverName
			*port = 8448
		}
	} else if *fqdn == "" {
		*fqdn = *serverName
	} else if *port == 0 {
		*port = 8448
	}

	if !strings.HasPrefix(*entryPoint, "#") {
		panic("Invalid entrypoint: " + *entryPoint)
	}

	homeserverURL := "http"
	if !*noTLS {
		homeserverURL = homeserverURL + "s"
	}
	homeserverURL = homeserverURL + "://" + *fqdn + ":" + strconv.Itoa(*port)

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

	println("Registered as " + (*resp).UserID)

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
	err = client.StateEvent(respJoin.RoomID, "m.room.aliases", *serverName, &content)
	regex := regexp.MustCompile("code=404")
	if err != nil && !regex.MatchString(err.Error()) {
		panic(err)
	}

	println("Fetched previous entrypoints for this homeserver")

	content.Aliases = append(content.Aliases, "#"+randSeq(10, false)+":"+*serverName)

	_, err = client.SendStateEvent(respJoin.RoomID, "m.room.aliases", *serverName, content)
	if err != nil {
		panic(err)
	}

	println("Entrypoint added")
}
