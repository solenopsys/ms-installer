# ms-installer module

# proto gen
protoc --proto_path=. --go_out=.  api.proto

# Env Vars
- zmq.SocketUrl = "tcp://*:5561"
- developerMode=true
- kubeconfigPath=C:\dev\sources\k3s.yaml

# Local tunnel
ssh -L 6443:127.0.0.1:6443 root@10.23.92.23
 