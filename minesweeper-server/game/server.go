package game

import (
	"log"
	"net"
)

func StartGameServer() {
	listener := listen("127.0.0.1:8080")
	defer closeListener(listener)

	log.Println("지뢰찾기 서버를 시작합니다.")

	serveClients(listener)
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

func serveClients(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("클라이언트 연결에 실패했습니다.", err)
			continue
		}

	}
}
