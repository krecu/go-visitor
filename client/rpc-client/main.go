package visitor_rpc_client

import (
	"encoding/json"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	pb "github.com/krecu/go-visitor/rpc"
	"log"
	"github.com/krecu/go-visitor/model"
	"math/rand"
	"time"
)

type Client struct {
	Addrs []string
}

func New(addrs []string) (*Client, error) {

	return &Client{
		Addrs: addrs,
	}, nil
}

func (v *Client) Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}

func (v *Client) RoundRobin() (addr string) {
	index := v.Random(0, len(v.Addrs))
	addr = v.Addrs[index]
	return
}

// получение на GRPC
func (v *Client) Get(id string, ip string, ua string, extra map[string]interface{}) (proto *model.Raw, err error) {

	var result *pb.VisitorReply

	conn, err := grpc.Dial(v.RoundRobin(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	chanel := pb.NewGreeterClient(conn)

	if extra != nil {
		extraJson, _ := json.Marshal(extra)
		result, err = chanel.GetVisitor(context.Background(), &pb.VisitorRequest{Ip: ip, Ua: ua, Id: id, Extra: string(extraJson)}); if err != nil {
			return
		}
	} else {
		result, err = chanel.GetVisitor(context.Background(), &pb.VisitorRequest{Ip: ip, Ua: ua, Id: id, Extra: ""}); if err != nil {
			return
		}
	}

	err = json.Unmarshal([]byte(result.GetBody()), &proto)

	return
}

// получение на GRPC
func (v *Client) Put(id string, ip string, ua string, extra map[string]interface{}) (proto *model.Raw, err error) {

	var result *pb.VisitorReply

	conn, err := grpc.Dial(v.RoundRobin(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	chanel := pb.NewGreeterClient(conn)

	if extra != nil {
		extraJson, _ := json.Marshal(extra)
		result, err = chanel.PutVisitor(context.Background(), &pb.VisitorRequest{Ip: ip, Ua: ua, Id: id, Extra: string(extraJson)}); if err != nil {
			return
		}
	} else {
		result, err = chanel.PutVisitor(context.Background(), &pb.VisitorRequest{Ip: ip, Ua: ua, Id: id, Extra: ""}); if err != nil {
			return
		}
	}

	err = json.Unmarshal([]byte(result.GetBody()), &proto)

	return
}