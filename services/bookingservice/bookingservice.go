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
	// bookings := bs.BookingRepo.Bookings
	bookings2 := bookingrepository.NewBoookingStore()
	fmt.Println("Your Events:")
	for _, booking := range bookings2.Bookings {
		if booking.UserID == config.GetUserID(ctx) {
			bs.printBooking(booking)
		}
	}
}

func (bs *BookingService) MakeBooking(ctx context.Context, userID string, showID string) {
	shows := bs.ShowRepo.Shows

	var requiredShow *models.Show
	var requiredIndex int
	for i := range shows {
		if shows[i].ID == showID {
			requiredShow = &shows[i]
			requiredIndex = i
			break
		}
	}
	if requiredShow == nil {
		color.Red("Error: Show not found.")
		return
	}

	venueRepo := venuerepository.NewVenueRepository()
	var requiredVenue models.Venue
	for _, venue := range venueRepo.Venues {
		if venue.ID == requiredShow.VenueID {
			requiredVenue = venue
			break
		}
	}

	var numTickets int
	for {
		fmt.Println("How many tickets do you want to book?")
		var err error
		numTickets, err = utils.TakeUserInput()
		if err != nil {
			color.Red("Invalid input. Please enter a valid number.")
			continue
		}
		if numTickets <= 0 || numTickets > 10 {
			color.Red("Number of tickets must be greater than zero.")
			continue
		}
		if numTickets > 10 {
			color.Red("A user can only book upto 10 tickets")
			continue
		}
		break
	}

	userTickets := make([]string, numTickets)
	totalPrice := numTickets * int(requiredShow.Price)

	reader := bufio.NewReader(os.Stdin)

	if !requiredVenue.IsSeatLayoutRequired {
		fmt.Printf("Total price: ₹%d\n", totalPrice)
		fmt.Println("Confirm booking? 1. Yes\n2. No: ")
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
		bs.ShowRepo.Shows[requiredIndex].BookedSeats = append(bs.ShowRepo.Shows[requiredIndex].BookedSeats, userTickets...)

		requiredShow.BookedSeats = append(requiredShow.BookedSeats, userTickets...)
		bs.ShowRepo.SaveShows(bs.ShowRepo.Shows)
	}

	newBooking := models.Booking{
		BookingID:         uuid.New().String(),
		UserID:            userID,
		ShowID:            showID,
		TimeBooked:        time.Now().Format(time.RFC3339),
		NumTickets:        numTickets,
		TotalBookingPrice: float64(totalPrice),
		Seats:             userTickets,
	}

	bs.BookingRepo.Bookings = append(bs.BookingRepo.Bookings, newBooking)
	bs.BookingRepo.SaveBookings(bs.BookingRepo.Bookings)

	color.Green("🎉 Congratulations! Your booking is confirmed. Here's your ticket:")
	bs.printBooking(newBooking)
}

func (bs *BookingService) printBooking(b models.Booking) {
	//i have booking, booking has show id, show id has event id and venue id so first required show needs to be caught

	shows := bs.ShowRepo.Shows
	var requiredShow models.Show
	for _, show := range shows {
		if show.ID == b.ShowID {
			requiredShow = show
		}
	}
	eventRepo := eventsrepository.NewEventRepository()
	venuRepository := venuerepository.NewVenueRepository()

	var requiredEvent models.Event
	for _, event := range eventRepo.Events {
		if event.ID == requiredShow.EventID {
			requiredEvent = event
			break
		}
	}

	var requiredVenue models.Venue
	for _, venue := range venuRepository.Venues {
		if venue.ID == requiredShow.VenueID {
			requiredVenue = venue
			break
		}
	}
	fmt.Println(config.Dash)
	fmt.Printf("Booking ID    \t: %s\n", b.BookingID)
	fmt.Printf("Event Name    \t: %s\n", requiredEvent.Name)
	fmt.Printf("Venue Name    \t: %s\n", requiredVenue.Name)
	fmt.Printf("Venue City    \t: %s\n", requiredVenue.City)
	fmt.Printf("Time Booked   \t: %s\n", b.TimeBooked)
	fmt.Printf("Tickets Booked\t: %d\n", b.NumTickets)
	fmt.Printf("Total Price   \t: ₹%.2f\n", b.TotalBookingPrice)
	fmt.Printf("Seats         \t: %v\n", b.Seats)
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
