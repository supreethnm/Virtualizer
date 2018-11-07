package logger

import (
	"fmt"
	"log" // "github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

func MonitoringResponseLogger(time time.Time, header http.Header, body string) {

	f, err := os.OpenFile("Monitor.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	log.SetOutput(f)
	headers := fmt.Sprintf("%v", header)

	headers = strings.Trim(headers, "map")

	log.Printf(
		"Response sent at %s\nHeaders: %s\nBody: %s\n\n",
		time,
		headers,
		body,
	)
}
func MonitoringRequestLogger(time time.Time, header http.Header, body string, operationname string, url string) {
	f, err := os.OpenFile("Monitor.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	log.SetOutput(f)
	headers := fmt.Sprintf("%v", header)

	headers = strings.Trim(headers, "map")

	log.Printf(
		"Requested recieved at %s\nURL: %s \nHeaders: %s\nBody: %s\nServed By Operation: %s\n\n",
		time,
		url,
		headers,
		body,
		operationname,
	)
}

func Logger(inner http.HandlerFunc, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner(w, r)

		f, err := os.OpenFile("Access.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}

		//defer to close when you're done with it, not because you think it's idiomatic!
		defer f.Close()

		//set output of logs to f
		log.SetOutput(f)

		log.Printf(
			"%s %s %f",
			r.Method,
			r.RequestURI,
			time.Since(start).Seconds(),
		)
	}
}
