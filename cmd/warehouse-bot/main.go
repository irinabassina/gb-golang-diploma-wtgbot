package main

import (
	"WarehouseTgBot/internal/commands"
	"WarehouseTgBot/internal/env"
	"WarehouseTgBot/internal/state"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := runMain(ctx); err != nil {
		log.Fatal(err)
	}
}

func runMain(ctx context.Context) error {
	e, err := env.Setup(ctx)
	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := e.TgBot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if e.StateMachine.GetCurrentChatState(update.Message.Chat.ID) == nil {
				e.StateMachine.SetCurrentChatState(update.Message.Chat.ID, state.Start)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			if update.Message.IsCommand() {
				commands.ProcessCommand(ctx, &update, &msg, e)
				if _, err = e.TgBot.Send(msg); err != nil {
					panic(err)
				}
				continue
			}

			err := e.StateMachine.HandleState(&update)
			if err != nil {
				msg.Text = err.Error()
				e.TgBot.Send(msg)
			}

		} else if update.CallbackQuery != nil {
			commands.HandleUpdateCallBackQuery(e, update)
		}
	}
	return nil
}
