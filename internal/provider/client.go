package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Internalclient struct {
	HostURL    *string
	HTTPClient *http.Client
}

func NewClient(url *string) (*Internalclient, error) {
	return &Internalclient{HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL: url}, nil
}

type NetworkR struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (n *Internalclient) GetNetwork(id types.Int64) (*NetworkR, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%v/api/network/%v", *n.HostURL, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := n.doRequest(request)
	if err != nil {
		return nil, err
	}
	var net NetworkR

	err = json.Unmarshal(body, &net)

	if err != nil {
		log.Fatal("Error GetNetwork", err)
	}
	return &net, nil
}

func (n *Internalclient) DeleteNetwork(id types.Int64) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%v/api/network/%v", *n.HostURL, id), nil)
	if err != nil {
		return err
	}
	body, err := n.doRequest(request)
	fmt.Print(body)
	if err != nil {
		return err
	}
	return nil
}

func (n *Internalclient) CreateNetwork(name string) (*NetworkR, error) {
	payload, err := json.Marshal(map[string]string{"name": name})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%v/api/network", *n.HostURL), strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	body, err := n.doRequest(request)
	if err != nil {
		return nil, err
	}
	var net NetworkR

	err = json.Unmarshal(body, &net)

	if err != nil {
		log.Fatal("Error CreateNetwork", err)
	}
	return &net, nil
}

func (n *Internalclient) doRequest(request *http.Request) ([]byte, error) {

	request.Header.Set("Content-Type", "application/json")
	response, err := n.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, Body: %s", response.StatusCode, body)
	}

	return body, nil
}
