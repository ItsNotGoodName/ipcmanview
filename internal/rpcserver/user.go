package rpcserver

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/build"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewUser(db sqlite.DB, dahuaStore *dahua.Store) *User {
	return &User{
		db:         db,
		dahuaStore: dahuaStore,
	}
}

type User struct {
	db         sqlite.DB
	dahuaStore *dahua.Store
}

func (u *User) GetHomePage(ctx context.Context, _ *emptypb.Empty) (*rpc.GetHomePageResp, error) {
	dbDevices, err := dahua.ListDevices(ctx, u.db)
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

	fileCount, err := dahua.CountFiles(ctx, u.db)
	if err != nil {
		return nil, err
	}

	eventCount, err := dahua.CountEvents(ctx, u.db)
	if err != nil {
		return nil, err
	}

	emailCount, err := dahua.CountEmails(ctx, u.db)
	if err != nil {
		return nil, err
	}

	latestFilesDTO, err := dahua.ListLatestFiles(ctx, u.db, 8)
	if err != nil {
		return nil, err
	}

	files := make([]*rpc.GetHomePageResp_File, 0, len(latestFilesDTO))
	for _, v := range latestFilesDTO {
		var thumbnailURL string
		if v.Type == models.DahuaFileTypeJPG {
			thumbnailURL = api.DahuaDeviceFileURI(v.DeviceID, v.FilePath)
		}
		files = append(files, &rpc.GetHomePageResp_File{
			Id:           v.ID,
			Url:          api.DahuaDeviceFileURI(v.DeviceID, v.FilePath),
			ThumbnailUrl: thumbnailURL,
			Type:         v.Type,
			StartTime:    timestamppb.New(v.StartTime.Time),
		})
	}

	latestEmailsDTO, err := dahua.ListLatestEmails(ctx, u.db, 5)
	if err != nil {
		return nil, err
	}

	emails := make([]*rpc.GetHomePageResp_Email, 0, len(latestEmailsDTO))
	for _, v := range latestEmailsDTO {
		emails = append(emails, &rpc.GetHomePageResp_Email{
			Id:              v.DahuaEmailMessage.ID,
			Subject:         v.DahuaEmailMessage.Subject,
			AttachmentCount: int32(v.AttachmentCount),
			CreatedAtTime:   timestamppb.New(v.DahuaEmailMessage.CreatedAt.Time),
		})
	}

	build := &rpc.GetHomePageResp_Build{
		Commit:     build.Current.Commit,
		Version:    build.Current.Version,
		Date:       timestamppb.New(build.Current.Date),
		RepoUrl:    build.Current.RepoURL,
		CommitUrl:  build.Current.CommitURL,
		LicenseUrl: build.Current.LicenseURL,
		ReleaseUrl: build.Current.ReleaseURL,
	}

	return &rpc.GetHomePageResp{
		Devices:    devices,
		FileCount:  fileCount,
		EventCount: eventCount,
		EmailCount: emailCount,
		Build:      build,
		Files:      files,
		Emails:     emails,
	}, nil
}

func (u *User) GetDevicesPage(ctx context.Context, req *rpc.GetDevicesPageReq) (*rpc.GetDevicesPageResp, error) {
	dbDevices, err := dahua.ListDevices(ctx, u.db)
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

	user, err := u.db.C().AuthGetUser(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	dbSessions, err := u.db.C().AuthListUserSessionsForUserAndNotExpired(ctx, repo.AuthListUserSessionsForUserAndNotExpiredParams{
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

	dbGroups, err := u.db.C().AuthListGroupsForUser(ctx, session.UserID)
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

func (u *User) GetEmailsPage(ctx context.Context, req *rpc.GetEmailsPageReq) (*rpc.GetEmailsPageResp, error) {
	page := parsePagePagination(req.Page)
	sort := parseSort(req.Sort).withDefaultOrder(rpc.Order_DESC)

	v, err := dahua.ListEmails(ctx, u.db, dahua.ListEmailsParams{
		Page:      page,
		Ascending: sort.Order == rpc.Order_ASC,
	})
	if err != nil {
		return nil, err
	}

	var emails []*rpc.GetEmailsPageResp_Email
	for _, v := range v.Items {
		emails = append(emails, &rpc.GetEmailsPageResp_Email{
			Id:              v.ID,
			DeviceId:        v.DeviceID,
			DeviceName:      v.DeviceName,
			From:            v.From,
			Subject:         v.Subject,
			AlarmEvent:      v.AlarmEvent,
			AttachmentCount: int32(v.AttachmentCount),
			CreatedAtTime:   timestamppb.New(v.CreatedAt.Time),
		})
	}

	return &rpc.GetEmailsPageResp{
		Emails:     emails,
		PageResult: encodePagePaginationResult(v.PageResult),
		Sort:       sort.encode(),
	}, nil
}

func (u *User) GetEmailsIDPage(ctx context.Context, req *rpc.GetEmailsIDPageReq) (*rpc.GetEmailsIDPageResp, error) {
	email, err := dahua.GetEmail(ctx, u.db, req.Id)
	if err != nil {
		return nil, err
	}

	emailAround, err := dahua.GetEmailAround(ctx, u.db, req.Id)
	if err != nil {
		return nil, err
	}

	attachments := make([]*rpc.GetEmailsIDPageResp_Attachment, 0, len(email.Attachments))
	for _, v := range email.Attachments {
		attachments = append(attachments, &rpc.GetEmailsIDPageResp_Attachment{
			Id:           v.DahuaEmailAttachment.ID,
			Name:         v.DahuaEmailAttachment.FileName,
			Url:          api.DahuaAferoFileURI(v.DahuaAferoFile.Name),
			ThumbnailUrl: api.DahuaAferoFileURI(v.DahuaAferoFile.Name),
			Size:         v.DahuaAferoFile.Size,
		})
	}

	emailCount, err := dahua.CountEmails(ctx, u.db)
	if err != nil {
		return nil, err
	}

	return &rpc.GetEmailsIDPageResp{
		Id:              email.Message.ID,
		DeviceId:        email.Message.DeviceID,
		From:            email.Message.From,
		Subject:         email.Message.Subject,
		To:              email.Message.To.Slice,
		Date:            timestamppb.New(email.Message.Date.Time),
		CreatedAtTime:   timestamppb.New(email.Message.CreatedAt.Time),
		Text:            email.Message.Text,
		Attachments:     attachments,
		NextEmailId:     email.NextEmailID,
		PreviousEmailId: emailAround.PreviousEmailID,
		EmailSeen:       emailAround.EmailSeen,
		EmailCount:      emailCount,
	}, nil
}

func (u *User) GetEventsPage(ctx context.Context, req *rpc.GetEventsPageReq) (*rpc.GetEventsPageResp, error) {
	page := parsePagePagination(req.Page)
	sort := parseSort(req.Sort).withDefaultOrder(rpc.Order_DESC)

	v, err := dahua.ListEvents(ctx, u.db, dahua.ListEventsParams{
		Page:      page,
		Ascending: sort.Order == rpc.Order_ASC,
	})
	if err != nil {
		return nil, err
	}

	var events []*rpc.GetEventsPageResp_Event
	for _, v := range v.Items {
		events = append(events, &rpc.GetEventsPageResp_Event{
			Id:            v.ID,
			DeviceId:      v.DeviceID,
			DeviceName:    v.DeviceName,
			Code:          v.Code,
			Action:        v.Action,
			Index:         v.Index,
			Data:          string(v.Data.RawMessage),
			CreatedAtTime: timestamppb.New(v.CreatedAt.Time),
		})
	}

	return &rpc.GetEventsPageResp{
		Events:     events,
		PageResult: encodePagePaginationResult(v.PageResult),
		Sort:       sort.encode(),
	}, nil
}

func (u *User) UpdateMyPassword(ctx context.Context, req *rpc.UpdateMyPasswordReq) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)

	dbUser, err := u.db.C().AuthGetUser(ctx, session.UserID)
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

	dbUser, err := u.db.C().AuthGetUser(ctx, session.UserID)
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
	client, err := u.dahuaStore.GetClient(ctx, req.Id)
	if err != nil {
		return nil, err
	}

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
	client, err := u.dahuaStore.GetClient(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	v, err := dahua.GetLicenseList(ctx, client.RPC)
	if err != nil {
		return nil, err
	}

	items := make([]*rpc.ListDeviceLicensesResp_License, 0, len(v))
	for _, v := range v {
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
	client, err := u.dahuaStore.GetClient(ctx, req.Id)
	if err != nil {
		return nil, err
	}

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
	client, err := u.dahuaStore.GetClient(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	v := dahua.GetRPCStatus(ctx, client.RPC)

	return &rpc.GetDeviceRPCStatusResp{
		Error:         v.Error,
		State:         v.State,
		LastLoginTime: timestamppb.New(v.LastLogin),
	}, nil
}

func (u *User) ListDeviceStorage(ctx context.Context, req *rpc.ListDeviceStorageReq) (*rpc.ListDeviceStorageResp, error) {
	client, err := u.dahuaStore.GetClient(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	v, err := dahua.GetStorage(ctx, client.RPC)
	if err != nil {
		return nil, err
	}

	items := make([]*rpc.ListDeviceStorageResp_Storage, 0, len(v))
	for _, v := range v {
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
