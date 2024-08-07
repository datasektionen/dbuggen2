package client

import (
	"database/sql"
	"dbuggen/server/database"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/h2non/gock"
)

func TestCoverpage(t *testing.T) {
	// Test case 1: Valid coverpage
	cp := sql.NullString{String: "https://example.com/cover.jpg", Valid: true}
	expected := `<img src="https://example.com/cover.jpg" style="max-width: 40vw;">`
	got := coverpage(cp)
	if string(got) != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}

	// Test case 2: Invalid coverpage
	cp = sql.NullString{String: "", Valid: false}
	expected = ""
	got = coverpage(cp)
	if string(got) != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}
}

func TestMemberpicture(t *testing.T) {
	// Test case 1: Valid coverpage
	mp := sql.NullString{String: "https://example.com/cover.jpg", Valid: true}
	expected := `<img src="https://example.com/cover.jpg" style="max-width: 10vw;">`
	got := memberpicture(mp)
	if string(got) != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}

	// Test case 2: Invalid coverpage
	mp = sql.NullString{String: "", Valid: false}
	expected = `<img src="/assets/default_member.svg" style="max-width: 10vw;">`
	got = memberpicture(mp)
	if string(got) != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}
}

func TestGetChefreds(t *testing.T) {
	getChefreds("https://dfunkt.datasektionen.se/")
}

func TestAuthortext(t *testing.T) {
	// Test case 1: Valid AuthorText
	authorText := sql.NullString{String: "Skriven av Test Testström", Valid: true}
	authors := []database.Author{{PreferedName: sql.NullString{String: "Ej Korrektström", Valid: true}, KthID: "testsupp"}}
	expected := "Skriven av Test Testström"
	got := authortext(authorText, authors)
	if got != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}

	// Test case 2: No AuthorText and empty list of authors
	authorText = sql.NullString{String: "", Valid: false}
	authors = []database.Author{}
	expected = "Skriven av redaqtionen"
	got = authortext(authorText, authors)
	if got != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}

	// Test case 3: No AuthorText with multiple authors
	authorText = sql.NullString{String: "", Valid: false}
	authors = []database.Author{
		{PreferedName: sql.NullString{String: "Skribent Skrivarsson", Valid: true}, KthID: "skribson"},
		{PreferedName: sql.NullString{String: "Skämt Skojsdotter", Valid: true}, KthID: "sskoj"},
		{PreferedName: sql.NullString{String: "", Valid: false}, KthID: "testsupp"},
	}
	expected = "Skriven av Skribent Skrivarsson, Skämt Skojsdotter och test support"
	got = authortext(authorText, authors)
	if got != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}

	// Test case 4: No AuthorText with single author
	authorText = sql.NullString{String: "", Valid: false}
	authors = []database.Author{{PreferedName: sql.NullString{String: "Testare #1", Valid: true}, KthID: "testsupp"}}
	expected = "Skriven av Testare #1"
	got = authortext(authorText, authors)
	if got != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}
}

func TestAuthorsName(t *testing.T) {
	// Test case 1: Preferred name is available
	author := database.Author{
		PreferedName: sql.NullString{String: "Testaren i dbuggen", Valid: true},
		KthID:        "testsupp",
	}
	expected := "Testaren i dbuggen"
	got := authorsName(author)
	if got != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}

	// Test case 2: Preferred name is not available, retrieve display name
	author = database.Author{
		PreferedName: sql.NullString{String: "", Valid: false},
		KthID:        "testsupp",
	}
	expected = "test support"
	got = authorsName(author)
	if got != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}

	// Test case 3: Failed to retrieve display name, return KTH ID
	author = database.Author{
		PreferedName: sql.NullString{String: "", Valid: false},
		KthID:        "jaghoppasingenpåkthhetersåhär",
	}
	expected = "jaghoppasingenpåkthhetersåhär"
	got = authorsName(author)
	if got != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}
}

func TestDarkmodeFalse(t *testing.T) {
	defer gock.Off()
	darkmodeURL := "http://darkmode.datasektionen.se"
	expected := false
	gock.New(darkmodeURL).
		Get("").
		Reply(http.StatusOK).
		JSON(strconv.FormatBool(expected))

	oldpoll := time.Date(1983, time.October, 7, 17, 0, 0, 0, time.Local)
	ds := DarkmodeStatus{
		Darkmode: true,
		LastPoll: oldpoll,
		Url:      darkmodeURL,
		Mutex:    sync.RWMutex{},
	}

	got := Darkmode(&ds)
	if got != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}
	newpoll := ds.LastPoll
	if newpoll == oldpoll {
		t.Errorf("The polling date of the struct has not been updated")
	}

	got2 := Darkmode(&ds)
	if got2 != expected {
		t.Errorf("got %v, wanted %v", got2, expected)
	}
	if ds.LastPoll != newpoll {
		t.Errorf("The darkmode status was polled again when it should not have been")
	}
}

func TestDarkmodeTrue(t *testing.T) {
	defer gock.Off()
	darkmodeURL := "http://darkmode.datasektionen.se"
	expected := true
	gock.New(darkmodeURL).
		Get("").
		Reply(http.StatusOK).
		JSON(strconv.FormatBool(expected))

	oldpoll := time.Date(1983, time.October, 7, 17, 0, 0, 0, time.Local)
	ds := DarkmodeStatus{
		Darkmode: false,
		LastPoll: oldpoll,
		Url:      darkmodeURL,
		Mutex:    sync.RWMutex{},
	}

	got := Darkmode(&ds)
	if got != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}
	newpoll := ds.LastPoll
	if newpoll == oldpoll {
		t.Errorf("The polling date of the struct has not been updated")
	}

	got2 := Darkmode(&ds)
	if got2 != expected {
		t.Errorf("got %v, wanted %v", got2, expected)
	}
	if ds.LastPoll != newpoll {
		t.Errorf("The darkmode status was polled again when it should not have been")
	}
}

func TestDarkmodeInvalid(t *testing.T) {
	defer gock.Off()
	darkmodeURL := "http://darkmode.datasektionen.se"
	gock.New(darkmodeURL).
		Get("").
		Reply(http.StatusOK).
		JSON("hehe, not a bool n00b")

	oldpoll := time.Date(1983, time.October, 7, 17, 0, 0, 0, time.Local)
	ds := DarkmodeStatus{
		Darkmode: false,
		LastPoll: oldpoll,
		Url:      darkmodeURL,
		Mutex:    sync.RWMutex{},
	}

	expected := true
	got := Darkmode(&ds)
	if got != expected {
		t.Errorf("got %v, wanted %v", got, expected)
	}
	if ds.LastPoll != oldpoll {
		t.Errorf("The darkmode status has been updated with an invalid url")
	}
}

func TestPathIntSeparator(t *testing.T) {
	// Test case 1: Valid inputs
	param := []string{
		"/123",
		"123",
		"0",
		"/1"}

	expected := []int{
		123,
		123,
		0,
		1}

	for i := 0; i < len(param); i++ {
		got, err := pathIntSeparator(param[i])
		if err != nil {
			t.Error(err)
		}
		if got != expected[i] {
			t.Errorf("got %v, wanted %v", got, expected[i])
		}
	}

	// Test case 2: Invalid inputs
	param2 := []string{
		"32a",
		"abcdef",
		"/1/1",
		"//1",
	}

	for i := 0; i < len(param2); i++ {
		_, err := pathIntSeparator(param2[i])
		if err == nil {
			t.Errorf(`expected error from string "%v"`, param2[i])
		}
	}
}
