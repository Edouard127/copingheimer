package http

import (
	"edouard127/copingheimer/src/intf"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"time"
)

type Dashboard struct {
	Engine    *gin.Engine
	Database  *intf.Database
	authUsers map[int32]User
}

func NewDashboard(mongo string) *Dashboard {
	return &Dashboard{
		Engine:    gin.Default(),
		Database:  intf.InitDatabase(mongo),
		authUsers: make(map[int32]User, 0),
	}
}

func (d *Dashboard) Start() {
	d.Engine.POST("/api/status", func(c *gin.Context) {
		// Parse the json data
		data := intf.StatusResponse{}
		if err := c.BindJSON(&data); err != nil {
			fmt.Println("Error while parsing the json data:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Received status:", data)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})

		if err := d.Database.Write(data.IP, &data); err != nil {
			fmt.Println("Error while writing the json data:", err)
			return
		}
	})
	if err := d.Engine.Run(":80"); err != nil {
		panic(err)
	}
}

type User struct {
	IP        net.IP
	Date      time.Time
	RefreshAt time.Time
}

func (d *Dashboard) register(ip net.IP, id int32, date time.Time) int32 {
	d.authUsers[id] = User{
		IP:        ip,
		Date:      date,
		RefreshAt: date.Add(time.Hour * 24),
	}
	return id
}

func (d *Dashboard) Delete(id int32) {
	delete(d.authUsers, id)
}

func (d *Dashboard) RefreshUser(id int32) {
	if user, ok := d.authUsers[id]; ok {
		user.RefreshAt = time.Now().Add(time.Hour * 24)
	}
}

func (d *Dashboard) FindUser(ip net.IP) (int32, bool) {
	for id, user := range d.authUsers {
		if user.IP.Equal(ip) {
			return id, true
		}
	}
	return 0, false
}

func (d *Dashboard) CreateNewUser(ip net.IP) int32 {
	n := int32(len(d.authUsers))
	date := time.Now()
	for i := 0; i < 4; i++ {
		n += int32(ip[i]) | int32(date.Nanosecond()>>(i*8))
	}

	return d.register(ip, n, date)
}
