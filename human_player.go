package main

// contains the game logic for the human player

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

// point{0, 0} --> A1, point{9, 9} --> J10
func pointToLetterNumber(pt point) (string, error) {
	if pt.invalid() {
		return "", fmt.Errorf("invalid point %v", pt)
	}
	return fmt.Sprintf("%c%d", letters[pt.X], pt.Y+1), nil
}

func letterNumberToPoint(s string) (pt point, err error) {
	s = strings.ToUpper(s)
	x := strings.Index(letters, s[:1])
	y := strings.Index("12345678910", s[1:])
	pt = point{y, x}
	if pt.invalid() {
		if s[:1] == "Q" {
			panic("User quit.")
		}
		err = errors.New("invalid: Try 'A1' or 'J10'")
	}
	return
}

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
