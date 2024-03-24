package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"net"
)

import proto "auth/pkg/user_v1"

const grpcPort = 50051

type Server struct {
	proto.UnimplementedUserV1Server
}

var users []*proto.User

func findUserById(id int64) *proto.User {
	for _, user := range users {
		if user.GetId() == id {
			return user
		}
	}
	return nil

}

func removeUserFromArray(id int64) {

	for i, user := range users {
		if user.GetId() == id {
			if (i + 1) == len(users) {
				users = users[:i]
			} else {
				users = append(users[:i], users[i+1:]...)
			}
		}
	}
}

func (s *Server) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	return &proto.GetResponse{
		User: findUserById(req.GetId()),
	}, nil

	//return &proto.GetResponse{
	//	User: &proto.User{
	//		Id: req.GetId(),
	//		Info: &proto.UserInfo{
	//			Name:            gofakeit.Name(),
	//			Email:           gofakeit.Email(),
	//			Password:        gofakeit.Password(true, false, false, false, false, 32),
	//			Role:            proto.UserInfo_Role(gofakeit.Number(0, 1)),
	//			PasswordConfirm: gofakeit.Password(true, false, false, false, false, 32),
	//		},
	//		CreatedAt: timestamppb.New(gofakeit.Date()),
	//		UpdatedAt: timestamppb.New(gofakeit.Date()),
	//	},
	//}, nil
}

func (s *Server) Create(ctx context.Context, req *proto.CreateRequest) (*proto.CreateResponse, error) {
	log.Printf("New user name: %s", req.GetInfo().GetName())
	id := rand.Int63()
	now := timestamppb.Now()

	users = append(users, &proto.User{
		Id:        id,
		Info:      req.GetInfo(),
		CreatedAt: now,
		UpdatedAt: now,
	})

	log.Printf("Users: %v", len(users))
	return &proto.CreateResponse{
		Id: id,
	}, nil
}

func (s *Server) Update(ctx context.Context, req *proto.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("edit user id: %s", req.GetInfo().GetId())

	user := findUserById(req.GetInfo().GetId())
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	user.Id = req.GetInfo().GetId()
	user.Info.Email = req.GetInfo().GetEmail().Value
	user.Info.Name = req.GetInfo().GetName().Value

	return &emptypb.Empty{}, nil
}

func (s *Server) Delete(ctx context.Context, req *proto.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("delete user id: %s", req.GetId())

	removeUserFromArray(req.GetId())

	return &emptypb.Empty{}, nil

}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	proto.RegisterUserV1Server(s, &Server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
