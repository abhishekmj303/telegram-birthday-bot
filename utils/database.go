package utils

import (
	"database/sql"
	"fmt"
	"strings"
)

var (
	DBConn *sql.DB
)

const (
	insertSQL         = "INSERT INTO birthdays (chat_id, name, day, month) VALUES (?, ?, ?, ?)"
	selectSQL         = "SELECT chat_id, name, day, month FROM birthdays"
	selectAllSQL      = "SELECT name, day, month FROM birthdays WHERE chat_id = ?"
	selectNameSQL     = "SELECT name, day, month FROM birthdays WHERE chat_id = ? AND name = ?"
	selectNameLikeSQL = "SELECT name, day, month FROM birthdays WHERE chat_id = ? AND LOWER(name) LIKE ?"
	selectMonthSQL    = "SELECT name, day, month FROM birthdays WHERE chat_id = ? AND month = ?"
	selectDaySQL      = "SELECT name, day, month FROM birthdays WHERE chat_id = ? AND day = ?"
	selectDateSQL     = "SELECT name, day, month FROM birthdays WHERE chat_id = ? AND day = ? AND month = ?"
	updateDateSQL     = "UPDATE birthdays SET day = ?, month = ? WHERE chat_id = ? AND name = ?"
	updateNameSQL     = "UPDATE birthdays SET name = ? WHERE chat_id = ? AND name = ?"
	deleteSQL         = "DELETE FROM birthdays WHERE chat_id = ? AND name = ?"
	orderbySQL        = " ORDER BY month, day"
)

func GetBirthday(chatID int64, name string) (BirthdayInfo, error) {
	var day, month int

	err := DBConn.QueryRow(selectNameSQL, chatID, name).
		Scan(&name, &day, &month)
	if err != nil {
		if err == sql.ErrNoRows {
			return BirthdayInfo{}, nil
		} else {
			return BirthdayInfo{}, err
		}
	}

	bd := BirthdayInfo{
		Name:  name,
		Day:   day,
		Month: month,
	}
	return bd, nil
}

func GetBirthdays() ([]BirthdayInfo, error) {
	rows, err := DBConn.Query(selectSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var info BirthdayInfo
	var infoList []BirthdayInfo
	for rows.Next() {
		err := rows.Scan(&info.ChatID, &info.Name, &info.Day, &info.Month)
		if err != nil {
			return nil, err
		}
		infoList = append(infoList, info)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return infoList, nil
}

func AddBirthday(bd *BirthdayInfo) error {
	_, err := DBConn.Exec(insertSQL, bd.ChatID, bd.Name, bd.Day, bd.Month)
	return err
}

func SearchBirthday(bd *BirthdayInfo, searchby string) ([]BirthdayInfo, error) {
	var rows *sql.Rows
	var err error

	switch searchby {
	case "all":
		rows, err = DBConn.Query(selectAllSQL+orderbySQL, bd.ChatID)
	case "name":
		likeExp := "%" + strings.ToLower(bd.Name) + "%"
		rows, err = DBConn.Query(selectNameLikeSQL+orderbySQL, bd.ChatID, likeExp)
	case "month":
		rows, err = DBConn.Query(selectMonthSQL+orderbySQL, bd.ChatID, bd.Month)
	case "day":
		rows, err = DBConn.Query(selectDaySQL+orderbySQL, bd.ChatID, bd.Day)
	case "date":
		rows, err = DBConn.Query(selectDateSQL+orderbySQL, bd.ChatID, bd.Day, bd.Month)
	default:
		return nil, fmt.Errorf("invalid searchby: %s", searchby)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	info := BirthdayInfo{ChatID: bd.ChatID}
	var infoList []BirthdayInfo
	for rows.Next() {
		err := rows.Scan(&info.Name, &info.Day, &info.Month)
		if err != nil {
			return nil, err
		}
		infoList = append(infoList, info)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return infoList, nil
}

func UpdateDate(bd *BirthdayInfo) error {
	_, err := DBConn.Exec(updateDateSQL, bd.Day, bd.Month, bd.ChatID, bd.Name)
	return err
}

func UpdateName(bd *BirthdayInfo, newName string) error {
	_, err := DBConn.Exec(updateNameSQL, newName, bd.ChatID, bd.Name)
	return err
}

func RemoveBirthday(bd *BirthdayInfo) error {
	_, err := DBConn.Exec(deleteSQL, bd.ChatID, bd.Name)
	return err
}
