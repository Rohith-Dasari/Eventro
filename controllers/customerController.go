package controllers

import (
	"context"
	"eventro2/config"
	"eventro2/services/bookingservice"
	"eventro2/services/searchevents"
	utils "eventro2/utils/userinput"
	"fmt"
)

type CustomerController struct {
	searchevents.SearchService
	bookingservice.BookingService
}

func NewCustomerController(se searchevents.SearchService, bs bookingservice.BookingService) *CustomerController {
	return &CustomerController{se, bs}
}

func (cc *CustomerController) ShowCustomerDashboard(ctx context.Context) {
	for {
		// showRepo := showrepository.NewShowRepository()
		// eventRepo := eventsrepository.NewEventRepository()
		// bookingRepo := bookingrepository.NewBoookingStore()
		// searchService := searchevents.NewSearchService(*eventRepo)
		// bookingService := bookingservice.NewBookingService(*bookingRepo, *showRepo)

		fmt.Println(config.CustomerDashboardMessage)
		fmt.Println(config.CustomerOptions)
		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()

		switch choice {
		case 1:
			cc.SearchService.Search(ctx)

		case 2:
			cc.BookingService.ViewBookingHistory(ctx)

		case 3:
			fmt.Println(config.LogoutMessage)
			return

		default:
			fmt.Println(config.DefaultChoiceMessage)
		}
	}
}
