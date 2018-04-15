package main

import (
	"math/rand"
)

// contains the game logic for the computer player

func randomBool() bool { return rand.Intn(2) == 1 }

func randomPoint() point { return point{rand.Intn(10), rand.Intn(10)} }

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

func compuTurn(opponent player) (gameOn bool) {
	ps := playableSquares(opponent)
	if len(ps) == 0 {
		panic("No playable squares!!")
	}
	sq := ps[rand.Intn(len(ps))]
	gameOn = dropABomb(opponent, sq)
	return
}
