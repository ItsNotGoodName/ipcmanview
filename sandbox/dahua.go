package sandbox

import (
	"context"
	"fmt"
	"os"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/internal/dahua"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/modules/global"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/modules/license"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/modules/magicbox"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/modules/storage"
)

func Dahua(ctx context.Context) {
	username, _ := os.LookupEnv("IPC_USERNAME")
	password, _ := os.LookupEnv("IPC_PASSWORD")
	ip, _ := os.LookupEnv("IPC_IP")

	cam := core.DahuaCamera{
		Address:  ip,
		Username: username,
		Password: password,
	}

	c := dahua.NewCameraActor(cam)
	defer c.Close(ctx)

	fmt.Println(global.GetCurrentTime(ctx, c))
	fmt.Println(magicbox.NeedReboot(ctx, c))
	fmt.Println(magicbox.GetCPUUsage(ctx, c))
	fmt.Println(magicbox.GetDeviceClass(ctx, c))
	fmt.Println(magicbox.GetDeviceType(ctx, c))
	fmt.Println(magicbox.GetHardwareVersion(ctx, c))
	fmt.Println(magicbox.GetMarketArea(ctx, c))
	fmt.Println(magicbox.GetMemoryInfo(ctx, c))
	fmt.Println(magicbox.GetProcessInfo(ctx, c))
	fmt.Println(magicbox.GetSerialNo(ctx, c))
	fmt.Println(magicbox.GetSoftwareVersion(ctx, c))
	fmt.Println(magicbox.GetUpTime(ctx, c))
	fmt.Println(magicbox.GetVendor(ctx, c))
	fmt.Println(license.GetLicenseInfo(ctx, c))
	fmt.Println(storage.GetDeviceAllInfo(ctx, c))

}
