package interceptor

import (
	"context"
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
)

type myClientStreamWrapper1 struct {
	grpc.ClientStream
}

func MyStreamClientInterceptor1(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
	method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// ストリームがopenされる前に行われる前処理
	log.Println("[pre] my stream client interceptor 1", method)

	stream, err := streamer(ctx, desc, cc, method, opts...)

	return &myClientStreamWrapper1{stream}, err
}

func (s *myClientStreamWrapper1) SendMsg(m interface{}) error {
	// リクエスト送信前に割り込ませる処理
	log.Println("[pre message] my stream client interceptor 1: ", m)

	// リクエスト送信
	return s.ClientStream.SendMsg(m)
}

func (s *myClientStreamWrapper1) RecvMsg(m interface{}) error {
	// レスポンス受信処理
	err := s.ClientStream.RecvMsg(m)

	// レスポンス受信後に割り込ませる処理
	if !errors.Is(err, io.EOF) {
		log.Println("[post message] my stream client interceptor 1: ", m)
	}
	return err
}

func (s *myClientStreamWrapper1) CloseSend() error {
	// ストリームをclose
	err := s.ClientStream.CloseSend()

	// ストリームがcloseされた後に行われる後処理
	log.Println("[post] my stream client interceptor 1")
	return err
}
