cd dev
set GOOS=linux
go build initialize.go
cd ..
docker build -f ./Dockerfile-deploy -t kademlia:latest .
pause
