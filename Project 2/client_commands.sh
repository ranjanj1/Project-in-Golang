./start_server.sh & sleep 10

grpcurl -plaintext localhost:50051 describe

echo "=== Create token ==="
grpcurl -plaintext -d  '{ "id": "1234" }' localhost:50051 token.v1.TokenService/CreateToken

echo "=== Write token ==="
grpcurl -plaintext -d  '{"token": {"id": "1234", "high": 10, "mid": 5, "low": 1, "name": "test_token"}}' localhost:50051 token.v1.TokenService/WriteToken

echo "=== Read token ==="
grpcurl -plaintext -d  '{ "id": "1234" }' localhost:50051 token.v1.TokenService/ReadToken

echo "=== Drop token ==="
grpcurl -plaintext -d  '{ "id": "1234" }' localhost:50051 token.v1.TokenService/DropToken

kill -9 $(lsof -t -i:50051)