package main

import (
	"fmt"
	"net"

	"github.com/keithalucas/jsonrpc/pkg/jsonrpc"
	"github.com/keithalucas/jsonrpc/pkg/spdk"
)

func main() {

	conn, err := net.Dial("unix", "/var/tmp/spdk.sock")

	if err != nil {
		fmt.Printf("Error opening socket: %v", err)
	}

	client := jsonrpc.NewClient(conn)

	errChan := client.Init()

	client.SendMsg("bdev_get_bdevs", nil)

	longhornCreate := spdk.LonghornVolumeCreate{
		Name: "demo",
		Replicas: []spdk.LonghornVolumeCreateReplica{
			spdk.LonghornVolumeCreateReplica{
				//Address: "127.0.0.1",
				Address: "lvs1/volume1",
				Port:    4420,
			},
			spdk.LonghornVolumeCreateReplica{
				//Address: "127.0.0.1",
				Address: "lvs1/volume2",
				Port:    4420,
			},
			spdk.LonghornVolumeCreateReplica{
				//Address: "127.0.0.1",
				Address: "lvs1/volume3",
				Port:    4420,
			},
		},
	}

	ext := spdk.NewLonghornSetExternalAddress("127.0.0.1")
	client.SendMsg(ext.GetMethod(), ext)

	aio := spdk.NewAioCreate("aio1", "/root/aio1.img", 4096)

	client.SendMsg(aio.GetMethod(), aio)
	//time.Sleep(time.Second)
	client.SendMsg("bdev_longhorn_create", longhornCreate)
	//time.Sleep(time.Second)
	client.SendMsg("bdev_get_bdevs", nil)

	lrc := spdk.NewLonghornCreateReplica("replica", 1024*1024*1024, "lvs1", "127.0.0.1", 4420)
	client.SendMsg(lrc.GetMethod(), lrc)

	<-errChan
}
