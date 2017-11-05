package studentlund

import (
	"time"

	"net/http"
	"io/ioutil"

	"github.com/laurent22/ical-go"
)

const (
	STUDENTLUND_DAILY	= "https://www.studentlund.se/event/idag/?ical=1&tribe_display=day&tribe-bar-date="
	STUDENTLUND_WEEKLY	= "https://www.studentlund.se/event/vecka/?ical=1&tribe_display=week&tribe-bar-date="
	STUDENTLUND_MONTHLY	= "https://www.studentlund.se/event/manad/?ical=1&tribe_display=month&tribe-bar-date="
)


func fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func translate(icalData string) ([]event, error) {
	calendar, err := ical.ParseCalendar(icalData)
	if err != nil {
		return nil, err
	}

	var events []event
	for _, node := range calendar.Children {
		if node.Type == 1 {
			_event := createEvent(node)
			if _event.Nation != "" {
				events = append(events, _event)
			}
		}
	}

	return events, nil
}

func convert(url string) ([]event, error) {
	icalData, err := fetch(url)
	if err != nil {
		return nil, err
	}

	events, err := translate(icalData)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func formatDate(date time.Time) string {
	return date.Format("2006-01-02")
}

func appendDate(url string, date time.Time) string {
	return url + formatDate(date)
}

func GetCurrentDay() ([]event, error) {
	return GetDay(time.Now())
}

func GetDay(date time.Time) ([]event, error) {
	url := appendDate(STUDENTLUND_DAILY, date)

	return convert(url)
}

func GetCurrentWeek() ([]event, error) {
	return GetWeek(time.Now())
}

func GetWeek(date time.Time) ([]event, error) {
	url := appendDate(STUDENTLUND_WEEKLY, date)

	return convert(url)
}

func GetCurrentMonth() ([]event, error) {
	return GetMonth(time.Now())
}

func GetMonth(date time.Time) ([]event, error) {
	url := appendDate(STUDENTLUND_MONTHLY, date)

	return convert(url)
}
