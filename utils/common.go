package utils

import (
	"fmt"
	"strings"
	"time"
)

type BirthdayInfo struct {
	ChatID int64
	Name   string
	Day    int
	Month  int
}

func (bd BirthdayInfo) Date() string {
	return time.Date(0, time.Month(bd.Month), bd.Day, 0, 0, 0, 0, time.UTC).Format("02-Jan")
}

func (bd BirthdayInfo) String() string {
	dateStr := bd.Date()
	return fmt.Sprintf("%s on %s", bd.Name, dateStr)
}

func BirthdayListStr(bdList []BirthdayInfo) string {
	var bdStrList []string
	for _, bd := range bdList {
		bdStrList = append(bdStrList, "  • "+bd.String())
	}
	return strings.Join(bdStrList, "\n")
}

const (
	InvalidAddReply = "Please send the name of the person along with the command\neg. /add John Doe"
)

func GetMsgData(msg string) (string, bool) {
	_, data, isfound := strings.Cut(msg, " ")
	return data, isfound
}

func RetryReply(command string) string {
	return fmt.Sprintf("Something went wrong!!\nPlease retry %s again.", command)
}
