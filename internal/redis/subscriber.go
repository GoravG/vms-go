package redis

import (
	"bufio"
	"net"
	"strings"
	"vms_go/internal/utils"
	"vms_go/internal/ws"
)

func SubscribeToRedis(hub *ws.Hub, addr, channel string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		utils.LogErrorf("Failed to connect to redis: %v", err)
	}
	utils.LogInfof("Connected to redis at %s, subscribing to %s", addr, channel)

	cmd := utils.EncodeCommand("SUBSCRIBE", channel)

	_, errr := conn.Write([]byte(cmd))
	if errr != nil {
		utils.LogErrorf("Failed to subscribe: %v", errr)
	}

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			utils.LogErrorf("Redis read error: %v", err)
		}

		if strings.Contains(line, "message") {
			_, _ = reader.ReadString('\n') // $<len of channel>
			_, _ = reader.ReadString('\n') // channel name
			_, _ = reader.ReadString('\n') // $<len of message>
			payload, _ := reader.ReadString('\n')
			payload = strings.TrimSpace(payload)
			utils.LogInfof("Received from redis: %s", payload)
			hub.Broadcast(payload)
		}
	}
}
