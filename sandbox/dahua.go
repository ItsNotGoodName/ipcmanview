package sandbox

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ItsNotGoodName/pkg/dahua"
	"github.com/ItsNotGoodName/pkg/dahua/auth"
	"github.com/ItsNotGoodName/pkg/dahua/modules/global"
)

func dahuaPrint(c *dahua.Conn) {
	fmt.Println("-", c.LastLogin)
}

func Dahua() {
	username, _ := os.LookupEnv("IPC_USERNAME")
	password, _ := os.LookupEnv("IPC_PASSWORD")
	ip, _ := os.LookupEnv("IPC_IP")

	c := dahua.NewConn(http.DefaultClient, ip)

	dahuaPrint(c)
	err := auth.Login(c, username, password)
	if err != nil {
		panic(err)
	}

	dahuaPrint(c)
	time, err := global.GetCurrentTime(c)
	if err != nil {
		panic(err)
	}
	fmt.Println(time)

	dahuaPrint(c)
	keep, err := global.KeepAlive(c)
	if err != nil {
		panic(err)
	}
	fmt.Println(keep)

	dahuaPrint(c)
	auth.Logout(c)
}
