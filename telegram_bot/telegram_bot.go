package telegram_bot

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// TelegramBot ..
type TelegramBot struct {
	*tgbotapi.BotAPI
	commands map[string]struct {
		i     int
		cmd   string
		intro string
		fn    func(*tgbotapi.Message) string
	}
	keyboard           tgbotapi.ReplyKeyboardMarkup
	inlineLinkKeyboard tgbotapi.InlineKeyboardMarkup
}

// NewTelegramBot ..
func NewTelegramBot(token, proxy string, debug bool) (*TelegramBot, error) {
	client := &http.Client{}
	if proxy != "" {
		url, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(url),
		}
	}

	bot, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		return nil, err
	}
	bot.Debug = debug

	imp := &TelegramBot{
		bot,
		make(map[string]struct {
			i     int
			cmd   string
			intro string
			fn    func(*tgbotapi.Message) string
		}),
		tgbotapi.ReplyKeyboardMarkup{},
		tgbotapi.InlineKeyboardMarkup{},
	}

	imp.RegisterCommand("links", "快捷跳转链接", func(msg *tgbotapi.Message) string {
		return fmt.Sprintf("%v", "快捷跳转链接")
	})

	imp.RegisterCommand("current_chat_id", "返回当前会话ID", func(msg *tgbotapi.Message) string {
		return fmt.Sprintf("%v", msg.Chat.ID)
	})
	return imp, nil
}

// RegisterCommand ..
func (bot *TelegramBot) RegisterCommand(command, intro string, fn func(*tgbotapi.Message) string) {

	bot.commands[command] = struct {
		i     int
		cmd   string
		intro string
		fn    func(*tgbotapi.Message) string
	}{i: len(bot.commands) + 1, cmd: command, intro: intro, fn: fn}
}

// RegisterKeyboard ..
func (bot *TelegramBot) RegisterKeyboard(lines [][]string) {
	var rows [][]tgbotapi.KeyboardButton

	for _, line := range lines {
		var buttons []tgbotapi.KeyboardButton
		for i := 0; i < len(line); i++ {
			buttons = append(buttons, tgbotapi.NewKeyboardButton(line[i]))
		}
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(buttons...))
	}
	rows = append(rows, []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("keyboard-close")})

	bot.keyboard = tgbotapi.NewReplyKeyboard(rows...)
}

// RegisterInlineLinkKeyboard ..
func (bot *TelegramBot) RegisterInlineLinkKeyboard(lines [][]struct{ Text, Link string }) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, line := range lines {
		var row []tgbotapi.InlineKeyboardButton
		for i := 0; i < len(line); i++ {
			row = append(row, tgbotapi.NewInlineKeyboardButtonURL(line[i].Text, line[i].Link))
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
	}

	bot.inlineLinkKeyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// Broadcast ..
func (bot *TelegramBot) Broadcast(cid int64, context string) {
	msg := tgbotapi.NewMessage(cid, context)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

var mode = "f"
var num = 0

// Run ..
func (bot *TelegramBot) Run() error {
	cfg := tgbotapi.NewUpdate(0)
	cfg.Timeout = 60

	updates, err := bot.GetUpdatesChan(cfg)
	if err != nil {
		log.Println(err)
	}

	for update := range updates {
		if ec := recover(); ec != nil {
			log.Fatalf("[TelegramBot]RUN EC:%v", ec)
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "set_mod":
				args := update.Message.CommandArguments()

				if args == "f" {
					mode = "f"
					bot.Broadcast(update.Message.Chat.ID, "设置解析F模式")
				}
				if args == "p" {
					mode = "p"
					bot.Broadcast(update.Message.Chat.ID, "设置解析P模式")
				}

			case "set_num":
				args := update.Message.CommandArguments()
				if mode == "p" {
					num, _ = strconv.Atoi(args)
					bot.Broadcast(update.Message.Chat.ID, fmt.Sprintf("设置NUM为%v", num))
				}
			}

		} else {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			erri := strings.Index(msg.Text, "err")
			if erri != -1 {
				if mode == "p" {
					if strings.Index(msg.Text, "[") != -1 && strings.Index(msg.Text, "]") != -1 {
						s := msg.Text[erri:]
						s = strings.Split(s, "[")[1]
						s = strings.Split(s, "]")[0]
						s = strings.TrimSpace(s)
						r := AnalyzePoker(s, num)
						bot.Send(tgbotapi.NewMessage(msg.ChatID, r))
					}
					continue
				}

				if strings.Index(msg.Text, "[") != -1 && strings.Index(msg.Text, "]") != -1 {
					s := msg.Text[erri:]
					s = strings.Split(s, "[")[1]
					s = strings.Split(s, "]")[0]
					s = strings.TrimSpace(s)
					r := AnalyzeFish(s)
					bot.Send(tgbotapi.NewMessage(msg.ChatID, r))
				}
			}
		}
	}

	return nil
}
