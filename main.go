package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"groupie-tracker/backend"
)

var logChan = make(chan string, 100)

func init() {
	go func() {
		logFile, err := os.OpenFile("history.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to open log file:", err)
			return
		}
		defer logFile.Close()

		logger := log.New(logFile, "", log.LstdFlags)
		for msg := range logChan {
			logger.Println(msg)
		}
	}()
}

func logHistory(message string) {
	logChan <- message
}

func main() {
	if len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, "check args!!!")
		return
	}

	fmt.Println("This is your port: http://localhost:8080/")

	startTime := time.Now()
	logHistory(fmt.Sprintf("Server started at %s", startTime.Format(time.RFC1123)))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logHistory(fmt.Sprintf("Accessed Home Page - %s", r.RemoteAddr))
		backend.HandleHome(w, r)
	})
	http.HandleFunc("/Artist/", func(w http.ResponseWriter, r *http.Request) {
		logHistory(fmt.Sprintf("Accessed Artist Page - %s", r.RemoteAddr))
		backend.HandlePage(w, r)
	})
	http.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		logHistory(fmt.Sprintf("404 Error Page Accessed - %s", r.RemoteAddr))
		backend.ErrorHandler(w, r)
	})
	http.HandleFunc("/frontend/css/", func(w http.ResponseWriter, r *http.Request) {
		logHistory(fmt.Sprintf("Accessed CSS Asset - %s", r.RemoteAddr))
		backend.CssHandler(w, r)
	})
	http.HandleFunc("/frontend/images/", func(w http.ResponseWriter, r *http.Request) {
		logHistory(fmt.Sprintf("Accessed Image Asset - %s", r.RemoteAddr))
		backend.ImageHandler(w, r)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logHistory(fmt.Sprintf("Server stopped due to error: %s", err))
		fmt.Fprintln(os.Stderr, "Server error:", err)
	}
}
