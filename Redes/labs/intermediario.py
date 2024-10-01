import socket

conecta4_direccion = ('localhost', 64511) #cambiar esto later
dir_intermediario = ("localhost", 8080) #cambiar esto later
server_intermediario = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_intermediario.bind(dir_intermediario)
server_intermediario.listen(1)

print("El servidor intermediario esta esperando una solicitud de partida")
client_socket, _ = server_intermediario.accept()
print("Cliente conectado")
server_conecta4 = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
server_conecta4.bind(('localhost', 0))
connect4_server_port = server_conecta4.getsockname()[1]
print("Se estableció conexión UDP con el servidor Conecta4 en el puerto", connect4_server_port)
server_conecta4.settimeout(10)

data = client_socket.recv(1024).decode()
print("Recibido desde servidor Cliente:", data)
if data == "1":
    server_conecta4.sendto(data.encode(), conecta4_direccion)
    print("Solicitando partida a servidor Conecta4")

    try:
        response, _ = server_conecta4.recvfrom(1024)
        response = response.decode()
        print("Received response from server Conecta4:", response)
        client_socket.send(response.encode())
    except socket.timeout:
        print("Timeout: No response received from server Conecta4")
    except Exception as e:
        print("An error occurred:", str(e))
