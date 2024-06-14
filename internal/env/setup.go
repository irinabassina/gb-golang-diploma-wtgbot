package env

import (
	"WarehouseTgBot/internal/ai"
	"WarehouseTgBot/internal/database/good"
	"WarehouseTgBot/internal/database/operation"
	"WarehouseTgBot/internal/database/user"
	"WarehouseTgBot/internal/service"
	"WarehouseTgBot/internal/state"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sashabaranov/go-openai"
	"github.com/sethvargo/go-envconfig"
	"log"
	"time"
)

type Env struct {
	Config           Config
	TgBot            *tgbotapi.BotAPI
	OpenaiClient     *openai.Client
	UserService      *service.UserService
	CategoryService  *service.CategoryService
	OperationService *service.OperationService
	StateMachine     *state.StateMachine
}

func Setup(ctx context.Context) (*Env, error) {
	var cfg Config
	env := &Env{}

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, fmt.Errorf("env processing: %w", err)
	}
	env.Config = cfg

	bot, err := tgbotapi.NewBotAPI(cfg.TokenConfig.TgToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = cfg.Logger.Debug
	log.Printf("Authorized on Telegram bot account %s", bot.Self.UserName)
	setBotCommands(bot)
	env.TgBot = bot

	env.OpenaiClient = ai.GetOpenAIClient(env.Config.TokenConfig.OpenAiToken)

	dbConn, err := pgxpool.Connect(ctx, cfg.PostgresConfig.ConnectionURL())
	if err != nil {
		return nil, fmt.Errorf("pgxpool Connect: %w", err)
	}

	const timeout = 5 * time.Second
	userService := service.NewUserService(user.NewRepository(dbConn, timeout), timeout)
	env.UserService = userService
	categoryService := service.NewCategoryService(good.NewRepository(dbConn, timeout), timeout)
	env.CategoryService = categoryService
	operationService := service.NewOperationService(operation.NewRepository(dbConn, timeout), timeout)
	env.OperationService = operationService

	env.StateMachine = state.NewStateMachine(ctx, env.TgBot, env.OpenaiClient, userService, categoryService, operationService)

	return env, nil
}

func setBotCommands(bot *tgbotapi.BotAPI) {
	setCommands := tgbotapi.NewSetMyCommands(tgbotapi.BotCommand{
		Command:     "director",
		Description: "Меню директора",
	},
		tgbotapi.BotCommand{
			Command:     "accountant",
			Description: "Меню бухгалтера",
		},
		tgbotapi.BotCommand{
			Command:     "help",
			Description: "Инструкция по работе с ботом",
		},
	)

	if _, err := bot.Request(setCommands); err != nil {
		log.Println("Unable to set commands")
	}
}
