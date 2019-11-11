package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/machinebox/graphql"
)

type UsersResponse struct {
	Users struct {
		Edges []struct {
			Node struct {
				ID           string
				Email        string
				CheckinDate  time.Time
				CheckoutDate time.Time
			}
		}
	}
}

type EventsResponse struct {
	Events []struct {
		User string
		Type string
	}
}

func getAllUsers(client *graphql.Client) (users UsersResponse) {
	req := graphql.NewRequest(`
    query {
        users{
			edges{
			  node{
				id
				email
				checkin_date
				checkout_date
			  }
			}
		}
	}`)

	req.Header.Set("Cache-Control", "no-cache")
	ctx := context.Background()

	if err := client.Run(ctx, req, &users); err != nil {
		log.Fatal(err)
	}
	return
}

func getEventsToday(client *graphql.Client, today string) (events EventsResponse) {
	req := graphql.NewRequest(`
    query ($filter: String) {
        events (filter:$filter) {
    		user
    		type
		}
	}`)
	req.Var("filter", fmt.Sprintf("{\"date_from\": \"%s 00:00:00\", \"date_to\": \"%s 23:59:59\"}", today, today))
	req.Header.Set("Cache-Control", "no-cache")
	ctx := context.Background()

	if err := client.Run(ctx, req, &events); err != nil {
		log.Fatal(err)
	}
	return
}

type EventResponse struct {
	CreateEvent struct {
		ID string
	}
}

func hasInEvent(id string, events_today EventsResponse) (res bool) {
	res = false
	for _, event := range events_today.Events {
		if event.User == id && event.Type == "IN" {
			res = true
			return
		}
	}
	return
}

func putAbsences(client *graphql.Client, users UsersResponse, events_today EventsResponse) (new_events []EventResponse) {
	for _, user := range users.Users.Edges {
		if !hasInEvent(user.Node.ID, events_today) {
			req := graphql.NewRequest(`
    			mutation ($new_event: NewEvent!) {
    			    createEvent (input:$new_event) {
						id
					}
				}`)
			req.Var("new_event", map[string]string{
				"user_email": user.Node.Email,
				"type":       "ABSENCE",
			})
			req.Header.Set("Cache-Control", "no-cache")
			ctx := context.Background()

			var new_event EventResponse
			if err := client.Run(ctx, req, &new_event); err != nil {
				log.Fatal(err)
			}
			new_events = append(new_events, new_event)
		}
	}
	return
}

func main() {
	client := graphql.NewClient("http://localhost:3000/gql") // TODO: put in env
	users := getAllUsers(client)
	today := time.Now()
	events_today := getEventsToday(client, today.Format("2006-01-02"))
	new_events := putAbsences(client, users, events_today)
	fmt.Printf("%d users with absences at %s\n", len(new_events), today.Format("2006-01-02"))
}
