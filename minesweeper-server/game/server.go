package game

import (
	"fmt"
	"log"
	"minesweeper-infrastructure/network"
	"minesweeper-infrastructure/protocol"
	"minesweeper-server/game/matchmaking"
	"minesweeper-server/game/room"
	"net"
)

func StartGameServer() {
	listener := listen("127.0.0.1:8080")
	defer closeListener(listener)

	log.Println("지뢰찾기 서버를 시작합니다.")

	matchMaker := matchmaking.NewMatchmaker()
	serveClients(listener, matchMaker)
}

func listen(address string) net.Listener {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("서버 시작을 실패했습니다.", err)
	}
	return listener
}

func closeListener(listener net.Listener) {
	err := listener.Close()
	if err != nil {
		log.Println("서버 종료를 실패했습니다.", err)
	}
}

func serveClients(listener net.Listener, matchMaker *matchmaking.Matchmaker) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("클라이언트 연결에 실패했습니다.", err)
			continue
		}

		go handleClient(network.NewConnection(conn), matchMaker)
	}
}

func handleClient(conn *network.Connection, matchMaker *matchmaking.Matchmaker) {
	defer closeConnection(conn)

	var gameRoom *room.Room

	for {
		message, err := conn.Receive()
		if err != nil {
			clientDisconnect(gameRoom, conn, matchMaker)
			return
		}

		handleMessage(conn, matchMaker, message, &gameRoom)
	}
}

func handleMessage(
	conn *network.Connection,
	matchMaker *matchmaking.Matchmaker,
	message protocol.Message,
	gameRoom **room.Room,
) {
	switch message.Type {
	case protocol.Join:
		handleJoin(conn, matchMaker, gameRoom)

	case protocol.Open:
		if *gameRoom != nil {
			(*gameRoom).HandleOpen(conn, message.Row, message.Col)
		}

	case protocol.Flag:
		if *gameRoom != nil {
			(*gameRoom).HandleFlag(conn, message.Row, message.Col)
		}

	case protocol.Exit:
		clientDisconnect(*gameRoom, conn, matchMaker)
	}
}

func handleJoin(conn *network.Connection, matchMaker *matchmaking.Matchmaker, gameRoom **room.Room) {
	room, id := matchMaker.FindOrCreateRoom(conn)
	*gameRoom = room

	response := protocol.Message{
		Type:     protocol.Joined,
		PlayerId: id,
		Message:  fmt.Sprintf("플레이어 %d로 입장했습니다.", id),
	}
	log.Println(response.Message)

	err := conn.Send(response)
	if err != nil {
		fmt.Println("메시지 전송을 실패했습니다.", err)
	}

	if room.IsFull() {
		room.StartGame()
	}
}

func closeConnection(conn *network.Connection) {
	err := conn.Close()
	if err != nil {
		log.Println("클라이언트 연결 종료를 실패했습니다.", err)
	}
}

func clientDisconnect(gameRoom *room.Room, conn *network.Connection, matchMaker *matchmaking.Matchmaker) {
	if gameRoom != nil {
		gameRoom.HandleDisconnect(conn)
		matchMaker.RemoveRoom(gameRoom.GetId())
	}
}
