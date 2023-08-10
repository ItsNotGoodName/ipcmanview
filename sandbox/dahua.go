package sandbox

//
// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"time"
//
// 	"github.com/ItsNotGoodName/ipcmanview/internal/core"
// 	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
// 	dh "github.com/ItsNotGoodName/ipcmanview/pkg/dahua"
// 	"github.com/ItsNotGoodName/ipcmanview/pkg/dahua/modules/global"
// 	"github.com/ItsNotGoodName/ipcmanview/pkg/dahua/modules/license"
// 	"github.com/ItsNotGoodName/ipcmanview/pkg/dahua/modules/magicbox"
// 	"github.com/ItsNotGoodName/ipcmanview/pkg/dahua/modules/mediafilefind"
// 	"github.com/ItsNotGoodName/ipcmanview/pkg/dahua/modules/storage"
// )
//
//
// 	json, err := json.MarshalIndent(data[0], "", "\t")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	fmt.Println(string(json))
// }
//
// func Dahua(ctx context.Context) {
// 	username, _ := os.LookupEnv("IPC_USERNAME")
// 	password, _ := os.LookupEnv("IPC_PASSWORD")
// 	ip, _ := os.LookupEnv("IPC_IP")
//
// 	cam := core.DahuaCamera{
// 		Address:  ip,
// 		Username: username,
// 		Password: password,
// 	}
//
// 	c := dahua.NewActorHandle(cam)
// 	defer c.Close(ctx)
//
// 	print(dahua.CameraDetailGet(ctx, c))
// 	print(dahua.CameraSoftwareVersionGet(ctx, c))
// 	print(dahua.CameraLicenseList(ctx, c))
//
// 	return
//
// 	fmt.Println(global.GetCurrentTime(ctx, c))
// 	fmt.Println(magicbox.NeedReboot(ctx, c))
// 	fmt.Println(magicbox.GetCPUUsage(ctx, c))
// 	fmt.Println(magicbox.GetDeviceClass(ctx, c))
// 	fmt.Println(magicbox.GetDeviceType(ctx, c))
// 	fmt.Println(magicbox.GetHardwareVersion(ctx, c))
// 	fmt.Println(magicbox.GetMarketArea(ctx, c))
// 	fmt.Println(magicbox.GetMemoryInfo(ctx, c))
// 	fmt.Println(magicbox.GetProcessInfo(ctx, c))
// 	fmt.Println(magicbox.GetSerialNo(ctx, c))
// 	fmt.Println(magicbox.GetSoftwareVersion(ctx, c))
// 	fmt.Println(magicbox.GetUpTime(ctx, c))
// 	fmt.Println(magicbox.GetVendor(ctx, c))
// 	fmt.Println(license.GetLicenseInfo(ctx, c))
// 	fmt.Println(storage.GetDeviceAllInfo(ctx, c))
//
// 	return
//
// 	{
// 		stream, err := mediafilefind.NewStream(
// 			ctx,
// 			c,
// 			mediafilefind.NewCondtion(
// 				dh.NewTimestamp(time.Now().Add(-30*24*time.Hour), time.Local),
// 				dh.NewTimestamp(time.Now(), time.Local),
// 			).Picture(),
// 		)
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		acum := 0
// 		for files, err := stream.Next(ctx, c); files != nil; files, err = stream.Next(ctx, c) {
// 			if err != nil {
// 				panic(err)
// 			}
//
// 			for _, fnfi := range files {
// 				acum += 1
// 				fmt.Printf("%d-------------%+v\n", acum, fnfi)
// 			}
// 		}
// 	}
// }
