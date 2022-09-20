package models

type Updates struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		MessageID int64 `json:"message_id"`
		From      struct {
			ID           int64  `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int64  `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

//type CallbackQuery struct {
//	ID   string `json:"id"`
//	From struct {
//		ID           int    `json:"id"`
//		IsBot        bool   `json:"is_bot"`
//		FirstName    string `json:"first_name"`
//		Username     string `json:"username"`
//		LanguageCode string `json:"language_code"`
//	} `json:"from"`
//	Message struct {
//		MessageID int `json:"message_id"`
//		From      struct {
//			ID        int64  `json:"id"`
//			IsBot     bool   `json:"is_bot"`
//			FirstName string `json:"first_name"`
//			Username  string `json:"username"`
//		} `json:"from"`
//		Chat struct {
//			ID        int    `json:"id"`
//			FirstName string `json:"first_name"`
//			Username  string `json:"username"`
//			Type      string `json:"type"`
//		} `json:"chat"`
//		Date        int    `json:"date"`
//		Text        string `json:"text"`
//		ReplyMarkup struct {
//			InlineKeyboard [][]struct {
//				Text         string `json:"text"`
//				CallbackData string `json:"callback_data"`
//			} `json:"inline_keyboard"`
//		} `json:"reply_markup"`
//	} `json:"message"`
//	ChatInstance string `json:"chat_instance"`
//	Data         string `json:"data"`
//}
