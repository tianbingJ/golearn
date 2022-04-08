package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

//Background()主要用于main函数、初始化代码中，作为context顶层的context
//TOD O() 使用场景不清晰的情况下使用这个

func worker(ch chan struct{}, ctx context.Context) {
	go worker2(ch, ctx)
Loop:
	for {
		//work
		time.Sleep(time.Microsecond * 10)
		select {
		case <-ctx.Done():
			break Loop
		default:
		}
	}
	ch <- struct{}{}
}

func worker2(ch chan<- struct{}, ctx context.Context) {
Loop:
	for {
		//work
		time.Sleep(time.Microsecond * 10)
		select {
		case <-ctx.Done():
			break Loop
		default:
		}
	}
	ch <- struct{}{}
}

func TestBasicUse(t *testing.T) {
	ch := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())
	go worker(ch, ctx)

	cancel()

	<-ch
	<-ch

	close(ch)
}

func TestWithTimeout(t *testing.T) {
	ch := make(chan struct{})
	t1 := time.Now()

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	go worker2(ch, ctx)

	<- ch
	t2 := time.Now()
	diff := t2.Sub(t1)
	if diff < time.Second {
		t.Fatalf("超时未生效")
	}
}

func TestWithValue(t *testing.T) {
	ctx := context.Background()
	key1 := 123
	val1 := 456
	v1Ctx := context.WithValue(ctx, key1, val1)
	v1 := v1Ctx.Value(123)
	if v1 != val1 {
		t.Fatalf("not equal %v %v", val1, v1)
	}

	key2 := "hello"
	val2 := "world"
	v2Ctx := context.WithValue(v1Ctx, key2, val2)
	v1 = v2Ctx.Value(key1)
	if v1 != val1 {
		t.Fatalf("not equal %v %v", val1, v1)
	}
	v2 := v2Ctx.Value(key2)
	if v2 != val2 {
		t.Fatalf("not equal %v %v", val2, v2)
	}
}

func TestDoneCancel(t *testing.T) {
	ctx := context.Background()
	ch := ctx.Done()
	fmt.Println(ch)
}
