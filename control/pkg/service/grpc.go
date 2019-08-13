package service

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

type GrpcObject struct {
	ImageConn     *grpc.ClientConn
	PatientConn   *grpc.ClientConn
	DiagnosisConn *grpc.ClientConn
	DateConn      *grpc.ClientConn
}

var GrpcConn GrpcObject

var opts []grpc.DialOption

var kacp = keepalive.ClientParameters{
	Time:                60 * time.Second, // send pings every 60 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

//func init() {
//	opts = append(opts, grpc.WithInsecure())
//	opts = append(opts, grpc.WithKeepaliveParams(kacp))
//	go func() {
//		for range time.Tick(time.Second * 60) {
//			GrpcConn.ConnectionCheck()
//		}
//	}()
//	return
//}
//func (g *GrpcObject) ConnectionCheck() {
//	{
//		conn, _ := g.GetImageConn()
//		if conn != nil {
//			d := DogfooterImage{}
//			_, err := d.KeepAlive(nil, Payload{}, UserObject{})
//			if err != nil {
//				fmt.Fprintf(os.Stderr, "DEBUG: ImageConn, This message is printed when rpc connection is wrong: %v\n", g.ImageConn.GetState())
//				g.ImageConn.Close()
//				g.ImageConn = nil
//			}
//		}
//	}
//	{
//		conn, _ := g.GetDiagnosisConn()
//		if conn != nil {
//			d := DogfooterDiagnosis{}
//			_, err := d.KeepAlive(nil, Payload{}, UserObject{})
//			if err != nil {
//				fmt.Fprintf(os.Stderr, "DEBUG: DiagnosisConn, This message is printed when rpc connection is wrong: %v\n", g.DiagnosisConn.GetState())
//				g.DiagnosisConn.Close()
//				g.DiagnosisConn = nil
//			}
//		}
//	}
//	{
//		conn, _ := g.GetDateConn()
//		if conn != nil {
//			d := DogfooterDate{}
//			_, err := d.KeepAlive(nil, Payload{}, UserObject{})
//			if err != nil {
//				fmt.Fprintf(os.Stderr, "DEBUG: DateConn, This message is printed when rpc connection is wrong: %v\n", g.DateConn.GetState())
//				g.DateConn.Close()
//				g.DateConn = nil
//			}
//		}
//	}
//	{
//		conn, _ := g.GetPatientConn()
//		if conn != nil {
//			d := DogfooterPatient{}
//			_, err := d.KeepAlive(nil, Payload{}, UserObject{})
//			if err != nil {
//				fmt.Fprintf(os.Stderr, "DEBUG: PatientConn, This message is printed when rpc connection is wrong: %v\n", g.PatientConn.GetState())
//				g.PatientConn.Close()
//				g.PatientConn = nil
//			}
//		}
//	}
//}
//
//func (g *GrpcObject) GetImageConn() (conn *grpc.ClientConn, err error) {
//	if g.ImageConn != nil {
//		if g.ImageConn.GetState() != connectivity.Ready {
//			fmt.Fprintf(os.Stderr, "Image - Warning!! %v\n", g.ImageConn.GetState())
//			g.ImageConn = nil
//		}
//	}
//	if g.ImageConn == nil {
//		addr := GetConfigClientImageGrpc()
//		g.ImageConn, err = grpc.Dial(addr, opts...)
//		if err != nil {
//			err = fmt.Errorf("fail to dial: %v", err)
//			return
//		}
//	}
//
//	conn = g.ImageConn
//
//	//defer Conn.Close()
//	return
//}
//func (g *GrpcObject) GetPatientConn() (conn *grpc.ClientConn, err error) {
//	if g.PatientConn != nil {
//		if g.PatientConn.GetState() != connectivity.Ready {
//			fmt.Fprintf(os.Stderr, "Warning!! %v\n", g.PatientConn.GetState())
//			g.PatientConn = nil
//		}
//	}
//	if g.PatientConn == nil {
//		addr := GetConfigClientPatientGrpc()
//		g.PatientConn, err = grpc.Dial(addr, opts...)
//		if err != nil {
//			err = fmt.Errorf("fail to dial: %v", err)
//			return
//		}
//	}
//	conn = g.PatientConn
//
//	//defer Conn.Close()
//	return
//}
//func (g *GrpcObject) GetDiagnosisConn() (conn *grpc.ClientConn, err error) {
//	if g.DiagnosisConn != nil {
//		if g.DiagnosisConn.GetState() != connectivity.Ready {
//			fmt.Fprintf(os.Stderr, "Warning!! %v\n", g.DiagnosisConn.GetState())
//			g.DiagnosisConn = nil
//		}
//	}
//	if g.DiagnosisConn == nil {
//		addr := GetConfigClientDiagnosisGrpc()
//		g.DiagnosisConn, err = grpc.Dial(addr, opts...)
//		if err != nil {
//			err = fmt.Errorf("fail to dial: %v", err)
//			return
//		}
//	}
//	conn = g.DiagnosisConn
//
//	//defer Conn.Close()
//	return
//}
//func (g *GrpcObject) GetDateConn() (conn *grpc.ClientConn, err error) {
//	if g.DateConn != nil {
//		if g.DateConn.GetState() != connectivity.Ready {
//			fmt.Fprintf(os.Stderr, "Warning!! %v\n", g.DateConn.GetState())
//			g.DateConn = nil
//		}
//	}
//	if g.DateConn == nil {
//		addr := GetConfigClientDateGrpc()
//		g.DateConn, err = grpc.Dial(addr, opts...)
//		if err != nil {
//			err = fmt.Errorf("fail to dial: %v", err)
//			return
//		}
//	}
//	conn = g.DateConn
//
//	//defer Conn.Close()
//	return
//}
