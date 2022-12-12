package main

import (
	"context"
	"errors"
	"fmt"
	"grpc-sample/cmd/server/interceptor"
	hellopb "grpc-sample/pkg/grpc"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	// 8080番portのLisnterを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// gRPCサーバーの作成
	// インターセプタの導入
	s := grpc.NewServer(
		// grpc.UnaryInterceptor(interceptor.MyUnaryServerInterceptor1),
		// 複数のインターセプタを使用するためにはChain~を使う
		// grpc.StreamInterceptor(interceptor.MyStreamServerInterceptor1),
		grpc.ChainUnaryInterceptor(
			interceptor.MyUnaryServerInterceptor1,
			interceptor.MyUnaryServerInterceptor2,
		),
		grpc.ChainStreamInterceptor(
			interceptor.MyStreamServerInterceptor1,
			interceptor.MyStreamServerInterceptor2,
		),
	)

	// gRPCサーバーをGreetingServiceに登録
	hellopb.RegisterGreetingServiceServer(s, NewMyServer())

	// サーバーリフレクションの設定
	reflection.Register(s)

	// 作成したgRPCサーバーを稼働させる
	go func() {
		log.Printf("start gRPC server port : %v", port)
		s.Serve(listener)
	}()

	// Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}

func (s *myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	//  エラー確認用
	// stat := status.New(codes.Unknown, "unknown error occurred")
	// stat, _ = stat.WithDetails(&errdetails.DebugInfo{
	// 	Detail: "detail reason of err",
	// })
	// err := stat.Err()
	// return nil, err

	// クライアントから送られたメタデータを取り出す
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println(md)
	}

	// HTTP/2ではヘッダーフレームを分割して送る
	// サーバーがクライアントに送る最初のヘッダーフレームのことをヘッダー、最後に送るヘッダーフレームのことをトレーサーと呼ぶ
	headerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "header"})
	if err := grpc.SetHeader(ctx, headerMD); err != nil {
		return nil, err
	}

	trailerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "trailer"})
	if err := grpc.SetTrailer(ctx, trailerMD); err != nil {
		return nil, err
	}

	// リクエストからnameフィールドを取り出して"Hello, [名前]!"というレスポンスを返す
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (s *myServer) HelloServerStream(req *hellopb.HelloRequest, stream hellopb.GreetingService_HelloServerStreamServer) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&hellopb.HelloResponse{
			Message: fmt.Sprintf("[%d] Hello %s", i, req.GetName()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

func (s *myServer) HelloClientStream(stream hellopb.GreetingService_HelloClientStreamServer) error {
	nameList := make([]string, 0)
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			message := fmt.Sprintf("Hello %v", nameList)
			return stream.SendAndClose(&hellopb.HelloResponse{
				Message: message,
			})
		}
		if err != nil {
			return err
		}
		nameList = append(nameList, req.GetName())
	}
}

func (s *myServer) HelloBiStreams(stream hellopb.GreetingService_HelloBiStreamsServer) error {
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		log.Println(md)
	}

	// (パターン1)すぐにヘッダーを送信したいならばこちら
	headerMD := metadata.New(map[string]string{"type": "stream", "from": "server", "in": "header"})
	if err := stream.SendHeader(headerMD); err != nil {
		return err
	}
	// (パターン2)本来ヘッダーを送るタイミングで送りたいならばこちら
	// if err := stream.SetHeader(headerMD); err != nil {
	// 	return err
	// }
	trailerMD := metadata.New(map[string]string{"type": "stream", "from": "server", "in": "trailer"})
	stream.SetTrailer(trailerMD)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		message := fmt.Sprintf("Hello, %v!", req.GetName())
		if err := stream.Send(&hellopb.HelloResponse{
			Message: message,
		}); err != nil {
			return err
		}
	}
}
