package server

import (
	"context"
	"fmt"
	"github.com/mmathys/acfts/common"
	pb "github.com/mmathys/acfts/proto"
	"github.com/mmathys/acfts/util"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type server struct {
	pb.UnimplementedACFTSServer
}

var id *common.Identity
var debug bool
var benchmarkMode bool
var SignedUTXO *sync.Map
var TxCounter *int32

func (s *server) Sign(ctx context.Context, in *pb.SignRequest) (*pb.SignReply, error) {
	log.Printf("got sign request: %v", in)

	if benchmarkMode {
		defer util.CountTx(TxCounter)
	}

	return nil, nil
}

func InitGRPC(port int, _id *common.Identity, _debug bool, _benchmark bool, _SignedUTXO *sync.Map, _TxCounter *int32) {
	id = _id
	debug = _debug
	benchmarkMode  = _benchmark
	SignedUTXO = _SignedUTXO
	TxCounter = _TxCounter

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterACFTSServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
