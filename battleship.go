///usr/bin/env go run *.go ${@} ; exit ${?}
// the line above is a shebang-like line for golang
// chmod +x battleship.go
// ./battleship.go

package main

/*
            H  u  m  a  n                  C  o  m  p  u  t  e  r
+------------------------------------------------------------------------+
|     A  B  C  D  E  F  G  H  I  J  ||  A  B  C  D  E  F  G  H  I  J     |
|  1                                 1                                 1 |
|  2        A  A  A  A  A            2                                 2 |
|  3                                 3                                 3 |
|  4           B  B  B  B            4                                 4 |
|  5                                 5                                 5 |
|  6  C  C  C                        6                                 6 |
|  7                                 7                                 7 |
|  8                       S  S  S   8                                 8 |
|  9                                 9                                 9 |
| 10              D  D              10                                10 |
|     A  B  C  D  E  F  G  H  I  J  ||  A  B  C  D  E  F  G  H  I  J     |
+------------------------------------------------------------------------+
*/

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"text/template"
	"time"
)

const port = "8080"

var ships = map[string]int{
	"Aircraft Carrier": 5,
	"Battleship":       4,
	"Cruiser":          3,
	"Submarine":        3,
	"Destroyer":        2,
}

type player struct {
	name  string
	board [][]string
}

var players = []player{
	player{"human", makeBoard()},
	player{"computer", makeBoard()},
}

// TODO: Is there a way to do this using html/template instead of text/template?
func htmlSubmitButton(y, x int) string {
	// return HTMLData(fmt.Sprintf("\"<button type='submit' name='%d,%d'>&nbsp;</button>\"", x, y))
	return fmt.Sprintf("<button type='submit' name='yx' value='%d,%d'>%[1]d,%d</button>", y, x)
}

func templateMap(players []player) map[string]string {
	const nbsp = string('\u00A0')
	m := map[string]string{
		"HomeStatus": players[0].name,
		"AwayStatus": players[1].name,
	}
	for i, player := range players {
		letter := string("HA"[i]) // Home v.s. Away
		for y, row := range player.board {
			for x, s := range row {
				if letter == "A" {
					if !strings.Contains(hitAndMiss, s) {
						// convert locs where human can drop bombs into html buttons
						s = htmlSubmitButton(y, x)
					}
				} else { // spaces --> nbsp to keep HTML table row heights constant
					s = strings.Replace(s, " ", nbsp, -1)
				}
				m[fmt.Sprintf("%s%d%d", letter, y, x)] = s
			}
		}
	}
	return m
}

var templates = template.Must(template.ParseFiles("battleship.html"))

type helloHandler struct{} // TODO: Is there a cleaner way to set up the handler?

func (h helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := templateMap(players)
	if err := templates.ExecuteTemplate(w, "battleship.html", m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	buttonPressed := r.FormValue("yx")
	if buttonPressed != "" {
		gameOn := dropABomb(players[1], coordsToPoint(buttonPressed))
		if gameOn {
			gameOn = compuTurn(players[0])
		}
		fmt.Println(boardDisplay(players))
		// TODO: Is there q way to force a refresh of the webpage without opening a new browser tab!!!
		openBrowser("http://localhost:" + port) // This opens a new browser tab!!!
		if !gameOn {
			panic("Game over man!")
		}
	}
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for shipName := range ships {
		compuPlaceShip(players[0], shipName) // TODO: Allow human to place ships
		compuPlaceShip(players[1], shipName)
	}
	fmt.Println(boardDisplay(players))
	url := "http://localhost:" + port
	fmt.Println("Point your browser to: " + url)
	openBrowser(url)
	// (d *Dialer) Dial(urlStr string, requestHeader http.Header) (*Conn, *http.Response, error)
	log.Fatal(http.ListenAndServe(":"+port, helloHandler{}))
}
