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
				CheckinDate  time.Time
				CheckoutDate time.Time
			}
		}
	}
}

type EventsResponse struct {
	Events []struct {
		Date string
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
			date
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

func getAbsences() {

	return
}

func getDelays() {

	return
}

func main() {
	client := graphql.NewClient("http://localhost:3000/gql") // TODO: put in env
	// users := getAllUsers(client)
	// fmt.Println(users)

	today := time.Now()
	events_today := getEventsToday(client, today.Format("2006-01-02"))
	fmt.Println(events_today)
}
