package protocol

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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
	id_len      int
	counter_len int
	domain      string
	domain_len  int
}

func (s Settings) PayloadLength() int {
	return 255 - s.id_len - s.counter_len - s.domain_len
}

var GlobalSettings = Settings{}

func init() {
	GlobalSettings.id_len = 16
	GlobalSettings.counter_len = 4
	GlobalSettings.domain = ".jacobseunglee.com"
	GlobalSettings.domain_len = len(GlobalSettings.domain)
}

func CalculateMaxFileSize() int {
	return (255 - GlobalSettings.PayloadLength()) * pow(16, GlobalSettings.counter_len)
}

func EncodePayload(id string, counter int, payload string) (string, error) {
	count := fmt.Sprintf("%0*x", GlobalSettings.counter_len, counter)

	return id + count + payload + GlobalSettings.domain, nil
}

type packet struct {
	id      string
	counter int
	payload string
}

func DecodePayload(payload string) (packet, error) {
	id := payload[:GlobalSettings.id_len]
	count := payload[GlobalSettings.id_len : GlobalSettings.id_len+GlobalSettings.counter_len]
	enc_payload := payload[GlobalSettings.id_len+GlobalSettings.counter_len:]

	dec_bytes, err := hex.DecodeString(enc_payload)
	if err != nil {
		return packet{}, err
	}

	counter, err := strconv.ParseInt(count, 16, 64)
	if err != nil {
		return packet{}, err
	}

	return packet{
		id:      id,
		counter: int(counter),
		payload: string(dec_bytes),
	}, nil
}

func CreateId() (string, error) {
	bytes := make([]byte, GlobalSettings.id_len/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
