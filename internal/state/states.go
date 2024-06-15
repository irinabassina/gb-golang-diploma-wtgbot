package state

import (
	"WarehouseTgBot/internal/service"
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
	"time"
)

const (
	Start = "start"

	AddEmployee     = "add_employee"
	DisableEmployee = "disable_employee"
	ShowEmployee    = "show_employee"
	AddCategory     = "add_category"

	ShowCategories = "show_categories"
	RemoveCategory = "remove_category"
	EditCategory   = "edit_category"

	GetBalance          = "get_balance"
	AddItems            = "add_items"
	PullItems           = "pull_items"
	RemoveLastOperation = "remove_last_op"

	RetroAI  = "retro_ai"
	FutureAI = "future_ai"
	CallAI   = "call_ai"
)

func GetStateInputRequest(state string) string {
	switch state {
	case Start:
		return "Введи новую команду из меню бота"
	case AddEmployee:
		return "Введите Telegram ID, имя, роль (director или storekeeper) нового сотрудника в формате \"id:имя:роль\""
	case DisableEmployee:
		return "Введите telegram ID сотрудника для удаления из системы"
	case AddCategory:
		return "Введите информацию о новой категории товара в формате \"название:описание:единица_измерения:цена\",\n" +
			"цена указывается в KZT, допустимые единицы измерения \"pcs\" или \"kg\""
	case RemoveCategory:
		return "Введите ID категории товара, которую хотите удалить из системы"
	case EditCategory:
		return "Введите новую информацию о категории товара в формате \"id_категории:название:описание:единица_измерения:цена\",\n" +
			"цена указывается в KZT, допустимые единицы измерения \"pcs\" или \"kg\""
	case CallAI:
		return "Зову ИИ в чат. Можете начать диалог."
	case AddItems:
		return "Введите ID категории и количество единиц товара к зачислению \"id:число_единиц\""
	case PullItems:
		return "Введите ID категории и количество единиц товара к списанию в формате \"id:число_единиц\""
	case RemoveLastOperation:
		return "Введите ID категории товара для которого нужно удалить последнюю операция прихода/расхода"
	case ShowEmployee, ShowCategories, GetBalance, RetroAI, FutureAI:
		return ""
	default:
		return "Неизвестное состояние"
	}
}

type State struct {
	name            string
	updateTimeStamp int64
}

type StateMachine struct {
	currentStates    map[int64]*State
	bot              *tgbotapi.BotAPI
	ctx              context.Context
	userService      *service.UserService
	categoryService  *service.CategoryService
	operationService *service.OperationService
	openai           *openai.Client
}

func NewStateMachine(ctx context.Context, bot *tgbotapi.BotAPI, openai *openai.Client, userService *service.UserService,
	categoryService *service.CategoryService, operationService *service.OperationService) *StateMachine {
	return &StateMachine{ctx: ctx, bot: bot, openai: openai, userService: userService, categoryService: categoryService,
		operationService: operationService, currentStates: make(map[int64]*State)}
}

func (s *StateMachine) GetCurrentChatState(chatId int64) *State {
	return s.currentStates[chatId]
}

func (s *StateMachine) SetCurrentChatState(chatId int64, stateName string) {
	s.currentStates[chatId] = &State{name: stateName, updateTimeStamp: time.Now().Unix()}
}

func (s *StateMachine) HandleState(update *tgbotapi.Update) error {
	var chatId int64

	if update.Message != nil {
		chatId = update.Message.Chat.ID
	} else {
		chatId = update.CallbackQuery.Message.Chat.ID
	}

	switch s.GetCurrentChatState(chatId).name {
	case Start:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, выберите пункт из главного меню")
		if _, err := s.bot.Send(msg); err != nil {
			return err
		}
		return nil
	case CallAI:
		return handleCallAi(update, s)
	case AddEmployee:
		return handleAddEmployee(update, s)
	case DisableEmployee:
		return handleRemoveEmployee(update, s)
	case ShowEmployee:
		return handleShowEmployee(chatId, s)
	case ShowCategories:
		return handleShowCategories(chatId, s)
	case AddCategory:
		return handleAddCategory(update, s)
	case RemoveCategory:
		return handleRemoveCategory(update, s)
	case EditCategory:
		return handleEditCategory(update, s)
	case AddItems:
		return handleAddItems(update, s)
	case PullItems:
		return handlePullItems(update, s)
	case RemoveLastOperation:
		return handleRemoveOperation(update, s)
	case GetBalance:
		return handleShowCurrentBalance(chatId, s)
	case RetroAI, FutureAI:
		//TODO implement me
		return errors.New("reports implementation on the way")
	default:
		return errors.New("not implemented")
	}
}
