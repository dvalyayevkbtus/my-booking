package booking

import "dvalyayevkbtu/my-booking/client"

type BookingRepr struct {
	Id        int64             `json:"id"`
	HotelName string            `json:"hotelName"`
	Price     string            `json:"price"`
	Currency  string            `json:"currency"`
	Client    client.ClientRepr `json:"client"`
	Paid      bool              `json:"paid"`
}

type Book struct {
	HotelName string `json:"hotelName"`
	Price     string `json:"price"`
	Currency  string `json:"currency"`
	ClientId  int64  `json:"clientId"`
}
