import socket

conecta4_direccion = ('localhost', 64511)
intermediario_direccion = ('localhost', 8080)

server_intermediario = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_intermediario.bind(intermediario_direccion)
server_intermediario.listen(1)

print("---El servidor intermediario esta esperando una solicitud de partida----")

try:
    while True:
        client_socket, _ = server_intermediario.accept()
        print("----Cliente conectado----")

        server_conecta4 = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        server_conecta4.bind(('localhost', 0))
        connect4_server_port = server_conecta4.getsockname()[1]
        print("----Se estableció conexión UDP con el servidor Conecta4 en el puerto", connect4_server_port, "----")

        while True:
            data = client_socket.recv(1024).decode()
            print("Recibido desde servidor Cliente:", data)
            if data == "1": # Cliente quiere jugar una partida
                server_conecta4.sendto(data.encode(), conecta4_direccion)
                print("----Solicitando partida a servidor Conecta4----")
                # Esperar la respuesta del servidor Conecta4
                response, _ = server_conecta4.recvfrom(1024)
                client_socket.send(response)
            elif data == "2": # Cliente quiere salir del juego
                client_socket.send("Saliendo del juego".encode())
                break
            else:
                server_conecta4.sendto(data.encode(), conecta4_direccion)
                response, _ = server_conecta4.recvfrom(1024)
                client_socket.send(response)
finally:
    client_socket.close()
    server_intermediario.close()
    server_conecta4.close()
