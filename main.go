package main

import (
	"context"
	"log"
	"time"

	"github.com/Victor2333/go-learn/observer"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
	"google.golang.org/grpc"
)

const prefix = "/election-demo"
const prop = "local"

func main() {
	doneChan := make(chan bool)
	flagChan := make(chan bool)
	sessionChan := make(chan *concurrency.Session)
	clientConfig := clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
		DialTimeout: 2 * time.Second,
	}

	cli, err := clientv3.New(clientConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	sub := observer.NewSubject("")
	obs := observer.NewObserver(func(s interface{}) {
		log.Println(s)
	})
	sub.Attach(obs)

	ctx, cancel := context.WithCancel(context.Background())

	go campaign(ctx, cli, prefix, prop, sessionChan, flagChan)

	go func() {
		for {
			select {
			case leaderFlag := <-flagChan:
				switch leaderFlag {
				case true:
					sub.SetState("master")
				case false:
					sub.SetState("slave")
				}
			}
		}
	}()

	<-doneChan
	cancel()
}

func campaign(ctx context.Context, cli *clientv3.Client, prefix string, prop string, sessionChan chan *concurrency.Session, flagC chan bool) {
	//TTL means how long the key will survey
	flagC <- false
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(3))
	if err != nil {
		log.Fatal(err)
	}
	election := concurrency.NewElection(session, prefix)

	//in election.go it will wait until the key is delete
	if err = election.Campaign(ctx, prop); err != nil {
		log.Fatal(err)
	}

	log.Println("Campaign: success")
	flagC <- true
	sessionChan <- session

	select {
	case <-session.Done():
		log.Println("elect: expired")
	}
}
