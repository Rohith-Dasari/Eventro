package controllers

import (
	"context"
	"eventro2/config"
	"eventro2/services/eventservice"
	privilegeservice "eventro2/services/priviligeservice"
	"eventro2/services/searchevents"
	"eventro2/services/showservice"
	"eventro2/services/userservice"
	utils "eventro2/utils/userinput"
	"fmt"
	"os"

	"github.com/fatih/color"
)

type AdminController struct {
	privilegeservice.PrivilegeService
	eventservice.EventService
	userservice.UserService
	showservice.ShowService
	searchevents.SearchService
}

func NewAdminController(p privilegeservice.PrivilegeService, e eventservice.EventService, u userservice.UserService, ss showservice.ShowService, se searchevents.SearchService) *AdminController {
	return &AdminController{p, e, u, ss, se}
}

func (ac *AdminController) ShowAdminDashboard(ctx context.Context) {
	for {
		fmt.Println(config.AdminDashboard)
		fmt.Println(config.AdminDashboardMSG)

		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()

		switch choice {
		case 1:
			ac.userModeration(ctx)
		case 2:
			ac.showModeration(ctx)
		case 3:
			ac.eventModeration(ctx)
		case 4:
			ac.PrivilegeService.EscalatePrivilege(ctx)
		case 5:
			ac.EventService.CreateNewEvent()
		case 6:
			ac.bookBehalfofUser(ctx)
		case 7:
			fmt.Println(config.LogoutMessage)
			os.Exit(0)
		}
	}
}

func (ac *AdminController) userModeration(ctx context.Context) {

	for {
		color.Blue(config.UserModeration)
		fmt.Println(config.UserModerationMenu)
		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()
		switch choice {
		case 1:
			ac.UserService.ModerateUser(ctx)
		case 2:
			ac.UserService.ViewBlockedUsers(ctx)
		case 3:
			return
		default:
			fmt.Println(config.InvalidMSG)
		}
	}
}

func (ac *AdminController) showModeration(ctx context.Context) {
	for {
		fmt.Println(config.ShowModeration)
		fmt.Println(config.ShowModerationMSG)
		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()
		switch choice {
		case 1:
			ac.ShowService.ModerateShow(ctx)
		case 2:
			ac.ShowService.ViewBlockedShows(ctx)
		case 3:
			return
		default:
			fmt.Println(config.InvalidMSG)
		}
	}
}
func (ac *AdminController) eventModeration(ctx context.Context) {
	for {
		fmt.Println(config.EventModeration)
		fmt.Println(config.EventModerationMSG)
		fmt.Println(config.ChoiceMessage)
		choice, _ := utils.TakeUserInput()
		switch choice {
		case 1:
			ac.EventService.ModerateEvents(ctx)
		case 2:
			ac.EventService.ViewBlockedEvents(ctx)
		case 3:
			return
		default:
			fmt.Println(config.InvalidMSG)
		}
	}
}

func (ac *AdminController) bookBehalfofUser(ctx context.Context) {
	ac.SearchService.Search(ctx)
}
