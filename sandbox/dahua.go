package sandbox

import (
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

func Dahua() {
	username, _ := os.LookupEnv("IPC_USERNAME")
	password, _ := os.LookupEnv("IPC_PASSWORD")
	ip, _ := os.LookupEnv("IPC_IP")

	c := dahua.NewConn(http.DefaultClient, dahua.NewCamera(ip))
	defer auth.Logout(c)

	dahuaPrint(c)
	err := auth.Login(c, username, password)
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

	fmt.Println(global.GetCurrentTime(c))

	fmt.Println(magicbox.NeedReboot(c))
	fmt.Println(magicbox.GetCPUUsage(c))
	fmt.Println(magicbox.GetDeviceClass(c))
	fmt.Println(magicbox.GetDeviceType(c))
	fmt.Println(magicbox.GetHardwareVersion(c))
	fmt.Println(magicbox.GetMarketArea(c))
	fmt.Println(magicbox.GetMemoryInfo(c))
	fmt.Println(magicbox.GetProcessInfo(c))
	fmt.Println(magicbox.GetSerialNo(c))
	fmt.Println(magicbox.GetSoftwareVersion(c))
	fmt.Println(magicbox.GetUpTime(c))
	fmt.Println(magicbox.GetVendor(c))

	dahuaPrint(c)
}
