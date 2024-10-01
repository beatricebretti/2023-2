package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	dir_conecta4 := "localhost:64511"
	udpAddr, _ := net.ResolveUDPAddr("udp", dir_conecta4)
	conn, _ := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	fmt.Println("El servidor Conecta4 se conecto en el puerto", dir_conecta4)

	buf := make([]byte, 1024)
	n, addr, _ := conn.ReadFromUDP(buf)
	mensaje := string(buf[:n])
	fmt.Printf("Recibido desde %s: %s\n", addr, mensaje)

	if mensaje == "1" {
		print("entre here")
		enviarMensaje(conn, "Disponible")
	}
}

func enviarMensaje(conn *net.UDPConn, mensaje string) {
	conn.Write([]byte(mensaje))
}