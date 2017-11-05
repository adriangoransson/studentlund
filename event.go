package studentlund

import (
	"regexp"
	"strings"
	"time"

	"github.com/laurent22/ical-go"
)

type event struct {
	Id			string		`json:"id"`
	Summary 	string		`json:"summary"`
	Description	string		`json:"description"`
	Url			string		`json:"url"`
	Location	string		`json:"location"`
	Nation		string		`json:"nation"`
	DateStart 	time.Time	`json:"date_start"`
	DateEnd 	time.Time	`json:"date_end"`
	LastUpdated	time.Time	`json:"last_updated"`
}

func resolveNation(summary string) string {
	// Because these just can't do it like the others
	if strings.Contains(summary, "Blekingska") {
		return "Blekingska Nationen"
	} else if strings.Contains(summary, "VG") {
		return "Västgöta Nation"
	}

	re := regexp.MustCompile(`([\wÅÄÖåäö]+) ([Nn]ation(?:en)?)`)
	matches := re.FindStringSubmatch(summary)
	// Find the name of the nation, the inflection of the word nation and capitalize both words
	if len(matches) > 1 {
		return (strings.Title(matches[1]) + " " + strings.Title(matches[2]))
	}

	return ""
}

func createEvent(node *ical.Node) event {
	summary := node.PropString("SUMMARY", "")
	description := node.PropString("DESCRIPTION", "")

	nation := resolveNation(summary)
	// If not found in summary, try the description
	if nation == "" {
		nation = resolveNation(description)
	}

	// Remove weird backslashes from the location field
	address := strings.Replace(node.PropString("LOCATION", ""), "\\", "", -1)

	return event{
		Id:				node.PropString("UID", ""),
		Summary: 		summary,
		Description: 	description,
		Url:			node.PropString("URL", ""),
		Nation:			nation,
		Location:		address,
		DateStart:		node.PropDate("DTSTART", time.Now()),
		DateEnd:		node.PropDate("DTEND", time.Now()),
		LastUpdated:	node.PropDate("LAST-MODIFIED", time.Now()),
	}
}
