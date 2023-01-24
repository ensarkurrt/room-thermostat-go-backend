package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	fromWifi := false
	isOpen := false
	temperature := 0.0
	desiredTempOnWifi := 0.0
	desiredTempOnBluetooth := 0.0

	fromWhatText := func() string {
		if fromWifi {
			return "wifi"
		} else {
			return "bt"
		}
	}

	desiredTempText := func() string {
		if fromWifi {
			return strconv.FormatFloat(desiredTempOnWifi, 'f', 2, 32)
		} else {
			return strconv.FormatFloat(desiredTempOnBluetooth, 'f', 2, 32)
		}
	}

	isOpenText := func() string {
		if isOpen {
			return "opened"
		} else {
			return "closed"
		}
	}

	statusResponse := func() string {
		return `{"fromWhat": "` + fromWhatText() + `","status": "` + isOpenText() + `" , "desiredTemp" : ` + desiredTempText() + ` ,  "temperature": ` + strconv.FormatFloat(temperature, 'f', 2, 32) + `}`
	}

	infoHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, statusResponse())
	}

	changeTypeHandler := func(w http.ResponseWriter, req *http.Request) {
		_fromWhat := req.URL.Query().Get("from-what")
		_desiredTemp := req.URL.Query().Get("desired-temp")

		fmt.Println(req.URL.Query())
		if _fromWhat == "wifi" {
			fromWifi = true
		} else {
			fromWifi = false
		}

		if s, err := strconv.ParseFloat(_desiredTemp, 32); err == nil {
			if fromWifi {
				desiredTempOnWifi = s
			} else {
				desiredTempOnBluetooth = s
			}
		}

		io.WriteString(w, statusResponse())
	}

	updateTemperatureHandler := func(w http.ResponseWriter, req *http.Request) {
		/* get temperature from url parameter */

		_temperature := req.URL.Query().Get("temperature")
		_desiredTempOnBluetooth := req.URL.Query().Get("desired-temp")
		_isOpen := req.URL.Query().Get("is-open")

		fmt.Println(req.URL.Query())

		if s, err := strconv.ParseFloat(_temperature, 32); err == nil {
			temperature = s
		}

		if s, err := strconv.ParseFloat(_desiredTempOnBluetooth, 32); err == nil {
			desiredTempOnBluetooth = s
		}

		if s, err := strconv.ParseInt(_isOpen, 10, 64); err == nil {
			if s == 1 {
				isOpen = true
			} else {
				isOpen = false
			}
		}

		io.WriteString(w, statusResponse())
	}

	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/change", changeTypeHandler)
	http.HandleFunc("/update-temp", updateTemperatureHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	for _, encodedRoute := range strings.Split(os.Getenv("ROUTES"), ",") {
		if encodedRoute == "" {
			continue
		}
		pathAndBody := strings.SplitN(encodedRoute, "=", 2)
		path, body := pathAndBody[0], pathAndBody[1]
		http.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, body)
		})
	}

	bindAddr := fmt.Sprintf(":%s", port)

	fmt.Println()
	fmt.Printf("==> Server listening at %s 🚀\n", bindAddr)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}
}
