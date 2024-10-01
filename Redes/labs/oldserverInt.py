import socket

def main():
    intermediary = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_address = ('localhost', 12345)
    intermediary.bind(server_address)
    intermediary.listen(1)

    print("Esperando la conexión del cliente...")
    client_conn, client_address = intermediary.accept()

    # Conectarse al servidor Conecta4 (UDP)
    connect4_server = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    connect4_server_address = ('localhost', 8080)  

    while True:
        data = client_conn.recv(1024).decode()
        if not data:
            break

        print(f"Recibido desde el servidor cliente: {data}")

        # Forward move to Connect 4 server (UDP)
        connect4_server.sendto(data.encode(), connect4_server_address)

        # Receive Connect 4 server response
        response, _ = connect4_server.recvfrom(1024)
        print(f"Recibido desde el servidor Conecta4: {response.decode()}")

        # Check for game over conditions and send back to client
        if "Game Over" in response.decode():
            client_conn.send(response)
            break
        else:
            client_conn.send(response.encode())

    intermediary.close()

if __name__ == "__main__":
    main()
