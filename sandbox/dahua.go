package sandbox

import (
	"fmt"
	"net/http"
	"os"
	"time"

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

	c := dahua.NewConn(http.DefaultClient, dahua.NewCamera(ip))

	dahuaPrint(c)
	err := auth.Login(c, username, password)
	if err != nil {
		panic(err)
	}

	dahuaPrint(c)
	date, err := global.GetCurrentTime(c)
	if err != nil {
		panic(err)
	}
	fmt.Println(date)

	time.Sleep(65 * time.Second)

	ok, err := auth.KeepAlive(c)
	if err != nil {
		panic(err)
	}
	fmt.Println("KeepAlive:", ok)

	dahuaPrint(c)
	auth.Logout(c)
}
