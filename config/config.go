package config

import (
	"eventro2/models"

	"github.com/fatih/color"
)

var CurrentUser *models.User

const (
	UsersFile    = "data/users.json"
	EventsFile   = "data/events.json"
	ShowsFile    = "data/shows.json"
	VenuesFile   = "data/venues.json"
	BookingsFile = "data/bookings.json"
)

var WelcomeMessage string = color.MagentaString("\t Welcome to Eventro!\t")
var ChoiceMessage string = color.GreenString("Enter Choice: ")
var CustomerDashboardMessage string = color.BlueString("Customer Dashboard")
var BookedSeat string = color.RedString("[  ] ")
var AvailableSeat string = color.GreenString("[%s] ")
var StartAppMessage string = color.WhiteString("1. Login\n2. Signup\n3. Exit")
var DefaultChoiceMessage string = color.RedString("Please enter out of given choices")
var LogoutMessage string = color.YellowString("Goodbye!")
