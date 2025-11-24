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

	log.Println("[Server] 지뢰찾기 서버를 시작합니다.")

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
		log.Printf("[Close] 서버 종료를 실패했습니다 - %v\n", err)
	}
}

func serveClients(listener net.Listener, matchMaker *matchmaking.Matchmaker) {
	for {
		log.Println("[Server] 클라이언트 연결 대기 중...")

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[Server] 클라이언트 연결 실패 - %v\n", err)
			continue
		}

		log.Printf("[Client %s] 연결됨\n", conn.RemoteAddr())

		go handleClient(network.NewConnection(conn), matchMaker)
	}
}

func handleClient(conn *network.Connection, matchMaker *matchmaking.Matchmaker) {
	defer closeConnection(conn)

	var gameRoom *room.Room

	for {
		message, err := conn.Receive()
		if err != nil {
			log.Printf("[Client %s] 연결 끊김 - %v\n", conn.Conn.RemoteAddr(), err)
			clientDisconnect(gameRoom, conn, matchMaker)
			return
		}
		log.Printf("[Client %s] 메시지 수신 - %s\n", conn.Conn.RemoteAddr(), message.Type)
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
		log.Printf("[Client %s] Exit 요청 수신\n", conn.Conn.RemoteAddr())
		clientDisconnect(*gameRoom, conn, matchMaker)
	}
}

func handleJoin(conn *network.Connection, matchMaker *matchmaking.Matchmaker, gameRoom **room.Room) {
	room, id := matchMaker.FindOrCreateRoom(conn)
	*gameRoom = room

	log.Printf("[Room %s][Join] Player %d 입장\n", room.GetId(), id)

	response := protocol.Message{
		Type:     protocol.Joined,
		PlayerId: id,
		Message:  fmt.Sprintf("플레이어 %d로 입장했습니다.", id),
	}
	log.Println(response.Message)

	err := conn.Send(response)
	if err != nil {
		log.Printf("[Room %s][Join] 메시지 전송 실패 - %v\n", room.GetId(), err)
	}

	if room.IsFull() {
		log.Printf("[Room %s] 두 플레이어 입장 완료 - 게임 시작\n", room.GetId())
		room.StartGame()
	}
}

func closeConnection(conn *network.Connection) {
	err := conn.Close()
	if err != nil {
		log.Printf("[Client %s] 연결 종료 실패 - %v\n", conn.Conn.RemoteAddr(), err)
	}
}

func clientDisconnect(gameRoom *room.Room, conn *network.Connection, matchMaker *matchmaking.Matchmaker) {
	if gameRoom != nil {
		log.Printf("[Room %s][Disconnect] 클라이언트 %s 처리\n", gameRoom.GetId(), conn.Conn.RemoteAddr())

		gameRoom.HandleDisconnect(conn)
		matchMaker.RemoveRoom(gameRoom.GetId())

		log.Printf("[Room %s][Disconnect] 룸 제거 완료\n", gameRoom.GetId())
	}
}
