protoc --proto_path=proto/order --proto_path=third-party --go_out=plugins=grpc:proto order/order.proto
protoc --proto_path=proto/payment --proto_path=third-party --go_out=plugins=grpc:proto payment/payment.proto
protoc --proto_path=proto/stock --proto_path=third-party --go_out=plugins=grpc:proto stock/stock.proto