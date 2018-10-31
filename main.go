package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

type withdraw_data struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}
type Data struct {
	Coin     string          `json:"coin"`
	Withdraw []withdraw_data `json:"withdraw_data"`
}

func httpReq(url string, data string, token string) string {
	payload := strings.NewReader(data)
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}

func jsonDecode(data string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	return result
}

func b64Encode(data string) string {
	encodeString := base64.StdEncoding.EncodeToString([]byte(data))
	return encodeString
}

func b64Decode(data string) []string {
	decodedString, _ := base64.StdEncoding.DecodeString(data)
	stringSplit := strings.SplitN(string(decodedString), ":", 2)
	return stringSplit
}

func loginApi(w http.ResponseWriter, r *http.Request) {
	data := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	decodedData := b64Decode(data[1])
	payloadData := fmt.Sprintf(`{"email" : "%s", "password" : "%s" }`, decodedData[0], decodedData[1])
	req := httpReq("http://159.65.13.106/web/public/api/login", payloadData, "")
	resp := jsonDecode(req)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(resp)
}

func newAddress(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new Address")
	token := r.Header.Get("X-Access-Token")
	if token == "" {
		http.Error(w, "Token Missing", 403)
	}
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	decodeData := jsonDecode(string(data))
	coin := decodeData["coin"]
	username := decodeData["username"]
	userid := decodeData["userid"]
	if coin == nil || coin == "" || username == nil || username == "" || userid == nil || userid == "" {
		http.Error(w, "Coin Or Username Missing", 403)
	}
	reqPayload := fmt.Sprintf(`{"coin": "%s", "username" : "%s", "userid" : "%s"}`, coin, username, fmt.Sprint(userid))
	req := httpReq("http://159.65.13.106/web/public/api/get_address", reqPayload, token)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println(req)
	resp := jsonDecode(req)
	json.NewEncoder(w).Encode(resp)
}

func getHotBalance(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Access-Token")
	if token == "" {
		http.Error(w, "Token Missing", 403)
	}
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	decodeData := jsonDecode(string(data))
	coin := decodeData["coin"]
	if coin == nil || coin == "" {
		http.Error(w, "Coin Missing", 403)
	}
	reqPayload := fmt.Sprintf(`{"coin" : "%s"}`, coin)
	req := httpReq("http://159.65.13.106/web/public/api/hot_balance", reqPayload, token)
	w.Header().Set("Content-Type", "application/json")
	resp := jsonDecode(req)
	json.NewEncoder(w).Encode(resp)
}

func getColdBalance(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Access-Token")
	if token == "" {
		http.Error(w, "Token Missing", 403)
	}
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	decodeData := jsonDecode(string(data))
	coin := decodeData["coin"]
	if coin == nil || coin == "" {
		http.Error(w, "Coin Missing", 403)
	}
	reqPayload := fmt.Sprintf(`{"coin" : "%s"}`, coin)
	req := httpReq("http://159.65.13.106/web/public/api/cold_balance", reqPayload, token)
	w.Header().Set("Content-Type", "application/json")
	resp := jsonDecode(req)
	json.NewEncoder(w).Encode(resp)
}

func createRawTx(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Access-Token")
	if token == "" {
		http.Error(w, "Token Missing", 403)
	}
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	decodeData := jsonDecode(string(data))
	coin := decodeData["coin"]
	address := decodeData["address"]
	amountraw := decodeData["amount"]
	if coin == nil || coin == "" || address == nil || address == "" || amountraw == "" || amountraw == nil {
		http.Error(w, "Data Missing", 403)
	}
	rawamount, _ := strconv.Atoi(fmt.Sprint(amountraw))
	amount := rawamount
	// NEED TO WORK BELOW
	var s withdraw_data
	s.Address = fmt.Sprint(address)
	s.Amount = fmt.Sprint(amount)
	wd := []withdraw_data{}
	wd = append(wd, s)
	full := &Data{Coin: "KMD", Withdraw: wd}
	jsonString, _ := json.Marshal(full)
	reqPayload := string(jsonString)
	// NEED TO WORK ABOVE
	fmt.Println(reqPayload)
	req := httpReq("http://159.65.13.106/web/public/api/createraw_hash", reqPayload, token)
	w.Header().Set("Content-Type", "application/json")
	resp := jsonDecode(req)
	json.NewEncoder(w).Encode(resp)
}

func broadcastTx(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Access-Token")
	if token == "" {
		http.Error(w, "Token Missing", 403)
	}
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	decodeData := jsonDecode(string(data))
	rawhash := decodeData["hash"]
	privateKey := "UuipWY1qWw7QmBfhtpLt62Z8td6LtMz67KeNMFcUcUXENJiNcEwn"
	if rawhash == nil || rawhash == "" {
		http.Error(w, "Data Missing", 403)
	}
	reqPayload := fmt.Sprintf(`{"raw_hash" : "%s", "coin" : "KMD", "private_key" : "%s"}`, rawhash, privateKey)
	req := httpReq("http://159.65.13.106/web/public/api/broadcast_hash", reqPayload, token)
	w.Header().Set("Content-Type", "application/json")
	resp := jsonDecode(req)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-Width", "Content-Type", "Authorization", "X-Access-Token"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	r.HandleFunc("/api/login", loginApi).Methods("POST")
	r.HandleFunc("/api/getnewaddress", newAddress).Methods("POST")
	r.HandleFunc("/api/gethotbalance", getHotBalance).Methods("POST")
	r.HandleFunc("/api/getcoldbalance", getColdBalance).Methods("POST")
	r.HandleFunc("/api/createrawtx", createRawTx).Methods("POST")
	r.HandleFunc("/api/broadcasttx", broadcastTx).Methods("POST")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(headers, methods, origins)(r)))

}
