///usr/bin/env go run ${0} ${@} ; exit ${?}
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
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const blanks = "          "
const hitAndMiss = "!."
const hitAndMissAndSpace = hitAndMiss + " "
const letters = "ABCDEFGHIJ"
const title = "            H  u  m  a  n                  C  o  m  p  u  t  e  r"

var ships = map[string]int{
	"Aircraft Carrier": 5,
	"Battleship":       4,
	"Cruiser":          3,
	"Submarine":        3,
	"Destroyer":        2,
}

type point struct {
	X, Y int
}

type player struct {
	name  string
	board []string
}

var players = []player{makePlayer("human"), makePlayer("computer")}

func invalidPoint(pt point) bool {
	return pt.X < 0 || pt.X > 9 || pt.Y < 0 || pt.Y > 9
}

func randomBool() bool { return rand.Intn(2) == 1 }

func randomPoint() point { return point{rand.Intn(10), rand.Intn(10)} }

func htmlSubmitButton(x, y int) string {
	// return HTMLData(fmt.Sprintf("\"<button type='submit' name='%d,%d'>&nbsp;</button>\"", x, y))
	return fmt.Sprintf("<button type='submit' name='xy' value='%d,%d'>%[1]d,%d</button>", x, y)
}

func template_map(homeTeam, awayTeam player) map[string]string {
	m := map[string]string{
		"HomeStatus": homeTeam.name,
		"AwayStatus": awayTeam.name,
	}
	teamBoards := map[string]([]string){
		"H": homeTeam.board,
		"A": awayTeam.board,
	}
	for letter, board := range teamBoards {
		for y, row := range board {
			for x, c := range row {
				s := string(c)
				if letter == "H" && strings.ContainsRune(hitAndMiss, c) == false {
					// convert locs where human can drop bombs into html buttons
					s = htmlSubmitButton(x, y)
				}
				m[fmt.Sprintf("%s%d%d", letter, x, y)] = s
			}
		}
	}
	return m
}

func neighbors(pt point) (pts []point) {
	u := point{pt.X - 1, pt.Y}
	d := point{pt.X + 1, pt.Y}
	l := point{pt.X, pt.Y - 1}
	r := point{pt.X, pt.Y + 1}
	for _, pt = range []point{u, d, l, r} {
		if invalidPoint(pt) == false {
			pts = append(pts, pt)
		}
	}
	return
}

func strReplaceRune(s string, pos int, r rune) string {
	// replace the rune at index pos with rune r
	return strings.Join([]string{s[:pos], s[pos+1:]}, string(r))
}

func pointsForShip(topLeft point, length int, across bool) (pts []point) {
	for i := 0; i < length; i++ {
		if across {
			pts = append(pts, point{topLeft.X + i, topLeft.Y})
		} else {
			pts = append(pts, point{topLeft.X, topLeft.Y + i})
		}
	}
	// if last point is invalid...
	if invalidPoint(pts[len(pts)-1]) {
		pts = []point{}
	}
	return
}

func pointToLetterNumber(pt point) (string, error) {
	if invalidPoint(pt) {
		return "", fmt.Errorf("invalid point %v", pt)
	}
	return fmt.Sprintf("%c%d", letters[pt.X], pt.Y+1), nil
}

func letterNumberToPoint(s string) (pt point, err error) {
	s = strings.ToUpper(s)
	x := strings.Index(letters, s[:1])
	y := strings.Index("12345678910", s[1:])
	pt = point{x, y}
	if invalidPoint(pt) {
		if s[:1] == "Q" {
			panic("User quit.")
		}
		err = errors.New("invalid: Try 'A1' or 'J10'")
	}
	return
}

func borderRow() string {
	return strings.Join([]string{"+", "+"}, strings.Repeat("-", 72))
}

func letterRow() string {
	letters := strings.Join(charsInStr(letters), "  ")
	return fmt.Sprintf("|     %s  ||  %s     |", letters, letters)
}

func charsInStr(str string) []string {
	// return [c for c in str]
	letters := []string{}
	for _, c := range str {
		letters = append(letters, string(c))
	}
	return letters
}

func clokeInStr(str string) []string {
	letters := []string{}
	for _, c := range str {
		if strings.ContainsRune(hitAndMiss, c) == false {
			c = ' ' // cloke the battleships!
		}
		letters = append(letters, string(c))
	}
	return letters
}

func formatRow(i int) string {
	return fmt.Sprintf("| %2d  %s  %2[1]d  %s  %2[1]d |", i, "%s")
}

func board(homeTeam, awayTeam player) string {
	rows := []string{title, borderRow(), letterRow()}
	for i := 0; i < 10; i++ {
		h := strings.Join(charsInStr(homeTeam.board[i]), "  ")
		a := strings.Join(clokeInStr(awayTeam.board[i]), "  ")
		rows = append(rows, fmt.Sprintf(formatRow(i+1), h, a))
	}
	return strings.Join(append(append(rows, rows[2]), rows[1]), "\n")
}

func makePlayer(name string) player {
	return player{name, []string{blanks, blanks, blanks, blanks, blanks,
		blanks, blanks, blanks, blanks, blanks}}
}

func playableSquares(opponent player) []point {
	squares := []point{}
	for y, row := range opponent.board {
		for x, c := range row {
			if strings.ContainsRune(hitAndMiss, c) == false {
				squares = append(squares, point{x, y})
			}
		}
	}
	fmt.Printf("%v\n", squares)
	return squares
}

func hasShip(p player, r rune) bool {
	return strings.ContainsRune(strings.Join(p.board, ""), r)
}

func hasAnyShips(p player) bool {
	for _, c := range strings.Join(p.board, "") {
		if strings.ContainsRune(hitAndMissAndSpace, c) == false {
			return true
		}
	}
	fmt.Printf("Game Over: Player %s has no ships!\n", p.name)
	return false
}

func askUser(prompt string) (string, error) {
	fmt.Print(prompt + " ")
	text, err := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimSpace(text), err
}

func askWhichSquare() string {
	text, _ := askUser("Enter a square between A1 and J10:")
	// pt, _ := letterNumberToPoint(text)
	// fmt.Printf("(%T) %s ==> %v|\n", text, text, pt)
	return text
}

func askIfAcross() bool {
	text, _ := askUser("[A]cross or [D]own:")
	return strings.ToUpper(text)[0] == 'A'
}

func dropABomb(opponent player, sq point) (gameOn bool) {
	gameOn = true
	oldRune := rune(opponent.board[sq.Y][sq.X])
	itsAHit := oldRune != ' '
	splash := '.'
	if itsAHit == true {
		splash = '!'
	}
	// Drop the bomb into the row
	opponent.board[sq.Y] = strReplaceRune(opponent.board[sq.Y], sq.X, splash)
	if itsAHit {
		if hasShip(opponent, rune(oldRune)) {
			fmt.Println("It's a hit!")
		} else {
			fmt.Printf("You sunk my battleship! %c\n", oldRune)
			gameOn = hasAnyShips(opponent)
		}
	}
	return
}

func humanTurn(opponent player) (gameOn bool) {
	sq := askWhichSquare()
	if sq[:1] == "Q" {
		gameOn = false // human wants out
	} else {
		pt, _ := letterNumberToPoint(sq)
		oldRune := rune(opponent.board[pt.Y][pt.X])
		if strings.ContainsRune(hitAndMiss, oldRune) {
			fmt.Println("You already tried that square.  Try different one:")
			gameOn = humanTurn(opponent)
		} else {
			gameOn = dropABomb(opponent, pt)
		}
	}
	return
}

func compuTurn(opponent player) (gameOn bool) {
	ps := playableSquares(opponent)
	if len(ps) == 0 {
		panic("No playable squares!!")
	}
	sq := ps[rand.Intn(len(ps))]
	gameOn = dropABomb(opponent, sq)
	return
}

func humanPlaceShip(p player, shipName string) {
	letter := shipName[0]
	length := ships[shipName]
	fmt.Printf("Placing %s (%c * %d)...\n", shipName, letter, length)
	topLeft, _ := letterNumberToPoint(askWhichSquare())
	across := askIfAcross()
	pts := pointsForShip(topLeft, length, across)
	if len(pts) == 0 {
		humanPlaceShip(p, shipName)
		return
	}
	for _, pt := range pts {
		oldRune := p.board[pt.Y][pt.X]
		if oldRune != ' ' {
			fmt.Printf("Placing %s failed: %v is %c\n", shipName, pt, oldRune)
			humanPlaceShip(p, shipName)
			return
		}
	}
	for _, pt := range pts {
		p.board[pt.Y] = strReplaceRune(p.board[pt.Y], pt.X, rune(letter))
	}
}

func compuPlaceShip(p player, shipName string) {
	letter := shipName[0]
	length := ships[shipName]
	topLeft := randomPoint()
	across := randomBool()
	pts := pointsForShip(topLeft, length, across)
	if len(pts) == 0 {
		compuPlaceShip(p, shipName)
		return
	}
	for _, pt := range pts {
		oldRune := p.board[pt.Y][pt.X]
		if oldRune != ' ' {
			// fmt.Printf("Placing %s failed: %v is %c\n", shipName, pt, oldRune)
			compuPlaceShip(p, shipName)
			return
		}
	}
	for _, pt := range pts {
		p.board[pt.Y] = strReplaceRune(p.board[pt.Y], pt.X, rune(letter))
	}
}

var templates = template.Must(template.ParseFiles("battleship.html"))

/*
func renderTemplate(w http.ResponseWriter, homeTeam, awayTeam player) {
	m := template_map(homeTeam, awayTeam)
	if err := templates.ExecuteTemplate(w, "battleship.html", m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
*/

type helloHandler struct{}

func (h helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
	/*
		   }

		   func battleshipHandler(w http.ResponseWriter, req *http.Request) {
		*-/
		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
		w.WriteHeader(http.StatusOK)

		io.WriteString(w, "Let's play battleship\n")
	*/
	// templates := template.Must(template.ParseFiles("battleship.html"))
	// templates := template.ParseFiles("battleship.html")
	m := template_map(players[0], players[1])
	if err := templates.ExecuteTemplate(w, "battleship.html", m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	buttonPressed := r.FormValue("xy")
	fmt.Println(buttonPressed)
	if buttonPressed != "" {
		x, _ := strconv.Atoi(buttonPressed[:1])
		y, _ := strconv.Atoi(buttonPressed[2:])
		dropABomb(players[0], point{x, y})
		browserReload()
	} else {
		fmt.Println(buttonPressed)
	}
	// r.ParseForm()
	// fmt.Println(r.Form)
}

func browserReload() {
	if runtime.GOOS == "darwin" {
		exec.Command("open", "http://localhost:8083").Start()
	}
}

func main() {
	// fmt.Println(htmlSubmitButton(point{1, 2}))
	rand.Seed(time.Now().UnixNano())
	humanPlayer := players[0] // makePlayer("human")
	compuPlayer := players[1] // makePlayer("computer")
	// fmt.Println(template_map(humanPlayer, compuPlayer))
	// panic("Dude.")
	for shipName := range ships {
		// fmt.Println(board(humanPlayer, compuPlayer))
		// humanPlaceShip(humanPlayer, shipName)
		compuPlaceShip(humanPlayer, shipName)
		compuPlaceShip(compuPlayer, shipName)
	}
	fmt.Println(board(humanPlayer, compuPlayer))
	// gameOn := hasAnyShips(humanPlayer) && hasAnyShips(compuPlayer)
	fmt.Println("Point your browser to: http://localhost:8080")
	// http.Handle("/", battleshipHandler)
	//log.Fatal(http.ListenAndServe(":8080", nil))
	browserReload()
	log.Fatal(http.ListenAndServe(":8083", helloHandler{}))
	/*
		for gameOn {
			gameOn = humanTurn(compuPlayer) && compuTurn(humanPlayer)
			fmt.Println(board(humanPlayer, compuPlayer))
		}
	*/
}
