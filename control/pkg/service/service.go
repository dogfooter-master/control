package service

import (
	"context"
	"dogfooter-control/control/pkg/grpc/pb"
	"errors"
	"fmt"
)

// ControlService describes the service.
type ControlService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	Api(ctx context.Context, req Payload) (res Payload, err error)
	Root(ctx context.Context, req *pb.RootRequest) (res *pb.RootReply, err error)
	File(ctx context.Context, req FileVars) (res FileVars, err error)
}

type basicControlService struct{}

// 같은방식인데 조회를 할 때 연결하느냐 아니면 업데이트를 할 때 연결하느냐의 차이라고 생각한다.
// 업데이트를 할 때 연결하는 것의 단점은 실제 조회시 싱크가 어긋날 가능성이 있고
// 조회를 할 때 연결하는 것은 성능상 느리다는 단점이 있다.
// 콘트롤 서버에서 각각의 데이터를 업데이트를 할 때 잘 연결하자.
func (b *basicControlService) Api(ctx context.Context, req Payload) (res Payload, err error) {
	// API Gateway
	switch req.Category {
	case "public":
		// 인증토큰 필요없는 서비스, 그 외에는 모두 인증토큰 필요
		s := DogfooterPublic{}
		res, err = s.Service(ctx, req)
	case "private":
		s := DogfooterPrivate{}
		res, err = s.Service(ctx, req)
	default:
		err = fmt.Errorf("unknown category: %v", req.Category)
	}

	return
}
func (b *basicControlService) File(ctx context.Context, req FileVars) (res FileVars, err error) {
	do := UserObject{
		SecretToken: SecretTokenObject{
			Token: req.AccessToken,
		},
	}
	var ro UserObject
	if ro, err = do.SecretToken.Authenticate(); err != nil {
		err = errors.New("unauthorized access")
		return
	} else {
		if len(req.Separators) <  3 {
			err = errors.New("invalid access")
			return
		}

		res = FileVars {
			Method: req.Method,
			R: req.R,
		}
		switch req.Separators[0] {
		case "i":
			res.StaticFilePath = StaticDataFilePath(ro.Login.Account, req.Id, req.Date, req.File)
		case "u":
			// 사용자 관련 파일들
			switch req.Separators[1] {
			case "p":
				switch req.Separators[2] {
				case "a":
					res.StaticFilePath = StaticProfileAvatarFilePath(ro.Login.Account, req.Date, req.File)
				}
			}
		}
	}
	return res, err
}
func (b *basicControlService) Root(ctx context.Context, req *pb.RootRequest) (res *pb.RootReply, err error) {
	return res, err
}

// NewBasicControlService returns a naive, stateless implementation of ControlService.
func NewBasicControlService() ControlService {
	return &basicControlService{}
}

// New returns a ControlService with all of the expected middleware wired in.
func New(middleware []Middleware) ControlService {
	var svc ControlService = NewBasicControlService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
