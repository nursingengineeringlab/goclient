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

type User struct {
	Name       string
	Age        int
	RoomNumber int //fixme: this should not be int
	Gender     int
	DeviceID   string
	DeviceType string
}

func UserPost(URL string, user string, passwd string, data []byte) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, bytes.NewReader(data))
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
func createUser() {
	userData, err := json.Marshal(User{Name: "John Doe", Age: 31, RoomNumber: 1, Gender: 0, DeviceID: "7C1A23F227B4", DeviceType: "RRI"})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := UserPost("http://"+*addr+"/seniors/", "test", "test", userData)
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

	createUser()

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
