package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	Rows     = 6
	Columns  = 6
	Empty    = " "
	Player   = "X"
	Computer = "O"
)

var board [Rows][Columns]string
var currentPlayer string

func main() {
	// Initialize the game board
	inicializarTablero()

	// Seed the random number generator for the computer's moves
	rand.Seed(time.Now().UnixNano())

	// Set the initial player to the player
	currentPlayer = Player

	// Start listening for UDP connections
	serverAddr := "localhost:64511" // Change to the desired port
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening on UDP:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connect 4 Server is listening on UDP port", serverAddr)

	for {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		message := string(buf[:n])
		fmt.Printf("Received message from %s: %s\n", addr, message)

		if message == "1" {
			// Start a new game
			inicializarTablero()
			currentPlayer = Player
			sendGameState(conn)
		} else if message == "2" {
			// End the game
			return
		} else if isValidMove(message) {
			// Handle the player's move
			makeMove(message)
			sendGameState(conn)

			// Check for a win or tie
			if checkWin(currentPlayer) {
				sendMessage(conn, "Player "+currentPlayer+" wins!")
			} else if isBoardFull() {
				sendMessage(conn, "It's a tie!")
			} else {
				// Switch to the computer's turn
				currentPlayer = Computer
				computerMove()
				sendGameState(conn)

				// Check for a win or tie
				if checkWin(currentPlayer) {
					sendMessage(conn, "Player "+currentPlayer+" wins!")
				} else if isBoardFull() {
					sendMessage(conn, "It's a tie!")
				} else {
					// Switch back to the player's turn
					currentPlayer = Player
				}
			}
		}
	}
}

func inicializarTablero() {
	for i := 0; i < Rows; i++ {
		for j := 0; j < Columns; j++ {
			board[i][j] = Empty
		}
	}
}

func sendGameState(conn *net.UDPConn) {
	gameState := "Current player: " + currentPlayer + "\n"
    print(gameState)

	for i := 0; i < Rows; i++ {
		for j := 0; j < Columns; j++ {
			gameState += board[i][j]
		}
		gameState += "\n"
	}

	conn.Write([]byte(gameState))
}

func sendMessage(conn *net.UDPConn, message string) {
	conn.Write([]byte(message))
}

func isValidMove(move string) bool {
	col, err := strconv.Atoi(move)
	if err != nil {
		return false
	}

	return col >= 1 && col <= Columns && board[0][col-1] == Empty
}

func makeMove(move string) {
	col, _ := strconv.Atoi(move)
	col--

	for i := Rows - 1; i >= 0; i-- {
		if board[i][col] == Empty {
			board[i][col] = currentPlayer
			break
		}
	}
}

func computerMove() {
	for {
		col := rand.Intn(Columns)
		if isValidMove(strconv.Itoa(col + 1)) {
			makeMove(strconv.Itoa(col + 1))
			break
		}
	}
}

func checkWin(player string) bool {
	// Check rows for a win
	for i := 0; i < Rows; i++ {
		for j := 0; j < Columns-3; j++ {
			if strings.Join(board[i][j:j+4], "") == player+player+player+player {
				return true
			}
		}
	}

	// Check columns for a win
	for i := 0; i < Rows-3; i++ {
		for j := 0; j < Columns; j++ {
			if board[i][j] == player && board[i+1][j] == player && board[i+2][j] == player && board[i+3][j] == player {
				return true
			}
		}
	}

	// Check diagonals (bottom-left to top-right) for a win
	for i := 3; i < Rows; i++ {
		for j := 0; j < Columns-3; j++ {
			if board[i][j] == player && board[i-1][j+1] == player && board[i-2][j+2] == player && board[i-3][j+3] == player {
				return true
			}
		}
	}

	// Check diagonals (top-left to bottom-right) for a win
	for i := 0; i < Rows-3; i++ {
		for j := 0; j < Columns-3; j++ {
			if board[i][j] == player && board[i+1][j+1] == player && board[i+2][j+2] == player && board[i+3][j+3] == player {
				return true
			}
		}
	}

	return false
}

func isBoardFull() bool {
	for i := 0; i < Rows; i++ {
		for j := 0; j < Columns; j++ {
			if board[i][j] == Empty {
				return false
			}
		}
	}
	return true
}
