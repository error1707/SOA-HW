package main

import (
	"encoding/json"
	"fmt"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"strconv"
)

var Channel *amqp.Channel

func main() {
	SetupQueue()
	err := SetupStorage()
	if err != nil {
		log.Panicf("Failed to connect to DB:\n%v", err)
	}
	r := gin.Default()
	makeHandler(r)
	err = r.Run(":8080")
	if err != nil {
		log.Panicf("Error on serving:\n%v", err)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func makeHandler(r *gin.Engine) {
	r.HTMLRender = ginview.Default()

	r.GET("/", func(ctx *gin.Context) {
		//render with master
		ctx.HTML(http.StatusOK, "index", gin.H{
			"title": "MafiaProfile",
		})
	})

	r.GET("/users", func(c *gin.Context) {
		users, err := DB.GetUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	r.POST("/user", func(c *gin.Context) {
		user := &UserProfile{}
		err := c.ShouldBindWith(user, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
			return
		}
		err = DB.CreateUser(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	})

	r.GET("/user/:username", func(c *gin.Context) {
		username, _ := c.Params.Get("username")
		user, err := DB.GetUser(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
			return
		}
		if user == nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, user)
	})

	r.PATCH("/user/:username", func(c *gin.Context) {
		username, _ := c.Params.Get("username")
		userPatch := &UserProfile{}
		err := c.ShouldBindWith(userPatch, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
			return
		}
		updated, err := DB.UpdateUser(username, userPatch)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
			return
		}
		if updated == 0 {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusOK)
	})

	r.DELETE("/user/:username", func(c *gin.Context) {
		username, _ := c.Params.Get("username")
		deleted, err := DB.DeleteUser(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
			return
		}
		if deleted == 0 {
			c.Status(http.StatusNoContent)
			return
		}
		c.Status(http.StatusOK)
	})

	r.POST("/user/:username/report", func(c *gin.Context) {
		username, _ := c.Params.Get("username")
		user, err := DB.GetUser(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
			return
		}
		if user == nil {
			c.Status(http.StatusNotFound)
			return
		}
		jobId, err := DB.CreateJob()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
			return
		}
		rawJob, err := json.Marshal(PDFJob{ID: jobId, User: user})
		err = Channel.Publish("",
			"requests",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        rawJob,
			})
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
			return
		}
		c.JSON(http.StatusOK, map[string]uint64{"job_id": jobId})
	})

	r.GET("/report/:job_id", func(c *gin.Context) {
		jobIdRaw, _ := c.Params.Get("job_id")
		jobId, err := strconv.Atoi(jobIdRaw)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{"bad job id"})
			return
		}
		pdf, err := DB.GetJob(jobId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
			return
		}
		if pdf == nil {
			c.Status(http.StatusNotFound)
			return
		}
		if pdf.FilePath == "" {
			c.JSON(http.StatusOK, map[string]string{"result": "Not ready yet"})
		}
		log.Printf("Get job!\n")
		c.Header("Content-Description", "Report")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+pdf.FilePath)
		c.Header("Content-Type", "application/octet-stream")
		c.File(fmt.Sprintf("/reports/%s", pdf.FilePath))
	})
}

func SetupQueue() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Error on connecting to rabbitmq: %v\n", err)
	}
	//defer conn.Close()
	Channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Error on opening a channel: %v\n", err)
	}
	//defer Channel.Close()

	_, err = Channel.QueueDeclare(
		"requests", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Error on declaring a queue: %v\n", err)
	}
	Responses, err := Channel.QueueDeclare(
		"responses", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Fatalf("Error on declaring a queue: %v\n", err)
	}
	go func() {
		msgs, err := Channel.Consume(
			Responses.Name, // queue
			"",             // consumer
			true,           // auto-ack
			false,          // exclusive
			false,          // no-local
			false,          // no-wait
			nil,            // args
		)
		if err != nil {
			log.Fatalf("Error on registering a consumer: %v\n", err)
		}
		for d := range msgs {
			var msg PDFJob
			json.Unmarshal(d.Body, &msg)
			log.Printf("PDF ready for: %d\n", msg.ID)
			found, err := DB.UpdateJob(msg)
			if err != nil {
				log.Printf("Error on updating job %d: %v\n", msg.ID, err)
			}
			if found == 0 {
				log.Printf("Job %d not found\n", msg.ID)
			}
		}
	}()
}

type PDFJob struct {
	ID       uint64       `json:"id" gorm:"id"`
	FilePath string       `json:"file_path" gorm:"file_path"`
	User     *UserProfile `json:"user" gorm:"-"`
}
