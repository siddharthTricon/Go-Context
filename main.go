package main

import(
	"context"
	"fmt"
	"time"
	"net/http"
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


	http.HandleFunc("/data", dataHandler)

	port := ":8080"
	fmt.Println("Server listening on port", port)

	if err := http.ListenAndServe(port, nil); err != nil{
		fmt.Println("Error starting server:", err)
	}
}


type DataRespone struct{
	Message string `json:"message"`
}

func fetchData(ctx context.Context) (string, error){
	select{
	case <- time.After(3 * time.Second):
		return "Data fetched successfully", nil
	case <- ctx.Done():
		return "", ctx.Err()
	}
}


func dataHandler(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(r.Context(), 2 * time.Second)
	defer cancel()

	_, err := fetchData(ctx)
	if err != nil{
		if err == context.DeadlineExceeded{
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		}else{
			http.Error(w, "Request cancelled", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
