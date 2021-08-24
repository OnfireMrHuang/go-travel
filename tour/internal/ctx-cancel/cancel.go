package ctx_cancel

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

// 使用http的一个上下文取消，<-ctx.Done表示监听上下文被取消事件，当我们关闭浏览器的时候serve内部会触发cancel。
func SimpleHttpCancelDemo() {
	// Create an HTTP server that listens on port 8000
	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// This prints to STDOUT to show that processing has started
		fmt.Fprint(os.Stdout, "processing request\n")
		// We use `select` to execute a peice of code depending on which
		// channel receives a message first
		select {
		case <-time.After(2 * time.Second):
			// If we receive a message after 2 seconds
			// that means the request has been processed
			// We then write this as the response
			w.Write([]byte("request processed"))
		case <-ctx.Done():
			// If the request gets cancelled, log it
			// to STDERR
			fmt.Fprint(os.Stderr, "request cancelled\n")
		}
	}))
}

// 上一个http的例子是serve内部触发了cancel，那如果我们的业务代码中要控制上下文的cancel动作怎么做呢？
// context.WithCancel会返回一个新的上下文和cancel函数，这个cancel函数就是把控制权给到我们，我们要取消的时候主动调用cancel就可以了。
func SimpleSubmitCancelDemo()  {
	ctx := context.Background()
	ctx,cancel := context.WithCancel(ctx)

	go func() {
		err := operation1(ctx)
		if err != nil {
			cancel()
		}
	}()
	operation2(ctx)
}

// 基于时间的一个超时控制取消
// 在我们调用http或rpc调用的时候，经常被用来控制超时异常
// 下面这个方法的意思是3秒后会被主动取消，如果要提前取消则直接调用返回的cancel
// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
// 和上面的方法意思类似，只是把时间间隔换成截止到参数里面的日期
// ctx, cancel := context.WithDeadline(ctx, time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))
func SimpleTimeoutCancelDemo()  {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)

	// Make a request, that will call the google homepage
	req, _ := http.NewRequest(http.MethodGet, "http://google.com", nil)
	// Associate the cancellable context we just created to the request
	req = req.WithContext(ctx)

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	// If the request failed, log to STDOUT
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	// Print the statuscode if the request succeeds
	fmt.Println("Response received, status code:", res.StatusCode)
}

func operation1(ctx context.Context) error {
	time.Sleep(100 * time.Millisecond)
	return errors.New("failed")
}

func operation2(ctx context.Context) {
	select {
		case <-time.After(500 * time.Millisecond):
			fmt.Println("done")
		case <- ctx.Done():
			fmt.Println("halted operation2")
	}
}

