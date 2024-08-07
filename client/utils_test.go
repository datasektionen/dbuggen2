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
	t.Run("valid picture", func(t *testing.T) {
		mp := sql.NullString{String: "https://example.com/cover.jpg", Valid: true}
		expected := `<img src="https://example.com/cover.jpg" style="max-width: 10vw;">`
		got := memberpicture(mp)
		if string(got) != expected {
			t.Errorf("got %v, wanted %v", got, expected)
		}
	})

	// Test case 2: Invalid picture
	t.Run("invalid picture", func(t *testing.T) {
		mp := sql.NullString{String: "", Valid: false}
		expected := `<img src="/assets/default_member.svg" style="max-width: 10vw;">`
		got := memberpicture(mp)
		if string(got) != expected {
			t.Errorf("got %v, wanted %v", got, expected)
		}
	})
}

func TestDisplaymemberize(t *testing.T) {
	t.Run("empty list of members", func(t *testing.T) {
		members := make([]database.Member, 0)
		got := displaymemberize(members)
		if len(got) != 0 {
			t.Errorf("length of displaymembers is %v, not 0", len(got))
		}
	})

	t.Run("list with multiple members", func(t *testing.T) {
		members := []database.Member{
			{
				KthID:        "testsupp",
				PreferedName: sql.NullString{Valid: true, String: "Testerino"},
				PictureURL:   sql.NullString{Valid: true, String: "rickroll.mp4"},
				Title:        "the cool one",
				Active:       true,
			},
			{
				KthID:        "test1",
				PreferedName: sql.NullString{Valid: true, String: "Test 1sson"},
				PictureURL:   sql.NullString{Valid: false, String: ""},
				Title:        "1ssons frestelse",
				Active:       true,
			},
			{
				KthID:        "test2",
				PreferedName: sql.NullString{Valid: true, String: "TE S. T"},
				PictureURL:   sql.NullString{Valid: true, String: "darknet.org/virus.exe"},
				Title:        "",
				Active:       true,
			},
		}

		expected := []displayMember{
			{
				KthID:   "redaqtionen/testsupp",
				Name:    "Testerino",
				Picture: memberpicture(members[0].PictureURL),
				Title:   "the cool one",
			},
			{
				KthID:   "redaqtionen/test1",
				Name:    "Test 1sson",
				Picture: memberpicture(members[1].PictureURL),
				Title:   "1ssons frestelse",
			},
			{
				KthID:   "redaqtionen/test2",
				Name:    "TE S. T",
				Picture: memberpicture(members[2].PictureURL),
				Title:   "",
			},
		}

		got := displaymemberize(members)
		if len(got) != len(expected) {
			t.Fatalf("list of display members is %v, instead of %v", len(got), len(expected))
		}

		for i, g := range got {
			if g != expected[i] {
				t.Errorf("gotten value of %v should be %v", g, expected[i])
			}
		}
	})
}

func TestGetChefreds(t *testing.T) {
	dfunktURL := "https://dfunkt.datasektionen.se"
	t.Run("not ok response", func(t *testing.T) {
		defer gock.Off()
		gock.New(dfunktURL).
			Get("api/role/chefred/current").
			Reply(http.StatusNotFound).
			JSON(`{
					"role": {
						"id": 1,
						"title": "Chefredaqtör",
						"description": "Är ordförande för Redaqtionen",
						"identifier": "chefred",
						"email": "chefred@dbu.gg",
						"active": true,
						"Group": {
						"name": "Funktionärer",
						"identifier": "dfunk"
						}
					},
					"mandates": [
						{
							"start": "2024-01-01",
							"end": "2024-12-31",
							"User": {
								"first_name": "Che",
								"last_name": "Fred",
								"email": null,
								"kthid": "chefen",
								"ugkthid": "chefefe"
							}
						}
					]
				}`)
		got := getChefreds(dfunktURL)
		if len(got) != 0 {
			t.Errorf("the result should be empty, but it is %v", got)
		}
	})

	t.Run("no current chefred", func(t *testing.T) {
		defer gock.Off()
		gock.New(dfunktURL).
			Get("api/role/chefred/current").
			Reply(http.StatusOK).
			JSON(`{
					"role": {
						"id": 1,
						"title": "Chefredaqtör",
						"description": "Är ordförande för Redaqtionen",
						"identifier": "chefred",
						"email": "chefred@dbu.gg",
						"active": true,
						"Group": {
						"name": "Funktionärer",
						"identifier": "dfunk"
						}
					},
					"mandates": []
				}`)

		got := getChefreds(dfunktURL)
		if len(got) != 0 {
			t.Errorf("the result should be empty, but it is %v", got)
		}
	})

	t.Run("single chefred", func(t *testing.T) {
		defer gock.Off()
		gock.New(dfunktURL).
			Get("api/role/chefred/current").
			Reply(http.StatusOK).
			JSON(`{
					"role": {
						"id": 1,
						"title": "Chefredaqtör",
						"description": "Är ordförande för Redaqtionen",
						"identifier": "chefred",
						"email": "chefred@dbu.gg",
						"active": true,
						"Group": {
						"name": "Funktionärer",
						"identifier": "dfunk"
						}
					},
					"mandates": [
						{
							"start": "2024-01-01",
							"end": "2024-12-31",
							"User": {
								"first_name": "Che",
								"last_name": "Fred",
								"email": null,
								"kthid": "chefen",
								"ugkthid": "chefefe"
							}
						}
					]
				}`)

		got := getChefreds(dfunktURL)
		expected := "chefen"
		if len(got) != 1 {
			t.Fatalf("there should only be a single chefred, there are %v many: %v", len(got), got)
		}

		if got[0] != expected {
			t.Errorf("result should be %v, but is %v", expected, got[0])
		}
	})

	t.Run("multiple chefreds", func(t *testing.T) {
		defer gock.Off()
		gock.New(dfunktURL).
			Get("api/role/chefred/current").
			Reply(http.StatusOK).
			JSON(`{
					"role": {
						"id": 1,
						"title": "Chefredaqtör",
						"description": "Är ordförande för Redaqtionen",
						"identifier": "chefred",
						"email": "chefred@dbu.gg",
						"active": true,
						"Group": {
						"name": "Funktionärer",
						"identifier": "dfunk"
						}
					},
					"mandates": [
						{
							"start": "2024-01-01",
							"end": "2024-12-31",
							"User": {
								"first_name": "Che",
								"last_name": "Fred",
								"email": null,
								"kthid": "chefen",
								"ugkthid": "chefefe"
							}
						},
						{
							"start": "2024-07-01",
							"end": "2025-12-31",
							"User": {
								"first_name": "Fred",
								"last_name": "Che",
								"email": null,
								"kthid": "bossen",
								"ugkthid": "bosososososo"
							}
						}
					]
				}`)

		got := getChefreds(dfunktURL)
		expected := []string{"chefen", "bossen"}
		if len(got) != len(expected) {
			t.Fatalf("there should only be 2 chefreds, there are %v many: %v", len(got), got)
		}

		for i := 0; i < len(got); i++ {
			if got[i] != expected[i] {
				t.Errorf("result should be %v, but is %v", expected[i], got[i])
			}
		}
	})
}

func TestRemoveDuplicateChefreds(t *testing.T) {
	t.Run("both lists empty", func(t *testing.T) {
		chefredsIDs := make([]string, 0)
		members := make([]database.Member, 0)

		chefreds, members := removeDuplicateChefreds(chefredsIDs, members)
		if len(chefreds) != 0 {
			t.Errorf("list of chefreds should be empty, is %v", len(members))
		}
		if len(members) != 0 {
			t.Errorf("list of members should be empty, is %v", len(members))
		}
	})

	t.Run("existing chefred with no members", func(t *testing.T) {
		chefredsIDs := []string{"chefen"}
		members := make([]database.Member, 0)

		chefreds, members := removeDuplicateChefreds(chefredsIDs, members)
		if len(chefreds) != 1 {
			t.Errorf("list of chefreds should 1, is %v", len(members))
		}
		if len(members) != 0 {
			t.Errorf("list of members should be empty, is %v", len(members))
		}
	})

	testsupp := database.Member{
		KthID:        "testsupp",
		PreferedName: sql.NullString{Valid: true, String: "Testerino"},
		PictureURL:   sql.NullString{Valid: true, String: "rickroll.mp4"},
		Title:        "the cool one",
		Active:       true,
	}
	test1 := database.Member{
		KthID:        "test1",
		PreferedName: sql.NullString{Valid: true, String: "Test 1sson"},
		PictureURL:   sql.NullString{Valid: false, String: ""},
		Title:        "1ssons frestelse",
		Active:       true,
	}
	test2 := database.Member{
		KthID:        "test2",
		PreferedName: sql.NullString{Valid: true, String: "TE S. T"},
		PictureURL:   sql.NullString{Valid: true, String: "darknet.org/virus.exe"},
		Title:        "",
		Active:       true,
	}

	chefen := database.Member{
		KthID:        "chefen",
		PreferedName: sql.NullString{Valid: false, String: ""},
		PictureURL:   sql.NullString{Valid: false, String: ""},
		Title:        "chefred",
		Active:       true,
	}
	bossen := database.Member{
		KthID:        "bossen",
		PreferedName: sql.NullString{Valid: false, String: ""},
		PictureURL:   sql.NullString{Valid: false, String: ""},
		Title:        "chefred",
		Active:       true,
	}

	t.Run("no chefred but multiple members", func(t *testing.T) {
		chefredsIDs := make([]string, 0)
		members := []database.Member{
			testsupp,
			test1,
			test2,
		}

		chefreds, gotMembers := removeDuplicateChefreds(chefredsIDs, members)
		if len(chefreds) != 0 {
			t.Errorf("list of chefreds should be empty, is %v", len(members))
		}
		if len(gotMembers) != len(members) {
			t.Fatalf("list of members be the same length as before, is %v; %v", len(gotMembers), gotMembers)
		}

		for i := 0; i < len(gotMembers); i++ {
			if gotMembers[i] != members[i] {
				t.Errorf("member %v should be %v", gotMembers[i], members[i])
			}
		}
	})

	t.Run("multiple chefreds and members without overlap", func(t *testing.T) {
		chefredsIDs := []string{"chefen", "bossen"}
		members := []database.Member{
			testsupp,
			test1,
			test2,
		}

		expectedChefreds := []database.Member{
			chefen,
			bossen,
		}

		gotChefreds, gotMembers := removeDuplicateChefreds(chefredsIDs, members)
		if len(gotChefreds) == len(expectedChefreds) {
			for i := 0; i < len(gotChefreds); i++ {
				if gotChefreds[i] != expectedChefreds[i] {
					t.Errorf("chefred %v should be %v", gotChefreds[i], expectedChefreds[i])
				}
			}
		} else {
			t.Errorf("length of resulting chefreds, %v, are not as expected %v", len(gotChefreds), len(expectedChefreds))
		}

		if len(gotMembers) == len(members) {
			for i := 0; i < len(gotMembers); i++ {
				if gotMembers[i] != members[i] {
					t.Errorf("member %v should be %v", gotMembers[i], members[i])
				}
			}
		} else {
			t.Errorf("length of resulting members, %v, should be unchanged %v", len(gotMembers), len(members))
		}
	})

	t.Run("multiple chefreds and members with some overlap", func(t *testing.T) {
		chefredsIDs := []string{"testsupp", "bossen"}
		members := []database.Member{
			testsupp,
			test1,
			test2,
		}

		expectedChefreds := []database.Member{
			testsupp,
			bossen,
		}

		expectedMembers := []database.Member{
			test1,
			test2,
		}

		gotChefreds, gotMembers := removeDuplicateChefreds(chefredsIDs, members)
		if len(gotChefreds) == len(expectedChefreds) {
			for i := 0; i < len(gotChefreds); i++ {
				if gotChefreds[i] != expectedChefreds[i] {
					t.Errorf("chefred %v should be %v", gotChefreds[i], expectedChefreds[i])
				}
			}
		} else {
			t.Errorf("length of resulting chefreds, %v, are not as expected %v", len(gotChefreds), len(expectedChefreds))
		}

		if len(gotMembers) == len(expectedMembers) {
			for i := 0; i < len(gotMembers); i++ {
				if gotMembers[i] != expectedMembers[i] {
					t.Errorf("member %v should be %v", gotMembers[i], expectedMembers[i])
				}
			}
		} else {
			t.Errorf("length of resulting members, %v, should be %v", len(gotMembers), len(expectedMembers))
		}
	})

	t.Run("multiple chefreds and members with full overlap", func(t *testing.T) {
		chefredsIDs := []string{"testsupp", "test1"}
		members := []database.Member{
			testsupp,
			test1,
		}

		expectedChefreds := []database.Member{
			testsupp,
			test1,
		}

		expectedMembers := make([]database.Member, 0)

		gotChefreds, gotMembers := removeDuplicateChefreds(chefredsIDs, members)
		if len(gotChefreds) == len(expectedChefreds) {
			for i := 0; i < len(gotChefreds); i++ {
				if gotChefreds[i] != expectedChefreds[i] {
					t.Errorf("chefred %v should be %v", gotChefreds[i], expectedChefreds[i])
				}
			}
		} else {
			t.Errorf("length of resulting chefreds, %v, are not as expected %v", len(gotChefreds), len(expectedChefreds))
		}

		if len(gotMembers) != len(expectedMembers) {
			t.Errorf("length of resulting members, %v, should be empty; %v", len(gotMembers), gotMembers)
		}
	})
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
