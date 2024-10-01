# Cliente
import socket

host_cliente = socket.gethostbyname(socket.gethostname())
puerto_cliente = 8080
msg = ''

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.connect((host_cliente, puerto_cliente))

    print("\n----- Bienvenid@ al juego 'Conecta 4' -----")
    while True:
        print("Seleccione una opcion:\n1- Jugar\n2- Salir")
        opcion = int(input('>> '))
        if opcion == 2:
            msg = '2'
            break
        else:
            msg = '1'
            continue
        s.sendall(opcion.encode())
        data = s.recv(1024)

print('Cerrando la partida ...')