package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/lqgl/grpc-server/pb"
)

const post = ":5001"

func main() {
	listen, err := net.Listen("tcp", post)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// 创建creds证书
	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatalln(err.Error())
	}
	// 通过grpc传递creds证书
	options := []grpc.ServerOption{grpc.Creds(creds)}
	// 创建server
	server := grpc.NewServer(options...)
	// 注册服务
	pb.RegisterEmployeeServiceServer(server, new(employeeService))

	fmt.Println("gRPC Server started ...", post)
	// 开启server监听listen端口号
	server.Serve(listen)
}

type employeeService struct {
	pb.UnimplementedEmployeeServiceServer
}

// GetByNo:通过员工编号找到员工
// 一元消息传递
func (s *employeeService) GetByNo(ctx context.Context, req *pb.GetByNoRequest) (*pb.EmployeeResponse, error) {
	for _, e := range employees {
		if req.Number == e.No {
			return &pb.EmployeeResponse{
				Employee: e,
			}, nil
		}
	}
	return nil, errors.New("employee not found")
}

// 二元消息传递
// 服务端会将数据以streaming的形式传回
func (s *employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {
	for _, e := range employees {
		stream.Send(&pb.EmployeeResponse{
			Employee: e,
		})
		time.Sleep(2 * time.Second)
	}
	return nil
}

// client 以stream的形式传输图片给服务端
func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		fmt.Printf("Employee: %s\n", md["no"])
	}
	var img []byte
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("File Size: %d\n", len(img))
			return stream.SendAndClose(&pb.AddPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}
		fmt.Printf("File received: %d\n", len(data.Data))
		img = append(img, data.Data...)
	}
}

func (s *employeeService) Save(context context.Context, stream *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	return nil, nil
}

// 双向传送stream
func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {

	for {
		empReq, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		employees = append(employees, empReq.Employee)
		stream.Send(&pb.EmployeeResponse{Employee: empReq.Employee})
	}
	for _, e := range employees {
		fmt.Println(e)
	}
	return nil
}
