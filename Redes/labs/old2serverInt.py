import socket

def main():
    intermediary_server_address = ('localhost', 8080)
    intermediary_server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    intermediary_server.bind(intermediary_server_address)
    intermediary_server.listen(1)
    connect4_server_address = ('localhost', 64511)

    print("El servidor intermediario esta esperando una conexion..")
    client_socket, _ = intermediary_server.accept()
    print(client_socket)

    while True:
        data = client_socket.recv(1024).decode()
        print(data)
        if not data:
            break

        print("Received message from Client:", data)

        if data == "1":
            connect4_server = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            connect4_server.sendto(data.encode(), connect4_server_address)
            response, _ = connect4_server.recvfrom(1024)
            connect4_server.close()
            client_socket.send(response.encode())
        elif data == "2":
            client_socket.send("Goodbye!".encode())
            break
        else:
            client_socket.send("Invalid choice. Please select 1 to play or 2 to exit.".encode())

    client_socket.close()
    intermediary_server.close()

if __name__ == "__main__":
    main()
