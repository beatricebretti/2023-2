package main

import (
    "fmt"
    "net"
    "math/rand"
    "time"
    "strconv"
)

const (
    Rows    = 6
    Columns = 6
    Empty   = " "
    Player1 = "X"
    Player2 = "O"
)

var board [Rows][Columns]string
var currentPlayer string

func main() {
    rand.Seed(time.Now().UnixNano())
    port := rand.Intn(65536-8000) + 8000
    address := fmt.Sprintf(":%d", port)

    // Open UDP connection for this game
    udpConn, err := net.ListenPacket("udp", address)
    if err != nil {
        fmt.Println("Error abriendo la conexion UDP:", err)
        return
    }
    defer udpConn.Close()

    // Connect to the intermediary server (UDP)
    intermediaryAddress := "localhost:8080"  // Replace with actual Intermediary Server address
    intermediaryConn, err := net.Dial("udp", intermediaryAddress)
    if err != nil {
        fmt.Println("Error en la coneccion con el Servidor Intermediario:", err)
        return
    }
    defer intermediaryConn.Close()

    fmt.Printf("Servidor Conecta4 esta escuchando en puerto UDP %d\n", port)

    // Initialize the game board
    initializeBoard()

    // Initialize the game state
    currentPlayer = Player1
    gameOver := false

    for !gameOver {
        // Send the current game state to the intermediary server
        sendGameState(intermediaryConn, currentPlayer)

        // Receive the player's move from the intermediary server
        move := receiveMove(intermediaryConn)
		// fmt.Println(move)

        // Apply the player's move to the game board
        if makeMove(move) {
            // Check for a win or draw
            if checkWin() {
                sendGameResult(intermediaryConn, currentPlayer+" gana!")
                gameOver = true
            } else if isBoardFull() {
                sendGameResult(intermediaryConn, "Es un empate!")
                gameOver = true
            } else {
                // Switch to the next player
                switchPlayer()
            }
        }
    }
}

func initializeBoard() {
    for i := 0; i < Rows; i++ {
        for j := 0; j < Columns; j++ {
            board[i][j] = Empty
        }
    }
}

// Get the height of a pile (column)
func getPileHeight(column int) int {
    height := 0
    for i := 0; i < Rows; i++ {
        if board[i][column] != Empty {
            height++
        }
    }
    return height
}

// Modify the sendGameState function to include pile information
func sendGameState(conn net.Conn, currentPlayer string) {
    gameState := "Jugador actual: " + currentPlayer + "\n"
	
    // Add pile information to the game state
    for i := 0; i < Rows; i++ {
        for j := 0; j < Columns; j++ {
            gameState += board[i][j]
        }
        gameState += "\n"
    }

    // Add pile information to the game state
    pileInfo := ""
    for j := 0; j < Columns; j++ {
        pileInfo += strconv.Itoa(getPileHeight(j)) + " "
    }
    gameState += pileInfo + "\n"

    conn.Write([]byte(gameState)) // Send game state with pile info
}



func receiveMove(conn net.Conn) int {
    buffer := make([]byte, 1024)
    n, _ := conn.Read(buffer)
    move, _ := strconv.Atoi(string(buffer[:n]))
    return move
}

func makeMove(column int) bool {
    if column < 0 || column >= Columns || board[0][column] != Empty {
        return false // Invalid move
    }

    // Find the first empty row in the selected column and place the player's piece
    for row := Rows - 1; row >= 0; row-- {
        if board[row][column] == Empty {
            board[row][column] = currentPlayer
            return true
        }
    }

    return false
}

func checkWin() bool {
    // Implement Connect 4 win conditions here
    // You need to check horizontally, vertically, and diagonally
    // Return true if there is a win, false otherwise
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

func switchPlayer() {
    if currentPlayer == Player1 {
        currentPlayer = Player2
    } else {
        currentPlayer = Player1
    }
}

func sendGameResult(conn net.Conn, result string) {
    conn.Write([]byte(result))
}



