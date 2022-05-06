package service

import (
	"course-work/app/consts"
	"course-work/app/redis"
	"course-work/app/types"
	"encoding/json"
	"math/rand"
	"strings"
)

func sendMessageChan(room, message string) {
	getChan(room) <- message
}

func getHistroyRedis(room string) string {
	data := redis.Session.Client.Get(redis.Session.Ctx, room+consts.ChatKey)
	if data.Val() == "" {
		return "[]"
	}

	return data.Val()
}

func saveMessageRedis(room, message string) {
	data := redis.Session.Client.Get(redis.Session.Ctx, room+consts.ChatKey)

	var messages []string = make([]string, 0)

	if data.Val() != "" {
		json.Unmarshal([]byte(data.Val()), &messages)
	}

	bytes, _ := json.Marshal(append(messages, message))

	redis.Session.Client.Set(redis.Session.Ctx, room+consts.ChatKey, string(bytes), 0)
}

func getAllUserInRoom(room string) ([]string, error) {
	data := redis.Session.Client.Get(redis.Session.Ctx, room+consts.UsersKey)

	var users []string = make([]string, 0)

	if data.Val() == "" {
		return users, nil
	}

	err := json.Unmarshal([]byte(data.Val()), &users)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func addUser(room, name string) error {
	users, err := getAllUserInRoom(room)

	if err != nil {
		return err
	}

	bytes, err := json.Marshal(append(users, name))

	if err != nil {
		return err
	}

	redis.Session.Client.Set(redis.Session.Ctx, room+consts.UsersKey, string(bytes), 0)

	return nil
}

func RemoveFromRoom(room string, user *types.User) error {

	users, err := getAllUserInRoom(room)

	if err != nil {
		return err
	}

	for i, userChat := range users {
		if userChat == user.Name {
			users = remove(i, users)
			break
		}
	}

	user.LoggedIn = false

	bytes, err := json.Marshal(users)

	if err != nil {
		return err
	}

	redis.Session.Client.Set(redis.Session.Ctx, room+consts.UsersKey, string(bytes), 0)

	admin := findAdmin(room)

	if admin == user.Name {
		if len(users) == 0 {
			setAdminRedis(room, "")
		} else {
			setAdminRedis(room, users[0])
			Pub(room, &types.MessageFront{
				Action: types.NewAdmin,
				Data:   users[0],
			})
		}
	}

	return nil
}

func findUser(room, name string) string {
	data := redis.Session.Client.Get(redis.Session.Ctx, room+consts.UsersKey)

	users := strings.Split(data.Val(), ",")

	for _, userChat := range users {
		if userChat == name {
			return name
		}
	}

	return ""
}

func findAdmin(room string) string {
	return redis.Session.Client.Get(redis.Session.Ctx, room+consts.AdminKey).Val()
}

func setAdminRedis(room, name string) {
	redis.Session.Client.Set(redis.Session.Ctx, room+consts.AdminKey, name, 0)
}

func remove(pos int, arr []string) []string {
	return append(arr[0:pos], arr[pos+1:]...)
}

func removeInt(pos int, arr []int) []int {
	return append(arr[0:pos], arr[pos+1:]...)
}

func getFieldRedis(room string) (string, error) {
	data := redis.Session.Client.Get(redis.Session.Ctx, room+consts.FieldKey)

	if data.Val() == "" {
		bytes, err := json.Marshal(createEmptyField())

		if err != nil {
			return "", err
		}

		redis.Session.Client.Set(redis.Session.Ctx, room+consts.FieldKey, string(bytes), 0)

		return string(bytes), nil

	}

	return data.Val(), nil
}

func getFieldInt(room string) ([][]int, error) {
	var field [][]int = make([][]int, 0)

	fieldStr, err := getFieldRedis(room)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(fieldStr), &field)

	if err != nil {
		return nil, err
	}

	return field, nil
}

func createEmptyField() [][]uint16 {

	var result [][]uint16 = make([][]uint16, 0, consts.Size)

	for i := 0; i < consts.Size; i++ {
		result = append(result, make([]uint16, consts.Size))
	}

	available := getAllPos()

	for i := 0; i < consts.AmountOf2; i++ {
		random := rand.Intn(len(available))

		position := available[random]

		result[position/consts.Size][position%consts.Size] = 2

		available = removeInt(random, available)
	}

	for i := 0; i < consts.AmountOf4; i++ {
		random := rand.Intn(len(available))

		position := available[random]

		result[position/consts.Size][position%consts.Size] = 4

		available = removeInt(random, available)

	}
	return result
}

func getAllPos() []int {
	result := make([]int, consts.Size*consts.Size)

	for i := 0; i < consts.Size*consts.Size; i++ {
		result[i] = i
	}

	return result
}

func saveField(room string, field [][]int) (string, error) {
	bytes, err := json.Marshal(field)

	if err != nil {
		return "", err
	}

	redis.Session.Client.Set(redis.Session.Ctx, room+consts.FieldKey, string(bytes), 0)
	return string(bytes), nil
}

func resetField(room string) (string, error) {
	bytes, err := json.Marshal(createEmptyField())

	if err != nil {
		return "", err
	}

	redis.Session.Client.Set(redis.Session.Ctx, room+consts.FieldKey, string(bytes), 0)

	return string(bytes), nil
}
