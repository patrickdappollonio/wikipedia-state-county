package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/franela/goreq"
)

const (
	WIKIPEDIA_URI = "https://en.wikipedia.org/wiki/List_of_counties_by_U.S._state"
	TABLE_STATE   = "states"
	TABLE_COUNTY  = "counties"
)

var createStatesTable = `drop table if exists %s;
create table %s (
	id int unsigned not null auto_increment,
	name varchar(50) not null,
	primary key(id)
) Engine = InnoDB;`

var createCountiesTable = `drop table if exists %s;
create table %s (
	state_id int unsigned not null,
	id int unsigned not null auto_increment,
	name varchar(50) not null,
	primary key (id, state_id),
	foreign key (state_id) references %s(id)
) Engine = InnoDB;`

var insertStateStatement = `insert into %s(id, name) values(%v, "%s");`
var insertCountyStatement = `insert into %s(state_id, name) values(%v, "%s");`

var (
	reTitles       = regexp.MustCompile(`<li><a href="(.*)" title="(.*)">(.*)<\/a>`)
	reAnchorText   = regexp.MustCompile(`">([^<]+)<\/a>`)
	anchorReplacer = strings.NewReplacer(`">`, "", `</a>`, "")
	strSeeAlso     = `<h2><span class="mw-headline" id="See_also">See also`
	strCommaSpace  = ", "
	stateCounty    = make(map[string][]string, 0)
	tmplIterator   = 0

	suffixes = []string{
		"County",
		"Census Area",
		"Borough",
	}

	prefixes = []string{
		"Consolidated Municipality of",
		"Municipality and County of",
		"City and County of",
		"City and Borough of",
		"Town and County of",
		"City of",
		"Municipality of",
	}
)

func main() {
	res, err := goreq.Request{Uri: WIKIPEDIA_URI}.Do()

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	str, err := res.Body.ToString()

	if err != nil {
		panic(err)
	}

	pair := strings.Split(str, strSeeAlso)
	str = pair[0]

	found := reTitles.FindAllString(str, -1)

	var state, county string
	var lastindex int
	for _, foundvalue := range found {
		foundvalue = reAnchorText.FindString(foundvalue)
		foundvalue = anchorReplacer.Replace(foundvalue)

		lastindex = strings.LastIndex(foundvalue, strCommaSpace)
		county = foundvalue[:lastindex]

		for _, prefix := range prefixes {
			county = strings.TrimPrefix(county, prefix)
		}
		for _, suffix := range suffixes {
			county = strings.TrimSuffix(county, suffix)
		}

		county = strings.TrimSpace(county)
		state = strings.TrimSpace(foundvalue[lastindex+len(strCommaSpace):])

		if _, exists := stateCounty[state]; !exists {
			stateCounty[state] = make([]string, 0)
		}

		stateCounty[state] = append(stateCounty[state], county)
	}

	fmt.Printf(createStatesTable, TABLE_STATE, TABLE_STATE)
	fmt.Printf("\n\n")

	fmt.Printf(createCountiesTable, TABLE_COUNTY, TABLE_COUNTY, TABLE_STATE)
	fmt.Printf("\n\n")

	iterator := 1
	for state, counties := range stateCounty {
		fmt.Printf(insertStateStatement, TABLE_STATE, iterator, state)
		fmt.Printf("\n")

		for _, county := range counties {
			fmt.Printf(insertCountyStatement, TABLE_COUNTY, iterator, county)
			fmt.Printf("\n")
		}

		iterator++
		fmt.Printf("\n")
	}
}
