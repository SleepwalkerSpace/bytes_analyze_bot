package main

import "bytes_analyze_bot/telegram_bot"

func main() {
	token := "6565489549:AAF0TZD_-_f_kPN6RKX2rIFvW869DCy1aKs"
	bot, err := telegram_bot.NewTelegramBot(token, "", true)
	if err != nil {
		panic(err)
	}

	bot.Run()

}
