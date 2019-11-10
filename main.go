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

func main() {
	client := graphql.NewClient("http://localhost:3000/gql") // TODO: put in env
	users := getAllUsers(client)
	fmt.Println(users)

}
