package commands

import (
	"WarehouseTgBot/internal/env"
	"WarehouseTgBot/internal/service"
	"WarehouseTgBot/internal/state"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var DirectorKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Показать сотрудников", state.ShowEmployee),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить сотрудника", state.AddEmployee),
		tgbotapi.NewInlineKeyboardButtonData("Отключить сотрудника", state.DisableEmployee),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Показать категории товаров", state.ShowCategories),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить категорию", state.AddCategory),
		tgbotapi.NewInlineKeyboardButtonData("Удалить категорию", state.RemoveCategory),
		tgbotapi.NewInlineKeyboardButtonData("Редактировать категорию", state.EditCategory),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удалить последнюю операцию по товару", state.RemoveLastOperation),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Общение с ИИ", state.CallAI),
	),
)

var AccountantKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Показать категории товаров", state.ShowCategories),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Внести приход товара", state.AddItems),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Внести расход товара", state.PullItems),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Показать текущий баланс склада", state.GetBalance),
	),
)

func ProcessCommand(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig, e *env.Env) {
	switch update.Message.Command() {
	case "help":
		msg.Text = "Допустимые команды:\n" +
			"1. /director - директорские операции (добавление новой категории товаров, добавление сотрудников, формирование отчетов и прогнозирование спроса) \n" +
			"2. /accountant - внести операции прихода/расхода товара."
	case "director":
		hasRole, err := e.UserService.UserHasRole(ctx, update.Message.From.ID, service.RoleDirector)
		if err != nil {
			msg.Text = "Ошибка проверки Вашей роли в системе"
			break
		}
		if !hasRole {
			msg.Text = fmt.Sprintf("Ваш пользователь (telegram id = %d) не является допустимым для этой роли управления ботом. Пожалуйста, обратитесь к администратору системы и оплатите подписку на диплом",
				update.SentFrom().ID)
		} else {
			msg.Text = "Пожалуйста, выберите команду директора"
			msg.ReplyMarkup = DirectorKeyboard
		}
	case "accountant":
		hasRole, err := e.UserService.UserHasRole(ctx, update.Message.From.ID, service.RoleDirector, service.RoleAccountant)
		if err != nil {
			msg.Text = "Ошибка проверки Вашей роли в системе"
			break
		}
		if !hasRole {
			msg.Text = fmt.Sprintf("Ваш пользователь (telegram id = %d) не является допустимым для этой роли управления ботом. Пожалуйста, обратитесь к администратору системы и оплатите подписку на диплом",
				update.SentFrom().ID)
		} else {
			msg.Text = "Пожалуйста, выберите команду бухгалтера"
			msg.ReplyMarkup = AccountantKeyboard
		}

	default:
		msg.Text = "Неизвестная команда. Список возможных команд: /help"
	}
}

func HandleUpdateCallBackQuery(e *env.Env, update tgbotapi.Update) {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := e.TgBot.Request(callback); err != nil {
		panic(err)
	}

	e.StateMachine.SetCurrentChatState(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	request := state.GetStateInputRequest(update.CallbackQuery.Data)
	if request != "" {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, request)
		if _, err := e.TgBot.Send(msg); err != nil {
			panic(err)
		}
	} else {
		err := e.StateMachine.HandleState(&update)
		if err != nil {
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, err.Error())
			if _, err := e.TgBot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}
