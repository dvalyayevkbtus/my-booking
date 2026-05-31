package client

import (
	"dvalyayevkbtu/my-booking/db"
	"dvalyayevkbtu/my-booking/utils"
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	db *db.BookingDb
}

func Init(db *db.BookingDb) *Client {
	return &Client{db}
}

func (c *Client) HandleClients(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		c.retrieveClients(rw, req)
		return
	}
	if req.Method == http.MethodPost {
		c.createClient(rw, req)
		return
	}
	utils.MethodNotAllowed(rw)
}

func (c *Client) retrieveClients(rw http.ResponseWriter, req *http.Request) {
	res, err := c.db.GetAllClients()
	if err != nil {
		log.Errorf("Cannot get all clients! %v", err)
		utils.InternalServerError(rw)
		return
	}

	clients := make([]ClientRepr, 0)
	for _, r := range res {
		clients = append(clients, ClientRepr{r.Id, r.FullName})
	}

	marshalled, err := json.Marshal(clients)
	if err != nil {
		log.Errorf("Cannot marshall all client! %v", err)
		utils.InternalServerError(rw)
		return
	}

	utils.SuccessString(rw, string(marshalled))
}

func (c *Client) createClient(rw http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Errorf("Cannot read body! %v", err)
		utils.InternalServerError(rw)
		return
	}
	var create CreateClient
	err = json.Unmarshal(body, &create)
	if err != nil {
		log.Errorf("Invalid body! %v", err)
		utils.BadRequest(rw)
		return
	}
	err = c.db.ClientInsert(create.FullName)
	if err != nil {
		log.Errorf("Cannot insert client! %v", err)
		utils.InternalServerError(rw)
		return
	}
	utils.Accepted(rw)
}
