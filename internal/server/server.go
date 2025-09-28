package server

import (
	"fmt"
	"net"
	"sneakygolem/internal/logger"
	"sneakygolem/internal/protocol"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

var (
	workers = make(map[string]*Worker)
)

func Run() {
	logger.Init()
	logger.Server.Info("Server is running...")
	// Set up DNS server

	udpWriter, err := net.ListenUDP("udp", nil)
	if err != nil {
		logger.Server.Error("Failed to set up UDP writer", "error", err)
		return
	}
	defer udpWriter.Close()

	dns.HandleFunc(".", handleDNSRequest)
	server := &dns.Server{Addr: ":53", Net: "udp"}
	err = server.ListenAndServe()
	if err != nil {
		logger.Server.Error("Failed to start DNS server", "error", err)
	}
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	logger.Server.Info("Received DNS request", "remote_addr", w.RemoteAddr().String())
	if len(r.Question) == 0 {
		logger.Server.Error("DNS request has no questions")
		return
	}
	sendDNSResponse(w, r)
	question := r.Question[0]
	logger.Server.Info("Processing question", "name", question.Name)
	payload := strings.Split(question.Name, ".")[0] // Remove trailing dot and domain length
	if len(payload) < protocol.GlobalSettings.IDLen+protocol.GlobalSettings.CounterLen {
		logger.Server.Error("Payload too short", "payload", payload)
		return
	}
	packet, err := protocol.DecodePayload(payload)
	if err != nil {
		logger.Server.Error("Failed to process payload", "error", err)
		return
	}
	// if packet.Counter == protocol.GlobalSettings.MaxCount-1 {

	// 	worker, exists := workers[packet.ID]
	// 	logger.Server.Info("Received finalization packet", "id", worker.done)
	// 	if exists {
	// 		logger.Server.Info("Finalizing worker", "id", packet.ID)
	// 		worker.queue <-
	// 	}
	// } else {
	worker := getOrCreateWorker(packet.ID)
	worker.queue <- packet
	// }
}

func sendDNSResponse(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Rcode = dns.RcodeSuccess
	m.Answer = make([]dns.RR, 1)
	a, _ := dns.NewRR(fmt.Sprintf("%s A %s", r.Question[0].Name, "1.2.3.4"))
	m.Answer[0] = a
	w.WriteMsg(m)
}

func getOrCreateWorker(id string) *Worker {
	worker, exists := workers[id]
	if !exists {
		worker = &Worker{
			queue: make(chan protocol.Packet, protocol.GlobalSettings.MaxCount),
			file:  fmt.Sprintf("output_%s.txt", id),
			mutex: sync.Mutex{},
			done:  false,
		}
		workers[id] = worker
		go worker.Run()
	}
	return worker
}
