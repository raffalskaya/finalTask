package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	pb "github.com/raffalskaya/finalTask/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

var expMap = NewExpressionsMap()

const hmacSampleSecret = "super_secret_signature"

type GRPCServer struct {
	pb.APIServiceServer // сервис из сгенерированного пакета
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

// Структура для чтения выражения из запроса
type ExpressionBody struct {
	Expression string `json:"expression"`
}

type TaskResultBody struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
}

func getUserIdToken(tokenString string) (string, error) {
	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		return "", err
		// log.Fatal(err)
	}

	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
		idValue, exists := claims["id"]
		if !exists {
			return "", err
		}

		// Преобразуем значение в строку
		idString := fmt.Sprintf("%v", idValue)
		return idString, nil
	} else {
		return "", err
	}
}

func createExpression(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header is missing"})
		return
	}

	// Проверяем, что заголовок начинается с "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(401, gin.H{"error": "Invalid authorization format"})
		return
	}

	// Извлекаем токен (убираем префикс "Bearer ")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	var userId, err = getUserIdToken(token)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}

	var expressionBody ExpressionBody
	// Привязываем тело запроса JSON к структуре
	if err := c.BindJSON(&expressionBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Проверим выражение на валидность
	_, valid := createStack(expressionBody.Expression)
	if !valid {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid expression"})
		return
	}
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var result, exp = expMap.AddExpression(ctx, db, expressionBody.Expression, userId)

	if result != http.StatusCreated {
		c.JSON(result, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(result, gin.H{"id": strconv.FormatInt(exp.Id, 10)})
}

func getExpressions(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header is missing"})
		return
	}

	// Проверяем, что заголовок начинается с "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(401, gin.H{"error": "Invalid authorization format"})
		return
	}

	// Извлекаем токен (убираем префикс "Bearer ")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	var userId, err = getUserIdToken(token)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}

	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	response, err := expMap.GetExpressions(ctx, db, userId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, gin.H{"expressions": response})
	}
}

func getExpression(c *gin.Context) {
	// Извлекаем параметр id из пути
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "id is no set"})
		return
	}
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	expression, status := expMap.GetExpression(ctx, db, id)

	c.JSON(status, gin.H{"expression": expression})
}

// Реализация метода SetTask по GRPC, установка результата вычисления
func (s *GRPCServer) SetTask(grpcctx context.Context, req *pb.TaskResult) (*emptypb.Empty, error) {
	taskBody := TaskResultBody{
		Id:     req.Id,
		Result: req.Result,
	}
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return &emptypb.Empty{}, err
	}
	defer db.Close()

	expMap.setTaskResult(ctx, db, taskBody)

	return &emptypb.Empty{}, nil
}

// Реализация метода GetTask по GRPC
func (s *GRPCServer) GetTask(ctx context.Context, req *emptypb.Empty) (*pb.TaskResponse, error) {
	task := expMap.getTask()
	if task == nil {
		return &pb.TaskResponse{
			Enabled: false,
			Task:    nil,
		}, nil
	}
	// uid, err := strconv.ParseInt(id, 16, 64)

	return &pb.TaskResponse{
		Enabled: true,
		Task: &pb.Task{
			Id:            task.Id.String(),
			ExpId:         strconv.FormatInt(task.expId, 10),
			Arg1:          task.Arg1,
			Arg2:          task.Arg2,
			Operation:     task.Operation,
			OperationTime: durationpb.New(task.OperationTime),
		},
	}, nil
}

func createToken(u User) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": u.Name,
		"id":   strconv.FormatInt(u.ID, 10),
		"nbf":  now.Unix(),
		"exp":  now.Add(60 * time.Minute).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		// panic(err)
		return "", err
	}
	return tokenString, nil
}

func login(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userFromDB, err := selectUser(ctx, db, user.Name)
	if err != nil {
		panic(err)
	}

	if nil == user.ComparePassword(userFromDB) {
		token, err := createToken(userFromDB)
		if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, gin.H{"token": token})
		}
	} else {
		c.Status(http.StatusForbidden)
	}
}

func register(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	password, err := generate(user.OriginPassword)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	dbUser := &User{
		Name:           user.Name,
		Password:       password,
		OriginPassword: user.OriginPassword,
	}
	userID, err := insertUser(ctx, db, dbUser)
	if err != nil {
		c.Status(http.StatusFound)
		log.Println("user already exists")
	} else {
		user.ID = userID
		c.Status(http.StatusOK)
	}
}

func main() {
	//-------------------------------------------------
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	if err = createUsersTable(ctx, db); err != nil {
		panic(err)
	}
	if err = createExpressionsTable(ctx, db); err != nil {
		panic(err)
	}
	//-------------------------------------------------

	expMap.ProcessEmptyExpressions(ctx, db)

	router := gin.Default()
	router.POST("/api/v1/calculate", createExpression)
	router.GET("/api/v1/expressions", getExpressions)
	router.GET("/api/v1/expressions/:id", getExpression)
	router.POST("/api/v1/login", login)
	router.POST("/api/v1/register", register)

	go router.Run(":8000")

	host := "localhost"
	port := "50051"

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", port)
	// создадим сервер grpc
	grpcServer := grpc.NewServer()
	// объект структуры, которая содержит реализацию
	// серверной части
	grpcServiceServer := NewGRPCServer()
	// зарегистрируем нашу реализацию сервера
	pb.RegisterAPIServiceServer(grpcServer, grpcServiceServer)
	// запустим grpc сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}
