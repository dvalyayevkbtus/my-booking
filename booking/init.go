package booking

import (
	"dvalyayevkbtu/my-booking/client"
	"dvalyayevkbtu/my-booking/db"
	"dvalyayevkbtu/my-booking/payment"
	"dvalyayevkbtu/my-booking/utils"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type Booking struct {
	db      *db.BookingDb
	payment *payment.Payment
}

func Init(db *db.BookingDb, payment *payment.Payment) *Booking {
	return &Booking{db, payment}
}

func (b *Booking) HandleBookings(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		b.retrieveBookings(rw, req)
		return
	}
	if req.Method == http.MethodPost {
		b.performBook(rw, req)
		return
	}
	utils.MethodNotAllowed(rw)
}

func (b *Booking) retrieveBookings(rw http.ResponseWriter, req *http.Request) {
	bookings, err := b.db.GetAllBookings()
	if err != nil {
		log.Debug("Cannot get all bookings!")
		utils.InternalServerError(rw)
		return
	}

	res := make([]BookingRepr, 0)
	for _, dbBook := range bookings {
		clientF, err := b.db.GetClient(dbBook.ClientId)
		if err != nil {
			log.Errorf("Cannot get client! %v", err)
			utils.InternalServerError(rw)
			return
		}

		payment, err := b.db.GetPayment(dbBook.Id)
		if err != nil {
			log.Errorf("Cannot get payment! %v", err)
			utils.InternalServerError(rw)
			return
		}

		repr := BookingRepr{dbBook.Id, dbBook.HotelName, dbBook.Price, dbBook.Currency,
			client.ClientRepr{Id: clientF.Id, FullName: clientF.FullName}, payment.Status == db.PaymentFulfilled}
		res = append(res, repr)
	}

	marshalled, err := json.Marshal(res)
	if err != nil {
		log.Errorf("Cannot marshal booking! %v", err)
		utils.InternalServerError(rw)
		return
	}
	utils.SuccessString(rw, string(marshalled))
}

func (b *Booking) performBook(rw http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Warnf("Cannot read body! %v", err)
		utils.BadRequest(rw)
		return
	}

	var book Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		log.Warnf("Cannot read body! %v", err)
		utils.BadRequest(rw)
		return
	}

	_, err = b.db.GetClient(book.ClientId)
	if err != nil {
		log.Warnf("Invalid client id! %v", err)
		utils.BadRequest(rw)
		return
	}

	id, err := b.db.RegisterBooking(book.HotelName, book.Price, book.Currency, book.ClientId)
	if err != nil {
		log.Errorf("Cannot save booking! %v", err)
		utils.InternalServerError(rw)
		return
	}
	paymentReference, err := b.db.CreatePayment(id)
	if err != nil {
		log.Errorf("Cannot create payment! %v", err)
		utils.InternalServerError(rw)
		return
	}
	reference := strconv.FormatInt(paymentReference, 10)
	err = b.payment.CreateInvoice(reference, book.Price, book.Currency)
	if err != nil {
		log.Errorf("Cannot create invoice! %v", err)
		utils.InternalServerError(rw)
		return
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			resp, err := b.payment.CheckPayment(reference)
			if err != nil {
				log.Error(err)
			}
			if resp {
				err = b.db.UpdatePayment(paymentReference, db.PaymentFulfilled)
				if err != nil {
					log.Error(err)
				}
				return
			}
		}
	}()
}

func (b *Booking) HandleBooking(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		utils.MethodNotAllowed(rw)
		return
	}

	id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		log.Debug("Invalid id!")
		utils.BadRequest(rw)
		return
	}

	booking, err := b.db.GetBooking(id)
	if err != nil {
		log.Errorf("Cannot get booking! %v", err)
		utils.InternalServerError(rw)
		return
	}

	clientF, err := b.db.GetClient(booking.ClientId)
	if err != nil {
		log.Errorf("Cannot get client! %v", err)
		utils.InternalServerError(rw)
		return
	}

	payment, err := b.db.GetPayment(booking.Id)
	if err != nil {
		log.Errorf("Cannot get payment! %v", err)
		utils.InternalServerError(rw)
		return
	}

	result := BookingRepr{booking.Id, booking.HotelName, booking.Price, booking.Currency,
		client.ClientRepr{Id: clientF.Id, FullName: clientF.FullName}, payment.Status == db.PaymentFulfilled}
	marshalled, err := json.Marshal(result)
	if err != nil {
		log.Errorf("Cannot marshal booking! %v", err)
		utils.InternalServerError(rw)
		return
	}
	utils.SuccessString(rw, string(marshalled))
}
