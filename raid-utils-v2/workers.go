/*

workers.go -
file that contains goroutine workers

*/

package main

import (
	// internals
	"fmt"
	// externals
	"github.com/bwmarrin/discordgo"
)

// the webhook worker
func webhookWorker(dg *discordgo.Session, channels []*discordgo.Channel, name, message, avatarUrl string) {

	// select a random channel to spam in
	channel := channels[randint(0, len(channels))]

	// create the initial webhook
	webhook, err := dg.WebhookCreate(channel.ID, name, "")
	if err != nil {

		// we were unable to create the webhook
		fmt.Printf("[err]: unable to create webhook...\n")
		fmt.Printf("       %v\n", err)
		return

	}

	// start the spam loop
	for {

		// send the message
		err = dg.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
			Content:   message,
			AvatarURL: avatarUrl,
		})
		if err != nil {

			// check if the webhook was deleted
			if err.(*discordgo.RESTError).Response.StatusCode == 404 {

				// if it was, recreate it
				webhook, err = dg.WebhookCreate(channel.ID, name, "")
				if err != nil {

					// we were unable to create the webhook
					fmt.Printf("[err]: unable to recreate webhook...\n")
					fmt.Printf("       %v\n", err)
					return

				}

			} else {

				// unknown error, so we just exit
				fmt.Printf("[err]: unknown error from api...\n")
				fmt.Printf("       %v\n", err)
				return

			}

		}

	}

	// delete the webhook
	err = dg.WebhookDelete(webhook.ID)

}
