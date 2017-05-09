#!/usr/bin/env python3

# func_dict = {func_name: (start_line, end_line, lines)}


# def get_lines(filename='battleship.go'):
#    with open(filename) as in_file:
#        return [line.strip() for line in in_file]
    
def funcs(filename='battleship.go'):
    # pass one
    with open(filename) as in_file:
        lines = []      # all lines in in_file
        func_dict = {}  # key = function name, value = [start_line, end_line]
        curr_func = ''  # will contain the current function with a trailing '('
        for i, line in enumerate(in_file):
            lines.append(line)
            if line.lstrip().startswith('func'):
                if curr_func:
                    if curr_func[0] == '(':
                        exit(i)
                    func_dict[curr_func].append(i - 1)  # record end_line
                curr_func = line.partition(' ')[-1].partition('(')[0] + '('
                if curr_func[0] == '(':
                        exit(i)
                func_dict[curr_func] = [i + 1]  # record start_line
    print(sorted(func_dict))
    assert False, 'Dude.'
    assert func_dict, 'No functions were found on the first pass!'
    func_dict[curr_func].append(i)  # record end_line
    print(func_dict, '\n')
    print()
    # pass two
    print(len(lines), len(func_dict))
    for curr_func, start_end_lines in func_dict.items():
        print(curr_func, start_end_lines)
        func_body = '\n'.join(lines[start_end_lines[0]:start_end_lines[1]])
        func_dict[curr_func] = []  # now record any functions that are called
        for called_func in func_dict:
            if called_func in func_body:
                func_dict[curr_func].append(called_func + ')')
    return {func + ')': tuple(called) for func, called in func_dict.items()}
        

print('=' * 25)
d = funcs()
for func in sorted(d):
    print('{}: {}'.format(func, d[func]))




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