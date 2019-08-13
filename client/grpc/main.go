package main

import (
	"context"
	"dogfooter-control/control/cmd/service"
	"dogfooter-control/control/pkg/grpc/pb"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"os"
	"time"
)

var fs = flag.NewFlagSet("dogfooter-control", flag.ExitOnError)
var gRpcAddr = fs.String("grpc-addr", "", "gRPC server address")

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Fprintf(os.Stderr, "%s took %v\n", name, elapsed.Seconds())
}
func main() {
	fs.Parse(os.Args[1:])
	defer timeTrack(time.Now(), "Test")

	var cfg = service.ServerConfig{}
	if err := cfg.ReadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		os.Exit(1)
	} else {
		*gRpcAddr = cfg.GrpcHosts
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*gRpcAddr, opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to dial: %v\n", err)
		return
	}
	defer conn.Close()

	c := pb.NewControlClient(conn)
	for i := 0; i < 5; i++ {
		if reply, err := c.Api(
			context.Background(),
			&pb.ApiRequest{
				Category: "patient",
			}); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		} else {
			if j, err2 := json.MarshalIndent(reply, "", " "); err2 != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err2)
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", string(j))
			}
		}
	}

	//if reply, err := c.ReadAllTimeOffHoliday(
	//	context.Background(),
	//	&pb.ReadAllTimeOffHolidayRequest{
	//	}); err != nil {
	//	fmt.Fprintf(os.Stderr, "%v\n", err)
	//	return
	//} else {
	//	if j, err2 := json.MarshalIndent(reply, "", " "); err2 != nil {
	//		fmt.Fprintf(os.Stderr, "%v\n", err2)
	//	} else {
	//		fmt.Fprintf(os.Stderr, "%v\n", string(j))
	//	}
	//}

	//if reply, err := c.ReadTimeOffHoliday(
	//	context.Background(),
	//	&pb.ReadTimeOffHolidayRequest{
	//		Id: "8ba562ce-3c1b-4374-96b7-0b870d3654dd",
	//	}); err != nil {
	//	fmt.Fprintf(os.Stderr, "%v\n", err)
	//	return
	//} else {
	//	if j, err2 := json.MarshalIndent(reply, "", " "); err2 != nil {
	//		fmt.Fprintf(os.Stderr, "%v\n", err2)
	//	} else {
	//		fmt.Fprintf(os.Stderr, "%v\n", string(j))
	//	}
	//}

	//if reply, err := c.ReadTimeOffHolidayBulk(
	//	context.Background(),
	//	&pb.ReadTimeOffHolidayBulkRequest{
	//		TimeOffHolidayIdList: []string{
	//			"8ba562ce-3c1b-4374-96b7-0b870d3654dd",
	//			"c84ca9fd-3e3a-45f3-b014-819000cf2737",
	//			"0489a159-654c-43e7-848f-a554990c1fab",
	//		},
	//	}); err != nil {
	//	fmt.Fprintf(os.Stderr, "%v\n", err)
	//	return
	//} else {
	//	if j, err2 := json.MarshalIndent(reply, "", " "); err2 != nil {
	//		fmt.Fprintf(os.Stderr, "%v\n", err2)
	//	} else {
	//		fmt.Fprintf(os.Stderr, "%v\n", string(j))
	//	}
	//}

	//if reply, err := c.UpdateTimeOffHoliday(
	//	context.Background(),
	//	&pb.UpdateTimeOffHolidayRequest{
	//		Id:        "8ba562ce-3c1b-4374-96b7-0b870d3654dd",
	//		Name:      "update_test_dogfooter-control_name",
	//		Year:      "Every year",
	//		FromMonth: 5,
	//		FromDay:   20,
	//		ToMonth:   5,
	//		ToDay:     25,
	//	}); err != nil {
	//	fmt.Fprintf(os.Stderr, "%v\n", err)
	//	return
	//} else {
	//	if j, err2 := json.MarshalIndent(reply, "", " "); err2 != nil {
	//		fmt.Fprintf(os.Stderr, "%v\n", err2)
	//	} else {
	//		fmt.Fprintf(os.Stderr, "%v\n", string(j))
	//	}
	//}
	//
	//if reply, err := c.DeleteTimeOffHoliday(
	//	context.Background(),
	//	&pb.DeleteTimeOffHolidayRequest{
	//		Id: "8ba562ce-3c1b-4374-96b7-0b870d3654dd",
	//	}); err != nil {
	//	fmt.Fprintf(os.Stderr, "%v\n", err)
	//	return
	//} else {
	//	if j, err2 := json.MarshalIndent(reply, "", " "); err2 != nil {
	//		fmt.Fprintf(os.Stderr, "%v\n", err2)
	//	} else {
	//		fmt.Fprintf(os.Stderr, "%v\n", string(j))
	//	}
	//}
}
