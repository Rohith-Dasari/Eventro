package app

import (
	"eventro2/handlers"
	"eventro2/middleware"
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
	"eventro2/services/showservice"
	"eventro2/services/userservice"
	"eventro2/services/venueservice"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, proceeding with environment variables")
	}

	// Get DB URL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Panic("DATABASE_URL is not set")
	}

	// Connect to DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("Failed to connect to database:", err)
	}

	// AutoMigrate
	err = db.AutoMigrate(
		&models.User{},
		&models.Booking{},
		&models.Venue{},
		&models.Event{},
		&models.Show{},
		&models.Artist{},
		&models.EventArtist{},
	)
	if err != nil {
		log.Panic("Migration failed:", err)
	}

	return db
}

func SetupServer(DB *gorm.DB) http.Handler {
	// repos
	showRepo := showrepository.NewShowRepositoryPG(DB)
	eventRepo := eventrepository.NewEventRepositoryPG(DB)
	bookingRepo := bookingrepository.NewBookingRepositoryPG(DB)
	userRepo := userrepository.NewUserRepositoryPG(DB)
	venueRepo := venuerepository.NewVenueRepositoryPG(DB)
	artistRepo := artistrepository.NewArtistRepositoryPG(DB)

	// services
	bookingService := bookingservice.NewBookingService(bookingRepo, showRepo, venueRepo, eventRepo)
	showService := showservice.NewShowService(showRepo, venueRepo, bookingRepo, eventRepo)
	eventService := eventservice.NewEventService(eventRepo)
	userService := userservice.NewUserService(userRepo)
	venueService := venueservice.NewVenueService(venueRepo)
	artistService := artistservice.NewArtistService(artistRepo)
	authService := authorisation.NewAuthService(userRepo)

	//handlers
	authHandler := &handlers.AuthHandler{
		AuthService: authService,
	}
	artistHandler := &handlers.ArtistHandler{
		ArtistService: &artistService,
	}
	eventHandler := &handlers.EventHandler{
		EventService: &eventService,
	}
	showHandler := &handlers.ShowHandler{
		ShowService: &showService,
	}
	venueHandler := &handlers.VenueHandler{
		VenueService: &venueService,
	}
	bookingHandler := &handlers.BookingHandler{
		BookingService: &bookingService,
	}
	userHandler := &handlers.UserHandler{
		UserService: &userService,
	}

	//setup mux
	mux := http.NewServeMux()

	//auth
	mux.HandleFunc("/api/v1/login", authHandler.Login)
	mux.HandleFunc("/api/v1/signup", authHandler.Signup)

	//artists
	mux.HandleFunc("GET /api/v1/artists", middleware.JWTAuth(artistHandler.BrowseArtists))          // GET
	mux.HandleFunc("POST /api/v1/artists", middleware.JWTAuth(artistHandler.CreateArtist))          // POST
	mux.Handle("DELETE /api/v1/artists/{artistID}", middleware.JWTAuth(artistHandler.DeleteArtist)) // DELETE

	//events
	mux.Handle("GET /api/v1/events", middleware.JWTAuth(eventHandler.BrowseEvents))             //browse
	mux.Handle("POST /api/v1/events", middleware.JWTAuth(eventHandler.CreateEvent))             //create
	mux.Handle("PATCH /api/v1/events/{eventID}", middleware.JWTAuth(eventHandler.UpdateEvent))  //patch
	mux.Handle("DELETE /api/v1/events/{eventID}", middleware.JWTAuth(eventHandler.DeleteEvent)) //delete
	mux.Handle("GET /api/v1/hosts/{hostID}/events", middleware.JWTAuth(eventHandler.EventsOfHost))

	//shows
	mux.Handle("GET /api/v1/shows", middleware.JWTAuth(showHandler.BrowseShows))            // browse
	mux.Handle("POST /api/v1/shows", middleware.JWTAuth(showHandler.CreateShow))            //create
	mux.Handle("PATCH /api/v1/shows/{showID}", middleware.JWTAuth(showHandler.UpdateShow))  // PATCH
	mux.Handle("DELETE /api/v1/shows/{showID}", middleware.JWTAuth(showHandler.DeleteShow)) // DELETE

	//venues
	mux.Handle("POST /api/v1/venues", middleware.JWTAuth(venueHandler.CreateVenue))
	mux.Handle("GET /api/v1/venues", middleware.JWTAuth(venueHandler.BrowseVenues))
	mux.Handle("DELETE /api/v1/venues/{venueID}", middleware.JWTAuth(venueHandler.DeleteVenue))
	mux.Handle("PATCH /api/v1/venues/{venueID}", middleware.JWTAuth(venueHandler.UpdateVenue))

	//users
	mux.Handle("GET /api/v1/users", middleware.JWTAuth(userHandler.BrowseUsers))
	mux.Handle("PATCH /api/v1/users/", middleware.JWTAuth(userHandler.UpdateUser))
	mux.Handle("GET /api/v1/{userID}/profile", middleware.JWTAuth(userHandler.GetUserProfile))
	mux.Handle("GET /api/v1/users/email/{email}", middleware.JWTAuth(userHandler.GetUserByMailID))

	//bookings
	mux.Handle("POST /api/v1/bookings", middleware.JWTAuth(bookingHandler.CreateBooking))
	mux.Handle("GET /api/v1/bookings", middleware.JWTAuth(bookingHandler.BrowseBookings))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})
	handler := c.Handler(mux)

	return handler
}

func Start() {
	DB := InitDB()
	mux := SetupServer(DB)

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
