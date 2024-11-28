package client

import (
	"database/sql"
	"dbuggen/server/database"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

// coverpage generates an HTML template for the cover image.
// If the coverpage is valid, it returns an HTML string with an image tag.
// If the coverpage is not valid, it returns an empty HTML string.
func coverpage(coverpage sql.NullString) template.HTML {
	if coverpage.Valid {
		return template.HTML(fmt.Sprintf(`<img src="%v" style="max-width: 40vw;">`, template.HTMLEscapeString(coverpage.String)))
	}
	return template.HTML("")
}

// Generates an html template for a member picture. If the member has
// an image that will be displayed, and otherwise it will show a
// default picture.
func memberpicture(picture sql.NullString) template.HTML {
	var pic string
	if picture.Valid {
		pic = picture.String
	} else {
		pic = "/public/default_member.svg"
	}
	return template.HTML(fmt.Sprintf(`<img src="%v" style="max-width: 10vw;">`, template.HTMLEscapeString(pic)))
}

// a struct for displaying members on the
// redaqtion page.
type displayMember struct {
	KthID   string
	Name    string
	Picture template.HTML
	Title   string
}

// creates a displaymember from a member struct, using the prefered
// name if there is any and a html template for the picture used.
func displaymemberize(members []database.Member) []displayMember {
	displaymembers := make([]displayMember, len(members))
	for i, member := range members {
		name := authorsName(database.Author{KthID: member.KthID, PreferedName: member.PreferedName})
		displaymembers[i] = displayMember{
			KthID:   fmt.Sprintf("redaqtionen/%v", member.KthID),
			Name:    name,
			Picture: memberpicture(member.PictureURL),
			Title:   member.Title,
		}
	}
	return displaymembers
}

// Gets a list of current chefreds kth ids from dfunkt
func getChefreds(DFUNKT_URL string) []string {
	type result struct { // "json"... more like "no, son"
		Mandates []struct { // "go"... more like "row".
			User struct { // the boat - pshshshchhhhh
				Kthid string `json:"kthid"`
			} `json:"user"`
		} `json:"mandates"`
	}

	if DFUNKT_URL[len(DFUNKT_URL)-1] != '/' { // no https://dfunkt.seapi/role/...
		DFUNKT_URL = fmt.Sprintf("%v/", DFUNKT_URL)
	}

	var chefreds []string

	resp, err := http.Get(fmt.Sprintf("%vapi/role/chefred/current", DFUNKT_URL))
	if err != nil {
		log.Println(err)
		return chefreds
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("unexpected http response from dfunkt: %v", resp.StatusCode)
		return chefreds
	}

	contents, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return chefreds
	}

	var res result
	err = json.Unmarshal([]byte(contents), &res)
	if err != nil {
		log.Println(err)
		return chefreds
	}

	for _, id := range res.Mandates {
		chefreds = append(chefreds, id.User.Kthid)
	}

	return chefreds
}

// Using a list of the kth ids of however many chefreds there are, it removes them
// the list of active members in case they are also present there as to not
// display duplicates. In case the chefred is not an active member it will
// create a default profile for them.
func removeDuplicateChefreds(chefredIDs []string, members []database.Member) ([]database.Member, []database.Member) {
	chefreds := make([]database.Member, len(chefredIDs))
	for i, chef := range chefredIDs {
		j := slices.IndexFunc(members, func(member database.Member) bool { return member.KthID == chef })
		if j == -1 {
			chefreds[i] = database.Member{
				KthID:        chef,
				PreferedName: sql.NullString{Valid: false, String: ""},
				PictureURL:   sql.NullString{Valid: false, String: ""},
				Title:        "chefred",
				Active:       true,
			}
			continue
		}

		chefreds[i] = members[j]
		members = slices.Delete(members, j, j+1)
	}
	return chefreds, members
}

// authortext returns the author text based on the given AuthorText and authors, wrapped
// in a html template.
// If AuthorText is valid, it returns the AuthorText string. Otherwise, it constructs
// the author text using the names of the authors.
func authortext(AuthorText sql.NullString, authors []database.Author) template.HTML {
	if AuthorText.Valid {
		return template.HTML(AuthorText.String)
	}

	if len(authors) == 0 {
		return "Skriven av redaqtionen"
	}

	var sb strings.Builder
	sb.WriteString("Skriven av ")

	sb.WriteString(string(authorURL(authorsName(authors[0]), authors[0].KthID)))
	if len(authors) == 1 {
		return template.HTML(sb.String())
	}

	for i := 1; i < len(authors)-1; i++ {
		sb.WriteString(fmt.Sprintf(", %v",
			authorURL(authorsName(authors[i]), authors[i].KthID)))
	}

	sb.WriteString(fmt.Sprintf(" och %v",
		authorURL(authorsName(authors[len(authors)-1]), authors[len(authors)-1].KthID)))
	return template.HTML(sb.String())
}

func authorURL(name string, kthID string) template.HTML {
	url := fmt.Sprintf(`<a href="/redaqtionen/%v">%v</a>`, kthID, name)
	return template.HTML(url)
}

// authorsName returns the preferred name of an author from the database.
// If the preferred name is not available, it retrieves the display name
// from hodis.datasektionen.se based on the author's KTH ID.
func authorsName(a database.Author) string {
	if a.PreferedName.Valid {
		return a.PreferedName.String
	}

	resp, err := http.Get(fmt.Sprintf("https://hodis.datasektionen.se/uid/%v", a.KthID))
	if err != nil {
		log.Println(err)
		return a.KthID
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("unexpected http GET status from hodis: %v", resp.StatusCode)
	}

	contents, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return a.KthID
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(contents, &m)
	if err != nil {
		log.Println(err)
		return a.KthID
	}

	displayName, ok := m["displayName"].(string)
	if !ok {
		log.Printf("Failed to convert displayName to string, from kthid %v", a.KthID)
		return a.KthID
	}
	return displayName
}

type DarkmodeStatus struct {
	Darkmode bool
	LastPoll time.Time
	Url      string
	Mutex    sync.RWMutex
}

// darkmode checks if the mörkläggning is active by making request to
// an external API. It parses and outputs the result as a bool.
// If any error occurs during the request or parsing the response, it returns
// the default dark mode status which is true.
func Darkmode(ds *DarkmodeStatus) bool {
	ds.Mutex.RLock()
	if time.Since(ds.LastPoll) <= time.Hour*24 {
		ds.Mutex.RUnlock()
		return ds.Darkmode
	}

	ds.Mutex.RUnlock()
	ds.Mutex.Lock()
	defer ds.Mutex.Unlock()
	defDarkmode := true

	resp, err := http.Get(ds.Url)
	if err != nil {
		log.Println(err)
		return defDarkmode
	}

	contents, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return defDarkmode
	}

	darkmodeStatus, err := strconv.ParseBool(strings.TrimSpace(string(contents)))
	if err != nil {
		log.Println(err)
		return defDarkmode
	}

	ds.LastPoll = time.Now()
	ds.Darkmode = darkmodeStatus
	return darkmodeStatus
}

// Function to remove the "/" before parameters, which was
// a problem. Turns "/123": string, into 123: int.
func pathIntSeparator(paramRaw string) (int, error) {
	_, paramLessRaw, found := strings.Cut(paramRaw, "/")

	if found {
		param, err := strconv.Atoi(paramLessRaw)
		if err != nil {
			return 0, err
		}
		return param, nil
	}

	param, err := strconv.Atoi(paramRaw)
	if err != nil {
		return 0, err
	}
	return param, nil
}
