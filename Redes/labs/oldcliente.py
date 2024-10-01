import socket

def main():
    client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_address = ('localhost', 12345)  # Intermediary Server address
    client.connect(server_address)

    try:
        while True:
            play_option = input("Enter 1 to play or 2 to exit: ")

            if play_option == '1':
                client.send("1".encode())  # Signal the desire to play
                break
            elif play_option == '2':
                client.send("2".encode())  # Signal the desire to exit
                print("Exiting the game.")
                return
            else:
                print("Invalid option. Please enter 1 to play or 2 to exit.")

        # Choose who plays first (1 for player, 2 for computer)
        first_player_choice = input("Choose who goes first (1 for player, 2 for computer): ")
        client.send(first_player_choice.encode())

        while True:
            message = input("Your move (1-6): ")
            client.send(message.encode())
            data = client.recv(1024).decode()
            print(data)

            # Extract and display the pile information
            pile_info = client.recv(1024).decode()
            print("Pile Information:", pile_info)

            if "Game Over" in data:
                break

            # Receive and display the game state for the next turn
            data = client.recv(1024).decode()
            print(data)

        new_game = input("Play again? (1 for Yes, 2 for No): ")
        if new_game == '1':
            client.send("1".encode())
        else:
            client.send("2".encode())

    finally:
        client.close()

if __name__ == "__main__":
    main()
