# Go-gRPC-Sample

## gRPCのサンプル

### 実行手順
1. [cmd/server](https://github.com/ChikaKakazu/Go-gRPC-Sample/tree/main/cmd/server)に移動する
    - ```cd cmd/server/```
1. ```go run main.go```で実行
1. 別でターミナルを起動し、[cmd/client](https://github.com/ChikaKakazu/Go-gRPC-Sample/tree/main/cmd/client)に移動する
    - ```cd cmd/client```
1. ```go run main.go```で実行
1. 以下のようなテキストが表示されるのでそれぞれ数字を入力する
    - ```start gRPC Client.
        1: send Request
        2: HelloServerStream
        3: HelloClientStream
        4: HelloBiStreams
        5: exit
      ```
1. 上から```Unary RPC```, ```Server Streaming RPC```, ```Client Streaming RPC```, ```Bidirectional Streaming RPC```となっている 
