package state

import (
	"WarehouseTgBot/internal/ai"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCallAi(update *tgbotapi.Update, s *StateMachine) error {
	answer, err := ai.AskGPT(s.ctx, s.openai, s.operationService, update)
	if err != nil {
		return errors.New("ошибка:" + err.Error())
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)
	if _, err = s.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func handleCloseAi(update *tgbotapi.Update, s *StateMachine) error {
	ai.CloseDialog(update.SentFrom().ID)

	msg := tgbotapi.NewMessage(update.FromChat().ID, "Диалог и его история закрыты")
	if _, err := s.bot.Send(msg); err != nil {
		return err
	}

	s.SetCurrentChatState(update.FromChat().ID, Start)
	return nil
}
