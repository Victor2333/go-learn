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
	var cancel func()
	var ctx context.Context
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
	obs := observer.NewObserver(func(oldState interface{}, newState interface{}) {
		log.Println("I am", newState)
		time.Sleep(5 * time.Second)
		if newState == "master" && oldState == "slave" {
			cancel()
		}
		if newState == "slave" {
			ctx, cancel = context.WithCancel(context.Background())
			go campaign(ctx, cli, prefix, prop, flagChan)
		}
	})
	sub.Attach(obs)

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

	flagChan <- false
	<-doneChan
}

func campaign(ctx context.Context, cli *clientv3.Client, prefix string, prop string, flagC chan bool) {
	//TTL means how long the key will survey
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

	for {
		select {
		case <-session.Done():
			log.Println("elect: expired")
			flagC <- false
			return
		case <-ctx.Done():
			session.Orphan()
		}
	}
}
