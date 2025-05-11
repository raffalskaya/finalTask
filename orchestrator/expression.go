package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Expression struct {
	Id          int64   `json:"id"`
	Status      string  `json:"status" binding:"oneof=active completed calculated"`
	Result      float64 `json:"result"`
	requestData string
	stack       []string
	tmpStack    []float64
}

type ExpressionsMap struct {
	m  map[int64]Expression
	mu sync.RWMutex
	tm TaskMap
}

func NewExpressionsMap() *ExpressionsMap {
	return &ExpressionsMap{
		m:  make(map[int64]Expression),
		tm: *NewTasksMap(),
	}
}

func (em *ExpressionsMap) GetExpression(ctx context.Context, db *sql.DB, id string) (*ExpressionDB, int) {
	// распарисм строку
	uid, err := strconv.ParseInt(id, 16, 64)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	var exp ExpressionDB
	var q = "SELECT id, expression, user_id FROM expressions WHERE id = $1"
	// var q = "SELECT id, name, password FROM users WHERE name=$1"
	err = db.QueryRowContext(ctx, q, uid).Scan(&exp.ID, &exp.Expression, &exp.UserID)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return &exp, http.StatusOK

	// uid, err := uuid.Parse(id)
	// if err != nil {
	// 	return nil, http.StatusInternalServerError
	// }
	// em.mu.Lock()
	// data, exists := em.m[uid]
	// em.mu.Unlock()
	// if !exists {
	// 	return nil, http.StatusNotFound
	// }
	// return &data, http.StatusOK
}

func (em *ExpressionsMap) setTaskResult(ctx context.Context, db *sql.DB, taskBody TaskResultBody) int {
	em.tm.mu.Lock()
	defer em.tm.mu.Unlock()
	uid, err := uuid.Parse(taskBody.Id)
	if err != nil {
		return http.StatusInternalServerError
	}
	task, exists := em.tm.m[uid]
	if !exists {
		return http.StatusNotFound
	}
	delete(em.tm.m, task.Id)

	em.mu.Lock()
	defer em.mu.Unlock()

	expression, ex := em.m[task.expId]
	if !ex {
		return http.StatusInternalServerError
	}
	expression.Status = "active"
	if len(expression.stack) > 0 {
		if len(expression.tmpStack) > 0 {
			expression.tmpStack = append(expression.tmpStack[:1], append([]float64{taskBody.Result}, expression.tmpStack[1:]...)...)
		} else {
			expression.tmpStack = append(expression.tmpStack[:0], append([]float64{taskBody.Result}, expression.tmpStack[0:]...)...)
		}
	} else {
		expression.Result = taskBody.Result
		expression.Status = "completed"
		err := updateExpression(ctx, db, expression.Id, taskBody.Result)
		if err != nil {
			return http.StatusInternalServerError
		}
	}

	em.m[expression.Id] = expression

	return http.StatusOK
}

func updateExpression(ctx context.Context, db *sql.DB, id int64, result float64) error {
	var q = "UPDATE expressions SET result = $1 WHERE id = $2"
	_, err := db.ExecContext(ctx, q, result, id)
	if err != nil {
		return err
	}

	return nil
}

func (em *ExpressionsMap) MoveTaskToStack(task Task) bool {
	em.tm.mu.Lock()
	defer em.tm.mu.Unlock()
	task, exists := em.tm.m[task.Id]
	if !exists {
		return false
	}
	delete(em.tm.m, task.Id)

	em.mu.Lock()
	defer em.mu.Unlock()
	expression, ex := em.m[task.expId]
	if !ex {
		return false
	}
	expression.Status = "active"

	expression.stack = append(expression.stack[:0], append([]string{task.Operation}, expression.stack[0:]...)...)
	expression.tmpStack = append(expression.tmpStack, task.Arg2)
	expression.tmpStack = append(expression.tmpStack, task.Arg1)

	em.m[expression.Id] = expression
	return true
}

func (em *ExpressionsMap) getTask() *Task {
	em.mu.Lock()
	defer em.mu.Unlock()

	for _, exp := range em.m {
		if exp.Status == "completed" || exp.Status == "calculated" {
			continue
		}
		for i, token := range exp.stack {
			if isMathOperator(token) {
				b := exp.tmpStack[len(exp.tmpStack)-1]
				a := exp.tmpStack[len(exp.tmpStack)-2]
				exp.tmpStack = exp.tmpStack[:len(exp.tmpStack)-2]

				for j := 0; j <= i; j++ {
					exp.stack = append(exp.stack[:0], exp.stack[1:]...)
				}

				exp.Status = "calculated"
				em.m[exp.Id] = exp

				em.tm.mu.Lock()
				defer em.tm.mu.Unlock()
				channel := make(chan Task)
				switch token {
				case "+":
					TIME_ADDITION_MS, exists := os.LookupEnv("TIME_ADDITION_MS")
					if !exists {
						TIME_ADDITION_MS = "1000"
					}
					TIME_ADDITION_MS_INT, err := time.ParseDuration(TIME_ADDITION_MS)
					if err != nil {
						TIME_ADDITION_MS_INT = 1000
					}
					task := Task{
						Id:            uuid.New(),
						expId:         exp.Id,
						Arg1:          a,
						Arg2:          b,
						Operation:     token,
						OperationTime: TIME_ADDITION_MS_INT,
					}
					em.tm.m[task.Id] = task

					go func() { channel <- task }()
					time.AfterFunc(TIME_ADDITION_MS_INT*time.Millisecond, func() {
						task := <-channel
						res := expMap.MoveTaskToStack(task)
						if res {
							fmt.Println("Время выполнения опрерации сложения истекло.")
						}
					})
					return &task
				case "-":
					TIME_SUBTRACTION_MS, exists := os.LookupEnv("TIME_SUBTRACTION_MS")
					if !exists {
						TIME_SUBTRACTION_MS = "1000"
					}
					TIME_SUBTRACTION_MS_INT, err := time.ParseDuration(TIME_SUBTRACTION_MS)
					if err != nil {
						TIME_SUBTRACTION_MS_INT = 1000
					}

					task := Task{
						Id:            uuid.New(),
						expId:         exp.Id,
						Arg1:          a,
						Arg2:          b,
						Operation:     token,
						OperationTime: TIME_SUBTRACTION_MS_INT,
					}
					em.tm.m[task.Id] = task
					go func() { channel <- task }()
					time.AfterFunc(TIME_SUBTRACTION_MS_INT*time.Millisecond, func() {
						task := <-channel
						res := expMap.MoveTaskToStack(task)
						if res {
							fmt.Println("Время выполнения опрерации вычитания истекло.")
						}
					})
					return &task
				case "*":
					TIME_MULTIPLICATIONS_MS, exists := os.LookupEnv("TIME_MULTIPLICATIONS_MS")
					if !exists {
						TIME_MULTIPLICATIONS_MS = "1000"
					}
					TIME_MULTIPLICATIONS_MS_INT, err := time.ParseDuration(TIME_MULTIPLICATIONS_MS)
					if err != nil {
						TIME_MULTIPLICATIONS_MS_INT = 1000
					}
					task := Task{
						Id:            uuid.New(),
						expId:         exp.Id,
						Arg1:          a,
						Arg2:          b,
						Operation:     token,
						OperationTime: TIME_MULTIPLICATIONS_MS_INT,
					}
					em.tm.m[task.Id] = task
					go func() { channel <- task }()
					time.AfterFunc(TIME_MULTIPLICATIONS_MS_INT*time.Millisecond, func() {
						task := <-channel
						res := expMap.MoveTaskToStack(task)
						if res {
							fmt.Println("Время выполнения опрерации умножения истекло.")
						}
					})
					return &task
				case "/":
					TIME_DIVISIONS_MS, exists := os.LookupEnv("TIME_DIVISIONS_MS")
					if !exists {
						TIME_DIVISIONS_MS = "1000"
					}
					TIME_DIVISIONS_MS_INT, err := time.ParseDuration(TIME_DIVISIONS_MS)
					if err != nil {
						TIME_DIVISIONS_MS_INT = 1000
					}
					task := Task{
						Id:            uuid.New(),
						expId:         exp.Id,
						Arg1:          a,
						Arg2:          b,
						Operation:     token,
						OperationTime: TIME_DIVISIONS_MS_INT,
					}
					em.tm.m[task.Id] = task
					go func() { channel <- task }()
					time.AfterFunc(TIME_DIVISIONS_MS_INT*time.Millisecond, func() {
						task := <-channel
						res := expMap.MoveTaskToStack(task)
						if res {
							fmt.Println("Время выполнения опрерации деления истекло.")
						}
					})
					return &task
				default:
					return nil
				}
			} else {
				value, _ := strconv.ParseFloat(token, 64)
				exp.tmpStack = append(exp.tmpStack, value)
			}
		}
	}
	return nil
}

func (em *ExpressionsMap) AddExpression(ctx context.Context, db *sql.DB, expression string, userId string) (int, *Expression) {

	allstack, success := createStack(expression)

	if !success {
		return http.StatusUnprocessableEntity, nil
	}

	var q = `
	INSERT INTO expressions (expression, user_id, result) values ($1, $2, "")
	`
	result, err := db.ExecContext(ctx, q, expression, userId)
	if err != nil {
		return http.StatusUnprocessableEntity, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return http.StatusUnprocessableEntity, nil
	}

	exp := Expression{
		Id:          id,
		Status:      "active",
		requestData: expression,
		stack:       allstack,
		tmpStack:    make([]float64, 0),
	}

	em.mu.Lock()
	em.m[exp.Id] = exp
	em.mu.Unlock()
	return http.StatusCreated, &exp
}

type (
	ExpressionDB struct {
		ID         int64
		Expression string
		UserID     int64
		Result     string
	}
)

func (em *ExpressionsMap) GetExpressions(ctx context.Context, db *sql.DB, userId string) ([]ExpressionDB, error) {
	var expressions []ExpressionDB
	var q = "SELECT id, expression, user_id, result FROM expressions WHERE user_id = $1"

	rows, err := db.QueryContext(ctx, q, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		e := ExpressionDB{}
		// err := rows.Scan(&e.ID)
		err := rows.Scan(&e.ID, &e.Expression, &e.UserID, &e.Result)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, e)
	}

	return expressions, nil
}

func (em *ExpressionsMap) ProcessEmptyExpressions(ctx context.Context, db *sql.DB) error {
	var q = "SELECT id, expression, user_id, result FROM expressions WHERE result = ''"

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		e := ExpressionDB{}
		err := rows.Scan(&e.ID, &e.Expression, &e.UserID, &e.Result)
		if err != nil {
			return err
		}

		allstack, success := createStack(e.Expression)

		if !success {
			continue
		}

		exp := Expression{
			Id:          e.ID,
			Status:      "active",
			requestData: e.Expression,
			stack:       allstack,
			tmpStack:    make([]float64, 0),
		}

		em.mu.Lock()
		em.m[exp.Id] = exp
		em.mu.Unlock()
	}

	return nil
}

func createExpressionsTable(ctx context.Context, db *sql.DB) error {
	const usersTable = `
	CREATE TABLE IF NOT EXISTS expressions(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		expression TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		result TEXT NOT NULL,
		FOREIGN KEY (user_id)  REFERENCES expressions (id)
	);`

	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	return nil
}
