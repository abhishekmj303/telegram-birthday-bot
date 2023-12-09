package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-telegram/bot"
)

var (
	tz *time.Location
)

func StartNotifier(ctx context.Context, b *bot.Bot) {
	tz, _ = time.LoadLocation("Asia/Kolkata")
	s := gocron.NewScheduler(tz)
	s.Every(1).Day().At("00:00").Do(notifyAllBefore, ctx, b, 0)
	s.Every(1).Day().At("20:00").Do(notifyAllBefore, ctx, b, 1)
	s.Every(1).Day().At("10:00").Do(notifyAllBefore, ctx, b, 7)
}

func notifyBefore(ctx context.Context, b *bot.Bot, bd BirthdayInfo, before int) {
	var text string

	switch before {
	case 0:
		text = "Today"
	case 1:
		text = "Tomorrow"
	default:
		text = fmt.Sprintf("%d days from now", before)
	}

	text += fmt.Sprintf(" (%s) is %s's birthday 🎂", bd.Date(), bd.Name)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: bd.ChatID,
		Text:   text,
	})
}

func notifyAllBefore(ctx context.Context, b *bot.Bot, before int) {
	bds, err := GetBirthdays()
	if err != nil {
		fmt.Println(err)
		return
	}

	now := time.Now().In(tz)
	day := now.Day() + before
	month := int(now.Month())

	for _, bd := range bds {
		if bd.Day == day && bd.Month == month {
			notifyBefore(ctx, b, bd, before)
		}
	}
}
