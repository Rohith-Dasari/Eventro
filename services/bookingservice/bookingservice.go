package bookingservice

import (
	"eventro/config"
	"eventro/models"
	"eventro/storage"
	"fmt"
)

func SeeBookingHistory() {
	bookings := storage.LoadBookings()
	fmt.Println("Your Events:")
	for _, booking := range bookings {
		if booking.UserID == config.CurrentUser.UserID {
			printBooking(booking)
		}
	}
}

func printBooking(b models.Booking) {
	fmt.Println("-------------")
	fmt.Printf("Booking ID        : %d\n", b.BookingID)
	fmt.Printf("Time Booked       : %s\n", b.TimeBooked)
	fmt.Printf("Tickets Booked    : %d\n", b.NumTickets)
	fmt.Printf("Total Price       : ₹%.2f\n", b.TotalBookingPrice)
	fmt.Printf("Seats             : %v\n", b.Seats)
	fmt.Println("-------------")
}
