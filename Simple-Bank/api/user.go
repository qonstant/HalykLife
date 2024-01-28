package api

import (
	db "Simple-Bank/db/sqlc"
	"Simple-Bank/util"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

type Message struct {
	IIN       uint64 `json:"iin"`
	Username  string `json:"username"`
	Action    string `json:"action"`
	Timestamp string `json:"timestamp"`
}

type createUserRequest struct {
	Iin      int64  `json:"iin" binding:"required,min=5`
	Username string `json:"username" binding:"required,alphanum"`
	Name     string `json:"name" binding:"required,min=5`
	Surname  string `json:"surname" binding:"required,min=5`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	UserRole string `json:"user_role"`
}

type userResponse struct {
	Iin            int64  `json:"iin"`
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	Email          string `json:"email"`
	UserRole       string `json:"user_role"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Iin:            user.Iin,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		Name:           user.Name,
		Surname:        user.Surname,
		Email:          user.Email,
		UserRole:       user.UserRole,
	}
}

type listUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	// Пример сообщения для отправки
	message := Message{
		IIN:       1234567890,
		Username:  "JohnDoe",
		Action:    "getListUsers",
		Timestamp: "2024-01-28T12:00:00Z",
	}
	var req listUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sendMessageToClickHouse(message)
	ctx.JSON(http.StatusOK, users)
}

type getUserRequestByIin struct {
	iin int64 `uri:"iin" binding:"required,min=1"`
}

func (server *Server) getQrForUser(ctx *gin.Context) {
	var req getUserRequestByIin
	var iin = ctx.Param("iin")
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fmt.Println(iin)
	var IntIin, err = strconv.ParseInt(iin, 10, 64)
	user, err := server.store.GetUser(ctx, IntIin)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	if user.UserRole != "user" {
		ctx.JSON(http.StatusUnauthorized, user)
		return
	}

	username := user.Username
	surname := user.Surname

	url := "mailto:info@halyklife.kz?subject=" + string(username) + " " + string(surname) + " IIN:" + string(iin)
	// Пример сообщения для отправки
	message := Message{
		IIN:       1234567890,
		Username:  "JohnDoe",
		Action:    "getQR",
		Timestamp: "2024-01-28T12:00:00Z",
	}
	sendMessageToClickHouse(message)

	ctx.JSON(http.StatusOK, url)
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Iin:            req.Iin,
		Username:       req.Username,
		Name:           req.Name,
		Surname:        req.Surname,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		UserRole:       string(req.UserRole),
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
			log.Println(pqErr.Code.Name())
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newUserResponse(user)
	// Пример сообщения для отправки
	message := Message{
		IIN:       1234567890,
		Username:  "JohnDoe",
		Action:    "userCreated",
		Timestamp: "2024-01-28T12:00:00Z",
	}
	sendMessageToClickHouse(message)
	ctx.JSON(http.StatusOK, response)
}

type loginUserRequest struct {
	Iin      int64  `json:"iin"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Iin)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	//err = util.CheckRole(user.UserRole, req.U)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
		user.UserRole,
	)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	// Пример сообщения для отправки
	message := Message{
		IIN:       1234567890,
		Username:  "JohnDoe",
		Action:    "userLogin",
		Timestamp: "2024-01-28T12:00:00Z",
	}
	sendMessageToClickHouse(message)
	ctx.JSON(http.StatusOK, rsp)
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, user)
}

func sendMessageToClickHouse(message Message) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092", // Указываем адрес брокера Kafka в идеале должен лежать в app.env но я не успел не обижайтесь :)
	}

	// Создаем нового producer
	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return
	}
	defer producer.Close()

	// Создаем канал для обработки результатов доставки
	deliveryChan := make(chan kafka.Event)
	// Обрабатываем результаты доставки в отдельной горутине
	handleDeliveryEvents(deliveryChan)

	topic := "baeldung" //это тоже :)

	err = SendMessageToKafkaTopicAsync(producer, topic, message, deliveryChan)
	if err != nil {
		fmt.Printf("Failed to send message to Kafka topic: %v\n", err)
		return
	}
	producer.Flush(30 * 1000)
}
func handleDeliveryEvents(deliveryChan <-chan kafka.Event) {
	go func() {
		for e := range deliveryChan {
			m := e.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		}
	}()
}

func SendMessageToKafkaTopicAsync(producer *kafka.Producer, topic string, message Message, deliveryChan chan kafka.Event) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %v", err) // сообщение идет на прямую в clickhouse там уже можно сделать крутые query для статистики
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          messageJSON,
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to produce message: %v", err)
	}

	return nil
}
