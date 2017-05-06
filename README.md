# [battleship](https://en.wikipedia.org/wiki/Battleship_(game))

This is my first app written in Go so I am sure that my future self will hate what is here but hey, this is the way that we learn...

TODOs:
- [ ] Web refresh is Mac-only
- [ ] Web refresh opens a __new tab__ after ever turn!
- [ ] Show the user better status messages in both UIs.
- [ ] Compartmentalize things by moving to multiple .go files
- [ ] Computer player needs state to follow up a hit by hitting the neighbors
- [ ] Tests!!!

## Quick start
* __$__ `go get github.com/cclauss/battleship`
* __$__ `cd ${GOPATH}/src/github.com/cclauss/battleship`
* __$__ `chmod +x battleship`
* __$__ `./battleship`

This will start the app running as a webserver on http://localhost:8080

If you are on a Mac, it will even open up a browser tab to allow you to play.  If you are not on a Mac then just click the URL above.

For best experience, place your windows so you can watch both the terminal and the browser at the same time.  ___Click buttons to drop bombs___ but watch out because the computer is _randomly_ dropping bombs on you as well.  You should really be worried if you lose.

![terminalAndBrowser](/images/terminalAndBrowser.png)

## History

I started out with something that looked like this:

```
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
```

Which helped me to understand strings and user input, etc.  After a while of that I started to think that gameplay and UX would be better if the go app was a webserver and the client could just click on buttons... 

![Battleship_web](/images/Battleship_web.png)

Which might work out better on mobile screens.

Next we will need a core battleship API that would drive both the text-based and the web-based ui.

At least it is up and stumbling...
