package Service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Map for storing websiteStatus of websites
var websiteStatus = map[string]string{}

// StatusChecker Introducing interfaces for loose coupling
type StatusChecker interface {
	Check(ctx http.ConnState, name string) (status bool, err error)
}

type httpChecker struct {
}

// Check -Function to check status of a website
func (checker *httpChecker) Check(ctx context.Context, name string) (status bool, err error) {
	resp, err := http.Get("http://" + name)
	if err != nil {
		return false, err
	}
	log.Print("Response Body returned for website "+name+" is: ", resp)
	return true, nil
}

func checkWebsiteStatus(website []string) {
	checker := httpChecker{}
	for {
		for _, val := range website {
			ctx := context.TODO()
			_, err := checker.Check(ctx, val)
			if err != nil {
				log.Printf("Error occurred=%s", err.Error())
				websiteStatus[val] = "DOWN"
			} else {
				websiteStatus[val] = "UP"
			}
		}
		time.Sleep(time.Minute)
	}
}

// HelloHandler -Handler function for handling incoming http request
func HelloHandler(writer http.ResponseWriter, userRequest *http.Request) {

	switch userRequest.Method {
	case "POST":
		var website []string
		json.NewDecoder(userRequest.Body).Decode(&website)
		//Go Routine
		go checkWebsiteStatus(website)

	case "GET":
		value := userRequest.URL.Query().Get("website")
		if value != "" {
			status, ok := websiteStatus[value]
			if !ok {
				writer.WriteHeader(http.StatusNotFound)
				return
			}
			jsonVal, err := json.Marshal(value + " : " + status)
			if err != nil {
				log.Printf("Error occurred=%s", err.Error())
				//Internal Server Error
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			writer.Write(jsonVal)
			return
		}
		jsonVal, err := json.Marshal(websiteStatus)
		if err != nil {
			log.Printf("Error occurred=%s", err.Error())
			//Internal Server Error
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Write(jsonVal)

	default:
		// Return 404
		writer.WriteHeader(http.StatusNotFound)
	}
}
