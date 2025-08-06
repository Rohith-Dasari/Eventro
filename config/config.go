package config

import (
	"eventro/models"

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

var WelcomeMessage string = color.MagentaString("🥳🥳🥳 \t Welcome to Eventro!\t🥳🥳🥳")
var ChoiceMessage string = color.GreenString("Enter Choice: ")
var CustomerDashboardMessage string = color.BlueString("Customer Dashboard")
var BookedSeat string = color.RedString("[❌] ")
var AvailableSeat string = color.GreenString("[%s] ")
