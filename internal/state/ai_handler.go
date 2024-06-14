package state

import (
	"WarehouseTgBot/internal/ai"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCallAi(update *tgbotapi.Update, s *StateMachine) error {
	answer, err := ai.AskGPT(s.openai, update)
	if err != nil {
		return errors.New("ошибка связи с ИИ")
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)
	if _, err = s.bot.Send(msg); err != nil {
		return err
	}
	return nil
}
