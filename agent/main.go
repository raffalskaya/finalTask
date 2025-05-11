package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	pb "github.com/raffalskaya/finalTask/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var sleepTime = 1 * time.Second

func calculate() {
	for {
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			fmt.Println("Оркестратор не запущен")
			time.Sleep(sleepTime)
			continue
		}
		defer conn.Close()

		// Создаем клиент
		client := pb.NewAPIServiceClient(conn)

		// Вызываем метод GetTask
		req := &emptypb.Empty{} // Пустой запрос
		resp, err := client.GetTask(context.Background(), req)
		if err != nil {
			fmt.Println("Ошибка получения задачи")
			time.Sleep(sleepTime)
			continue
		}

		if resp.Enabled != true {
			fmt.Println("Нет задач")
			time.Sleep(sleepTime)
			continue
		}

		var result float64
		switch resp.Task.Operation {
		case "+":
			result = resp.Task.Arg1 + resp.Task.Arg2
		case "-":
			result = resp.Task.Arg1 - resp.Task.Arg2
		case "*":
			result = resp.Task.Arg1 * resp.Task.Arg2
		case "/":
			if resp.Task.Arg2 == 0 {
				result = 0
			}
			result = resp.Task.Arg1 / resp.Task.Arg2
		default:
			result = 0
		}

		taskRes := &pb.TaskResult{
			Id:     resp.Task.Id,
			Result: result,
		}

		setStatusResp, err := client.SetTask(context.Background(), taskRes)
		if err != nil {
			fmt.Println("Ошибка отправки результата задачи", setStatusResp)
			time.Sleep(sleepTime)
			continue
		}
	}
}

func main() {
	COMPUTING_POWER, exists := os.LookupEnv("COMPUTING_POWER")
	if !exists {
		COMPUTING_POWER = "1"
	}

	computing_power_int, err := strconv.Atoi(COMPUTING_POWER)

	if err != nil {
		computing_power_int = 1
	}
	for i := 0; i < computing_power_int; i++ {
		go calculate()
	}

	select {}
}
