package main

import (
	"time"

	"github.com/lqgl/grpc-server/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var employees = []*pb.Employee{
	{
		Id:        1,
		No:        1995,
		FirstName: "Chandler",
		LastName:  "Bing",
		MonthSalary: &pb.MonthSalary{
			Basic: 5000,
			Bonus: 200,
		},
		Status: pb.EmployeeStatus_NORMAL,
		LastModified: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
	{
		Id:        2,
		No:        1996,
		FirstName: "Rachel",
		LastName:  "Green",
		MonthSalary: &pb.MonthSalary{
			Basic: 5500,
			Bonus: 300,
		},
		Status: pb.EmployeeStatus_ON_VACATION,
		LastModified: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
	{
		Id:        3,
		No:        1997,
		FirstName: "Jimmy",
		LastName:  "Chen",
		MonthSalary: &pb.MonthSalary{
			Basic: 6000,
			Bonus: 200,
		},
		Status: pb.EmployeeStatus_RESIGNED,
		LastModified: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
}
