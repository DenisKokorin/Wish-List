package group

import (
	"context"

	proto "github.com/DenisKokorin/WishListProto/gen/go/group"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Group interface {
	CreateGroup(ctx context.Context, userId int64, title string) (*proto.GroupResponse, error)
	GetGroup(ctx context.Context, groupId int64) (*proto.GroupResponse, error)
	GetAllGroups(ctx context.Context, userId int64) (*proto.AllGroupsResponse, error)
	InviteUser(ctx context.Context, groupId int64, username string) (string, error)
	AcceptInvite(ctx context.Context, token string) (string, error)
	DeleteGroup(ctx context.Context, groupId int64) (string, error)
	LeaveGroup(ctx context.Context, userId int64, groupId int64) (string, error)
}

type serverAPI struct {
	group Group
	proto.UnimplementedGroupServiceServer
}

func Register(gRPC *grpc.Server, group Group) {
	proto.RegisterGroupServiceServer(gRPC, &serverAPI{group: group})
}

func (s *serverAPI) CreateGroup(ctx context.Context, req *proto.CreateGroupRequest) (*proto.GroupResponse, error) {
	res, err := s.group.CreateGroup(ctx, req.GetUserId(), req.GetTitle())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (s *serverAPI) GetGroup(ctx context.Context, req *proto.GetGroupRequest) (*proto.GroupResponse, error) {
	res, err := s.group.GetGroup(ctx, req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (s *serverAPI) GetAllGroups(ctx context.Context, req *proto.AllGroupsRequest) (*proto.AllGroupsResponse, error) {
	res, err := s.group.GetAllGroups(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (s *serverAPI) InviteUser(ctx context.Context, req *proto.InviteRequest) (*proto.Invite, error) {
	res, err := s.group.InviteUser(ctx, req.GetGroupId(), req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.Invite{Token: res}, nil
}

func (s *serverAPI) AcceptInvite(ctx context.Context, req *proto.Invite) (*proto.AcceptResponse, error) {
	res, err := s.group.AcceptInvite(ctx, req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.AcceptResponse{Success: res}, nil
}

func (s *serverAPI) DeleteGroup(ctx context.Context, req *proto.DeleteGroupRequest) (*proto.DeleteGroupResponse, error) {
	res, err := s.group.DeleteGroup(ctx, req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.DeleteGroupResponse{Success: res}, nil
}

func (s *serverAPI) LeaveGroup(ctx context.Context, req *proto.LeaveRequest) (*proto.LeaveResponse, error) {
	res, err := s.group.LeaveGroup(ctx, req.GetUserId(), req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.LeaveResponse{Success: res}, nil
}
