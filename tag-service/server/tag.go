package server

import (
	"context"
	"encoding/json"
	"go-travel/tag-service/pkg/bapi"
	"go-travel/tag-service/pkg/errcode"
	pb "go-travel/tag-service/proto"
)

type TagServer struct {
	pb.UnimplementedTagServiceServer
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, in *pb.GetTagListRequest) (*pb.GetTagListReply, error) {

	//clientConn,err := GetClientConn(ctx,"localhost:8004",[]grpc.DialOption{
	//	grpc.WithUnaryInterceptor(
	//		grpc_middleware.ChainUnaryClient(
	//			middleware.UnaryContextTimeout(),
	//			middleware.ClientTracing(),
	//			)),
	//})
	//if err != nil {
	//	return nil,errcode.TogRPCError(errcode.Fail)
	//}
	//defer clientConn.Close()
	//tagServiceClient := pb.NewTagServiceClient(clientConn)
	//resp,err := tagServiceClient.GetTagList(ctx,&pb.GetTagListRequest{Name: "GO"})
	//if err != nil {
	//	return nil,errcode.TogRPCError(errcode.Fail)
	//}


	api := bapi.NewAPI("http://127.0.0.1:8000")
	body, err := api.GetTagList(ctx, in.GetName())
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ErrorGetTagListFail)
	}

	tagList := pb.GetTagListReply{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, err
	}
	return &tagList, nil
}
