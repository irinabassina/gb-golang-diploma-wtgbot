package ai

import (
	"WarehouseTgBot/internal/service"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
	"log"
)

var conversationMessages = make(map[int64][]openai.ChatCompletionMessage)

func GetOpenAIClient(token string) *openai.Client {
	return openai.NewClient(token)
}

func AskGPT(ctx context.Context, openaiClient *openai.Client, operationService *service.OperationService,
	update *tgbotapi.Update) (string, error) {

	message := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: update.Message.Text,
	}
	err := addMessagesContext(ctx, operationService, update.SentFrom().ID, message)
	if err != nil {
		return "", err
	}

	resp, err := openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: conversationMessages[update.SentFrom().ID],
		},
	)
	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}
	err = addMessagesContext(nil, operationService, update.SentFrom().ID, resp.Choices[0].Message)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func addMessagesContext(ctx context.Context, operationService *service.OperationService, userID int64, msg openai.ChatCompletionMessage) error {
	msgs := conversationMessages[userID]
	if msgs == nil || len(msgs) == 0 {
		err := fillDialogStart(ctx, operationService, userID)
		if err != nil {
			return err
		}
	}
	conversationMessages[userID] = append(conversationMessages[userID], msg)
	return nil
}

func fillDialogStart(ctx context.Context, operationService *service.OperationService, userID int64) error {
	msgs := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "you are a warehouse assistant",
		},
	}
	historyCSV, err := operationService.GetOperationsHistory(ctx)
	if err != nil {
		return err
	}
	msgs = append(msgs, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: historyCSV,
	})
	conversationMessages[userID] = msgs
	return nil
}

func CloseDialog(userID int64) {
	conversationMessages[userID] = nil
}
