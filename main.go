package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Struct type to marsha and unmarshal request
type ReqMessage struct {
	Message string `json:"message"`
}

//Map for saving message and hash of message
var Hashs = make(map[string]string)

func main() {

	http.HandleFunc("/messages", message)
	http.HandleFunc("/messages/", getMessage)
	log.Fatal(http.ListenAndServeTLS(":5000", "localhost.crt", "localhost.key", nil))

}

// handler func for POST req to generate a hash
func message(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		decoder := json.NewDecoder(r.Body)

		var h ReqMessage
		err := decoder.Decode(&h)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		// Encrypt message and assign to msg variable
		msg := encrypt(h.Message)
		digest := map[string]string{"digest": msg}

		//Saving hash and message into map and in memory for search by hash key
		Hashs[msg] = h.Message

		resp, err := json.Marshal(digest)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, "%s", string(resp))
	} else {
		p := map[string]string{"err_msg": "Method Not Allowed"}
		resp, _ := json.Marshal(p)
		http.Error(w, string(resp), http.StatusMethodNotAllowed)
	}
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-type", "application/json")
		msg := r.URL.Path[len("/messages/"):]
		if msg == "" {
			m := map[string]string{"err_msg": "Empty arguments"}
			resp, _ := json.Marshal(m)
			fmt.Println(string(resp))
			http.Error(w, string(resp), http.StatusNotFound)

		} else if val, ok := Hashs[msg]; ok {
			m2 := &ReqMessage{
				Message: val,
			}
			resp, _ := json.Marshal(m2)
			fmt.Fprintf(w, "%s", string(resp))
		} else {
			m3 := map[string]string{"Message": "Message not found"}
			resp, _ := json.Marshal(m3)
			http.Error(w, string(resp), http.StatusNotFound)
		}
	}
}

func encrypt(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
