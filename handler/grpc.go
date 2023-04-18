package handler

import (
	"context"
	"github.com/GarnBarn/common-go/proto"
	"github.com/GarnBarn/gb-tag-service/service"
)

type Grpc struct {
	proto.UnimplementedTagServer
	tagService service.Tag
}

func NewGrpcServer(tagService service.Tag) *Grpc {
	return &Grpc{
		tagService: tagService,
	}
}

func (g *Grpc) GetTag(c context.Context, tr *proto.TagRequest) (*proto.TagPublic, error) {
	tagPublic, err := g.tagService.GetTagById(int(tr.TagId), tr.ConsealPrivateKey)
	if err != nil {
		return nil, err
	}
	protoTag := tagPublic.ToProtoTag()
	return &protoTag, nil
}

func (g *Grpc) IsTagExists(c context.Context, tr *proto.TagRequest) (*proto.TagExistsResponse, error) {
	isTagExist := g.tagService.IsTagExist(int(tr.TagId))
	return &proto.TagExistsResponse{
		IsExists: isTagExist,
	}, nil
}
