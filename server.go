package zebra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"os/user"
)

type Server interface {
	Validate(ctx context.Context) error
	ProcessServerStatus() ServStatus
}

type ServStatus struct {
	Fault Fault `json:"fault,omitempty"`
	State State `json:"state,omitempty"`
}

type StatusFactory interface {
	New(string) Server
	Add(Type, TypeConstructor) StatusFactory
	Types() []Type
	Type(string) (Type, bool)
	Constructor(string) (TypeConstructor, bool)
}

type ServerStore struct {
	factory  StatusFactory
	servData map[string]Server
}

type ServerMap struct {
	fact       StatusFactory
	ServerInfo map[string]*StatusList
}

type StatusList struct {
	ctr      TypeConstructor
	Statuses []Server
}

func NewServerStore(serverData *ServerMap) *ServerStore {
	ids := &ServerStore{
		factory: ServFactory(),
		servData: func() map[string]Server {
			resMap := make(map[string]Server)
			for _, l := range serverData.ServerInfo {
				for _, res := range l.Statuses {
					resMap["server-status"] = res
				}
			}

			return resMap
		}(),
	}

	return ids
}

func ServFactory() StatusFactory {
	return typeMap{}
}

func ProcessServerStatus(rawResult []byte) map[string]string {
	c := make(map[string]json.RawMessage)

	// unmarschal JSON
	e := json.Unmarshal(rawResult, &c)

	// panic on error
	if e != nil {
		panic(e)
	}

	// a string slice to hold the keys
	results := make(map[string]string)

	// iteration counter
	i := 0

	// copy c's keys into k
	for s, val := range c {
		results[s] = string(val)
		i++
	}

	return results
}

var (
	SSH_USER, err = user.Current()
	REQUEST       = "curl "
	SERVER_HOST   = "https://" + string(SSH_USER.Username) + "@" + string(SSH_USER.Username) + "cs-bld.insieme.local:8000"
	REDFISH_AGENT = "/redfish/v1/Systems/ "
)

type redfishAPI struct {
	user        *user.User
	curl        string
	requestType string
	flags       string
	greps       string
}

func NewRedfish(thisType string, flag string, grep string) *redfishAPI {
	return &redfishAPI{
		user:        SSH_USER,
		curl:        REQUEST + SERVER_HOST + REDFISH_AGENT,
		requestType: thisType,
		flags:       flag,
		greps:       grep,
	}

}

// function for the server health.
func ServerStatus() (map[string]string, error) {
	redfish := NewRedfish("Systems/System.Embedded.1", "--user root:password ", "| jq .Status")

	executableCommand := redfish.curl + redfish.requestType + redfish.flags + redfish.greps

	erred := new(error)

	cmd := exec.Command(executableCommand)

	ret, err := exec.Command(executableCommand).Output()

	err2 := cmd.Run()

	data := ProcessServerStatus(ret)

	if err == nil && err2 == nil {
		return data, nil
	} else if err != nil {
		erred = &err
	} else if err2 != nil {
		erred = &err2
	}

	return nil, *erred
}

/*
// function for the server health.
func ServerHealthStatus() (string, error) {
	curl := "curl https://<OOB>/redfish/v1/Systems/System.Embedded.1 --user root:password | jq .Status"

	erred := new(error)

	cmd := exec.Command(curl)

	ret, err := exec.Command(curl).Output()

	err2 := cmd.Run()

	if err == nil && err2 == nil {
		// log.Fatal(err)
		return string(ret), nil
	} else if err != nil {
		erred = &err
	} else if err2 != nil {
		erred = &err2
	}

	return "", *erred
}
*/

func StatusCode(PAGE string, AUTH string) (r string) {
	// Setup the request.
	req, err := http.NewRequest("GET", PAGE, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", AUTH)

	// Execute the request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err.Error()
	}

	// Close response body as required.
	defer resp.Body.Close()

	fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

	return resp.Status
	// or fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
}

/*
func redfishConfig() *redfish.Configuration {

	return &redfish.Configuration{
		BasePath:      "http://localhost:8000",
		DefaultHeader: make(map[string]string),
		UserAgent:     "go-redfish/client",
	}
}

var redCfg = redfishConfig()
var redfishApi = redfish.NewAPIClient(redCfg).DefaultApi


type APIClient struct {
	DefaultApi *DefaultApiService
	// contains filtered or unexported fields
}

// function for the server health.
func (redfishApi) ServerStatus() (string, error) {
	erred := new(error)

	/*	curl := "curl https://<OOB>/redfish/v1/Systems/System.Embedded.1 --user root:password | jq .Status"

		erred := new(error)

		cmd := exec.Command(curl)

		ret, err := exec.Command(curl).Output()

		err2 := cmd.Run()

		if err == nil && err2 == nil {
			// log.Fatal(err)
			return string(ret), nil
		} else if err != nil {
			erred = &err
		} else if err2 != nil {
			erred = &err2
		}*/
/*
		sl, _, _ := redfishApi.ListSystems(context.Background())

		a := redfishApi.Status.State
		return *sl.Description, *erred
	}
*/

/*

For the server, we can use page and api data to get the status.

Multiple functions can be used to achieve various status tasks for the server, which can be structured based on
what is needed. They can be used to execute and retreive the specific information that can later be processed as desired.
The results can be kept in store until a later iteration.

*/

/*

Hello. Here are some ideas that we can utilize for the server status. Multiple functions can be used to achieve various status tasks, which can be structured based on what is needed. These functions can be used to execute, retrieve, and return specific information that can later be processed as desired.
The results can be kept in store until a later iteration.
*/

/*
Hello. Here are some ideas that we can utilize for the server status. Multiple functions can be used to achieve various status tasks, based on what is needed. These functions can be used to return specific information that can later be processed as desired.
The results can also be kept in store.
*/
