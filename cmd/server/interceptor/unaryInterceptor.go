package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

// Unary RPCのインターセプタ
func MyUnaryServerInterceptor1(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// ハンドラの前に割り込ませる前処理
	log.Println("[pre] my unary server interceptor 1 : ", info.FullMethod)

	// 本来の処理
	res, err := handler(ctx, req)

	// ハンドラの後に割り込ませる後処理
	log.Println("[post] my unary server interceptor 1 : ", res)

	return res, err
}

func MyUnaryServerInterceptor2(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("[pre] my unary server interceptor 2: ", info.FullMethod, req)
	res, err := handler(ctx, req) // 本来の処理
	log.Println("[post] my unary server interceptor 2: ", res)
	return res, err
}
