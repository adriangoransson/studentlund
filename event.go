package studentlund

import (
	"regexp"
	"strings"
	"time"

	"net/url"

	ical "github.com/adriangoransson/ical-go"
)

type Organizer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Dates struct {
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	LastUpdated time.Time `json:"last_updated"`
}

type Event struct {
	Id          string    `json:"id"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	ImageUrl    string    `json:"image_url"`
	Location    string    `json:"location"`
	Date        Dates     `json:"date"`
	Organizer   Organizer `json:"organizer"`
}

type ByDate []Event

func (e ByDate) Len() int {
	return len(e)
}

func (e ByDate) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e ByDate) Less(i, j int) bool {
	// If the start dates are equal, shortest event goes first
	if e[i].Date.Start.Equal(e[j].Date.Start) {
		return e[i].Date.End.Before(e[j].Date.End)
	}

	return e[i].Date.Start.Before(e[j].Date.Start)
}

func stripMailto(email string) string {
	if strings.HasPrefix(email, "MAILTO:") {
		return email[7:]
	}

	return email
}

func resolveOrganizer(node *ical.Node) (Organizer, error) {
	// Try the ORGANIZER field, CN parameter first
	organizer := node.ChildByName("ORGANIZER")

	if organizer != nil {
		// URL decode the CN field
		organizerName, err := url.PathUnescape(organizer.Parameter("CN", ""))
		organizerName = strings.Title(organizerName)

		if err != nil {
			return Organizer{}, err
		}

		// Remove quotes surrounding the organizer name
		organizerName = strings.Trim(organizerName, "\"")

		// strip MAILTO: before email if it exists
		return Organizer{
			Name:  organizerName,
			Email: stripMailto(organizer.Value),
		}, nil
	}

	// ORGANIZER field didn't exist in calendar node
	// Try parsing the organizer name from the SUMMARY field instead
	organizerName := resolveNationByText(node.PropString("SUMMARY", ""))

    if organizerName == "" {
        organizerName = resolveNationByText(node.PropString("DESCRIPTION", ""))
    }

    if organizerName == "" {
        organizerName = resolveNationByText(node.PropString("LOCATION", ""))
    }

    return Organizer{
        Name:  organizerName,
        Email: "",
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
		return strings.Title(matches[1] + " " + matches[2])
	}

	return ""
}

func createEvent(node *ical.Node) (Event, error) {
	organizer, err := resolveOrganizer(node)

	if err != nil {
		return Event{}, err
	}

	date := Dates{
		Start:       node.PropDate("DTSTART", time.Now()),
		End:         node.PropDate("DTEND", time.Now()),
		LastUpdated: node.PropDate("LAST-MODIFIED", time.Now()),
	}

	// Remove weird backslashes from the location field
	address := strings.Replace(node.PropString("LOCATION", ""), "\\", "", -1)

	return Event{
		Id:          node.PropString("UID", ""),
		Summary:     node.PropString("SUMMARY", ""),
		Description: strings.TrimSpace(node.PropString("DESCRIPTION", "")),
		Url:         node.PropString("URL", ""),
		ImageUrl:    node.PropString("ATTACH", ""),
		Organizer:   organizer,
		Location:    address,
		Date:        date,
	}, nil
}
