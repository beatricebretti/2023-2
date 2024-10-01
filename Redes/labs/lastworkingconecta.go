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
var mensaje2 string

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
			//estadoPartida(conn)
			mensaje2 = "Disponible"
			conn.Write([]byte(mensaje2))
		} else if mensaje == "2" {
			return
		} else if movimientoValido(mensaje) {
			llenarPila(mensaje)
			estadoPartida(conn)

			// Ver si gano o empate
			if verificarGanador(JugadorActual) {
				mensaje = "Jugador "+JugadorActual+" gana!"
				conn.Write([]byte(mensaje))
			} else if tableroLleno() {
				mensaje =  "Empate"
				conn.Write([]byte(mensaje))
			} else {
				// le toca a pc
				JugadorActual = Bot
				BotMove()
				estadoPartida(conn)

				// ver si gano o empate
				if verificarGanador(JugadorActual) {
					mensaje =  "Jugador "+JugadorActual+" gana!"
					conn.Write([]byte(mensaje))
				} else if tableroLleno() {
					mensaje =  "Empate"
					conn.Write([]byte(mensaje))
				} else {
					// volver a jugadorr
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
	estadoJuego := "Jugador actual: " + JugadorActual + "\n"
    print(estadoJuego)

	for i := 0; i < Filas; i++ {
		for j := 0; j < Columnas; j++ {
			estadoJuego += tablero[i][j]
		}
		estadoJuego += "\n"
	}

	conn.Write([]byte(estadoJuego))
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
	// ver filas
	for i := 0; i < Filas; i++ {
		for j := 0; j < Columnas-3; j++ {
			if strings.Join(tablero[i][j:j+4], "") == jugador+jugador+jugador+jugador {
				return true
			}
		}
	}

	// ver columnas
	for i := 0; i < Filas-3; i++ {
		for j := 0; j < Columnas; j++ {
			if tablero[i][j] == jugador && tablero[i+1][j] == jugador && tablero[i+2][j] == jugador && tablero[i+3][j] == jugador {
				return true
			}
		}
	}

	// ver diagonal abajo izq arriba der
	for i := 3; i < Filas; i++ {
		for j := 0; j < Columnas-3; j++ {
			if tablero[i][j] == jugador && tablero[i-1][j+1] == jugador && tablero[i-2][j+2] == jugador && tablero[i-3][j+3] == jugador {
				return true
			}
		}
	}

	// ver diagonal abjo der arriba izq
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
