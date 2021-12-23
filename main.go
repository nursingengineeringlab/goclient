package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var addr = flag.String("addr", "127.0.0.1:8000", "http service address")

var httpPrefix = "http://"

var basicAuthUser = "test"
var basicAuthPass = "test"

type UserStruct struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PatientPageList struct {
	Count int `json:"count"`
	//Next     int       `json:"next"`
	//Previous int       `json:"previous"`
	Results []Patient `json:"results"`
}

type Patient struct {
	Name       string     `json:"name"`
	Age        int        `json:"age"`
	RoomNumber int        `json:"room_no"` //fixme: this should not be int
	Gender     string     `json:"gender"`
	DeviceID   string     `json:"device_id"`
	DeviceType string     `json:"device_type"`
	User       UserStruct `json:"user"`
}

func userPost(URL string, user string, passwd string, data []byte) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewReader(data))
	req.SetBasicAuth(user, passwd)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	fmt.Println(s)

	return resp, err
}

func getUserList() []Patient {
	client := &http.Client{}
	URL := httpPrefix + *addr + "/seniors/"
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	req.SetBasicAuth(basicAuthUser, basicAuthPass)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	var pList PatientPageList
	err = json.Unmarshal(bodyText, &pList)
	if err != nil {
		log.Fatal(err)
	}
	return pList.Results
}

func deleteUser(deviceID string) {
	client := &http.Client{}
	URL := httpPrefix + *addr + "/seniors/" + deviceID + "/" //https://docs.djangoproject.com/en/2.2/ref/settings/#append-slash
	fmt.Println(URL)

	req, err := http.NewRequest(http.MethodDelete, URL, nil)
	req.SetBasicAuth(basicAuthUser, basicAuthPass)
	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(respBody))
}

func createUser() {
	userData, err := json.Marshal(Patient{Name: "John Doe1", Age: 31, RoomNumber: 1, Gender: "M", DeviceID: "7C1A23F227B4", DeviceType: "RRI", User: UserStruct{Username: "test88", Email: "test@test.com", Password: "test23"}})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := userPost(httpPrefix+*addr+"/seniors/", basicAuthUser, basicAuthPass, userData)
	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println(res["json"])

}
func main() {
	flag.Parse()
	log.SetFlags(0)

	userList := getUserList()
	for _, val := range userList {
		fmt.Println("deleting user: ", val)
		deleteUser(val.DeviceID)
	}

	//createUser()

	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)
	//
	//u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws/sensor/RR"}
	//log.Printf("connecting to %s", u.String())
	//
	//c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	//if err != nil {
	//	log.Fatal("dial:", err)
	//}
	//defer c.Close()
	//
	//done := make(chan struct{})
	//
	//go func() {
	//	defer close(done)
	//	for {
	//		_, message, err := c.ReadMessage()
	//		if err != nil {
	//			log.Println("read:", err)
	//			return
	//		}
	//		log.Printf("recv: %s", message)
	//	}
	//}()
	//
	//ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()
	//
	//for {
	//	select {
	//	case <-done:
	//		return
	//	case t := <-ticker.C:
	//		err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
	//		if err != nil {
	//			log.Println("write:", err)
	//			return
	//		}
	//	case <-interrupt:
	//		log.Println("interrupt")
	//
	//		// Cleanly close the connection by sending a close message and then
	//		// waiting (with timeout) for the server to close the connection.
	//		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	//		if err != nil {
	//			log.Println("write close:", err)
	//			return
	//		}
	//		select {
	//		case <-done:
	//		case <-time.After(time.Second):
	//		}
	//		return
	//	}
	//}
}
