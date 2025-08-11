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

var (
	WelcomeMessage           = color.MagentaString("==== Welcome to Eventro! ====")
	ChoiceMessage            = color.CyanString("Enter Choice: ")
	CustomerDashboardMessage = color.BlueString("==== Customer Dashboard ====")
	BookedSeat               = color.RedString("[  ] ")
	AvailableSeat            = color.GreenString("[%s] ")
	StartAppMessage          = color.WhiteString("1. Login\n2. Signup\n3. Exit")
	DefaultChoiceMessage     = color.RedString("Please chooose one of the given options")
	LogoutMessage            = color.YellowString("Goodbye!")
	LoginErrorMessage        = color.RedString("\nYou no longer have access to this account. Please contact Admin")
	AccessMessage            = color.RedString("Access Denied")
	WelcomeBack              = color.GreenString("\nWelcome back %s! You are logged in as: %s\n")
	CustomerOptions          = "1. Search\n2. Booking History\n3. Logout"
	SearchEventsMessage      = color.CyanString("Select how you want to search")
	InvalidMSG               = color.RedString("Invalid choice.")
	Dash                     = color.MagentaString("-------------")
	AdminDashboard           = color.BlueString("\n==== Admin Dashboard ====")
	HostDashboard            = color.BlueString("\n==== Host Dashboard ====")
	ShowModerationMSG        = "1. Block/Unblock Show\n2. View Blocked Shows\n3. Back"
	EventModeration          = "1. Block/Unblock Event\n2. View Blocked Events\n3. Back"
)
