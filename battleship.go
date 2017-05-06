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
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const PORT = "8080"
const hit = "!"
const miss = "."
const hitAndMiss = hit + miss
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
	Y, X int // order is 00, 01, 02 [new row] 10, 11, 12
}

type player struct {
	name  string
	board [][]string
}

func coordsToPoint(yCommaX string) point {
	strs := strings.SplitN(yCommaX, ",", 2)
	y, _ := strconv.Atoi(strings.TrimSpace(strs[0]))
	x, _ := strconv.Atoi(strings.TrimSpace(strs[1]))
	return point{y, x}
}

var players = []player{
	player{"human", makeBoard()},
	player{"computer", makeBoard()},
}

// create a 10 x 10 matrix of strings with each string set to " "
func makeBoard() [][]string {
	board := [][]string{}
	for i := 0; i < 10; i++ {
		row := []string{" ", " ", " ", " ", " ", " ", " ", " ", " ", " "}
		board = append(board, row)
	}
	return board
}

// concatinate all the strings together to ease the finding of ships, etc.
func boardToStr(board [][]string) string {
	rowStrings := []string{}
	for _, row := range board {
		rowStrings = append(rowStrings, strings.Join(row, ""))
	}
	return strings.Join(rowStrings, "")
}

func invalidPoint(pt point) bool {
	return pt.X < 0 || pt.X > 9 || pt.Y < 0 || pt.Y > 9
}

func randomBool() bool { return rand.Intn(2) == 1 }

func randomPoint() point { return point{rand.Intn(10), rand.Intn(10)} }

func htmlSubmitButton(y, x int) string {
	// return HTMLData(fmt.Sprintf("\"<button type='submit' name='%d,%d'>&nbsp;</button>\"", x, y))
	return fmt.Sprintf("<button type='submit' name='yx' value='%d,%d'>%[1]d,%d</button>", y, x)
}

func template_map(homeTeam, awayTeam player) map[string]string {
	m := map[string]string{
		"HomeStatus": homeTeam.name,
		"AwayStatus": awayTeam.name,
	}
	teamBoards := map[string]([][]string){
		"H": homeTeam.board,
		"A": awayTeam.board,
	}
	for letter, board := range teamBoards {
		for y, row := range board {
			for x, s := range row {
				if letter == "H" && !strings.Contains(hitAndMiss, s) {
					// convert locs where human can drop bombs into html buttons
					s = htmlSubmitButton(y, x)
				}
				m[fmt.Sprintf("%s%d%d", letter, y, x)] = s
			}
		}
	}
	return m
}

// once computer has a hit it will want to hit the neghbors next
func neighbors(pt point) (pts []point) {
	u := point{pt.Y, pt.X - 1}
	d := point{pt.Y, pt.X + 1}
	l := point{pt.Y - 1, pt.X}
	r := point{pt.Y + 1, pt.X}
	for _, pt = range []point{u, d, l, r} {
		if invalidPoint(pt) == false {
			pts = append(pts, pt)
		}
	}
	return
}

/*
func strReplaceRune(s string, pos int, r rune) string {
	// replace the rune at index pos with rune r
	return strings.Join([]string{s[:pos], s[pos+1:]}, string(r))
}
*/

// Useful for the intitial placement of ships on the board
func pointsForShip(topLeft point, length int, across bool) (pts []point) {
	for i := 0; i < length; i++ {
		if across {
			pts = append(pts, point{topLeft.Y, topLeft.X + i})
		} else {
			pts = append(pts, point{topLeft.Y + i, topLeft.X})
		}
	}
	// if last point is invalid...
	if invalidPoint(pts[len(pts)-1]) {
		pts = []point{}
	}
	return
}

// point{0, 0} --> A1, point{9, 9} --> J10
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
	pt = point{y, x}
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

/*
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
*/

func formatRow(i int) string {
	return fmt.Sprintf("| %2d  %s  %2[1]d  %s  %2[1]d |", i, "%s")
}

// Don't allow those cheating humans to see the computer's ships!
func clokeRow(row []string) []string {
	newRow := []string{}
	for _, s := range row {
		if !strings.Contains(hitAndMiss, s) {
			s = " "
		}
		newRow = append(newRow, s)
	}
	return newRow
}

func boardDisplay(players []player) string {
	rows := []string{title, borderRow(), letterRow()}
	for i := 0; i < 10; i++ {
		homeTeam := strings.Join(players[0].board[i], "  ")
		awayTeam := strings.Join(clokeRow(players[1].board[i]), "  ")
		rows = append(rows, fmt.Sprintf(formatRow(i+1), homeTeam, awayTeam))
	}
	return strings.Join(append(append(rows, rows[2]), rows[1]), "\n")
}

func playableSquares(opponent player) []point {
	squares := []point{}
	for y, row := range opponent.board {
		for x, c := range row {
			if strings.Contains(hitAndMiss, c) == false {
				squares = append(squares, point{y, x})
			}
		}
	}
	// fmt.Printf("%v\n", squares)
	return squares
}

func hasShip(p player, s string) bool {
	return strings.Contains(boardToStr(p.board), s)
}

func hasAnyShips(p player) bool {
	for _, r := range boardToStr(p.board) {
		if !strings.ContainsRune(hitAndMissAndSpace, r) {
			return true
		}
	}
	fmt.Printf("Game Over: Player %s has no ships!\n", p.name)
	return false
}

/*
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
*/

func dropABomb(opponent player, sq point) (gameOn bool) {
	gameOn = true
	oldStr := opponent.board[sq.Y][sq.X]
	itsAHit := oldStr != " "
	splash := miss
	if itsAHit == true {
		splash = hit
	}
	// Drop the bomb into the square
	opponent.board[sq.Y][sq.X] = splash
	if itsAHit {
		if hasShip(opponent, oldStr) {
			fmt.Println("It's a hit!")
		} else {
			fmt.Printf("You sunk my battleship! %s\n", oldStr)
			gameOn = hasAnyShips(opponent)
		}
	}
	return
}

/*
func humanTurn(opponent player) (gameOn bool) {
	sq := askWhichSquare()
	if sq[:1] == "Q" {
		gameOn = false // human wants out
	} else {
		pt, _ := letterNumberToPoint(sq)
		if strings.Contains(hitAndMiss, opponent.board[pt.Y][pt.X]) {
			fmt.Println("You already tried that square.  Try different one:")
			gameOn = humanTurn(opponent)
		} else {
			gameOn = dropABomb(opponent, pt)
		}
	}
	return
}
*/

func compuTurn(opponent player) (gameOn bool) {
	ps := playableSquares(opponent)
	if len(ps) == 0 {
		panic("No playable squares!!")
	}
	sq := ps[rand.Intn(len(ps))]
	gameOn = dropABomb(opponent, sq)
	return
}

/*
func humanPlaceShip(p player, shipName string) {
	letter := string(shipName[0])
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
		oldStr := p.board[pt.Y][pt.X]
		if oldStr != " " {
			fmt.Printf("Placing %s failed: %v is %s\n", shipName, pt, oldStr)
			humanPlaceShip(p, shipName)
			return
		}
	}
	for _, pt := range pts {
		p.board[pt.Y][pt.X] = letter
	}
}
*/

func compuPlaceShip(p player, shipName string) {
	letter := string(shipName[0])
	length := ships[shipName]
	topLeft := randomPoint()
	across := randomBool()
	pts := pointsForShip(topLeft, length, across)
	if len(pts) == 0 { // the ship would have gone outside the matrix :-(
		compuPlaceShip(p, shipName) // recursive call
		return
	}
	for _, pt := range pts {
		oldStr := p.board[pt.Y][pt.X]
		if oldStr != " " { // square is already occupied
			compuPlaceShip(p, shipName) // recursive call
			return
		}
	}
	// all clear, put the ship on the board
	for _, pt := range pts {
		p.board[pt.Y][pt.X] = letter
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
	m := template_map(players[1], players[0])
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
		browserReload(PORT) // TODO: This open a new browser tab!!!
		if !gameOn {
			panic("Game over man!")
		}
	}
}

func browserReload(port string) {
	if runtime.GOOS == "darwin" {
		exec.Command("open", "http://localhost:"+port).Start()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for shipName := range ships {
		compuPlaceShip(players[0], shipName) // TODO: Allow human to place ships
		compuPlaceShip(players[1], shipName)
	}
	fmt.Println(boardDisplay(players))
	fmt.Println("Point your browser to: http://localhost:" + PORT)
	browserReload(PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, helloHandler{}))
	/*
		for gameOn {
			gameOn = humanTurn(compuPlayer) && compuTurn(humanPlayer)
			fmt.Println(board(humanPlayer, compuPlayer))
		}
	*/
}
