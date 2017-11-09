package visitor_rpc_client

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/krecu/go-visitor/model"
	pb "github.com/krecu/go-visitor/protoc/visitor"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	return rand.Intn(max-min) + min
}

func (v *Client) RoundRobin() (addr string) {
	index := v.Random(0, len(v.Addrs))
	addr = v.Addrs[index]
	return
}

// получение на GRPC
func (v *Client) Get(id string) (proto *model.Raw, err error) {

	conn, err := grpc.Dial(v.RoundRobin(), grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	chanel := pb.NewGreeterClient(conn)

	if result, err := chanel.Get(context.Background(), &pb.GetRequest{Id: id}); err == nil {
		err = json.Unmarshal([]byte(result.GetBody()), &proto)
	}

	return
}

// создание на GRPC
func (v *Client) Post(id string, ip string, ua string, extra map[string]interface{}) (proto *model.Raw, err error) {

	var (
		result    *pb.Reply
		extraJson []byte
	)
	conn, err := grpc.Dial(v.RoundRobin(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	chanel := pb.NewGreeterClient(conn)

	if extra != nil {
		extraJson, err = json.Marshal(extra)
		if err != nil {
			return
		}
	}

	result, err = chanel.Post(context.Background(), &pb.PostRequest{Ip: ip, Ua: ua, Id: id, Extra: string(extraJson)})
	if err != nil {
		return
	} else {
		err = json.Unmarshal([]byte(result.GetBody()), &proto)
	}

	return
}
