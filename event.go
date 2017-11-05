package studentlund

import (
	"regexp"
	"strings"
	"time"
	"net/url"

	"github.com/laurent22/ical-go"
)

type organizer struct {
	Name	string	`json:"name"`
	Email	string	`json:"email"`
}

type dates struct {
	Start		time.Time	`json:"start"`
	End			time.Time	`json:"end"`
	LastUpdated	time.Time	`json:"last_updated"`
}

type event struct {
	Id			string		`json:"id"`
	Summary 	string		`json:"summary"`
	Description	string		`json:"description"`
	Url			string		`json:"url"`
	ImageUrl	string		`json:"image_url"`
	Location	string		`json:"location"`
	Date 		dates		`json:"date"`
	Organizer	organizer	`json:"organizer"`
}

func stripMailto(email string) string {
	if strings.HasPrefix(email, "MAILTO:") {
		return email[7:]
	}

	return email
}

func resolveOrganizer(node *ical.Node) (organizer, error) {
	// Try the ORGANIZER field, CN parameter first
	organizerField := node.ChildByName("ORGANIZER")

	if organizerField != nil {
		// URL decode the CN field
		organizerName, err := url.PathUnescape(organizerField.Parameter("CN", ""))
		if err != nil {
			return organizer{}, err
		}

		// Remove quotes surrounding the organizer name
		organizerName = strings.Trim(organizerName, "\"")

		// strip MAILTO: before email if it exists
		return organizer{
			Name:	organizerName,
			Email:	stripMailto(organizerField.Value),
		}, nil
	}

	// ORGANIZER field didn't exist in calendar node
	// Try parsing the organizer name from the SUMMARY field instead
	organizerName := resolveNationByText(node.PropString("SUMMARY", ""))
	if organizerName != "" {
		return organizer{
			Name:	organizerName,
			Email:	"",
		}, nil
	}

	// No match in SUMMARY. Try DESCRIPTION instead
	return organizer{
		Name:	resolveNationByText(node.PropString("DESCRIPTION", "")),
		Email:	"",
	}, nil
}

func resolveNationByText(text string) string {
	// Because these just can't do it like the others
	if strings.Contains(text, "Blekingska") {
		return "Blekingska Nationen"
	} else if strings.Contains(text, "VG") {
		return "Västgöta Nation"
	}

	re := regexp.MustCompile(`([\wÅÄÖåäö]+) ([Nn]ation(?:en)?)`)
	matches := re.FindStringSubmatch(text)
	// Find the name of the nation, the inflection of the word nation and capitalize both words
	if len(matches) > 1 {
		return (strings.Title(matches[1]) + " " + strings.Title(matches[2]))
	}

	return ""
}

func createEvent(node *ical.Node) (event, error) {
	organizerData, err := resolveOrganizer(node)

	if err != nil {
		return event{}, err
	}

	date := dates{
		Start:			node.PropDate("DTSTART", time.Now()),
		End:			node.PropDate("DTEND", time.Now()),
		LastUpdated:	node.PropDate("LAST-MODIFIED", time.Now()),
	}

	// Remove weird backslashes from the location field
	address := strings.Replace(node.PropString("LOCATION", ""), "\\", "", -1)

	return event{
		Id:				node.PropString("UID", ""),
		Summary: 		node.PropString("SUMMARY", ""),
		Description: 	node.PropString("DESCRIPTION", ""),
		Url:			node.PropString("URL", ""),
		ImageUrl:		node.PropString("ATTACH", ""),
		Organizer:		organizerData,
		Location:		address,
		Date:			date,
	}, nil
}
