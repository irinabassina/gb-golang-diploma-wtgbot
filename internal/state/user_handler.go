package state

import (
	"WarehouseTgBot/internal/database"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

func handleShowEmployee(chatId int64, s *StateMachine) error {
	users, err := s.userService.FindAllActiveUsers(s.ctx)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatId, s.userService.ConvertToHTML(users))
	msg.ParseMode = tgbotapi.ModeHTML

	if _, err = s.bot.Send(msg); err != nil {
		return err
	}
	s.SetCurrentChatState(chatId, Start)

	return nil
}

func handleRemoveEmployee(update *tgbotapi.Update, s *StateMachine) error {
	i, err := strconv.ParseInt(update.Message.Text, 10, 64)
	if err != nil {
		return errors.New("ошибка: неверный Telegram ID сотрудника")
	}
	err = s.userService.DisableUser(s.ctx, i)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сотрудник исключен из Вашей организации")
	if _, err := s.bot.Send(msg); err != nil {
		return err
	}
	s.SetCurrentChatState(update.Message.Chat.ID, Start)
	return nil
}

func handleAddEmployee(update *tgbotapi.Update, s *StateMachine) error {
	split := strings.Split(update.Message.Text, ":")
	if len(split) != 3 {
		return errors.New("ошибка: неверный формат информации о сотруднике")
	}
	i, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return errors.New("ошибка: неверный Telegram ID сотрудника")
	}
	user := database.User{
		ID:        i,
		Name:      split[1],
		Role:      split[2],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Enabled:   true,
	}
	err = s.userService.AddUser(s.ctx, user)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сотрудник добавлен в Вашу организацию")
	if _, err := s.bot.Send(msg); err != nil {
		return err
	}
	s.SetCurrentChatState(update.Message.Chat.ID, Start)
	return nil
}
