package rpcserver

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewAdmin(db repo.DB, bus *event.Bus) *Admin {
	return &Admin{
		db:  db,
		bus: bus,
	}
}

type Admin struct {
	db  repo.DB
	bus *event.Bus
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
			sb = sb.OrderBy(parseOrderSQL("name", req.Sort.GetOrder()))
		case "url":
			sb = sb.OrderBy(parseOrderSQL("url", req.Sort.GetOrder()))
		case "createdAt":
			sb = sb.OrderBy(parseOrderSQL("created_at", req.Sort.GetOrder()))
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
		return nil, err
	}

	count, err := func() (int64, error) {
		var row struct{ Count int64 }
		return row.Count, ssq.QueryOne(ctx, a.db, &row, sq.
			Select("COUNT(*) AS count").
			From("dahua_devices"))
	}()
	if err != nil {
		return nil, err
	}

	return &rpc.GetAdminDevicesPageResp{
		Items:      items,
		PageResult: parsePagePaginationResult(page.Result(int(count))),
		Sort:       req.Sort,
	}, nil
}

func (a *Admin) GetAdminDevicesIDPage(ctx context.Context, req *rpc.GetAdminDevicesIDPageReq) (*rpc.GetAdminDevicesIDPageResp, error) {
	v, err := a.db.DahuaGetFatDevice(ctx, repo.DahuaFatDeviceParams{IDs: []int64{req.Id}})
	if err != nil {
		return nil, err
	}

	return &rpc.GetAdminDevicesIDPageResp{
		Device: &rpc.GetAdminDevicesIDPageResp_Device{
			Id:             v.DahuaDevice.ID,
			Name:           v.DahuaDevice.Name,
			Url:            v.DahuaDevice.Url.String(),
			Username:       v.DahuaDevice.Username,
			Disabled:       v.DisabledAt.Valid,
			Location:       v.DahuaDevice.Location.String(),
			CreatedAtTime:  timestamppb.New(v.DahuaDevice.CreatedAt.Time),
			UpdatedAtTime:  timestamppb.New(v.DahuaDevice.UpdatedAt.Time),
			DisabledAtTime: timestamppb.New(v.DahuaDevice.DisabledAt.Time),
			Features:       dahua.FeatureToStrings(v.Feature),
		},
	}, nil

}

func (a *Admin) GetDevice(ctx context.Context, req *rpc.GetDeviceReq) (*rpc.GetDeviceResp, error) {
	v, err := a.db.DahuaGetFatDevice(ctx, repo.DahuaFatDeviceParams{IDs: []int64{req.Id}})
	if err != nil {
		return nil, err
	}

	return &rpc.GetDeviceResp{
		Id:       v.ID,
		Name:     v.Name,
		Url:      v.Url.String(),
		Username: v.Username,
		Location: v.Location.String(),
		Features: dahua.FeatureToStrings(v.Feature),
	}, nil
}

func (a *Admin) CreateDevice(ctx context.Context, req *rpc.CreateDeviceReq) (*rpc.CreateDeviceResp, error) {
	urL, err := url.Parse(req.Url)
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation(req.Location)
	if err != nil {
		return nil, err
	}

	id, err := dahua.CreateDevice(ctx, a.db, a.bus, dahua.CreateDeviceParams{
		Name:     req.Name,
		URL:      urL,
		Username: req.Username,
		Password: req.Password,
		Location: loc,
		Feature:  dahua.FeatureFromStrings(req.Features),
	})
	if err != nil {
		return nil, err
	}

	return &rpc.CreateDeviceResp{
		Id: id,
	}, nil
}

func (a *Admin) UpdateDevice(ctx context.Context, req *rpc.UpdateDeviceReq) (*emptypb.Empty, error) {
	urL, err := url.Parse(req.Url)
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation(req.Location)
	if err != nil {
		return nil, err
	}

	dbDevice, err := a.db.DahuaGetFatDevice(ctx, repo.DahuaFatDeviceParams{IDs: []int64{req.Id}})
	if err != nil {
		return nil, err
	}

	err = dahua.UpdateDevice(ctx, a.db, a.bus, dbDevice.DahuaDevice, dahua.UpdateDeviceParams{
		Name:        req.Name,
		URL:         urL,
		Username:    req.Username,
		NewPassword: req.NewPassword,
		Location:    loc,
		Feature:     dahua.FeatureFromStrings(req.Features),
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) DeleteDevice(ctx context.Context, req *rpc.DeleteDeviceReq) (*emptypb.Empty, error) {
	for _, id := range req.Ids {
		err := dahua.DeleteDevice(ctx, a.db, a.bus, id)
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}

func (a *Admin) SetDeviceDisable(ctx context.Context, req *rpc.SetDeviceDisableReq) (*emptypb.Empty, error) {
	for _, v := range req.Items {
		err := dahua.UpdateDeviceDisabled(ctx, a.db, a.bus, v.Id, v.Disable)
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
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
			sb = sb.OrderBy(parseOrderSQL("username", req.Sort.GetOrder()))
		case "email":
			sb = sb.OrderBy(parseOrderSQL("email", req.Sort.GetOrder()))
		case "createdAt":
			sb = sb.OrderBy(parseOrderSQL("users.created_at", req.Sort.GetOrder()))
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
		return nil, err
	}

	count, err := func() (int64, error) {
		var row struct{ Count int64 }
		return row.Count, ssq.QueryOne(ctx, a.db, &row, sq.
			Select("COUNT(*) AS count").
			From("users"))
	}()
	if err != nil {
		return nil, err
	}

	return &rpc.GetAdminUsersPageResp{
		Items:      items,
		PageResult: parsePagePaginationResult(page.Result(int(count))),
		Sort:       req.Sort,
	}, nil

}

func (a *Admin) CreateUser(ctx context.Context, req *rpc.CreateUserReq) (*emptypb.Empty, error) {
	_, err := auth.CreateUser(ctx, a.db, auth.CreateUserParams{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Admin:    req.Admin,
		Disabled: req.Disabled,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) GetUser(ctx context.Context, req *rpc.GetUserReq) (*rpc.GetUserResp, error) {
	v, err := a.db.AuthGetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &rpc.GetUserResp{
		Id:       v.ID,
		Username: v.Username,
		Email:    v.Email,
	}, nil
}

func (a *Admin) UpdateUser(ctx context.Context, req *rpc.UpdateUserReq) (*emptypb.Empty, error) {
	dbUser, err := a.db.AuthGetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	err = auth.UpdateUser(ctx, a.db, dbUser, auth.UpdateUserParams{
		Email:    req.Email,
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) DeleteUser(ctx context.Context, req *rpc.DeleteUserReq) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)
	for _, id := range req.Ids {
		if id != session.UserID {
			err := auth.DeleteUser(ctx, a.db, id)
			if err != nil {
				return nil, err
			}
		}
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) SetUserDisable(ctx context.Context, req *rpc.SetUserDisableReq) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)
	for _, item := range req.Items {
		if item.Id != session.UserID {
			err := auth.UpdateUserDisabled(ctx, a.db, item.Id, item.Disable)
			if err != nil {
				return nil, err
			}
		}
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) SetUserAdmin(ctx context.Context, req *rpc.SetUserAdminReq) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)
	if req.Id == session.UserID {
		return nil, fmt.Errorf("Cannot modify current user.")
	}

	err := auth.UpdateUserAdmin(ctx, a.db, req.Id, req.Admin)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) ResetUserPassword(ctx context.Context, req *rpc.ResetUserPasswordReq) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)

	dbUser, err := a.db.AuthGetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if err := auth.UpdateUserPassword(ctx, a.db, dbUser, auth.UpdateUserPasswordParams{
		NewPassword:      req.NewPassword,
		CurrentSessionID: session.SessionID,
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
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
			sb = sb.OrderBy(parseOrderSQL("name", req.Sort.GetOrder()))
		case "userCount":
			sb = sb.OrderBy(parseOrderSQL("user_count", req.Sort.GetOrder()))
		case "createdAt":
			sb = sb.OrderBy(parseOrderSQL("groups.created_at", req.Sort.GetOrder()))
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
		return nil, err
	}

	count, err := func() (int64, error) {
		var row struct{ Count int64 }
		return row.Count, ssq.QueryOne(ctx, a.db, &row, sq.
			Select("COUNT(*) AS count").
			From("groups"))
	}()
	if err != nil {
		return nil, err
	}

	return &rpc.GetAdminGroupsPageResp{
		Items:      items,
		PageResult: parsePagePaginationResult(page.Result(int(count))),
		Sort:       req.Sort,
	}, nil
}

func (a *Admin) GetAdminGroupsIDPage(ctx context.Context, req *rpc.GetAdminGroupsIDPageReq) (*rpc.GetAdminGroupsIDPageResp, error) {
	dbUsers, err := a.db.AuthListUsersByGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	users := make([]*rpc.GetAdminGroupsIDPageResp_User, 0, len(dbUsers))
	for _, v := range dbUsers {
		users = append(users, &rpc.GetAdminGroupsIDPageResp_User{
			Id:       v.ID,
			Username: v.Username,
		})
	}

	v, err := a.db.AuthGetGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &rpc.GetAdminGroupsIDPageResp{
		Group: &rpc.GetAdminGroupsIDPageResp_Group{
			Id:             v.ID,
			Name:           v.Name,
			Description:    v.Description,
			Disabled:       v.DisabledAt.Valid,
			DisabledAtTime: timestamppb.New(v.DisabledAt.Time),
			CreatedAtTime:  timestamppb.New(v.CreatedAt.Time),
			UpdatedAtTime:  timestamppb.New(v.UpdatedAt.Time),
		},
		Users: users,
	}, nil
}

func (a *Admin) GetGroup(ctx context.Context, req *rpc.GetGroupReq) (*rpc.GetGroupResp, error) {
	v, err := a.db.AuthGetGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &rpc.GetGroupResp{
		Id:          v.ID,
		Name:        v.Name,
		Description: v.Description,
	}, nil
}

func (a *Admin) CreateGroup(ctx context.Context, req *rpc.CreateGroupReq) (*rpc.CreateGroupResp, error) {
	id, err := auth.CreateGroup(ctx, a.db, auth.CreateGroupParams{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &rpc.CreateGroupResp{
		Id: id,
	}, nil
}

func (a *Admin) UpdateGroup(ctx context.Context, req *rpc.UpdateGroupReq) (*emptypb.Empty, error) {
	dbGroup, err := a.db.AuthGetGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	err = auth.UpdateGroup(ctx, a.db, dbGroup, auth.UpdateGroupParams{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) DeleteGroup(ctx context.Context, req *rpc.DeleteGroupReq) (*emptypb.Empty, error) {
	for _, id := range req.Ids {
		err := auth.DeleteGroup(ctx, a.db, id)
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}

func (a *Admin) SetGroupDisable(ctx context.Context, req *rpc.SetGroupDisableReq) (*emptypb.Empty, error) {
	for _, item := range req.Items {
		err := auth.UpdateGroupDisable(ctx, a.db, item.Id, item.Disable)
		if err != nil {
			return nil, err
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
	features := make([]*rpc.ListDeviceFeaturesResp_Item, 0, len(dahua.FeatureList))
	for _, v := range dahua.FeatureList {
		features = append(features, &rpc.ListDeviceFeaturesResp_Item{
			Name:        v.Name,
			Value:       v.Value,
			Description: v.Description,
		})
	}
	listDeviceFeaturesResp = &rpc.ListDeviceFeaturesResp{Features: features}
}

func (*Admin) ListDeviceFeatures(context.Context, *emptypb.Empty) (*rpc.ListDeviceFeaturesResp, error) {
	return listDeviceFeaturesResp, nil
}
