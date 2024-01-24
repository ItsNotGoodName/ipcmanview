package rpcserver

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	sq "github.com/Masterminds/squirrel"
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

	groups, err := func() ([]*rpc.Group, error) {
		// SELECT ...
		sb := sq.
			Select(
				"groups.*",
				"COUNT(group_users.group_id) AS user_count",
			).
			From("groups").
			LeftJoin("group_users ON group_users.group_id = groups.id").
			GroupBy("groups.id")
		// ORDER BY
		switch req.Sort.GetField() {
		case "name":
			sb = sb.OrderBy(orderBySQL("name", req.Sort.GetOrder()))
		case "userCount":
			sb = sb.OrderBy(orderBySQL("user_count", req.Sort.GetOrder()))
		case "createdAt":
			sb = sb.OrderBy(orderBySQL("groups.created_at", req.Sort.GetOrder()))
		}
		// OFFSET ...
		sb = sb.
			Offset(uint64(page.Offset())).
			Limit(uint64(page.Limit()))

		rows, scanner, err := ssq.QueryRows(ctx, a.db, sb)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var groups []*rpc.Group
		for rows.Next() {
			var v struct {
				repo.Group
				UserCount int64
			}
			err := scanner.Scan(&v)
			if err != nil {
				return nil, err
			}

			groups = append(groups, &rpc.Group{
				Id:            v.ID,
				Name:          v.Name,
				Description:   v.Description,
				UserCount:     v.UserCount,
				CreatedAtTime: timestamppb.New(v.CreatedAt.Time),
				UpdatedAtTime: timestamppb.New(v.UpdatedAt.Time),
			})
		}

		return groups, nil
	}()
	if err != nil {
		return nil, NewError(err).Internal()
	}

	count, err := a.db.CountGroup(ctx)
	if err != nil {
		return nil, NewError(err).Internal()
	}

	return &rpc.ListGroupsResp{
		Items:      groups,
		PageResult: convertPagePaginationResult(page.Result(int(count))),
		Sort:       req.Sort,
	}, nil
}
