import random
import socket

flag=True
while flag:
    print(" --- NUEVO INTENTO -- \n")

    # Definir host y port
    host = 'jdiaz.inf.santiago.usm.cl'
    puertos = [50006, 50007, 50008, 50009, 50010]
    # Elegir un puerto aleatorio para realizar la conexión
    port = random.choice(puertos)  

    # Mensaje a enviar para tener los datos de la imagen
    message = 'GET NEW IMG DATA'
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    client_socket.sendto(message.encode(), (host, port))

    response, server_address = client_socket.recvfrom(1024)
    response_str = response.decode()
    print('MENSAJE RECIBIDO:    ' + response_str)
    response_parts = response_str.split()

    image_id = response_parts[0].split(':')[1]
    width = int(response_parts[1].split(':')[1])
    height = int(response_parts[2].split(':')[1])
    tcp_port = int(response_parts[3].split(':')[1])
    udp_port_1 = int(response_parts[4].split(':')[1])
    udp_port_2 = None
    udp_port_3 = None
    if len(response_parts) > 5:
        udp_port_2 = int(response_parts[5].split(':')[1])
        Y = 2
    if len(response_parts) > 6:
        udp_port_3 = int(response_parts[6].split(':')[1])
        Y = 3
    pv_port = int(response_parts[-1].split(':')[1])

    # Calcular tamanio del buffer
    buffer_size = width * height * 3

    # Crear un socket TCP para la primera parte de la imagen
    tcp_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    tcp_socket.connect((host, tcp_port))
    print("Conectando al primer server TCP: "+str(tcp_port))
    tcp_socket.sendall('GET 1/{} IMG ID:{}'.format(Y, image_id).encode())
    image_data_1 = tcp_socket.recv(buffer_size)
    print("Primera parte de la imagen obtenida! ")
    tcp_socket.close()

    # Crear un socket UDP para la segunda parte de la imagen
    udp_socket_1 = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    print("Conectando al segundo server UDP: "+str(udp_port_1))
    udp_socket_1.sendto('GET 2/{} IMG ID:{}'.format(Y, image_id).encode(), (host, udp_port_1))
    image_data_2, _ = udp_socket_1.recvfrom(buffer_size)
    print("Segunda parte recibida! ")
    udp_socket_1.close()

    # Crear un socket UDP para la tercera parte de la imagen si es que esta
    if udp_port_3 is not None:
        udp_socket_2 = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        print("Conectando al tercer server UDP: "+str(udp_port_2))
        udp_socket_2.sendto('GET 3/3 IMG ID:{}'.format(image_id).encode(), (host, udp_port_2))
        image_data_3, _ = udp_socket_2.recvfrom(buffer_size)
        print("Tercera parte recibida! ")
        udp_socket_2.close()

    print("Checking Bytes...    ")
    # Combinar la data de la imagen en una unica variable
    if udp_port_3 is not None:
        image_data = image_data_1 + image_data_2 + image_data_3
    else:
        image_data = image_data_1 + image_data_2

    # Crear un socket TCP para el puerto PV
    pv_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    pv_socket.connect((host, pv_port))
    print("Conectando al puerto PV: "+str(pv_port))

    # Enviar la data de la imagen al puerto PV
    pv_socket.sendall(image_data)
    response = pv_socket.recv(1024)
    response_str = response.decode()

    # Checkear si la foto es correcta 
    if response_str=="200: SUCCESS":
        #print(response_str) 
        print("Todo está correcto! Escribiendo imagen...    ")

        # Reconstruir la imagen
        with open(f'{image_id}.png', 'wb') as f:
            f.write(image_data)

        print("Imagen escrita! \n")
        print("-- FIN --")
        break
    else:
        print("Posible error!, intentando nuevamente \n")
        