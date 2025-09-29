package protocol

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sneakygolem/internal/logger"
	"strconv"
)

func pow(a, b int) int {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}

type Settings struct {
	IDLen      int
	CounterLen int
	Domain     string
	DomainLen  int
	MaxCount   int
}

func (s Settings) PayloadLength() int {
	return 63 - s.IDLen - s.CounterLen - s.DomainLen
}

var GlobalSettings = Settings{}

func init() {
	GlobalSettings.IDLen = 16
	GlobalSettings.CounterLen = 4
	GlobalSettings.Domain = ".log.jacobseunglee.com"
	GlobalSettings.DomainLen = len(GlobalSettings.Domain)
	GlobalSettings.MaxCount = pow(16, GlobalSettings.CounterLen)
	logger.Init()
}

func CalculateMaxFileSize() int {
	return (63 - GlobalSettings.PayloadLength()) * GlobalSettings.MaxCount
}

func EncodePayload(id string, counter int, payload string) (string, error) {
	count := fmt.Sprintf("%0*x", GlobalSettings.CounterLen, counter)

	return id + count + payload + GlobalSettings.Domain, nil
}

type Packet struct {
	ID      string
	Counter int
	Payload string
}

// Redo this packet section to use struct methods and process things in object oriented way

func DecodePayload(payload string) (Packet, error) {
	id := payload[:GlobalSettings.IDLen]
	count := payload[GlobalSettings.IDLen : GlobalSettings.IDLen+GlobalSettings.CounterLen]
	enc_payload := payload[GlobalSettings.IDLen+GlobalSettings.CounterLen:]

	counter, err := strconv.ParseInt(count, 16, 64)
	if err != nil {
		return Packet{}, err
	}

	return Packet{
		ID:      id,
		Counter: int(counter),
		Payload: enc_payload,
	}, nil
}

func FinalizePayload(id string) string {
	return id + fmt.Sprintf("%0*x", GlobalSettings.CounterLen, pow(16, GlobalSettings.CounterLen)-1) + GlobalSettings.Domain
}

func CreateId() (string, error) {
	bytes := make([]byte, GlobalSettings.IDLen/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
