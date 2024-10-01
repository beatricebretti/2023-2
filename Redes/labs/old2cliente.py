import socket

def main():
    server_address = ('localhost', 64511)  # Adjust to the Connect 4 Server's address and port
    #server_address = ('localhost', 8080) 
    client = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    print("---------------- Bienvenido a Connecta4 ----------------")

    while True:
        print("- Seleccione una opcion")
        print("1- Jugar")
        print("2- Salir")
        choice = input()

        if choice == '1':
            client.sendto(choice.encode(), server_address)
            while True:
                data, server = client.recvfrom(1024)
                print(data.decode())
                move = input("Enter your move (1-6): ")
                client.sendto(move.encode(), server_address)
        elif choice == '2':
            client.sendto(choice.encode(), server_address)
            break
        else:
            print("Invalid choice. Please select 1 to play or 2 to exit.")

    client.close()

if __name__ == "__main__":
    main()
