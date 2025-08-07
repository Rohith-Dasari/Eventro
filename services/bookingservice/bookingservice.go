package bookingservice

import (
	"bufio"
	"eventro/config"
	"eventro/models"
	"eventro/storage"
	utils "eventro/utils/userinput"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
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
	fmt.Printf("Booking ID        : %s\n", b.BookingID)
	fmt.Printf("Time Booked       : %s\n", b.TimeBooked)
	fmt.Printf("Tickets Booked    : %d\n", b.NumTickets)
	fmt.Printf("Total Price       : ₹%.2f\n", b.TotalBookingPrice)
	fmt.Printf("Seats             : %v\n", b.Seats)
	fmt.Println("-------------")
}

func MakeBooking(UserID string, showID string) {
	var requiredShow *models.Show
	shows := storage.LoadShows()
	for _, show := range shows {
		if show.ID == showID {
			requiredShow = &show
			break
		}
	}
	if requiredShow == nil {
		fmt.Println("Error: Show not found")
		return
	}
	bookedTickets := requiredShow.BookedSeats
	var numTickets int
	var err error
	for {
		fmt.Println("How many tickets do you want to book?")
		numTickets, err = utils.TakeUserInput()
		if err != nil {
			continue
		}
		break
	}
	var userTickets []string = make([]string, numTickets)
	for i := 0; i < numTickets; {
		fmt.Println("Enter seat ", i+1)
		var userTicket string
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		userTicket = strings.TrimSpace(input)
		isValid := ValidateTicket(userTicket, bookedTickets)
		if !isValid {
			fmt.Println("Entered ticket is either invalid or already booked")
		} else {
			userTickets[i] = userTicket
			i++
		}
	}
	fmt.Println("your selected seats are: ", userTickets)
	fmt.Println("total price: ", numTickets*int(requiredShow.Price))
	fmt.Println("confirm booking? y/n")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(input)
	if choice == "y" {
		//unmarshal shows, append list to ticket seats
		requiredShow.BookedSeats = append(requiredShow.BookedSeats, userTickets...)
		//marshall and write to file
		storage.SaveShows(shows)
		//unmarshall, append, marshall bookings
		bookings := storage.LoadBookings()
		newBooking := &models.Booking{
			BookingID:         uuid.New().String(),
			UserID:            UserID,
			ShowID:            showID,
			TimeBooked:        time.Now().Format(time.RFC3339),
			NumTickets:        numTickets,
			TotalBookingPrice: float64(numTickets) * requiredShow.Price,
			Seats:             userTickets,
		}
		bookings = append(bookings, *newBooking)
		storage.SaveBookings(bookings)
		fmt.Println("Congratulations! your booking is confirmed here is your ticket!")
		printBooking(*newBooking)
		//display qr code
		//fmt.Println("do you want to create a ics file?")
	}

}

func ValidateTicket(userTicket string, bookedTickets []string) bool {
	//check if userTicket is a alphabet followe by a number, and then check if it is not present in bookedTicket

	userTicket = strings.ToUpper(userTicket)

	matched, err := regexp.MatchString(`(A|B|C|D|E|F|G|H|I|J)(10|[1-9])$`, userTicket)
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
