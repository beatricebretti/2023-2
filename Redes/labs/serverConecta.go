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
	Filas     = 6
	Columnas  = 6
	Empty    = " "
	Jugador   = "X"
	Bot = "O"
)

var tablero [Filas][Columnas]string
var JugadorActual string

func main() {
	inicializarTablero()
	rand.Seed(time.Now().UnixNano())
	JugadorActual = Jugador

	// Abrir conexiones UDP
	intermediario_direccion := "localhost:64511"
	udpAddr, _ := net.ResolveUDPAddr("udp", intermediario_direccion)
	conn, _ := net.ListenUDP("udp", udpAddr)
	defer conn.Close()

	fmt.Println("El servidor Conecta4 se conecto al servidor intermediario en el puerto", intermediario_direccion)

	for {
		buf := make([]byte, 1024)
		n, addr, _ := conn.ReadFromUDP(buf)
		mensaje := string(buf[:n])
		fmt.Printf("Recibido desde %s: %s\n", addr, mensaje)

		if mensaje == "1" {
			inicializarTablero()
			JugadorActual = Jugador
			estadoPartida(conn)
			enviarMensaje(conn, "Disponible")
		} else if mensaje == "2" {
			return
		} else if movimientoValido(mensaje) {
			llenarPila(mensaje)
			estadoPartida(conn)

			// Check for a win or tie
			if verificarGanador(JugadorActual) {
				enviarMensaje(conn, "Jugador "+JugadorActual+" gana!")
			} else if tableroLleno() {
				enviarMensaje(conn, "Empate")
			} else {
				// Switch to the Bot's turn
				JugadorActual = Bot
				BotMove()
				estadoPartida(conn)

				// Check for a win or tie
				if verificarGanador(JugadorActual) {
					enviarMensaje(conn, "Jugador "+JugadorActual+" gana!")
				} else if tableroLleno() {
					enviarMensaje(conn, "Empate")
				} else {
					// Switch back to the jugador's turn
					JugadorActual = Jugador
				}
			}
		}
	}
}

func inicializarTablero() {
	for i := 0; i < Filas; i++ {
		for j := 0; j < Columnas; j++ {
			tablero[i][j] = Empty
		}
	}
}

func estadoPartida(conn *net.UDPConn) {
	estadoJuego := "Current jugador: " + JugadorActual + "\n"
    print(estadoJuego)

	for i := 0; i < Filas; i++ {
		for j := 0; j < Columnas; j++ {
			estadoJuego += tablero[i][j]
		}
		estadoJuego += "\n"
	}

	conn.Write([]byte(estadoJuego))
}

func enviarMensaje(conn *net.UDPConn, mensaje string) {
	conn.Write([]byte(mensaje))
}

func movimientoValido(move string) bool {
	col, _ := strconv.Atoi(move)

	return col >= 1 && col <= Columnas && tablero[0][col-1] == Empty
}

func llenarPila(move string) {
	col, _ := strconv.Atoi(move)
	col--

	for i := Filas - 1; i >= 0; i-- {
		if tablero[i][col] == Empty {
			tablero[i][col] = JugadorActual
			break
		}
	}
}

func BotMove() {
	for {
		col := rand.Intn(Columnas)
		if movimientoValido(strconv.Itoa(col + 1)) {
			llenarPila(strconv.Itoa(col + 1))
			break
		}
	}
}

func verificarGanador(jugador string) bool {
	// Check rows for a win
	for i := 0; i < Filas; i++ {
		for j := 0; j < Columnas-3; j++ {
			if strings.Join(tablero[i][j:j+4], "") == jugador+jugador+jugador+jugador {
				return true
			}
		}
	}

	// Check columns for a win
	for i := 0; i < Filas-3; i++ {
		for j := 0; j < Columnas; j++ {
			if tablero[i][j] == jugador && tablero[i+1][j] == jugador && tablero[i+2][j] == jugador && tablero[i+3][j] == jugador {
				return true
			}
		}
	}

	// Check diagonals (bottom-left to top-right) for a win
	for i := 3; i < Filas; i++ {
		for j := 0; j < Columnas-3; j++ {
			if tablero[i][j] == jugador && tablero[i-1][j+1] == jugador && tablero[i-2][j+2] == jugador && tablero[i-3][j+3] == jugador {
				return true
			}
		}
	}

	// Check diagonals (top-left to bottom-right) for a win
	for i := 0; i < Filas-3; i++ {
		for j := 0; j < Columnas-3; j++ {
			if tablero[i][j] == jugador && tablero[i+1][j+1] == jugador && tablero[i+2][j+2] == jugador && tablero[i+3][j+3] == jugador {
				return true
			}
		}
	}

	return false
}

func tableroLleno() bool {
	for i := 0; i < Filas; i++ {
		for j := 0; j < Columnas; j++ {
			if tablero[i][j] == Empty {
				return false
			}
		}
	}
	return true
}
