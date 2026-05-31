package main

import (
	"dvalyayevkbtu/my-booking/booking"
	"dvalyayevkbtu/my-booking/client"
	"dvalyayevkbtu/my-booking/config"
	"dvalyayevkbtu/my-booking/db"
	"dvalyayevkbtu/my-booking/logging"
	"dvalyayevkbtu/my-booking/payment"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	logging.SetupLogger()

	conf, cErr := config.InitConfig()
	if cErr != nil {
		log.Errorf("Error on loading config! %v", cErr)
		panic(cErr)
	}
	db, dbErr := db.InitDatabase(conf.DB)
	if dbErr != nil {
		log.Errorf("Error on initializing database! %v", dbErr)
		panic(dbErr)
	}
	payment := payment.CreatePayment(conf.Payment)

	cli := client.Init(db)
	http.HandleFunc("/client", cli.HandleClients)

	book := booking.Init(db, payment)
	http.HandleFunc("/booking", book.HandleBookings)
	http.HandleFunc("/booking/{id}", book.HandleBooking)

	log.Info("My booking started!")
	http.ListenAndServe(":8080", nil)
}
