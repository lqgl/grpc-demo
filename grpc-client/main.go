package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/lqgl/grpc-client/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const address = "localhost:5001"

func main() {
	creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		log.Fatalln(err.Error())
	}
	options := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial(address, options...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close()
	client := pb.NewEmployeeServiceClient(conn)
	// getByNo(client)
	// getAll(client)
	// addPhoto(client)
	saveAll(client)
}

func getByNo(client pb.EmployeeServiceClient) {
	res, err := client.GetByNo(context.Background(), &pb.GetByNoRequest{Number: 1996})
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(res.Employee)
}

func getAll(client pb.EmployeeServiceClient) {
	stream, err := client.GetAll(context.Background(), &pb.GetAllRequest{})
	if err != nil {
		log.Fatalln(err.Error())
	}

	for {
		resp, err := stream.Recv()
		// 如果服务端数据发送结束，则为EOF
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println(resp.Employee)
	}
}

func addPhoto(client pb.EmployeeServiceClient) {
	imgFile, err := os.Open("ya.jpg")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer imgFile.Close()
	// metadata相当于报文的header，我们只需要把用户number放在header传输一次就可以了
	md := metadata.New(map[string]string{"no": "1996"})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.AddPhoto(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// 循环分块传输数据
	for {
		chunk := make([]byte, 128*1024)
		chunkSize, err := imgFile.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		if chunkSize < len(chunk) {
			chunk = chunk[:chunkSize]
		}
		// 开始分块发送数据
		stream.Send(&pb.AddPhotoRequest{Data: chunk})
	}

	// CloseAndRecv 会向客户端发送一个信号EOF，等待服务端发回一个响应
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(res.IsOk)
}

func saveAll(client pb.EmployeeServiceClient) {
	employees := []*pb.Employee{
		{
			Id:        200,
			No:        201,
			FirstName: "xx",
			LastName:  "xx1",
			MonthSalary: &pb.MonthSalary{
				Basic: 200,
				Bonus: 125.5,
			},
			Status: pb.EmployeeStatus_NORMAL,
			LastModified: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
		{
			Id:        300,
			No:        301,
			FirstName: "asd",
			LastName:  "wefewf",
			MonthSalary: &pb.MonthSalary{
				Basic: 300,
				Bonus: 5.5,
			},
			Status: pb.EmployeeStatus_NORMAL,
			LastModified: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
		{
			Id:        400,
			No:        401,
			FirstName: "www",
			LastName:  "wwwwq",
			MonthSalary: &pb.MonthSalary{
				Basic: 4566,
				Bonus: 100,
			},
			Status: pb.EmployeeStatus_NORMAL,
			LastModified: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
	}

	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}
	// 我们不知道什么时候服务器会把数据发回，我们不能在这阻塞，采用goroutine
	finishChannel := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				finishChannel <- struct{}{}
				break
			}
			if err != nil {
				log.Fatalln(err.Error())
			}
			fmt.Println(res.Employee)
		}
	}()

	for _, e := range employees {
		err := stream.Send(&pb.EmployeeRequest{
			Employee: e,
		})
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	stream.CloseSend()
	<-finishChannel
}
