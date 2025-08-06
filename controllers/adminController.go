package controllers

import (
	"bufio"
	"context"
	"eventro/config"
	"eventro/models"
	"eventro/services/eventservice"
	"eventro/services/showservice"
	"eventro/services/userservice"
	"eventro/storage"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ShowAdminDashboard(ctx context.Context) {
	fmt.Println("\n=== Admin Dashboard ===")
	//view users, block users, unblock users, assign new host--done
	//view show, block show/event, unblock show/event,-done
	//book a show behalf of user
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("1. User Moderation 2. Show Moderation 3. Event Moderation 4. Add a new host 5. Booking Request/Book behalf of user")
	var choice int
	fmt.Println(config.ChoiceMessage)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	choice, _ = strconv.Atoi(input)

	switch choice {
	case 1:
		//create class of userservice class
		fmt.Println("1. block/unlock user 2. view blocked user ")
		fmt.Println(config.ChoiceMessage)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, _ = strconv.Atoi(input)
		switch choice {
		case 1:
			userservice.ModerateUser(ctx)
		case 2:
			userservice.ViewBlockedUsers(ctx)
		}
	case 2:
		fmt.Println("1. block/unblock show 2. view blocked show")
		fmt.Println(config.ChoiceMessage)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, _ = strconv.Atoi(input)
		switch choice {
		case 1:
			showservice.ModerateShow(ctx)
		case 2:
			showservice.ViewBlockedShows(ctx)
		}
	case 3:
		fmt.Println("1. block/unblock event 2. view blocked event")
		fmt.Println(config.ChoiceMessage)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, _ = strconv.Atoi(input)
		switch choice {
		case 1:
			eventservice.ModerateEvents(ctx)
		case 2:
			eventservice.ViewBlockedEvents(ctx)
		}
	case 4:
		fmt.Println("Enter user mail id whose privilige you want to escalate ")
		var usermail string
		fmt.Scanf("%s", &usermail)
		users := storage.LoadUsers()
		var requiredUser *models.User
		for _, user := range users {
			if user.Email == usermail {
				requiredUser = &user
			}
		}
		requiredUser.Role = models.Host
	}

}
