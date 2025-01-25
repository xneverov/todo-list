package tasks

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(nowDate, date, repeat string) (string, error) {
	errIncorrectRepeat := fmt.Errorf("неверный формат правил повторения (repeat): %s", repeat)

	if repeat == "" {
		return "", fmt.Errorf("нет правил повторения")
	}

	now, err := time.Parse(dateFormat, nowDate)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты (now): %s", nowDate)
	}

	nextDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты (date): %s", date)
	}

	repeatParts := strings.Split(repeat, " ")

	switch repeatParts[0] {
	case "d":
		if len(repeatParts) != 2 {
			return "", errIncorrectRepeat
		}

		nextDate, err = calcDays(now, nextDate, repeatParts[1])
		if err != nil {
			return "", errIncorrectRepeat
		}

	case "w":
		if len(repeatParts) != 2 {
			return "", errIncorrectRepeat
		}

		nextDate, err = calcDaysOfWeek(now, nextDate, repeatParts[1])
		if err != nil {
			return "", errIncorrectRepeat
		}

	case "m":
		if len(repeatParts) < 2 || len(repeatParts) > 3 {
			return "", errIncorrectRepeat
		}

		nextDate, err = calcDaysOfMonth(now, nextDate, strings.Join(repeatParts[1:], " "))
		if err != nil {
			return "", errIncorrectRepeat
		}

	case "y":
		if len(repeatParts) != 1 {
			return "", errIncorrectRepeat
		}

		nextDate = nextDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}

	default:
		return "", errIncorrectRepeat
	}

	return nextDate.Format(dateFormat), nil
}

func calcDays(now, date time.Time, repeat string) (time.Time, error) {
	days, err := strconv.Atoi(repeat)
	if err != nil {
		return time.Time{}, err
	}

	if days < 1 || days > 400 {
		return time.Time{}, errors.New("the number exceeds the limits")
	}

	date = date.AddDate(0, 0, days)
	for date.Before(now) {
		date = date.AddDate(0, 0, days)
	}

	return date, nil
}

func calcDaysOfWeek(now, date time.Time, repeat string) (time.Time, error) {
	availableDays, err := stringToIntMap(repeat, 7, 1)
	if err != nil {
		return time.Time{}, err
	}

	if _, ok := availableDays[7]; ok {
		availableDays[0] = struct{}{}
		delete(availableDays, 7)
	}

	for {
		day := int(date.Weekday())
		if _, ok := availableDays[day]; ok && now.Before(date) {
			break
		}

		date = date.AddDate(0, 0, 1)
	}

	return date, nil
}

func calcDaysOfMonth(now, date time.Time, repeat string) (time.Time, error) {
	repeatParts := strings.Split(repeat, " ")

	availableDays, err := stringToIntMap(repeatParts[0], 31, -2)
	if err != nil {
		return time.Time{}, err
	}

	availableMonths := make(map[int]struct{})
	if len(repeatParts) != 2 {
		for i := 1; i <= 12; i++ {
			availableMonths[i] = struct{}{}
		}
	} else {
		availableMonths, err = stringToIntMap(repeatParts[1], 12, 1)
		if err != nil {
			return time.Time{}, err
		}

		if err := validateDaysInMonths(availableDays, availableMonths); err != nil {
			return time.Time{}, err
		}
	}

	for {
		var lastDay int
		var penultimateDay int

		if _, ok := availableDays[-1]; ok {
			lastDay = daysInMonth(date)
		}

		if _, ok := availableDays[-2]; ok {
			penultimateDay = daysInMonth(date) - 1
		}

		day := int(date.Day())
		month := int(date.Month())

		if _, ok := availableDays[day]; ok || day == lastDay || day == penultimateDay {
			if _, ok := availableMonths[month]; ok {
				if now.Before(date) {
					break
				}
			}
		}

		date = date.AddDate(0, 0, 1)
	}

	return date, nil
}

func stringToIntMap(repeat string, max, min int) (map[int]struct{}, error) {
	m := make(map[int]struct{})
	arr := strings.Split(repeat, ",")

	for _, i := range arr {
		n, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}

		if n > max || n < min {
			return nil, errors.New("the number exceeds the limits")
		}

		m[n] = struct{}{}
	}

	return m, nil
}

func daysInMonth(date time.Time) int {
	lastDay := date.AddDate(0, 1, -date.Day())
	return lastDay.Day()
}

func validateDaysInMonths(days, months map[int]struct{}) error {
	var maxDay int

	for month := range months {
		date := time.Date(time.Now().Year(), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		lastDay := daysInMonth(date)

		if month == 2 {
			lastDay = 29
		}

		if lastDay > maxDay {
			maxDay = lastDay
		}
	}

	for day := range days {
		if day > maxDay {
			return errors.New("day exceeds the maximum day of the months")
		}
	}

	return nil
}
