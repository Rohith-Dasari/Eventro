package bookingservice

import (
	"bufio"
	"context"
	"eventro2/config"
	"eventro2/models"
	bookingrepository "eventro2/repository/booking_repository"
	showrepository "eventro2/repository/show_repository"
	utils "eventro2/utils/userinput"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type BookingService struct {
	BookingRepo bookingrepository.BookingRepository
	ShowRepo    showrepository.ShowRepository
}

func NewBookingService(bRepo bookingrepository.BookingRepository, sRepo showrepository.ShowRepository) *BookingService {
	return &BookingService{
		BookingRepo: bRepo,
		ShowRepo:    sRepo,
	}
}

func (bs *BookingService) ViewBookingHistory(ctx context.Context) {
	bookings := bs.BookingRepo.Bookings
	fmt.Println("Your Events:")
	for _, booking := range bookings {
		if booking.UserID == config.GetUserID(ctx) {
			bs.printBooking(booking)
		}
	}
}

func (bs *BookingService) MakeBooking(ctx context.Context, userID string, showID string) {
	shows := bs.ShowRepo.Shows

	var requiredShow *models.Show
	for i := range shows {
		if shows[i].ID == showID {
			requiredShow = &shows[i]
			break
		}
	}
	if requiredShow == nil {
		fmt.Println("Error: Show not found")
		return
	}

	bookedTickets := requiredShow.BookedSeats

	var numTickets int
	for {
		fmt.Println("How many tickets do you want to book?")
		var err error
		numTickets, err = utils.TakeUserInput()
		if err == nil {
			break
		}
	}

	userTickets := make([]string, numTickets)
	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < numTickets; {
		fmt.Printf("Enter seat %d: ", i+1)
		input, _ := reader.ReadString('\n')
		userTicket := strings.TrimSpace(input)

		if !bs.validateTicket(userTicket, bookedTickets) {
			fmt.Println("Entered ticket is either invalid or already booked.")
			continue
		}

		userTickets[i] = userTicket
		i++
	}

	fmt.Println("Your selected seats:", userTickets)
	totalPrice := numTickets * int(requiredShow.Price)
	fmt.Printf("Total price: ₹%d\n", totalPrice)

	fmt.Print("Confirm booking? (y/n): ")
	choice, _ := reader.ReadString('\n')
	if strings.TrimSpace(choice) != "y" {
		fmt.Println("Booking canceled.")
		return
	}

	requiredShow.BookedSeats = append(requiredShow.BookedSeats, userTickets...)
	bs.ShowRepo.SaveShows(shows)

	bookings := bs.BookingRepo.Bookings
	newBooking := models.Booking{
		BookingID:         uuid.New().String(),
		UserID:            userID,
		ShowID:            showID,
		TimeBooked:        time.Now().Format(time.RFC3339),
		NumTickets:        numTickets,
		TotalBookingPrice: float64(totalPrice),
		Seats:             userTickets,
	}
	bookings = append(bookings, newBooking)
	bs.BookingRepo.SaveBookings(bookings)

	fmt.Println("🎉 Congratulations! Your booking is confirmed. Here's your ticket:")
	bs.printBooking(newBooking)
}

func (bs *BookingService) printBooking(b models.Booking) {
	fmt.Println("-------------")
	fmt.Printf("Booking ID        : %s\n", b.BookingID)
	fmt.Printf("Time Booked       : %s\n", b.TimeBooked)
	fmt.Printf("Tickets Booked    : %d\n", b.NumTickets)
	fmt.Printf("Total Price       : ₹%.2f\n", b.TotalBookingPrice)
	fmt.Printf("Seats             : %v\n", b.Seats)
	fmt.Println("-------------")
}

func (bs *BookingService) validateTicket(userTicket string, bookedTickets []string) bool {
	userTicket = strings.ToUpper(userTicket)

	matched, err := regexp.MatchString(`^[A-J](10|[1-9])$`, userTicket)
	if err != nil || !matched {
		return false
	}

	for _, ticket := range bookedTickets {
		if strings.ToUpper(ticket) == userTicket {
			return false
		}
	}
	return true
}
