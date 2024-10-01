import socket

intermediario_direccion = ('localhost', 8080)
cliente = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
try:
    cliente.connect(intermediario_direccion)
    print("----Bienvenido a Conecta 4 (conectado a servidor intermediario)----")
    while True:
        print("-Seleccione una opcion")
        print("1- Jugar")
        print("2- Salir")
        opcion = input()
        if opcion == '1':
            print("entre al if opcion == 1")
            cliente.sendall(opcion.encode())
            disponibilidad = cliente.recv(1024).decode()
            print("Respuesta de disponibilidad:", disponibilidad)  # Aquí recibes la respuesta del servidor intermediario

            if disponibilidad == "Disponible":
                print("llegue aqui c:")
                while True:
                    print("llegue aqui c: 2")
                    data = cliente.recv(1024).decode()
                    print(data)
                    if data.startswith("Jugador") or data == "Empate":
                        break
                    move = input("Ingresa tu movimiento (1-6): ")
                    cliente.sendall(move.encode())
            else:
                print("El servidor Conecta4 no está disponible en este momento.")
        elif opcion == '2':
            cliente.sendall(opcion.encode())
            break
finally:
    cliente.close()
