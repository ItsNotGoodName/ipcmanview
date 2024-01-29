package repo

import "github.com/ItsNotGoodName/ipcmanview/internal/models"

// this is stupid

func (c DahuaFile) Convert() models.DahuaFile {
	return models.DahuaFile{
		ID:          c.ID,
		DeviceID:    c.DeviceID,
		Channel:     int(c.Channel),
		StartTime:   c.StartTime.Time,
		EndTime:     c.EndTime.Time,
		Length:      int(c.Length),
		Type:        c.Type,
		FilePath:    c.FilePath,
		Duration:    int(c.Duration),
		Disk:        int(c.Disk),
		VideoStream: c.VideoStream,
		Flags:       c.Flags.Slice,
		Events:      c.Events.Slice,
		Cluster:     int(c.Cluster),
		Partition:   int(c.Partition),
		PicIndex:    int(c.PicIndex),
		Repeat:      int(c.Repeat),
		WorkDir:     c.WorkDir,
		WorkDirSN:   c.WorkDirSn,
		Storage:     c.Storage,
	}
}

func (c DahuaStream) Convert(embedURL string) models.DahuaStream {
	return models.DahuaStream{
		ID:           c.ID,
		DeviceID:     c.DeviceID,
		Name:         c.Name,
		Channel:      int(c.Channel),
		Subtype:      int(c.Subtype),
		MediamtxPath: c.MediamtxPath,
		EmbedURL:     embedURL,
	}
}

func (c DahuaStorageDestination) Convert() models.DahuaStorageDestination {
	return models.DahuaStorageDestination{
		ID:              c.ID,
		Name:            c.Name,
		Storage:         c.Storage,
		ServerAddress:   c.ServerAddress,
		Port:            c.Port,
		Username:        c.Username,
		Password:        c.Password,
		RemoteDirectory: c.RemoteDirectory,
	}
}

func (c User) Convert() models.User {
	return models.User{
		ID:       c.ID,
		Email:    c.Email,
		Username: c.Username,
		Password: c.Password,
	}
}

func (c Group) Convert() models.Group {
	return models.Group{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
	}
}
