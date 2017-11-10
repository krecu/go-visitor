package visitor_rpc_client

import (
	"encoding/json"
	"math/rand"
	"time"

	"fmt"

	"github.com/krecu/go-visitor/model"
	pb "github.com/krecu/go-visitor/protoc/visitor"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Client struct {
	Addrs []string
	conn  []*grpc.ClientConn
}

func New(addrs []string) (proto *Client, err error) {

	proto = &Client{
		Addrs: addrs,
	}

	for _, h := range addrs {
		if conn, err := grpc.Dial(h, grpc.WithInsecure()); err == nil {
			proto.conn = append(proto.conn, conn)
		}
	}

	if len(proto.conn) == 0 {
		err = fmt.Errorf("Not enable connect to visitor hosts ")
	}

	return
}

func (v *Client) Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func (v *Client) RoundRobin() (index int) {
	index = v.Random(0, len(v.conn))
	return
}

// получение на GRPC
func (v *Client) Get(id string) (proto *model.Raw, err error) {

	chanel := pb.NewGreeterClient(v.conn[v.RoundRobin()])

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

	chanel := pb.NewGreeterClient(v.conn[v.RoundRobin()])

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

// создание на GRPC
func (v *Client) Patch(id string, fields map[string]interface{}) (proto *model.Raw, err error) {

	var (
		result    *pb.Reply
		filedJson []byte
	)

	chanel := pb.NewGreeterClient(v.conn[v.RoundRobin()])

	filedJson, err = json.Marshal(fields)
	if err != nil {
		return
	}

	result, err = chanel.Patch(context.Background(), &pb.PatchRequest{Id: id, Fields: string(filedJson)})
	if err != nil {
		return
	} else {
		err = json.Unmarshal([]byte(result.GetBody()), &proto)
	}

	return
}
