package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	// "time"

	pb "webserver/biz/pb"

	"google.golang.org/grpc"
)

type businessServer struct {
	pb.UnimplementedBusinessServer
}

type User struct {
	id int32
	name string
	family string
	age int32
	sex string
	created_at string
}

var db = connectToDB()

func (s *businessServer) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	user_id := req.UserId
	var rows *sql.Rows
		
	if user_id > 0 {
		rows, _ = db.Query("SELECT id, name, family, age, sex, created_at FROM users where id = $1", user_id)
		fmt.Println("getting an specific user")
	} else {
		rows, _ = db.Query("SELECT id, name, family, age, sex, created_at FROM users limit 100")
		fmt.Println("in 100 users")
	}
	print(rows)
	pbUsers := []*pb.User{}
	defer rows.Close()
	for rows.Next() {
		var user User
		fmt.Println(rows.Columns())
	
		_ = rows.Scan(&user.id, &user.name, &user.family, &user.age, &user.sex, &user.created_at)
 
    	// fmt.Printf("%s %s %d %d %s %s\n", user.name, user.family, user.id, user.age, user.sex, user.created_at)
		// fmt.Println("here")
		pbUser := &pb.User{Name: user.name, Family: user.family, Id: user.id, Age: user.age, Sex: user.sex, CreatedAt: user.created_at}
		pbUsers = append(pbUsers, pbUser)
	}
	return &pb.GetUsersResponse{Users: pbUsers}, nil
}

func (s *businessServer) GetUsersWithSQLInject(ctx context.Context, req *pb.GetUsersWithSQLInjectRequest) (*pb.GetUsersResponse, error) {
	return &pb.GetUsersResponse{}, nil
}

func main() {
	defer db.Close()
	lis, err := net.Listen("tcp", ":5062")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	bizServer := &businessServer{}
	pb.RegisterBusinessServer(s, bizServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
