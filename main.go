package main

import(
	"context"
	"fmt"
	"time"
)


func performTask(ctx context.Context, taskname string){
	select{
		case <- time.After(3 * time.Second):
			fmt.Println(taskname, "completed successfully")
		case <- ctx.Done():
			fmt.Println(taskname,"cancelled:", ctx.Err())
	}
}

func main(){
	fmt.Println("Starting main...")

	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	go performTask(ctx, "Task1")

	time.Sleep(4 * time.Second)
	fmt.Println("Exiting main...")
}