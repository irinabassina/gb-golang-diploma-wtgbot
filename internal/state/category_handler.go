package state

import (
	"WarehouseTgBot/internal/database"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

func handleRemoveCategory(update *tgbotapi.Update, s *StateMachine) error {
	i, err := strconv.ParseInt(update.Message.Text, 10, 64)
	if err != nil {
		return errors.New("ошибка: неверный формат ID категории")
	}
	err = s.categoryService.DisableCategory(s.ctx, i)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Категория удалена из Вашей организации")
	if _, err := s.bot.Send(msg); err != nil {
		return err
	}
	s.SetCurrentChatState(update.Message.Chat.ID, Start)
	return nil
}

func handleEditCategory(update *tgbotapi.Update, s *StateMachine) error {
	split := strings.Split(update.Message.Text, ":")
	if len(split) != 5 {
		return errors.New("ошибка: неверный формат информации о категории товара")
	}

	id, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {

	}

	cost, err := strconv.ParseFloat(split[4], 64)
	if err != nil {
		return errors.New("ошибка: неверный формат цены категории товара")
	}

	gc := database.GoodCategory{
		ID:          id,
		Name:        split[1],
		Description: split[2],
		Unit:        split[3],
		Cost:        cost,
		CreatedBy:   update.SentFrom().ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Enabled:     true,
	}

	err = s.categoryService.UpdateCategory(s.ctx, gc)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Категория товара отредактирована")
	if _, err := s.bot.Send(msg); err != nil {
		return err
	}
	s.SetCurrentChatState(update.Message.Chat.ID, Start)
	return nil
}

func handleAddCategory(update *tgbotapi.Update, s *StateMachine) error {
	split := strings.Split(update.Message.Text, ":")
	if len(split) != 4 {
		return errors.New("ошибка: неверный формат информации о категории товара")
	}

	cost, err := strconv.ParseFloat(split[3], 64)
	if err != nil {
		return errors.New("ошибка: неверный формат цены категории товара")
	}

	gc := database.GoodCategory{
		Name:        split[0],
		Description: split[1],
		Unit:        split[2],
		Cost:        cost,
		CreatedBy:   update.SentFrom().ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Enabled:     true,
	}

	err = s.categoryService.AddCategory(s.ctx, gc)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Категория товара добавлена в Вашу организацию")
	if _, err := s.bot.Send(msg); err != nil {
		return err
	}
	s.SetCurrentChatState(update.Message.Chat.ID, Start)
	return nil
}

func handleShowCategories(chatId int64, s *StateMachine) error {
	cats, err := s.categoryService.FindAllActiveCategories(s.ctx)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatId, s.categoryService.ConvertToHTML(cats))
	msg.ParseMode = tgbotapi.ModeHTML

	if _, err = s.bot.Send(msg); err != nil {
		return err
	}
	s.SetCurrentChatState(chatId, Start)
	return nil
}
