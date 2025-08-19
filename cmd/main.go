package main

import (
	"context"
	"eventro2/config"
	"eventro2/controllers"
	"eventro2/models"
	artistrepository "eventro2/repository/artists_repository"
	bookingrepository "eventro2/repository/booking_repository"
	eventrepository "eventro2/repository/event_repository"
	showrepository "eventro2/repository/show_repository"
	userrepository "eventro2/repository/user_repository"
	venuerepository "eventro2/repository/venue_repository"
	"eventro2/services/artistservice"
	"eventro2/services/authorisation"
	"eventro2/services/bookingservice"
	"eventro2/services/eventservice"
	privilegeservice "eventro2/services/priviligeservice"
	"eventro2/services/searchevents"
	"eventro2/services/showservice"
	"eventro2/services/userservice"
	"eventro2/services/venueservice"
	utils "eventro2/utils/userinput"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("env not found")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	db.DisableForeignKeyConstraintWhenMigrating = true
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("failed to migrate users: %v", err)
	}

	err = db.AutoMigrate(&models.Venue{})
	if err != nil {
		log.Fatalf("failed to migrate venues: %v", err)
	}

	err = db.AutoMigrate(&models.Event{})
	if err != nil {
		log.Fatalf("failed to migrate events: %v", err)
	}

	err = db.AutoMigrate(&models.Show{})
	if err != nil {
		log.Fatalf("failed to migrate shows: %v", err)
	}

	err = db.AutoMigrate(&models.Booking{})
	if err != nil {
		log.Fatalf("failed to migrate bookings: %v", err)
	}

	err = db.AutoMigrate(&models.Artist{})
	if err != nil {
		log.Fatalf("failed to migrate artists: %v", err)
	}

	err = db.AutoMigrate(&models.EventArtist{})
	if err != nil {
		log.Fatalf("failed to migrate eventartists: %v", err)
	}

	fmt.Println(config.WelcomeMessage)
	ctx := context.Background()

	for {
		ctx = startApp(ctx, db)
		if config.GetUserID(ctx) == "" {
			fmt.Println(config.LoginErrorMessage)
			continue
		}
		break
	}
	launchDashboard(ctx, db)
}

func startApp(ctx context.Context, db *gorm.DB) context.Context {
	userRepo := userrepository.NewUserRepositoryPG(db)
	authService := authorisation.NewAuthService(userRepo)
	authController := controllers.NewAuthController(*authService)

	fmt.Println(config.StartAppMessage)
	for {
		fmt.Print(config.ChoiceMessage)

		choice, err := utils.TakeUserInput()
		if err != nil {
			continue
		}
		switch choice {
		case 1:
			return authController.LoginFlow(ctx)
		case 2:
			return authController.SignupFlow(ctx)
		case 3:
			fmt.Println(config.LogoutMessage)
			os.Exit(0)
		default:
			fmt.Println(config.DefaultChoiceMessage)
		}
	}
}

func launchDashboard(ctx context.Context, db *gorm.DB) {
	// repos
	showRepo := showrepository.NewShowRepositoryPG(db)
	eventRepo := eventrepository.NewEventRepositoryPG(db)
	bookingRepo := bookingrepository.NewBookingRepositoryPG(db)
	userRepo := userrepository.NewUserRepositoryPG(db)
	venueRepo := venuerepository.NewVenueRepositoryPG(db)
	artistRepo := artistrepository.NewArtistRepositoryPG(db)

	// services
	searchService := searchevents.NewSearchService(eventRepo, venueRepo, showRepo, bookingRepo, artistRepo)
	bookingService := bookingservice.NewBookingService(bookingRepo, showRepo, venueRepo, eventRepo)
	showService := showservice.NewShowService(showRepo, venueRepo, bookingRepo, eventRepo)
	eventService := eventservice.NewEventService(eventRepo)
	priviligeService := privilegeservice.NewPrivilegeService(userRepo)
	userService := userservice.NewUserService(userRepo)
	venueService := venueservice.NewVenueService(venueRepo)
	artistService := artistservice.NewArtistService(artistRepo)

	// controllers
	role := config.GetUserRole(ctx)
	switch role {
	case models.Admin:
		adminController := controllers.NewAdminController(priviligeService, eventService, userService, showService, searchService, artistService)
		adminController.ShowAdminDashboard(ctx)
	case models.Host:
		hostController := controllers.NewHostController(showService, venueService)
		hostController.ShowHostDashboard(ctx)
	case models.Customer:
		customerController := controllers.NewCustomerController(searchService, bookingService)
		customerController.ShowCustomerDashboard(ctx)
	default:
		fmt.Println(config.AccessMessage)
	}
}
