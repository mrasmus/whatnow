package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	ical "github.com/laurent22/ical-go"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		callback := r.FormValue("jsonp")
		calid := r.FormValue("calid")
		res, _ := http.Get("https://www.google.com/calendar/ical/" + calid + "%40group.calendar.google.com/public/basic.ics")

		scanner := bufio.NewScanner(res.Body)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		node, _, _, _ := ical.ParseCalendarNode(lines, 0)

		var summaries []string
		now := time.Now()

		for _, child := range node.ChildrenByName("VEVENT") {
			start := child.PropDate("DTSTART", now)
			end := child.PropDate("DTEND", now)
			if now.After(start) && now.Before(end) {
				summaries = append(summaries, child.PropString("SUMMARY", ""))
			}
		}
		jsonBytes, _ := json.Marshal(summaries)
		if callback != "" {
			fmt.Fprintf(w, "%s(%s)", callback, string(jsonBytes))
		} else {
			fmt.Fprintf(w, "%s", string(jsonBytes))
		}
	})

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}

}
