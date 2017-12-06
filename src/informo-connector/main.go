package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Informo/goutils"
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

const informoRoomID = "!xkMuBYHNWUOLHIoOEw:matrix.org"

func main() {
	flag.Parse()

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	if *fqdn == "" || *port == 0 {
		success, f, p := goutils.Lookup(*serverName)
		if !success {
			redefinePreserveContent(fqdn, port, *serverName, 8448)
		} else {
			redefinePreserveContent(fqdn, port, f, p)
		}
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
	regClient.Client.Transport = tr
	resp, err := regClient.RegisterDummy(&reqReq)
	if err != nil {
		panic(err)
	}

	println("Registered as " + (*resp).UserID)

	client, err := gomatrix.NewClient(homeserverURL, (*resp).UserID, (*resp).AccessToken)
	if err != nil {
		panic(err)
	}
	client.Client.Transport = tr

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

	var newAlias = "#" + randSeq(10, false) + ":" + *serverName

	content.Aliases = append(content.Aliases, newAlias)

	if success, err := createAliasOnServer(newAlias, homeserverURL, (*resp).AccessToken); err != nil {
		panic(err)
	} else if !success {
		return
	}

	_, err = client.SendStateEvent(respJoin.RoomID, "m.room.aliases", *serverName, content)
	if err != nil {
		panic(err)
	}

	println("Entrypoint added")
}

func redefinePreserveContent(fqdnBuf *string, portBuf *int, fqdn string, port int) {
	if *fqdnBuf == "" {
		*fqdnBuf = fqdn
	}

	if *portBuf == 0 {
		*portBuf = port
	}
}

func createAliasOnServer(alias string, homeserverURL string, accessToken string) (bool, error) {
	url := homeserverURL + "/_matrix/client/r0/directory/room/" + url.PathEscape(alias) + "?access_token=" + accessToken
	content := []byte(`{
		"room_id": "` + informoRoomID + `"
	}`)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(content))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		fmt.Printf("Can't create alias, got %d status\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Response body is %s\n", string(body))
		return false, err
	}

	println("Alias created")
	return true, err
}
