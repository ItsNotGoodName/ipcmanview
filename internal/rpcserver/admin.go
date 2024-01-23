package rpcserver

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewAdmin(db repo.DB) *Admin {
	return &Admin{
		db: db,
	}
}

type Admin struct {
	db repo.DB
}

func (a *Admin) CreateGroup(ctx context.Context, req *rpc.CreateGroupReq) (*rpc.CreateGroupResp, error) {
	id, err := auth.CreateGroup(ctx, a.db, models.Group{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		if errs, ok := asValidationErrors(err); ok {
			return nil, NewError(err, "Failed to create group.").Validation(errs, [][2]string{
				{"Name", "name"},
				{"Description", "description"},
			})
		}

		if constraintErr, ok := sqlite.AsConstraintError(err, sqlite.CONSTRAINT_UNIQUE); ok {
			return nil, NewError(err, "Failed to create group.").Constraint(constraintErr, [][3]string{
				{"groups.name", "name", "Name already taken."},
			})
		}
	}

	return &rpc.CreateGroupResp{
		Id: id,
	}, nil
}

func (a *Admin) DeleteGroup(ctx context.Context, req *rpc.DeleteGroupReq) (*emptypb.Empty, error) {
	err := a.db.DeleteGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) GetGroup(ctx context.Context, req *rpc.GetGroupReq) (*rpc.GetGroupResp, error) {
	v, err := a.db.GetGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &rpc.GetGroupResp{
		Group: &rpc.Group{
			Id:            v.ID,
			Name:          v.Name,
			Description:   v.Description,
			CreatedAtTime: timestamppb.New(v.CreatedAt.Time),
			UpdatedAtTime: timestamppb.New(v.UpdatedAt.Time),
		},
	}, nil
}

// UpdateGroup implements rpc.Admin.
func (*Admin) UpdateGroup(context.Context, *rpc.UpdateGroupReq) (*emptypb.Empty, error) {
	panic("unimplemented")
}

func (a *Admin) ListGroups(ctx context.Context, req *rpc.ListGroupsReq) (*rpc.ListGroupsResp, error) {
	page := parsePagePagination(req.Page)

	dbGroups, err := a.db.ListGroup(ctx, repo.ListGroupParams{
		Limit:  int64(page.Limit()),
		Offset: int64(page.Offset()),
	})
	if err != nil {
		return nil, NewError(err).Internal()
	}

	count, err := a.db.CountGroup(ctx)
	if err != nil {
		return nil, NewError(err).Internal()
	}

	groups := make([]*rpc.Group, 0, len(dbGroups))
	for _, v := range dbGroups {
		groups = append(groups, &rpc.Group{
			Id:            v.ID,
			Name:          v.Name,
			Description:   v.Description,
			UserCount:     v.UserCount,
			CreatedAtTime: timestamppb.New(v.CreatedAt.Time),
			UpdatedAtTime: timestamppb.New(v.UpdatedAt.Time),
		})
	}

	return &rpc.ListGroupsResp{
		Groups:     groups,
		PageResult: convertPagePaginationResult(page.Result(int(count))),
	}, nil
}
