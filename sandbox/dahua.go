package sandbox

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ItsNotGoodName/pkg/dahua"
	"github.com/ItsNotGoodName/pkg/dahua/auth"
	"github.com/ItsNotGoodName/pkg/dahua/modules/global"
	"github.com/ItsNotGoodName/pkg/dahua/modules/magicbox"
)

func dahuaPrint(c *dahua.Conn) {
	fmt.Println("-", c.LastLogin)
}

func Dahua(ctx context.Context) {
	username, _ := os.LookupEnv("IPC_USERNAME")
	password, _ := os.LookupEnv("IPC_PASSWORD")
	ip, _ := os.LookupEnv("IPC_IP")

	c := dahua.NewConn(http.DefaultClient, dahua.NewCamera(ip))
	defer auth.Logout(context.Background(), c)

	dahuaPrint(c)
	err := auth.Login(ctx, c, username, password)
	if err != nil {
		panic(err)
	}

	// ----------------------------- global
	// time.Sleep(65 * time.Second)
	//
	// ok, err := auth.KeepAlive(c)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("KeepAlive:", ok)

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

	dahuaPrint(c)
}
