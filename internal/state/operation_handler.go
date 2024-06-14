package state

import (
	"WarehouseTgBot/internal/database"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

func handleAddItems(update *tgbotapi.Update, s *StateMachine) error {
	operation, err := getOperation(update)
	if err != nil {
		return err
	}

	err = s.operationService.AddOperation(s.ctx, *operation)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Товары зачислены на склад")
	if _, err := s.bot.Send(msg); err != nil {
		return err
	}
	s.SetCurrentChatState(update.Message.Chat.ID, Start)
	return nil
}

func getOperation(update *tgbotapi.Update) (*database.Operation, error) {
	split := strings.Split(update.Message.Text, ":")
	if len(split) != 2 {
		return nil, errors.New("ошибка: неверный формат информации о товаре")
	}

	id, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return nil, errors.New("ошибка: неверный формат ID категории товара")
	}

	value, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return nil, errors.New("ошибка: неверный формат единиц товара")
	}

	operation := database.Operation{
		CategoryID: id,
		Value:      value,
		CreatedBy:  update.SentFrom().ID,
	}
	return &operation, nil
}

func handlePullItems(update *tgbotapi.Update, s *StateMachine) error {
	operation, err := getOperation(update)
	if err != nil {
		return err
	}
	operation.Value = -operation.Value
	err = s.operationService.AddOperation(s.ctx, *operation)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Товары списаны на склада")
	if _, err := s.bot.Send(msg); err != nil {
		return err
	}
	s.SetCurrentChatState(update.Message.Chat.ID, Start)
	return nil
}

func handleRemoveOperation(update *tgbotapi.Update, s *StateMachine) error {
	id, err := strconv.ParseInt(update.Message.Text, 10, 64)
	if err != nil {
		return errors.New("ошибка: неверный формат ID категории товара")
	}

	operation, err := s.operationService.RemoveLastOperation(s.ctx, id)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("последняя операция удалена успешно. дата последней операции была %s, значение %f",
			operation.CreatedAt, operation.Value))
	if _, err := s.bot.Send(msg); err != nil {
		return err
	}
	s.SetCurrentChatState(update.Message.Chat.ID, Start)

	return nil
}

func handleShowCurrentBalance(chatId int64, s *StateMachine) error {
	cats, err := s.operationService.ShowCurrentBalance(s.ctx)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatId, s.operationService.ConvertToHTML(cats))
	msg.ParseMode = tgbotapi.ModeHTML

	if _, err = s.bot.Send(msg); err != nil {
		return err
	}
	s.SetCurrentChatState(chatId, Start)
	return nil
}
