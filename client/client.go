package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/surajjyoti/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func Sum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting Sum service..")
	req := calculatorpb.SumRequest{
		FirstNum:  15,
		SecondNum: 55,
	}
	resp, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling Sum grpc unary call: %v", err)
	}

	log.Printf("Response from Sum Unary Call : %v", resp.Result)

}

func Prime(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting Prime number service..")
	req := calculatorpb.PrimeRequest{
		Num: 20,
	}
	resStream, err := c.Prime(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling Prime server-side streaming grpc : %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//we have reached to the end of the file
			break
		}

		if err != nil {
			log.Fatalf("error while receving server stream : %v", err)
		}

		fmt.Println("Response From Prime Server : ", msg.Result)
	}

}

func Average(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting Average service")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client-side streaming : %v", err)
	}
	req := []*calculatorpb.AverageRequest{
		&calculatorpb.AverageRequest{Num: 5},
		&calculatorpb.AverageRequest{Num: 10},
		&calculatorpb.AverageRequest{Num: 15},
		&calculatorpb.AverageRequest{Num: 20},
		&calculatorpb.AverageRequest{Num: 25},
		&calculatorpb.AverageRequest{Num: 30},
	}
	for _, val := range req {
		fmt.Println("\nSending number.... : ", val)
		stream.Send(val)
		time.Sleep(1 * time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from server : %v", err)
	}
	fmt.Println("\n****Response From Server : ", resp.Result)

}

func Maxnumber(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting Maxnumber service")
	req := []*calculatorpb.MaxnumberRequest{
		&calculatorpb.MaxnumberRequest{Num: 1},
		&calculatorpb.MaxnumberRequest{Num: 3},
		&calculatorpb.MaxnumberRequest{Num: 7},
		&calculatorpb.MaxnumberRequest{Num: 9},
		&calculatorpb.MaxnumberRequest{Num: 2},
		&calculatorpb.MaxnumberRequest{Num: 5},
		&calculatorpb.MaxnumberRequest{Num: 22},
		&calculatorpb.MaxnumberRequest{Num: 15},
		&calculatorpb.MaxnumberRequest{Num: 21},
		&calculatorpb.MaxnumberRequest{Num: 19},
	}
	stream, err := c.Maxnumber(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client side streaming : %v", err)
	}

	waitchan := make(chan int32)

	go func(req []*calculatorpb.MaxnumberRequest) {
		for _, val := range req {
			fmt.Println("\nSending number... : ", val.Num)
			err := stream.Send(val)
			if err != nil {
				log.Fatalf("error while sending request to Maxnumber service : %v", err)
			}
			time.Sleep(1000 * time.Millisecond)

		}
		stream.CloseSend()
	}(req)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitchan)
				return
			}

			if err != nil {
				log.Fatalf("Error receiving response from server : %v", err)
			}

			fmt.Printf("\nResponse From Server : %v", resp.Result)
		}

	}()
	<-waitchan

}

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//Sum(c)
	//Prime(c)
	//Average(c)
	Maxnumber(c)
}
