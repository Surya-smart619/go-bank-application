package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PersonalDetails struct {
	Name string
	Age  int
}

type Account struct {
	ID              string
	Type            string
	Balance         float64
	PersonalDetails PersonalDetails
	Transcations    []Transcation
}

type Transcation struct {
	Credit         float64
	Debit          float64
	ClosingBalance float64
}

type TcpRequest struct {
	Method string
	Url    string
	Data   string
	Params struct {
		ID string
	}
}

func main() {
	router := gin.Default()
	router.POST("/accounts", openAccount)
	router.GET("/accounts", getAccounts)
	router.GET("/accounts/:id", getAccountsByID)
	router.PATCH("/accounts/:id/profile", updateProfileByID)
	router.POST("/accounts/:id/transcations", makeTranscation)
	router.Run("localhost:8080")
}

func getConnection() (net.Conn, error) {
	return net.Dial("tcp", "localhost:8000")
}

func openAccount(c *gin.Context) {
	var newAccount Account

	if err := c.BindJSON(&newAccount); err != nil {
		return
	}
	conn, err := getConnection()
	if err != nil {
		fmt.Println("Error while get connection", err.Error())
	}
	accountJson, err := json.Marshal(newAccount)
	tcpReq := &TcpRequest{Method: "POST", Url: "/accounts", Data: string(accountJson)}
	err = gob.NewEncoder(conn).Encode(tcpReq)
	if err != nil {
		fmt.Println("Error while encoding", err)
	}

	//Reads Data from Server
	reader := bufio.NewReader(conn)
	data := &Account{}
	gob.NewDecoder(reader).Decode(data)
	c.IndentedJSON(http.StatusCreated, data)
}

func getAccounts(c *gin.Context) {
	conn, err := getConnection()
	if err != nil {
		fmt.Println("Error while get connection", err.Error())
	}
	tcpReq := &TcpRequest{Method: "GET", Url: "/accounts"}
	gob.NewEncoder(conn).Encode(tcpReq)

	//Reads Data from Server
	reader := bufio.NewReader(conn)
	data := []Account{}
	gob.NewDecoder(reader).Decode(&data)
	c.IndentedJSON(http.StatusOK, data)
}

func getAccountsByID(c *gin.Context) {
	conn, err := getConnection()
	if err != nil {
		fmt.Println("Error while get connection", err.Error())
	}
	tcpReq := &TcpRequest{Method: "GET", Url: "/accounts"}
	tcpReq.Params.ID = c.Param("id")
	gob.NewEncoder(conn).Encode(tcpReq)

	//Reads Data from Server
	reader := bufio.NewReader(conn)
	data := &Account{}
	gob.NewDecoder(reader).Decode(data)
	c.IndentedJSON(http.StatusCreated, data)
}

func updateProfileByID(c *gin.Context) {
	conn, err := getConnection()
	if err != nil {
		fmt.Println("Error while get connection", err.Error())
	}
	var profile PersonalDetails
	if err := c.BindJSON(&profile); err != nil {
		return
	}
	profileJson, err := json.Marshal(profile)
	tcpReq := &TcpRequest{Method: "PATCH", Url: "/profile", Data: string(profileJson)}
	tcpReq.Params.ID = c.Param("id")
	err = gob.NewEncoder(conn).Encode(tcpReq)
	if err != nil {
		fmt.Println("Error while encoding", err)
	}

	//Reads Data from Server
	reader := bufio.NewReader(conn)
	data := &Account{}
	gob.NewDecoder(reader).Decode(data)
	c.IndentedJSON(http.StatusOK, data)
}

func makeTranscation(c *gin.Context) {
	conn, err := getConnection()
	if err != nil {
		fmt.Println("Error while get connection", err.Error())
	}
	var transcation Transcation
	if err := c.BindJSON(&transcation); err != nil {
		return
	}
	transcationJson, err := json.Marshal(transcation)
	tcpReq := &TcpRequest{Method: "POST", Url: "/transcations", Data: string(transcationJson)}
	tcpReq.Params.ID = c.Param("id")
	err = gob.NewEncoder(conn).Encode(tcpReq)
	if err != nil {
		fmt.Println("Error while encoding", err)
	}

	//Reads Data from Server
	reader := bufio.NewReader(conn)
	data := &Account{}
	gob.NewDecoder(reader).Decode(data)
	c.IndentedJSON(http.StatusCreated, data)
}
