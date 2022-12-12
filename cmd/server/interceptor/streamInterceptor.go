package interceptor

import (
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
)

// grpc.ServerStreamインターフェースを満たす独自構造体
type myServerStreamWrapper1 struct {
	grpc.ServerStream
}

func MyStreamServerInterceptor1(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// ストリームがopenされたときに行われる前処理
	log.Println("[pre stream] my stream server interceptor 1 : ", info.FullMethod)

	err := handler(srv, &myServerStreamWrapper1{ss})

	// ストリームがcloseされるときに行われる後処理
	log.Println("[post stream] my stream server interceptor 1: ")
	return err
}

// リクエストを受信するinterfaceメソッド
func (s *myServerStreamWrapper1) RecvMsg(m interface{}) error {
	// ストリームから、リクエストを受信
	err := s.ServerStream.RecvMsg(m)
	// 受信したリクエストを、ハンドラで処理する前に差し込む前処理
	if !errors.Is(err, io.EOF) {
		log.Println("[pre message] my stream server interceptor 1: ", m)
	}
	return err
}

// リクエストを送信するinterfaceメソッド
func (s *myServerStreamWrapper1) SendMsg(m interface{}) error {
	// ハンドラで作成したレスポンスを、ストリームから返信する直前に差し込む後処理
	log.Println("[post message] my stream server interceptor 1: ", m)
	return s.ServerStream.SendMsg(m)
}

// 2つ目のインターセプタ
func MyStreamServerInterceptor2(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("[pre stream] my stream server interceptor 2: ", info.FullMethod)
	err := handler(srv, &myServerStreamWrapper2{ss}) // 本来のストリーム処理
	log.Println("[post stream] my stream server interceptor 2: ")
	return err
}

type myServerStreamWrapper2 struct {
	grpc.ServerStream
}

func (s *myServerStreamWrapper2) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		log.Println("[pre message] my stream server interceptor 2: ", m)
	}
	return err
}

func (s *myServerStreamWrapper2) SendMsg(m interface{}) error {
	log.Println("[post message] my stream server interceptor 2: ", m)
	return s.ServerStream.SendMsg(m)
}
