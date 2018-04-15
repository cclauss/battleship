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
	"fmt"
	"strconv"
	"strings"
)

const hit = "!"
const miss = "."
const hitAndMiss = hit + miss
const hitAndMissAndSpace = hitAndMiss + " "
const letters = "ABCDEFGHIJ"
const title = "            H  u  m  a  n                  C  o  m  p  u  t  e  r"

type point struct {
	Y, X int // order is 00, 01, 02 [new row] 10, 11, 12
}

func (pt point) invalid() bool {
	return pt.X < 0 || pt.X > 9 || pt.Y < 0 || pt.Y > 9
}

// once computer has a hit it will want to hit the neghbors next
func neighbors(pt point) (pts []point) {
	u := point{pt.Y, pt.X - 1}
	d := point{pt.Y, pt.X + 1}
	l := point{pt.Y - 1, pt.X}
	r := point{pt.Y + 1, pt.X}
	for _, pt = range []point{u, d, l, r} {
		if pt.invalid() == false {
			pts = append(pts, pt)
		}
	}
	return
}

func coordsToPoint(yCommaX string) point {
	strs := strings.SplitN(yCommaX, ",", 2)
	y, _ := strconv.Atoi(strings.TrimSpace(strs[0]))
	x, _ := strconv.Atoi(strings.TrimSpace(strs[1]))
	return point{y, x}
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
	if pts[len(pts)-1].invalid() {
		pts = []point{}
	}
	return
}

func charsInStr(str string) []string {
	// Python: return [c for c in str]
	letters := []string{}
	for _, c := range str {
		letters = append(letters, string(c))
	}
	return letters
}

func formatRow(i int) string {
	return fmt.Sprintf("| %2d  %s  %2[1]d  %s  %2[1]d |", i, "%s")
}

func borderRow() string {
	return strings.Join([]string{"+", "+"}, strings.Repeat("-", 72))
}

func letterRow() string {
	letters := strings.Join(charsInStr(letters), "  ")
	return fmt.Sprintf("|     %s  ||  %s     |", letters, letters)
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
