package main

import (
	"github.com/jeremywohl/flatten"
	"github.com/logzio/logzio-go"
	"github.com/nats-io/stan.go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logzioToken := os.Getenv("LOGZIO_TOKEN")
	clusterID := os.Getenv("CLUSTER_ID")
	natsURL := os.Getenv("NATS_URL")
	subject := os.Getenv("SUBJECT")
	queueGroup := os.Getenv("QUEUE_GROUP")
	durableName := queueGroup

	l, err := logzio.New(
		logzioToken,
		logzio.SetDebug(os.Stderr),
		logzio.SetUrl("https://listener.logz.io:8071"),
		logzio.SetDrainDuration(time.Second*10),
		logzio.SetTempDirectory("/tmp/logzio"),
		logzio.SetDrainDiskThreshold(99),
	)

	if err != nil {
		panic(err)
	}

	var clientID string
	hostname, err := os.Hostname()
	if err != nil {
		clientID = hostname
	} else {
		clientID = "client"
	}


	qcb := func (m *stan.Msg) {
		flatMsg, err := flatten.FlattenString(string(m.Data), "", flatten.DotStyle)
		if err != nil {
			println(err.Error())
			return
		}
		if err := l.Send([]byte(flatMsg)); err != nil {
			println(err.Error())
		}
	}

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		panic(err.Error())
	}

	sub, err := sc.QueueSubscribe(subject, queueGroup, qcb, stan.DurableName(durableName))
	if err != nil {
		panic(err.Error())
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	sub.Close()
	sc.Close()
	l.Stop()
}
