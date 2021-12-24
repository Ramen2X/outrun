package orpc

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/Ramen2X/outrun/config"
	"github.com/Ramen2X/outrun/orpc/rpcobj"
)

func Start() {
	rpc.Register(new(rpcobj.Toolbox))
	rpc.Register(new(rpcobj.Config))
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":"+config.CFile.RPCPort)
	if err != nil {
		log.Fatal("error listening in ORPC:", err)
	}
	log.Println("Starting ORPC server on port " + config.CFile.RPCPort)
	go http.Serve(listener, nil)
}
