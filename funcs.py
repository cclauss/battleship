#!/usr/bin/env python3


def funcs(filename='battleship.go'):
    with open(filename) as in_file:
        for line in in_file:
            line = line.strip()
            if line.startswith('func'):
                yield line.rstrip(' {')


print('\n'.join(funcs()))
"""
func coordsToPoint(yCommaX string) point
func makeBoard() [][]string
func boardToStr(board [][]string) string
func invalidPoint(pt point) bool
func randomBool() bool { return rand.Intn(2) == 1 }
func randomPoint() point { return point{rand.Intn(10), rand.Intn(10)} }
func htmlSubmitButton(y, x int) string
func templateMap(players []player) map[string]string
func neighbors(pt point) (pts []point)
func pointsForShip(topLeft point, length int, across bool) (pts []point)
func pointToLetterNumber(pt point) (string, error)
func letterNumberToPoint(s string) (pt point, err error)
func borderRow() string
func letterRow() string
func charsInStr(str string) []string
func formatRow(i int) string
func clokeRow(row []string) []string
func boardDisplay(players []player) string
func playableSquares(opponent player) []point
func hasShip(p player, s string) bool
func hasAnyShips(p player) bool
func askUser(prompt string) (string, error)
func askWhichSquare() string
func askIfAcross() bool
func dropABomb(opponent player, sq point) (gameOn bool)
func humanTurn(opponent player) (gameOn bool)
func compuTurn(opponent player) (gameOn bool)
func humanPlaceShip(p player, shipName string)
func compuPlaceShip(p player, shipName string)
func renderTemplate(w http.ResponseWriter, homeTeam, awayTeam player)
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc
func (h helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)
func openBrowser(url string) error
func main()
"""
