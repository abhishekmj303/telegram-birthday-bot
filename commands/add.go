package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abhishekmj303/telegram-birthday-bot/ui"
	"github.com/abhishekmj303/telegram-birthday-bot/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var (
	addInfo = make(map[int64]utils.BirthdayInfo)
)

func AddHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	msgID := update.Message.ID
	var text string
	var kb models.ReplyMarkup

	name, isfound := utils.GetMsgData(update.Message.Text)
	if !isfound {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           chatID,
			Text:             utils.InvalidAddReply,
			ReplyToMessageID: msgID,
		})
		return
	}

	bd, err := utils.GetBirthday(chatID, name)
	if err != nil {
		text = utils.RetryReply("/add")
	} else {
		if bd.Name == "" {
			text = fmt.Sprintf("Select the Birthday of '%s'", name)
			addInfo[chatID] = utils.BirthdayInfo{
				ChatID: chatID,
				Name:   name,
			}
			kb = ui.Datepicker(b, datepickerAddHandler)
		} else {
			text = fmt.Sprintf("Name '%s' already exists.\nDo you want to /remove it?", name)
			kb = ui.Birthdaypicker(b, birthdaypickerRemoveHandler, []utils.BirthdayInfo{bd})
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           chatID,
		Text:             text,
		ReplyMarkup:      kb,
		ReplyToMessageID: msgID,
	})
}

func datepickerAddHandler(ctx context.Context, b *bot.Bot, mes *models.Message, date time.Time) {
	chatID := mes.Chat.ID
	var text string

	info, ok := addInfo[chatID]
	if !ok {
		fmt.Println("info not found")
		text = utils.RetryReply("/add")
	} else {
		info.Day = date.Day()
		info.Month = int(date.Month())
		text = "<b>Added Birthday</b> for " + info.String()
	}

	err := utils.AddBirthday(&info)
	if err != nil {
		log.Println(err)
		text = utils.RetryReply("/add")
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}
