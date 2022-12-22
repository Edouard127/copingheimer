package http

import (
	"edouard127/copingheimer/src/intf"
	"flag"
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

func NewDashboard() *Dashboard {
	mongo := flag.String("mongo", "mongodb://localhost:27017", "MongoDB connection string")
	flag.Parse()
	return &Dashboard{
		Engine:    gin.Default(),
		Database:  intf.InitDatabase(*mongo),
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
		bytes, err := data.Json(data.IP)
		if err != nil {
			fmt.Println("Error while marshalling the json data:", err)
			return
		}
		if err := d.Database.Write(data.IP, bytes); err != nil {
			fmt.Println("Error while writing the json data:", err)
			return
		}
	})
	if err := d.Engine.Run(":80"); err != nil {
		panic(err)
	}
}

type User struct {
	Date      time.Time
	RefreshAt time.Time
}

func (d *Dashboard) register(id int32, date time.Time) int32 {
	d.authUsers[id] = User{
		Date:      date,
		RefreshAt: date.Add(time.Minute * 5),
	}
	return id
}

func (d *Dashboard) CreateNewUser(ip net.IP) int32 {
	n := int32(len(d.authUsers))
	date := time.Now()
	for i := 0; i < 4; i++ {
		n += int32(ip[i]) | int32(date.Nanosecond()>>(i*8))
	}

	return d.register(n, date)
}
