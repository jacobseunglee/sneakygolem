package server

import (
	"fmt"
	"sneakygolem/internal/logger"
	"sneakygolem/internal/protocol"
	"strings"
	"sync"

	"github.com/mr-tron/base58"
)

type Worker struct {
	queue    chan protocol.Packet
	file     string
	received map[int]protocol.Packet
	done     bool
	mutex    sync.Mutex
}

func checkAllRecieved(received map[int]protocol.Packet) bool {
	for count := range len(received) {
		_, ok := received[count]
		if !ok {
			logger.Server.Error("Missing packet", count)
			return false
		}
	}
	return true
}

func (w *Worker) Run() {
	w.received = make(map[int]protocol.Packet)
	for len(w.queue) > 0 && !w.done {
		packet := <-w.queue
		if packet.Counter == protocol.GlobalSettings.MaxCount-1 {
			logger.Server.Info("Worker", "file", w.file, "log", "received finalization packet")
			w.mutex.Lock()
			w.done = true
			w.mutex.Unlock()
		}
		logger.Server.Info("Worker", "file", w.file, "received packet", packet.Counter)
		w.mutex.Lock()
		w.received[packet.Counter] = packet
		w.mutex.Unlock()
	}
	logger.Server.Info("Worker", "file", w.file, "log", "finalizing")
	w.mutex.Lock()
	if checkAllRecieved(w.received) {
		err := w.processAll()
		if err != nil {
			logger.Server.Error("Worker", "file", w.file, "failed to process all packets", err)
		}
	} else {
		logger.Server.Error("Worker", "file", w.file, "error", "not all packets received")
	}
	w.mutex.Unlock()
}

func (w *Worker) processAll() error {
	payload_list := []string{}
	for count := range len(w.received) {
		logger.Server.Info("Worker", "file", w.file, "processing packet", count)
		packet, ok := w.received[count]
		if !ok {
			logger.Server.Error("Worker", "file", w.file, "missing packet", count)
			return fmt.Errorf("missing packet %d", count)
		}
		payload_list = append(payload_list, packet.Payload)
	}
	payload := strings.Join(payload_list, "")
	dec_payload, err := processPayload(payload)
	if err != nil {
		logger.Server.Error("Worker", "file", w.file, "failed to process payload", err)
		return fmt.Errorf("failed to process payload: %w", err)
	}
	err = protocol.AppendBytesToFile(dec_payload, w.file)
	if err != nil {
		logger.Server.Error("Worker", "file", w.file, "failed to write", err)
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func processPayload(payload string) ([]byte, error) {
	// Simulate payload processing
	logger.Server.Info("Processing payload", "payload", payload)
	dec_payload, err := base58.Decode(payload)
	if err != nil {
		logger.Server.Error("Failed to decode payload", "error", err)
		return nil, err
	}
	return dec_payload, nil
}
