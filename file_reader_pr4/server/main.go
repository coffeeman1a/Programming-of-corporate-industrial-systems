package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const (
	port       = ":9000"
	codeLen    = 4
	minPlayers = 2
	maxPlayers = 2
)

type Player struct {
	conn   net.Conn
	name   string
	reader *bufio.Reader
	writer *bufio.Writer
}

type PlayerResult struct {
	XMLName  xml.Name `xml:"player"`
	Name     string   `xml:"name,attr"`
	Attempts int      `xml:"attempts"`
}

type RoundResult struct {
	XMLName   xml.Name       `xml:"round"`
	StartTime string         `xml:"start"`
	EndTime   string         `xml:"end"`
	Code      string         `xml:"code"`
	Players   []PlayerResult `xml:"players>player"`
	Winner    string         `xml:"winner"`
}

func main() {
	// start the server
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Listen: %v", err)
	}
	defer ln.Close()
	log.Println("Server is running on: ", port)

	players := waitForPlayers(ln)

	for round := 1; ; round++ {
		start := time.Now()
		secret := genCode(codeLen)
		log.Println("Secret code: ", secret)

		broadcast(players, fmt.Sprintf("\n=== New round %d! ===\n", round))
		broadcast(players, fmt.Sprintf("Game is starting! Code length: %d.\n", codeLen))

		attempts, winner, disconnected := playRoundWithStats(players, secret)
		if disconnected {
			break
		}
		end := time.Now()
		result := RoundResult{
			StartTime: start.Format(time.RFC3339),
			EndTime:   end.Format(time.RFC3339),
			Code:      secret,
			Winner:    winner,
		}
		for _, p := range players {
			result.Players = append(result.Players, PlayerResult{Name: p.name, Attempts: attempts[p.name]})
		}
		if err := saveRoundResult(result, round); err != nil {
			log.Printf("Error saving XML: %v", err)
		}

	}

}

func waitForPlayers(ln net.Listener) []*Player {

	players := make([]*Player, 0, maxPlayers)
	for len(players) < maxPlayers {
		conn, _ := ln.Accept()
		p := &Player{
			conn:   conn,
			name:   conn.RemoteAddr().String(),
			reader: bufio.NewReader(conn),
			writer: bufio.NewWriter(conn),
		}
		players = append(players, p)
		p.writer.WriteString("Welcome, player " + p.name + "\n")
		p.writer.Flush()
		log.Println("New connection: ", p.name)
	}

	return players
}

func playRoundWithStats(players []*Player, secret string) (map[string]int, string, bool) {
	attempts := make(map[string]int)
	current := 0
	for {
		p := players[current]
		p.writer.WriteString("Your turn: GUESS:xxxx\n")
		p.writer.Flush()

		line, err := p.reader.ReadString('\n')
		if err != nil {
			log.Printf("Read error from %s: %v", p.name, err)
			return attempts, "", true
		}
		guess := strings.TrimSpace(line)
		if !strings.HasPrefix(guess, "GUESS:") {
			p.writer.WriteString("Invalid format, try GUESS:1234\n")
			p.writer.Flush()
			continue
		}
		code := strings.TrimPrefix(guess, "GUESS:")
		b, w := evaluate(secret, code)
		res := fmt.Sprintf("%dB%dW\n", b, w)
		attempts[p.name]++

		broadcast(players,
			fmt.Sprintf("%s guessed %s → %s", p.name, code, res),
		)

		if b == codeLen {
			broadcast(players, "WIN! Winner: "+p.name+"\n")
			return attempts, p.name, false
		}
		current = (current + 1) % len(players)
	}
}

func broadcast(players []*Player, msg string) {
	for _, p := range players {
		p.writer.WriteString(msg)
		p.writer.Flush()
	}
}

func genCode(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = byte('0' + rand.Intn(10))
	}
	return string(s)
}

func evaluate(secret, guess string) (b, w int) {
	usedS := make([]bool, len(secret))
	usedG := make([]bool, len(guess))
	// black
	for i := range secret {
		if guess[i] == secret[i] {
			b++
			usedS[i], usedG[i] = true, true
		}
	}
	// white
	for i := range secret {
		if usedS[i] {
			continue
		}
		for j := range guess {
			if usedG[j] || secret[i] != guess[j] {
				continue
			}
			w++
			usedS[i], usedG[j] = true, true
			break
		}
	}
	return
}

func saveRoundResult(result RoundResult, round int) error {
	filename := fmt.Sprintf("round_%d.xml", round)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")
	if err := enc.Encode(result); err != nil {
		return err
	}
	return nil
}
