package bookingservice

import (
	"bufio"
	"context"
	"eventro2/config"
	"eventro2/models"
	bookingrepository "eventro2/repository/booking_repository"
	eventsrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	venuerepository "eventro2/repository/venue_repository"
	utils "eventro2/utils/userinput"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

type BookingService struct {
	BookingRepo bookingrepository.BookingRepository
	ShowRepo    showrepository.ShowRepository
	venueRepo   venuerepository.VenueRepository
	eventRepo   eventsrepository.EventRepository
}

func NewBookingService(bRepo bookingrepository.BookingRepository, sRepo showrepository.ShowRepository, vRepo venuerepository.VenueRepository, eRepo eventsrepository.EventRepository) BookingService {
	return BookingService{
		BookingRepo: bRepo,
		ShowRepo:    sRepo,
		venueRepo:   vRepo,
		eventRepo:   eRepo,
	}
}

func (bs *BookingService) ViewBookingHistory(ctx context.Context) {
	userID := config.GetUserID(ctx)

	bookings, err := bs.BookingRepo.ListByUser(userID)
	if err != nil {
		fmt.Println("Error fetching bookings:", err)
		return
	}

	if len(bookings) == 0 {
		fmt.Println("No bookings found.")
		return
	}

	fmt.Println("Your Events:")
	for _, booking := range bookings {
		bs.printBooking(booking)
	}
}

func (bs *BookingService) MakeBooking(ctx context.Context, userID string, showID string) {
	requiredShow, err := bs.ShowRepo.GetByID(showID)
	if err != nil || requiredShow == nil {
		color.Red("Error: Show not found.")
		return
	}

	requiredVenue, err := bs.venueRepo.GetByID(requiredShow.VenueID)
	if err != nil || requiredVenue == nil {
		color.Red("Error: Venue not found.")
		return
	}

	var numTickets int
	for {
		fmt.Println("How many tickets do you want to book?")
		numTickets, err = utils.TakeUserInput()
		if err != nil {
			color.Red("Invalid input. Please enter a valid number.")
			continue
		}
		if numTickets <= 0 || numTickets > 10 {
			color.Red("Number of tickets must be between 1 and 10.")
			continue
		}
		break
	}

	userTickets := make([]string, numTickets)
	totalPrice := numTickets * int(requiredShow.Price)

	reader := bufio.NewReader(os.Stdin)

	if !requiredVenue.IsSeatLayoutRequired {
		fmt.Printf("Total price: ₹%d\n", totalPrice)
		fmt.Println("Confirm booking?\n1. Yes\n2. No: ")
		choice, _ := utils.TakeUserInput()
		if choice != 1 {
			color.Red("Booking canceled.")
			fmt.Println("But it is more fun to book :)")
			return
		}
	} else {
		bookedTickets := requiredShow.BookedSeats

		for i := 0; i < numTickets; {
			fmt.Printf("Enter seat %d: ", i+1)
			input, _ := reader.ReadString('\n')
			userTicket := strings.ToUpper(strings.TrimSpace(input))

			if !bs.validateTicket(userTicket, bookedTickets) {
				color.Red("Entered ticket is either invalid or already booked.")
				continue
			}

			userTickets[i] = userTicket
			i++
		}

		fmt.Println("Your selected seats:", userTickets)
		fmt.Printf("Total price: ₹%d\n", totalPrice)
		fmt.Println("Confirm booking? \n1. Yes\n2. No: ")
		choice, _ := utils.TakeUserInput()
		if choice != 1 {
			color.Red("Booking canceled.")
			fmt.Println("But it is more fun to book :)")
			return
		}

		requiredShow.BookedSeats = append(requiredShow.BookedSeats, userTickets...)
		if err := bs.ShowRepo.Update(requiredShow); err != nil {
			color.Red("Error updating show seats:", err)
			return
		}
	}
	newBooking := &models.Booking{
		UserID:            userID,
		ShowID:            showID,
		NumTickets:        numTickets,
		TotalBookingPrice: float64(totalPrice),
		Seats:             userTickets,
	}

	if err := bs.BookingRepo.Create(newBooking); err != nil {
		color.Red("Error creating booking:", err)
		return
	}

	color.Green("🎉 Congratulations! Your booking is confirmed. Here's your ticket:")
	bs.printBooking(*newBooking)
}

func (bs *BookingService) printBooking(b models.Booking) {
	requiredShow, err := bs.ShowRepo.GetByID(b.ShowID)
	if err != nil || requiredShow == nil {
		color.Red("Error: Show not found for booking.")
		return
	}

	requiredEvent, err := bs.eventRepo.GetByID(requiredShow.EventID)
	if err != nil || requiredEvent == nil {
		color.Red("Error: Event not found for show.")
		return
	}

	requiredVenue, err := bs.venueRepo.GetByID(requiredShow.VenueID)
	if err != nil || requiredVenue == nil {
		color.Red("Error: Venue not found for show.")
		return
	}

	fmt.Println(config.Dash)
	fmt.Printf("%-16s : %s\n", "Booking ID", b.BookingID)
	fmt.Printf("%-16s : %s\n", "Event Name", requiredEvent.Name)
	fmt.Printf("%-16s : %s\n", "Venue Name", requiredVenue.Name)
	fmt.Printf("%-16s : %s\n", "Venue City", requiredVenue.City)
	fmt.Printf("%-16s : %s\n", "Time Booked", b.TimeBooked.Format(time.RFC1123))
	fmt.Printf("%-16s : %d\n", "Tickets Booked", b.NumTickets)
	fmt.Printf("%-16s : ₹%.2f\n", "Total Price", b.TotalBookingPrice)
	fmt.Printf("%-16s : %v\n", "Seats", b.Seats)
	fmt.Println(config.Dash)
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
