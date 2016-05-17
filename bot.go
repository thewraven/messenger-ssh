package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/frodsan/fbot"
)

func main() {
	bot := fbot.NewBot(fbot.Config{
		AccessToken: os.Getenv("EAALVTDYz2rcBACQkrtrSon8eNXZCtJl40HA1bSwfKpOMV977JFsB7EXuD3hrYxnugW9s837KSxe5w6JscQSwMsbjPZA62ZBIqyyAl2ArsmDlJtxnfEGREvffwJBsefbxs24jyPJTyhdYLvmC65CU7YT8jtyoJMeEL4lqmlGXQZDZD"),
		AppSecret:   os.Getenv("2032fab87b7b8197b430cf5a65c9b639"),
		VerifyToken: os.Getenv("0b8baf8b94e97c2416584afc2a8e9016"),
	})

	bot.On(fbot.EventMessage, func(event *fbot.Event) {
		fmt.Println(event.Message.Text)

		bot.Deliver(fbot.DeliverParams{
			Recipient: event.Sender,
			Message: &fbot.Message{
				Text: event.Message.Text,
			},
		})
	})

	http.Handle("/bot", fbot.Handler(bot))

	http.ListenAndServe(":4567", nil)
}