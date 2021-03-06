package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alecthomas/template"
	"github.com/davidbanham/human_duration"
	"github.com/davidbanham/required_env"

	rpio "github.com/stianeikeland/go-rpio/v4"
)

var limit = 0

type pinPair struct {
	Pin    rpio.Pin  `json:"pin"`
	StopAt time.Time `json:"stop_at"`
}

type PayloadPin struct {
	StopAt time.Time `json:"stop_at"`
	Name   string    `json:"name"`
	State  bool      `json:"state"`
	Id     int       `json:"id"`
}

var pins = []pinPair{}

var numbers = []int{
	2, 3, 4, 17, 27, 22, 10, 9, 11, 5, 6, 13, 19, 26, 14, 15, 18, 23, 24, 25, 8, 7, 12, 16, 20, 21,
}

func init() {
	required_env.Ensure(map[string]string{
		"NUM_RELAYS": "4",
	})

	if parsed, err := strconv.Atoi(os.Getenv("NUM_RELAYS")); err != nil {
		log.Fatal("env input err", err)
	} else {
		limit = parsed
	}
}

func main() {
	err := rpio.Open()
	if err != nil {
		log.Println(err)
		log.Fatal("error opening gpio memory")
	}

	for i, num := range numbers {
		if i >= limit {
			continue
		}
		pin := rpio.Pin(num)
		pin.Output()
		pin.High()
		pins = append(pins, pinPair{
			Pin: pin,
		})
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/spa", spaIndex)
	http.HandleFunc("/css/main.css", css)

	http.HandleFunc("/pins", currentState)

	p := 0
	for p < limit {
		http.HandleFunc("/pins/"+strconv.Itoa(p)+"/toggle", toggler(p))
		http.HandleFunc("/pins/"+strconv.Itoa(p)+"/on", onner(p))
		http.HandleFunc("/pins/"+strconv.Itoa(p)+"/off", offer(p))
		http.HandleFunc("/pins/"+strconv.Itoa(p), stater(p)) // Get current state
		p++
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func toggler(i int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pins[i].Pin.Toggle()
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func lengthFromForm(input string) time.Duration {
	lengthStr := input + "m"
	length, err := time.ParseDuration(lengthStr)
	if err != nil {
		return 0
	}
	return length
}

func timer(i int, length time.Duration, shouldTurnOn bool) {
	stopAt := time.Now().Add(length)

	pins[i].StopAt = stopAt

	timer := time.NewTimer(length)
	<-timer.C

	if shouldTurnOn {
		pins[i].Pin.Low()
	} else {
		pins[i].Pin.High()
	}
}

func onner(i int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		length := lengthFromForm(r.FormValue("length"))
		if length > 0 {
			go timer(i, length, false)
		}
		pins[i].Pin.Low()
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func offer(i int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pins[i].Pin.High()
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func stater(i int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := pins[i].Pin.Read()

		w.WriteHeader(http.StatusOK)

		if state == rpio.Low {
			w.Write([]byte("{\"on\": true}"))
		} else {
			w.Write([]byte("{\"on\": false}"))
		}

		return
	}
}

type pinsPayload struct {
	Pins []PayloadPin `json:"pins"`
}

func currentState(w http.ResponseWriter, r *http.Request) {
	payload := []PayloadPin{}
	for i, pin := range pins {
		pp := PayloadPin{
			Name:   strconv.Itoa(i),
			State:  pin.Pin.Read() == rpio.Low,
			StopAt: pin.StopAt,
			Id:     i,
		}
		payload = append(payload, pp)
	}

	json.NewEncoder(w).Encode(pinsPayload{
		Pins: payload,
	})
}

type PageData struct {
	Pins  []pinPair
	Title string
}

var funcMap = template.FuncMap{
	"howLong": func(t time.Time) string {
		answer := human_duration.String(t.Sub(time.Now()), "minute")
		if answer == "" {
			return human_duration.String(t.Sub(time.Now()), "second")
		}
		return answer
	},
}

func index(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Pins:  pins,
		Title: "Relay Runnner",
	}
	tmpl := template.Must(template.New("index").Funcs(funcMap).Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta http-equiv="refresh" content="30">
<title>{{.Title}}</title>
<meta http-equiv="X-UA-Compatible" content="IE=Edge">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="description" content="Turn things on. Or off. I don't mind.">
<style type="text/css">
.wrapper {
	display: grid;
	grid-template-columns: repeat( auto-fit, minmax(250px, 1fr) );
	grid-gap: 10px;
	grid-auto-rows: minmax(100px, auto);

	text-align: center;
}
@media only screen and (min-width: 1000px) {
	.wrapper {
		margin: 20%;
		margin-top: 0px;
	}
}
 /* The switch - the box around the slider */
.switch {
	position: relative;
	display: inline-block;
	width: 60px;
	height: 34px;
}

/* Hide default HTML checkbox */
.switch input {
	opacity: 0;
	width: 0;
	height: 0;
}

/* The slider */
.slider {
	position: absolute;
	cursor: pointer;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background-color: #ccc;
	-webkit-transition: .4s;
	transition: .4s;
}

.slider:before {
	position: absolute;
	content: "";
	height: 26px;
	width: 26px;
	left: 4px;
	bottom: 4px;
	background-color: white;
	-webkit-transition: .4s;
	transition: .4s;
}

input:checked + .slider {
	background-color: #2196F3;
}

input:focus + .slider {
	box-shadow: 0 0 1px #2196F3;
}

input:checked + .slider:before {
	-webkit-transform: translateX(26px);
	-ms-transform: translateX(26px);
	transform: translateX(26px);
}

/* Rounded sliders */
.slider.round {
	border-radius: 34px;
}

.slider.round:before {
	border-radius: 50%;
}

.timeinput {
	margin-top: 20px;
	max-width: 50px;
}
</style>
</head>
<body>
	<div class="wrapper">
	{{ range $i, $pair := .Pins }}
	<div>
		<h4>Relay {{$i}}</h4>
		{{ if (eq $pair.Pin.Read 1)}}
		<form method="post" action="/pins/{{$i}}/on" onchange="this.submit()">
			<label class="switch">
				<input type="checkbox">
				<span class="slider round"></span>
			</label>
			<br>
			<input class="timeinput" type="number" step="1" name="length">
		</form>
		{{ else }}
		<form method="post" action="/pins/{{$i}}/off" onchange="this.submit()">
			<label class="switch">
				<input type="checkbox" checked>
				<span class="slider round"></span>
			</label>
		</form>
		<p>Off in {{ howLong $pair.StopAt }}</p>
		{{ end }}
	</div>
	{{ end }}
	</div>
</body>
</html>
	`))
	if err := tmpl.Execute(w, data); err != nil {
		log.Fatal("template error", err)
	}
}

func spaIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(indexPage))
}

func css(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(cssPage))
}
