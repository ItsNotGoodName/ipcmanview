package sandbox

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ItsNotGoodName/pkg/dahua"
	"github.com/ItsNotGoodName/pkg/dahua/client"
	"github.com/ItsNotGoodName/pkg/dahua/modules/global"
)

func Dahua() {
	username, _ := os.LookupEnv("IPC_USERNAME")
	password, _ := os.LookupEnv("IPC_PASSWORD")
	ip, _ := os.LookupEnv("IPC_IP")

	c := client.Client{
		Username: username,
		Password: password,
		Camera:   dahua.NewCamera(ip),
		Conn:     dahua.NewConn(http.DefaultClient),
	}

	err := client.Login(c.Conn, c, c.Username, c.Password)
	if err != nil {
		panic(err)
	}

	time, err := global.GetCurrentTime(c)
	if err != nil {
		panic(err)
	}
	fmt.Println(time)

	keep, err := global.KeepAlive(c)
	if err != nil {
		panic(err)
	}
	fmt.Println(keep)

	client.Logout(c.Conn, c)
}
