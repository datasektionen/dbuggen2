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
	"strconv"
	"strings"
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

// authortext returns the author text based on the given AuthorText and authors.
// If AuthorText is valid, it returns the AuthorText string. Otherwise, it constructs
// the author text using the names of the authors.
func authortext(AuthorText sql.NullString, authors []database.Author) string {
	if AuthorText.Valid {
		return AuthorText.String
	}

	if len(authors) == 0 {
		return "Skriven av redaqtionen"
	}

	var sb strings.Builder
	sb.WriteString("Skriven av ")

	sb.WriteString(authorsName(authors[0]))
	if len(authors) == 1 {
		return sb.String()
	}

	for i := 1; i < len(authors)-1; i++ {
		sb.WriteString(fmt.Sprintf(", %v", authorsName(authors[i])))
	}

	sb.WriteString(fmt.Sprintf(" och %v", authorsName(authors[len(authors)-1])))
	return sb.String()
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
	errg := json.Unmarshal(contents, &m)
	if errg != nil {
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
}

// darkmode checks if the mörkläggning is active by making request to
// an external API. It parses and outputs the result as a bool.
// If any error occurs during the request or parsing the response, it returns
// the default dark mode status which is true.
func Darkmode(ds *DarkmodeStatus, url string) bool {
	if time.Now().Sub(ds.LastPoll) <= time.Hour*24 {
		return ds.Darkmode
	}

	defDarkmode := true

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return defDarkmode
	}

	contents, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Println(err)
		return defDarkmode
	}

	darkmodeStatus, err := strconv.ParseBool(string(contents))
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
