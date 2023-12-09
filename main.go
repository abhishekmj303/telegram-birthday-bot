package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/abhishekmj303/telegram-birthday-bot/commands"
	"github.com/abhishekmj303/telegram-birthday-bot/utils"
)

func main() {
	envs, err := godotenv.Read(".env")
	if err != nil {
		panic(err)
	}

	utils.DBConn, err = sql.Open("mysql", envs["DB_DSN"])
	if err != nil {
		panic(err)
	}
	defer utils.DBConn.Close()

	err = utils.DBConn.Ping()
	if err != nil {
		panic(err)
	}
	utils.DBConn.SetConnMaxLifetime(time.Minute * 3)
	utils.DBConn.SetMaxOpenConns(10)
	utils.DBConn.SetMaxIdleConns(10)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		// bot.WithDebug(),
	}

	b, err := bot.New(envs["BOT_TOKEN"], opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/add", bot.MatchTypePrefix, commands.AddHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/showall", bot.MatchTypeExact, commands.ShowallHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/search", bot.MatchTypeExact, commands.SearchHandler)
	// b.RegisterHandler(bot.HandlerTypeMessageText, "/edit", bot.MatchTypeExact, commands.SearchHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/remove", bot.MatchTypeExact, commands.SearchHandler)

	utils.StartNotifier(ctx, b)
	b.Start(ctx)
}

func sendHelpCommand(ctx context.Context, b *bot.Bot, chatID int64) {
	msg := "Hi, I am a birthday reminder bot. I can help you remember birthdays of your friends and family.\n\n"
	msg += "You can use the following commands to interact with me:\n"
	msg += "  • /add <name> - Add a new birthday\n"
	msg += "  • /showall - Show all birthdays\n"
	msg += "  • /search - Search for birthdays\n"
	msg += "  • /remove - Remove a birthday\n"
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	})
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		fmt.Println("update.Message is nil")
		return
	}

	if update.Message.ReplyToMessage == nil {
		sendHelpCommand(ctx, b, update.Message.Chat.ID)
		return
	}

	cmd := strings.Split(update.Message.ReplyToMessage.Text, "/")
	if len(cmd) != 2 {
		sendHelpCommand(ctx, b, update.Message.Chat.ID)
		return
	}
	if !slices.Contains([]string{"search", "edit", "remove"}, cmd[1]) {
		sendHelpCommand(ctx, b, update.Message.Chat.ID)
		return
	}

	update.Message.ReplyToMessage.Text = "/" + cmd[1]
	bd := utils.BirthdayInfo{
		ChatID: update.Message.Chat.ID,
		Name:   update.Message.Text,
	}
	commands.Search(ctx, b, update.Message, &bd, "name")
}
