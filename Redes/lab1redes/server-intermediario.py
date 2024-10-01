# Servidor intermediario
import socket

puerto_intermediario = 8080
host_intermediario = socket.gethostbyname(socket.gethostname())
print("---El servidor intermediario esta esperando una solicitud de partida----")

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.bind((host_intermediario, puerto_intermediario))
    s.listen()
    conn, addr = s.accept()

    server_conecta4 = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    server_conecta4.bind(('localhost', 0))
    connect4_server_puerto = server_conecta4.getsockname()[1]
    print("----Se estableció conexión UDP con el servidor Conecta4 en el puerto", connect4_server_puerto, "----")

    with conn:
        try:
            print('Se establecio conexion TCP con cliente', addr)
            while True:
                data = conn.recv(1024).decode()
                print("Recibido desde servidor cliente:", data)
                if not data:
                    print("No hay mas data")
                    break
                conn.sendall(data)
        except ValueError:
            print('Se perdio conexion con el cliente')
