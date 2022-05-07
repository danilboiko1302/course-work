package service

import (
	"bytes"
	"course-work/app/consts"
	"course-work/app/nats"
	"course-work/app/types"
	"course-work/app/vocabulary"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"
)

func Sub(name string) (func() error, error) {
	return nats.Connection.Sub(name, func(data []byte) {
		go func(r io.Reader) {
			var msg types.Message

			err := gob.NewDecoder(r).Decode(&msg)

			if err != nil {
				log.Println("error decoding message")
				return
			}

		}(bytes.NewReader(data))
	})
}

func Pub(name string, msg *types.MessageFront) error {
	bytes, err := json.Marshal(msg)

	if err != nil {
		log.Println("error json.Marshal: " + err.Error())
		return err
	}

	return nats.Connection.Pub(name, bytes)
}

func HandleMessage(room string, msg types.Message, user *types.User) (*types.MessageFront, *types.MessageFront, error) {
	switch msg.Action {
	case types.ErrorAction:
		{
			return nil, nil, errors.New(vocabulary.ERROR_ACTION)
		}
	case types.Login:
		{
			return login(room, msg.Data, user)
		}
	case types.Logout:
		{
			return logout(room, msg.Data, user)
		}
	case types.GetUsers:
		{
			return getUsers(room, user)
		}
	case types.SetAdmin:
		{
			return setAdmin(room, msg.Data, user)
		}
	case types.GetField:
		{
			return getField(room, user)
		}
	case types.Up:
		{
			return move(room, user, types.MoveUp)
		}
	case types.Down:
		{
			return move(room, user, types.MoveDown)
		}
	case types.Left:
		{
			return move(room, user, types.MoveLeft)
		}
	case types.Right:
		{
			return move(room, user, types.MoveRight)
		}
	case types.SendMessage:
		{
			return sendMessage(room, msg.Data, user)
		}
	case types.GetHistory:
		{
			return getHistroy(room, user)
		}
	}

	return nil, nil, errors.New(vocabulary.UNKNOWN_ACTION)
}

func getHistroy(room string, user *types.User) (*types.MessageFront, *types.MessageFront, error) {
	if !user.LoggedIn {
		return nil, nil, errors.New("you are not logined in")
	}

	return &types.MessageFront{
		Action: types.GetMessages,
		Data:   getHistroyRedis(room),
	}, nil, nil
}

func sendMessage(room, data string, user *types.User) (*types.MessageFront, *types.MessageFront, error) {
	if !user.LoggedIn {
		return nil, nil, errors.New("you are not logined in")
	}

	sendMessageChan(room, user.Name+": "+data)

	return &types.MessageFront{
		Action: types.GetMessage,
		Data:   user.Name + ": " + data,
	}, nil, nil
}

func setAdmin(room, name string, user *types.User) (*types.MessageFront, *types.MessageFront, error) {
	if !user.LoggedIn {
		return nil, nil, errors.New("you are not logined in")
	}

	admin := findAdmin(room)

	if admin != user.Name {
		return nil, nil, errors.New("you are not admin")
	}

	if admin == name {
		return nil, nil, errors.New("you are already admin")
	}

	setAdminRedis(room, name)

	return &types.MessageFront{
		Action: types.NewAdmin,
		Data:   name,
	}, nil, nil
}

func move(room string, user *types.User, direction types.Move) (*types.MessageFront, *types.MessageFront, error) {
	if !user.LoggedIn {
		return nil, nil, errors.New("you are not logined in")
	}

	admin := findAdmin(room)

	if admin != user.Name {
		return nil, nil, errors.New("you are not admin")
	}

	field, err := getFieldInt(room)

	if err != nil {
		return nil, nil, err
	}

	switch direction {
	case types.MoveUp:
		{
			field = rotateLeft(field)
		}
	case types.MoveDown:
		{
			field = rotateRight(field)
		}
	case types.MoveRight:
		{
			field = rotateLeft(field)
			field = rotateLeft(field)
		}
	}
	newField := step(field)

	if compare(field, newField) {
		return nil, nil, errors.New("you can`t do this move")
	}

	newField = add(newField)

	switch direction {
	case types.MoveUp:
		{
			newField = rotateRight(newField)
		}
	case types.MoveDown:
		{
			newField = rotateLeft(newField)
		}
	case types.MoveRight:
		{
			newField = rotateLeft(newField)
			newField = rotateLeft(newField)
		}
	}

	copy := make([][]int, 0, len(field))

	for _, row := range newField {
		copy = append(copy, append([]int{}, row...))
	}

	if check(copy) {
		reset, err := resetField(room)

		if err != nil {
			return nil, nil, err
		}

		bytes, err := json.Marshal(newField)

		if err != nil {
			return nil, nil, err
		}

		Pub(room, &types.MessageFront{
			Action: types.Field,
			Data:   string(bytes),
		})
		time.Sleep(time.Millisecond * 500)

		Pub(room, &types.MessageFront{
			Action: types.Lost,
		})
		time.Sleep(time.Millisecond * 500)

		Pub(room, &types.MessageFront{
			Action: types.Field,
			Data:   reset,
		})
		return nil, nil, nil
	}

	data, err := saveField(room, newField)

	if err != nil {
		return nil, nil, err
	}

	return &types.MessageFront{
		Action: types.Field,
		Data:   data,
	}, nil, nil
}

func check(field [][]int) bool {
	res := []bool{checkField(field)}

	for i := 0; i < 3; i++ {
		field = rotateLeft(field)
		res = append(res, checkField(field))

	}

	for _, check := range res {
		if !check {
			return false
		}
	}

	return true
}

func checkField(field [][]int) bool {
	newField := step(field)

	return compare(field, newField)
}

func compare(field, newField [][]int) bool {
	for i := 0; i < consts.Size; i++ {
		for j := 0; j < consts.Size; j++ {
			if field[i][j] != newField[i][j] {
				return false
			}
		}
	}

	return true
}

func add(field [][]int) [][]int {
	random := rand.Intn(100)

	num := 2

	if random < consts.AmountOf4 {
		num = 4
	}

	available := findPos(field)

	random = rand.Intn(len(available))

	pos := available[random]

	field[pos/consts.Size][pos%consts.Size] = num

	return field
}

func findPos(field [][]int) []int {
	result := make([]int, 0)

	for i, row := range field {
		for j, elem := range row {
			if elem == 0 {
				result = append(result, i*consts.Size+j)
			}
		}
	}

	return result
}

func step(field [][]int) [][]int {
	result := make([][]int, 0, len(field))

	for _, row := range field {
		result = append(result, stepRow(append([]int{}, row...)))
	}

	return result
}
func stepRow(row []int) []int {

	return moveRow(connect(moveRow(row)))
}

func moveRow(row []int) []int {
	for i := 0; i < len(row)-1; i++ {
		if row[i] == 0 && row[i+1] != 0 {
			row[i] = row[i+1]
			row[i+1] = 0
			return moveRow(row)
		}
	}

	return row
}

func connect(row []int) []int {
	for i := 0; i < len(row)-1; i++ {
		if row[i] != 0 && row[i] == row[i+1] {
			row[i] = row[i] * 2
			row[i+1] = 0
		}
	}
	return row
}

func rotateLeft(field [][]int) [][]int {
	for i := 0; i < consts.Size/2; i++ {
		for j := i; j < consts.Size-i-1; j++ {
			temp := field[i][j]
			field[i][j] = field[j][consts.Size-1-i]
			field[j][consts.Size-1-i] = field[consts.Size-1-i][consts.Size-1-j]
			field[consts.Size-1-i][consts.Size-1-j] = field[consts.Size-1-j][i]
			field[consts.Size-1-j][i] = temp
		}
	}

	return field
}

func rotateRight(field [][]int) [][]int {
	return rotateLeft(rotateLeft(rotateLeft(field)))
}

func getField(room string, user *types.User) (*types.MessageFront, *types.MessageFront, error) {
	if !user.LoggedIn {
		return nil, nil, errors.New("you are not logined in")
	}

	field, err := getFieldRedis(room)

	if err != nil {
		return nil, nil, err
	}

	return nil, &types.MessageFront{
		Action: types.Field,
		Data:   field,
	}, nil
}

func getUsers(room string, user *types.User) (*types.MessageFront, *types.MessageFront, error) {
	if !user.LoggedIn {
		return nil, nil, errors.New("you are not logined in")
	}

	users, err := getAllUserInRoom(room)

	if err != nil {
		return nil, nil, err
	}

	bytes, err := json.Marshal(users)

	if err != nil {
		return nil, nil, err
	}

	return nil, &types.MessageFront{
		Action: types.AllUsers,
		Data:   string(bytes),
	}, nil
}

func logout(room string, name string, user *types.User) (*types.MessageFront, *types.MessageFront, error) {
	if !user.LoggedIn {
		return nil, nil, errors.New("you are not logined in")
	}

	RemoveFromRoom(room, user)

	user.LoggedIn = false

	return &types.MessageFront{
		Action: types.UserLeft,
		Data:   name,
	}, nil, nil

}

func login(room string, name string, user *types.User) (*types.MessageFront, *types.MessageFront, error) {
	users, err := getAllUserInRoom(room)

	if err != nil {
		return nil, nil, err
	}

	for _, user := range users {
		if user == name {
			return nil, nil, fmt.Errorf(vocabulary.LOGIN_ERROR, "user already exists")
		}
	}

	addUser(room, name)

	user.LoggedIn = true
	user.Name = name

	admin := findAdmin(room)

	if admin == "" {
		setAdminRedis(room, name)
	}

	return &types.MessageFront{
			Action: types.NewUser,
			Data:   name,
		}, &types.MessageFront{
			Action: types.LogedIn,
		}, nil
}
