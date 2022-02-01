package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type PersonalDetails struct {
	Name string
	Age  int
}

type TcpRequest struct {
	Method string
	Url    string
	Data   string
	Params struct {
		ID string
	}
}

type Account struct {
	ID              string
	Type            string
	Balance         float64
	Transcations    []Transcation
	PersonalDetails *PersonalDetails
}

type Transcation struct {
	Credit         float64
	Debit          float64
	ClosingBalance float64
}

var accounts = []Account{}

func init() {
	log.SetPrefix("Bank core: ")
}

func main() {

	ln, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatalln("Error listening:", err.Error())
	}
	defer ln.Close()
	for {
		log.Println("Server listening: localhost:8000")
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln("Error accepting:", err.Error())
		}
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	defer conn.Close()
	handleRoutes(conn)
}

func handleRoutes(conn net.Conn) {

	tcpReq := &TcpRequest{}
	gob.NewDecoder(conn).Decode(tcpReq)
	log.Printf("Received : %+v", tcpReq)
	endpoint := tcpReq.Url
	switch method := tcpReq.Method; method {
	case "GET":
		if endpoint == "/accounts" {
			if tcpReq.Params.ID != "" {
				getAccountByID(conn, tcpReq.Params.ID)
			} else {
				getAccounts(conn)
			}
			return
		}
	case "POST":
		log.Println("POST METHOD", tcpReq)
		if endpoint == "/accounts" {
			openAccount(conn, tcpReq.Data)
			return
		} else if endpoint == "/transcations" && tcpReq.Params.ID != "" {
			makeTranscation(conn, tcpReq.Params.ID, tcpReq.Data)
			return
		}
	case "PATCH":
		if endpoint == "/profile" && tcpReq.Params.ID != "" {
			updateProfileByID(conn, tcpReq.Params.ID, tcpReq.Data)
			return
		}
	}

}

func openAccount(c net.Conn, dataString string) {
	var newAccount Account
	err := json.Unmarshal([]byte(dataString), &newAccount)
	if err != nil {
		fmt.Println("Error in unmarshal", err.Error())
	}
	accounts = append(accounts, newAccount)
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(newAccount)
	c.Write(buf.Bytes())
}

func getAccounts(c net.Conn) {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(accounts)
	c.Write(buf.Bytes())
}

func getAccountByID(c net.Conn, id string) {
	account := findAcountByID(id)
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(account)
	c.Write(buf.Bytes())
}

func findAcountByID(id string) *Account {
	for i, account := range accounts {
		if account.ID == id {
			return &accounts[i]
		}
	}
	return nil
}

func updateProfileByID(c net.Conn, id string, dataString string) {
	var profile PersonalDetails
	err := json.Unmarshal([]byte(dataString), &profile)
	if err != nil {
		fmt.Println("Error in unmarshal", err.Error())
	}
	account := findAcountByID(id)
	account.PersonalDetails = &profile
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(account)
	c.Write(buf.Bytes())
}

func makeTranscation(c net.Conn, id string, dataString string) {
	var transcation Transcation
	err := json.Unmarshal([]byte(dataString), &transcation)
	if err != nil {
		fmt.Println("Error in unmarshal", err.Error())
	}
	account := findAcountByID(id)
	account.Balance += transcation.Credit
	account.Balance -= transcation.Debit
	transcation.ClosingBalance = account.Balance
	account.Transcations = append(account.Transcations, transcation)
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(account)
	c.Write(buf.Bytes())
}
