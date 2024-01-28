package rpcserver

import (
	"context"
	"net/url"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertCreateUpdateGroupError(msg string, err error) error {
	if errs, ok := asValidationErrors(err); ok {
		return NewError(err, msg).Validation(errs, [][2]string{
			{"name", "Name"},
			{"description", "Description"},
		})
	}

	if constraintErr, ok := asConstraintError(err); ok {
		return NewError(err, msg).Constraint(constraintErr, [][3]string{
			{"name", "groups.name", "Name already taken."},
		})
	}

	return check(err)
}

func convertCreateUpdateDeviceError(msg string, err error) error {
	if errs, ok := asValidationErrors(err); ok {
		return NewError(err, msg).Validation(errs, [][2]string{
			{"name", "Name"},
			{"description", "Description"},
			{"location", "Location"},
		})
	}

	if constraintErr, ok := asConstraintError(err); ok {
		return NewError(err, msg).Constraint(constraintErr, [][3]string{
			{"name", "dahua_devices.name", "Name already taken."},
			{"url", "dahua_devices.ip", "URL already taken."},
		})
	}

	return check(err)
}

func NewAdmin(db repo.DB, bus *core.Bus) *Admin {
	return &Admin{
		db:  db,
		bus: bus,
	}
}

type Admin struct {
	db  repo.DB
	bus *core.Bus
}

// ---------- Device

func (a *Admin) GetAdminDevicesPage(ctx context.Context, req *rpc.GetAdminDevicesPageReq) (*rpc.GetAdminDevicesPageResp, error) {
	page := parsePagePagination(req.Page)

	items, err := func() ([]*rpc.GetAdminDevicesPageResp_Device, error) {
		var row struct {
			repo.DahuaDevice
		}
		// SELECT ...
		sb := sq.
			Select(
				"dahua_devices.*",
			).
			From("dahua_devices")
		// ORDER BY
		switch req.Sort.GetField() {
		case "name":
			sb = sb.OrderBy(convertOrderToSQL("name", req.Sort.GetOrder()))
		case "createdAt":
			sb = sb.OrderBy(convertOrderToSQL("created_at", req.Sort.GetOrder()))
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

		var items []*rpc.GetAdminDevicesPageResp_Device
		for rows.Next() {
			err := scanner.Scan(&row)
			if err != nil {
				return nil, err
			}

			items = append(items, &rpc.GetAdminDevicesPageResp_Device{
				Id:             row.ID,
				Name:           row.Name,
				Url:            row.Url.String(),
				Username:       row.Username,
				Disabled:       row.DisabledAt.Valid,
				DisabledAtTime: timestamppb.New(row.DisabledAt.Time),
				CreatedAtTime:  timestamppb.New(row.CreatedAt.Time),
			})
		}

		return items, nil
	}()
	if err != nil {
		return nil, check(err)
	}

	count, err := func() (int64, error) {
		var row struct{ Count int64 }
		return row.Count, ssq.QueryOne(ctx, a.db, &row, sq.
			Select("COUNT(*) AS count").
			From("dahua_devices"))
	}()
	if err != nil {
		return nil, check(err)
	}

	return &rpc.GetAdminDevicesPageResp{
		Items:      items,
		PageResult: convertPagePaginationResult(page.Result(int(count))),
		Sort:       req.Sort,
	}, nil
}

func (a *Admin) GetAdminDevicesIDPage(ctx context.Context, req *rpc.GetAdminDevicesIDPageReq) (*rpc.GetAdminDevicesIDPageResp, error) {
	dbDevice, err := a.db.GetDahuaDevice(ctx, req.Id)
	if err != nil {
		return nil, check(err)
	}

	return &rpc.GetAdminDevicesIDPageResp{
		Device: &rpc.GetAdminDevicesIDPageResp_Device{
			Id:             dbDevice.ID,
			Name:           dbDevice.Name,
			Url:            dbDevice.Url.String(),
			Username:       dbDevice.Username,
			Disabled:       false,
			Location:       dbDevice.Location.String(),
			CreatedAtTime:  timestamppb.New(dbDevice.CreatedAt.Time),
			UpdatedAtTime:  timestamppb.New(dbDevice.UpdatedAt.Time),
			DisabledAtTime: timestamppb.New(dbDevice.DisabledAt.Time),
			Features:       []string{},
		},
	}, nil

}

func (*Admin) GetDevice(context.Context, *rpc.GetDeviceReq) (*rpc.GetDeviceResp, error) {
	return nil, errNotImplemented
}

func (a *Admin) CreateDevice(ctx context.Context, req *rpc.CreateDeviceReq) (*rpc.CreateDeviceResp, error) {
	urL, err := url.Parse(req.Url)
	if err != nil {
		return nil, NewError(nil, "URL is invalid.").Field("url")
	}
	loc, err := time.LoadLocation(req.Location)
	if err != nil {
		return nil, NewError(nil, "Location is invalid.").Field("location")
	}

	res, err := dahua.CreateDevice(ctx, a.db, a.bus, models.DahuaDevice{
		Name:     req.Name,
		Url:      urL,
		Username: req.Username,
		Password: req.Password,
		Location: loc,
		Feature:  dahua.FeatureFromStrings(req.Features),
	})
	if err != nil {
		return nil, convertCreateUpdateDeviceError("Failed to create device.", err)
	}

	return &rpc.CreateDeviceResp{
		Id: res.DahuaDevice.ID,
	}, nil
}

func (a *Admin) UpdateDevice(ctx context.Context, req *rpc.UpdateDeviceReq) (*emptypb.Empty, error) {
	return nil, errNotImplemented
}

func (a *Admin) DeleteDevice(ctx context.Context, req *rpc.DeleteDeviceReq) (*emptypb.Empty, error) {
	for _, id := range req.Ids {
		err := dahua.DeleteDevice(ctx, a.db, a.bus, id)
		if err != nil {
			return nil, check(err)
		}
	}
	return &emptypb.Empty{}, nil
}

func (*Admin) SetDeviceDisable(context.Context, *rpc.SetDeviceDisableReq) (*emptypb.Empty, error) {
	return nil, errNotImplemented
}

// ---------- User

func (a *Admin) GetAdminUsersPage(ctx context.Context, req *rpc.GetAdminUsersPageReq) (*rpc.GetAdminUsersPageResp, error) {
	page := parsePagePagination(req.Page)

	items, err := func() ([]*rpc.GetAdminUsersPageResp_User, error) {
		var row struct {
			repo.User
			Admin bool
		}
		// SELECT ...
		sb := sq.
			Select(
				"users.*",
				"admins.user_id IS NOT NULL as 'admin'",
			).
			From("users").
			LeftJoin("admins ON admins.user_id = users.id")
		// ORDER BY
		switch req.Sort.GetField() {
		case "username":
			sb = sb.OrderBy(convertOrderToSQL("username", req.Sort.GetOrder()))
		case "email":
			sb = sb.OrderBy(convertOrderToSQL("email", req.Sort.GetOrder()))
		case "createdAt":
			sb = sb.OrderBy(convertOrderToSQL("users.created_at", req.Sort.GetOrder()))
		default:
			sb = sb.OrderBy("admin DESC")
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

		var items []*rpc.GetAdminUsersPageResp_User
		for rows.Next() {
			err := scanner.Scan(&row)
			if err != nil {
				return nil, err
			}

			items = append(items, &rpc.GetAdminUsersPageResp_User{
				Id:             row.ID,
				Username:       row.Username,
				Email:          row.Email,
				Disabled:       row.DisabledAt.Valid,
				Admin:          row.Admin,
				DisabledAtTime: timestamppb.New(row.DisabledAt.Time),
				CreatedAtTime:  timestamppb.New(row.CreatedAt.Time),
			})
		}

		return items, nil
	}()
	if err != nil {
		return nil, check(err)
	}

	count, err := func() (int64, error) {
		var row struct{ Count int64 }
		return row.Count, ssq.QueryOne(ctx, a.db, &row, sq.
			Select("COUNT(*) AS count").
			From("users"))
	}()
	if err != nil {
		return nil, check(err)
	}

	return &rpc.GetAdminUsersPageResp{
		Items:      items,
		PageResult: convertPagePaginationResult(page.Result(int(count))),
		Sort:       req.Sort,
	}, nil

}

func (*Admin) DeleteUser(context.Context, *rpc.DeleteUserReq) (*emptypb.Empty, error) {
	return nil, errNotImplemented
}

func (a *Admin) SetUserDisable(ctx context.Context, req *rpc.SetUserDisableReq) (*emptypb.Empty, error) {
	authSession := useAuthSession(ctx)
	for _, item := range req.Items {
		if item.Id != authSession.UserID {
			err := auth.UpdateUserDisable(ctx, a.db, item.Id, item.Disable)
			if err != nil {
				return nil, check(err)
			}
		}
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) SetUserAdmin(ctx context.Context, req *rpc.SetUserAdminReq) (*emptypb.Empty, error) {
	authSession := useAuthSession(ctx)
	if req.Id == authSession.UserID {
		return nil, NewError(nil, "Cannot modify current user.").Field("item")
	}

	err := auth.UpdateUserAdmin(ctx, a.db, req.Id, req.Admin)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ResetUserPassword implements rpc.Admin.
func (*Admin) ResetUserPassword(context.Context, *rpc.ResetUserPasswordReq) (*emptypb.Empty, error) {
	return nil, errNotImplemented
}

// ---------- Group

func (a *Admin) GetAdminGroupsPage(ctx context.Context, req *rpc.GetAdminGroupsPageReq) (*rpc.GetAdminGroupsPageResp, error) {
	page := parsePagePagination(req.Page)

	items, err := func() ([]*rpc.GetAdminGroupsPageResp_Group, error) {
		var row struct {
			repo.Group
			UserCount int64
		}
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
			sb = sb.OrderBy(convertOrderToSQL("name", req.Sort.GetOrder()))
		case "userCount":
			sb = sb.OrderBy(convertOrderToSQL("user_count", req.Sort.GetOrder()))
		case "createdAt":
			sb = sb.OrderBy(convertOrderToSQL("groups.created_at", req.Sort.GetOrder()))
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

		var items []*rpc.GetAdminGroupsPageResp_Group
		for rows.Next() {
			err := scanner.Scan(&row)
			if err != nil {
				return nil, err
			}

			items = append(items, &rpc.GetAdminGroupsPageResp_Group{
				Id:             row.ID,
				Name:           row.Name,
				UserCount:      row.UserCount,
				Disabled:       row.DisabledAt.Valid,
				DisabledAtTime: timestamppb.New(row.DisabledAt.Time),
				CreatedAtTime:  timestamppb.New(row.CreatedAt.Time),
			})
		}

		return items, nil
	}()
	if err != nil {
		return nil, check(err)
	}

	count, err := func() (int64, error) {
		var row struct{ Count int64 }
		return row.Count, ssq.QueryOne(ctx, a.db, &row, sq.
			Select("COUNT(*) AS count").
			From("groups"))
	}()
	if err != nil {
		return nil, check(err)
	}

	return &rpc.GetAdminGroupsPageResp{
		Items:      items,
		PageResult: convertPagePaginationResult(page.Result(int(count))),
		Sort:       req.Sort,
	}, nil
}

func (a *Admin) GetAdminGroupsIDPage(ctx context.Context, req *rpc.GetAdminGroupsIDPageReq) (*rpc.GetAdminGroupsIDPageResp, error) {
	dbGroup, err := a.db.GetGroup(ctx, req.Id)
	if err != nil {
		return nil, check(err)
	}

	dbUsers, err := a.db.GetUserByGroup(ctx, req.Id)
	if err != nil {
		return nil, check(err)
	}

	users := make([]*rpc.GetAdminGroupsIDPageResp_User, 0, len(dbUsers))
	for _, v := range dbUsers {
		users = append(users, &rpc.GetAdminGroupsIDPageResp_User{
			Id:       v.ID,
			Username: v.Username,
		})
	}

	return &rpc.GetAdminGroupsIDPageResp{
		Group: &rpc.GetAdminGroupsIDPageResp_Group{
			Id:             dbGroup.ID,
			Name:           dbGroup.Name,
			Description:    dbGroup.Description,
			Disabled:       dbGroup.DisabledAt.Valid,
			DisabledAtTime: timestamppb.New(dbGroup.DisabledAt.Time),
			CreatedAtTime:  timestamppb.New(dbGroup.CreatedAt.Time),
			UpdatedAtTime:  timestamppb.New(dbGroup.UpdatedAt.Time),
		},
		Users: users,
	}, nil
}

func (a *Admin) GetGroup(ctx context.Context, req *rpc.GetGroupReq) (*rpc.GetGroupResp, error) {
	dbGroup, err := a.db.GetGroup(ctx, req.Id)
	if err != nil {
		return nil, check(err)
	}
	return &rpc.GetGroupResp{
		Model: &rpc.GroupModel{
			Id:          req.Id,
			Name:        dbGroup.Name,
			Description: dbGroup.Description,
		},
	}, nil
}

func (a *Admin) CreateGroup(ctx context.Context, req *rpc.CreateGroupReq) (*rpc.CreateGroupResp, error) {
	id, err := auth.CreateGroup(ctx, a.db, models.Group{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, convertCreateUpdateGroupError("Failed to create group.", err)
	}

	return &rpc.CreateGroupResp{
		Id: id,
	}, nil
}

func (a *Admin) UpdateGroup(ctx context.Context, req *rpc.UpdateGroupReq) (*emptypb.Empty, error) {
	dbGroup, err := a.db.GetGroup(ctx, req.Model.GetId())
	if err != nil {
		return nil, check(err)
	}
	group := dbGroup.Convert()

	group.Name = req.Model.GetName()
	group.Description = req.Model.GetDescription()

	_, err = auth.UpdateGroup(ctx, a.db, group)
	if err != nil {
		return nil, convertCreateUpdateGroupError("Failed to update group.", err)
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) DeleteGroup(ctx context.Context, req *rpc.DeleteGroupReq) (*emptypb.Empty, error) {
	for _, id := range req.Ids {
		err := auth.DeleteGroup(ctx, a.db, id)
		if err != nil {
			return nil, check(err)
		}
	}
	return &emptypb.Empty{}, nil
}

func (a *Admin) SetGroupDisable(ctx context.Context, req *rpc.SetGroupDisableReq) (*emptypb.Empty, error) {
	for _, item := range req.Items {
		err := auth.UpdateGroupDisable(ctx, a.db, item.Id, item.Disable)
		if err != nil {
			return nil, check(err)
		}
	}
	return &emptypb.Empty{}, nil
}

func (*Admin) ListLocations(context.Context, *emptypb.Empty) (*rpc.ListLocationsResp, error) {
	return &rpc.ListLocationsResp{
		Locations: core.Locations,
	}, nil
}

var listDeviceFeaturesResp *rpc.ListDeviceFeaturesResp

func init() {
	res := make([]*rpc.ListDeviceFeaturesResp_Item, 0, len(dahua.FeatureList))
	for _, v := range dahua.FeatureList {
		res = append(res, &rpc.ListDeviceFeaturesResp_Item{
			Name:        v.Name,
			Value:       v.Value,
			Description: v.Description,
		})
	}
	listDeviceFeaturesResp = &rpc.ListDeviceFeaturesResp{Features: res}
}

func (*Admin) ListDeviceFeatures(context.Context, *emptypb.Empty) (*rpc.ListDeviceFeaturesResp, error) {
	return listDeviceFeaturesResp, nil
}
