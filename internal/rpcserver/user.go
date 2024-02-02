package rpcserver

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewUser(db repo.DB, dahuaStore *dahua.Store) *User {
	return &User{
		db:         db,
		dahuaStore: dahuaStore,
	}
}

type User struct {
	db         repo.DB
	dahuaStore *dahua.Store
}

func (u *User) GetHomePage(ctx context.Context, _ *emptypb.Empty) (*rpc.GetHomePageResp, error) {
	dbDevices, err := u.db.DahuaListFatDevices(ctx)
	if err != nil {
		return nil, err
	}

	devices := make([]*rpc.GetHomePageResp_Device, 0, len(dbDevices))
	for _, v := range dbDevices {
		devices = append(devices, &rpc.GetHomePageResp_Device{
			Id:   v.ID,
			Name: v.Name,
		})
	}

	fileCount, err := u.db.DahuaCountDahuaFiles(ctx)
	if err != nil {
		return nil, err
	}

	eventCount, err := u.db.DahuaCountDahuaEvents(ctx)
	if err != nil {
		return nil, err
	}

	emailCount, err := u.db.DahuaCountDahuaEmails(ctx)
	if err != nil {
		return nil, err
	}

	return &rpc.GetHomePageResp{
		Devices:    devices,
		FileCount:  fileCount,
		EventCount: eventCount,
		EmailCount: emailCount,
	}, nil
}

func (u *User) GetDevicesPage(ctx context.Context, req *rpc.GetDevicesPageReq) (*rpc.GetDevicesPageResp, error) {
	dbDevices, err := u.db.DahuaListFatDevices(ctx)
	if err != nil {
		return nil, err
	}

	devices := make([]*rpc.GetDevicesPageResp_Device, 0, len(dbDevices))
	for _, v := range dbDevices {
		devices = append(devices, &rpc.GetDevicesPageResp_Device{
			Id:            v.ID,
			Name:          v.Name,
			Url:           v.Url.String(),
			Username:      v.Username,
			CreatedAtTime: timestamppb.New(v.CreatedAt.Time),
			Disabled:      v.DisabledAt.Valid,
		})
	}

	return &rpc.GetDevicesPageResp{
		Devices: devices,
	}, nil
}

func (u *User) GetProfilePage(ctx context.Context, _ *emptypb.Empty) (*rpc.GetProfilePageResp, error) {
	session := useAuthSession(ctx)

	user, err := u.db.AuthGetUser(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	dbSessions, err := u.db.AuthListUserSessionsForUserAndNotExpired(ctx, repo.AuthListUserSessionsForUserAndNotExpiredParams{
		UserID: session.UserID,
		Now:    types.NewTime(time.Now()),
	})
	if err != nil {
		return nil, err
	}

	activeCutoff := time.Now().Add(-24 * time.Hour)
	sessions := make([]*rpc.GetProfilePageResp_Session, 0, len(dbSessions))
	for _, v := range dbSessions {
		sessions = append(sessions, &rpc.GetProfilePageResp_Session{
			Id:             v.ID,
			UserAgent:      v.UserAgent,
			Ip:             v.Ip,
			LastIp:         v.LastIp,
			LastUsedAtTime: timestamppb.New(v.LastUsedAt.Time),
			CreatedAtTime:  timestamppb.New(v.CreatedAt.Time),
			Active:         v.LastUsedAt.After(activeCutoff),
			Current:        v.ID == session.SessionID,
		})
	}

	dbGroups, err := u.db.AuthListGroupsForUser(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	groups := make([]*rpc.GetProfilePageResp_Group, 0, len(dbGroups))
	for _, v := range dbGroups {
		groups = append(groups, &rpc.GetProfilePageResp_Group{
			Id:           v.ID,
			Name:         v.Name,
			Description:  v.Description,
			JoinedAtTime: timestamppb.New(v.JoinedAt.Time),
		})
	}

	return &rpc.GetProfilePageResp{
		Username:      user.Username,
		Email:         user.Email,
		Admin:         session.Admin,
		CreatedAtTime: timestamppb.New(user.CreatedAt.Time),
		UpdatedAtTime: timestamppb.New(user.UpdatedAt.Time),
		Sessions:      sessions,
		Groups:        groups,
	}, nil
}

func (u *User) UpdateMyPassword(ctx context.Context, req *rpc.UpdateMyPasswordReq) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)

	dbUser, err := u.db.AuthGetUser(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	if err := auth.CheckUserPassword(dbUser.Password, req.OldPassword); err != nil {
		return nil, err
	}

	if err := auth.UpdateUserPassword(ctx, u.db, dbUser, auth.UpdateUserPasswordParams{
		NewPassword:      req.NewPassword,
		CurrentSessionID: session.SessionID,
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (u *User) UpdateMyUsername(ctx context.Context, req *rpc.UpdateMyUsernameReq) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)

	dbUser, err := u.db.AuthGetUser(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	if err := auth.UpdateUserUsername(ctx, u.db, dbUser, req.NewUsername); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (u *User) RevokeAllMySessions(ctx context.Context, rCreateUpdateGroupeq *emptypb.Empty) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)

	err := auth.DeleteOtherUserSessions(ctx, u.db, session.UserID, session.SessionID)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (u *User) RevokeMySession(ctx context.Context, req *rpc.RevokeMySessionReq) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)

	if err := auth.DeleteUserSession(ctx, u.db, session.UserID, req.SessionId); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (u *User) GetDeviceDetail(ctx context.Context, req *rpc.GetDeviceDetailReq) (*rpc.GetDeviceDetailResp, error) {
	device, err := u.db.DahuaGetFatDevice(ctx, repo.DahuaFatDeviceParams{IDs: []int64{req.Id}})
	if err != nil {
		return nil, err
	}
	client := u.dahuaStore.Client(ctx, dahua.ConnFrom(device))

	v, err := dahua.GetDahuaDetail(ctx, client.RPC)
	if err != nil {
		return nil, err
	}

	return &rpc.GetDeviceDetailResp{
		Sn:               v.SN,
		DeviceClass:      v.DeviceClass,
		DeviceType:       v.DeviceType,
		HardwareVersion:  v.HardwareVersion,
		MarketArea:       v.MarketArea,
		ProcessInfo:      v.ProcessInfo,
		Vendor:           v.Vendor,
		OnvifVersion:     v.OnvifVersion,
		AlgorithmVersion: v.AlgorithmVersion,
	}, nil
}

func (u *User) ListDeviceLicenses(ctx context.Context, req *rpc.ListDeviceLicensesReq) (*rpc.ListDeviceLicensesResp, error) {
	device, err := u.db.DahuaGetFatDevice(ctx, repo.DahuaFatDeviceParams{IDs: []int64{req.Id}})
	if err != nil {
		return nil, err
	}
	client := u.dahuaStore.Client(ctx, dahua.ConnFrom(device))

	vv, err := dahua.GetLicenseList(ctx, client.RPC)
	if err != nil {
		return nil, err
	}

	items := make([]*rpc.ListDeviceLicensesResp_License, 0, len(vv))
	for _, v := range vv {
		items = append(items, &rpc.ListDeviceLicensesResp_License{
			AbroadInfo:    v.AbroadInfo,
			AllType:       v.AllType,
			DigitChannel:  int32(v.DigitChannel),
			EffectiveDays: int32(v.EffectiveDays),
			EffectiveTime: timestamppb.New(v.EffectiveTime),
			LicenseId:     int32(v.LicenseID),
			ProductType:   v.ProductType,
			Status:        int32(v.Status),
			Username:      v.Username,
		})
	}

	return &rpc.ListDeviceLicensesResp{
		Items: items,
	}, nil
}

func (u *User) GetDeviceSoftwareVersion(ctx context.Context, req *rpc.GetDeviceSoftwareVersionReq) (*rpc.GetDeviceSoftwareVersionResp, error) {
	device, err := u.db.DahuaGetFatDevice(ctx, repo.DahuaFatDeviceParams{IDs: []int64{req.Id}})
	if err != nil {
		return nil, err
	}
	client := u.dahuaStore.Client(ctx, dahua.ConnFrom(device))

	v, err := dahua.GetSoftwareVersion(ctx, client.RPC)
	if err != nil {
		return nil, err
	}

	return &rpc.GetDeviceSoftwareVersionResp{
		Build:                   v.Build,
		BuildDate:               v.BuildDate,
		SecurityBaseLineVersion: v.SecurityBaseLineVersion,
		Version:                 v.Version,
		WebVersion:              v.WebVersion,
	}, nil
}

func (u *User) GetDeviceRPCStatus(ctx context.Context, req *rpc.GetDeviceRPCStatusReq) (*rpc.GetDeviceRPCStatusResp, error) {
	device, err := u.db.DahuaGetFatDevice(ctx, repo.DahuaFatDeviceParams{IDs: []int64{req.Id}})
	if err != nil {
		return nil, err
	}
	client := u.dahuaStore.Client(ctx, dahua.ConnFrom(device))

	v := dahua.GetRPCStatus(ctx, client.RPC)

	return &rpc.GetDeviceRPCStatusResp{
		Error:         v.Error,
		State:         v.State,
		LastLoginTime: timestamppb.New(v.LastLogin),
	}, nil
}

func (u *User) ListDeviceStorage(ctx context.Context, req *rpc.ListDeviceStorageReq) (*rpc.ListDeviceStorageResp, error) {
	device, err := u.db.DahuaGetFatDevice(ctx, repo.DahuaFatDeviceParams{IDs: []int64{req.Id}})
	if err != nil {
		return nil, err
	}
	client := u.dahuaStore.Client(ctx, dahua.ConnFrom(device))

	vv, err := dahua.GetStorage(ctx, client.RPC)
	if err != nil {
		return nil, err
	}

	items := make([]*rpc.ListDeviceStorageResp_Storage, 0, len(vv))
	for _, v := range vv {
		items = append(items, &rpc.ListDeviceStorageResp_Storage{
			Name:       v.Name,
			State:      v.State,
			Path:       v.Path,
			Type:       v.Type,
			TotalBytes: v.TotalBytes,
			UsedBytes:  v.UsedBytes,
			IsError:    v.IsError,
		})
	}

	return &rpc.ListDeviceStorageResp{
		Items: items,
	}, nil
}
