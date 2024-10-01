import socket

dir_intermediario = ("localhost", 8080) #cambiar esto later
server_cliente = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

server_cliente.connect(dir_intermediario)
print("Bienvenido a Conecta4 (conectado a servidor intermediario)")
print("-Seleccione una opcion")
print("1- Jugar")
print("2- Salir")
opcion = input()

if opcion == "1":
    cliente.sendall(opcion.encode())
    disponibilidad = cliente.recv(1024).decode()

    if disponibilidad == "NO":
        print("No hay partidas disponibles")
    elif disponibilidad == "SI":
        print("Si puedes jugar :D")
