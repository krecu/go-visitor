package visitor_rpc_client

import (
	"encoding/json"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	pb "github.com/krecu/go-visitor/client/rpc-client/rpc"
	"log"
	"github.com/krecu/go-visitor/model"
)

type Client struct {
	Addr string
}

type Body struct {
	Ip string
	Ua string
	Id string
}

func New(addr string) (*Client, error) {

	return &Client{
		Addr: addr,
	}, nil
}

// получение на GRPC
func (v *Client) Rpc(id string, ip string, ua string) (proto *model.Raw, err error) {

	conn, err := grpc.Dial(v.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	chanel := pb.NewGreeterClient(conn)

	result, err := chanel.GetVisitor(context.Background(), &pb.VisitorRequest{Ip: ip, Ua: ua, Id: id}); if err != nil {
		return
	}

	err = json.Unmarshal([]byte(result.GetBody()), &proto)

	return
}

func main()  {
	
}