package commands

import (
	"WarehouseTgBot/internal/env"
	"WarehouseTgBot/internal/service"
	"WarehouseTgBot/internal/state"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var AiKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Ретроспективный ИИ анализ", state.RetroAI),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Предиктивный ИИ анализ", state.FutureAI),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Общение с ИИ", state.CallAI),
	),
)

var OperationsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

var CategoriesKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Показать категории товаров", state.ShowCategories),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить категорию", state.AddCategory)),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Деактивировать категорию", state.RemoveCategory),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Редактировать категорию", state.EditCategory),
	),
)

var EmployeesKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Показать сотрудников", state.ShowEmployee),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить сотрудника", state.AddEmployee),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отключить сотрудника", state.DisableEmployee),
	),
)

func ProcessCommand(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig, e *env.Env) {
	switch update.Message.Command() {
	case "operations":
		checkRoleAndShowKB(ctx, update, msg, e, "Пожалуйста, выберите операцию управления товарами на складе",
			OperationsKeyboard, service.RoleDirector, service.RoleStorekeeper)
	case "categories":
		checkRoleAndShowKB(ctx, update, msg, e, "Пожалуйста, выберите команду управления категориями товаров",
			CategoriesKeyboard, service.RoleDirector)
	case "employees":
		checkRoleAndShowKB(ctx, update, msg, e, "Пожалуйста, выберите команду управления сотрудниками организации",
			EmployeesKeyboard, service.RoleDirector)
	case "ai":
		checkRoleAndShowKB(ctx, update, msg, e, "Пожалуйста, выберите команду по ИИ анализу склада",
			AiKeyboard, service.RoleDirector)
	default:
		msg.Text = "Неизвестная команда."
	}
}

func checkRoleAndShowKB(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig, e *env.Env,
	kbDescription string, keyboard tgbotapi.InlineKeyboardMarkup, roles ...string) {
	hasRole, err := e.UserService.UserHasRole(ctx, update.Message.From.ID, roles...)
	if err != nil {
		msg.Text = "Ошибка проверки Вашей роли в системе"
		return
	}
	if !hasRole {
		msg.Text = fmt.Sprintf("Ваш пользователь (telegram id = %d) не является допустимым для этой роли управления ботом. Пожалуйста, обратитесь к администратору системы и оплатите подписку на диплом",
			update.SentFrom().ID)
	} else {
		msg.Text = kbDescription
		msg.ReplyMarkup = keyboard
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
