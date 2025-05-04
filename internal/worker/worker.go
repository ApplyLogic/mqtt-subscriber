package worker

import (
	"github.com/ApplyLogic/mqtt-subscriber/internal/metrics"
	"github.com/ApplyLogic/mqtt-subscriber/internal/mqtt"
	"github.com/ApplyLogic/mqtt-subscriber/internal/parser"
	"log"
	"net"
	"time"
)

type Worker struct {
	Host       string
	Port       string
	parser     *parser.Parser
	Subscriber *mqtt.Subscriber
}

func NewWorker(host, port string, parser *parser.Parser, subscriber *mqtt.Subscriber) *Worker {
	return &Worker{
		Host:       host,
		Port:       port,
		parser:     parser,
		Subscriber: subscriber,
	}
}

func (g *Worker) Start() error {

}
