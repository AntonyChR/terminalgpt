package openaiservice

import (
	"os"
	"testing"
)

func TestSaveConversationAsTxtFile(t *testing.T) {
	chat := NewChat(ApiConfiguration{
		Apikey: "",
		ApiUrl: "",
	})

	chat.AddMessageAsUser("Hello")
	chat.AddMessageAsAssistant("Hi there!")

	chat.AddMessageAsUser("How are you?")
	chat.AddMessageAsAssistant("I'm doing well, thank you!")

	chat.AddMessageAsUser("Goodbye")
	chat.AddMessageAsAssistant("Goodbye!")

	CONVERSATIONS_TEST_DIR := "./conversations_test"

	if err := chat.Save("test", "test.txt", CONVERSATIONS_TEST_DIR); err != nil {
		t.Errorf("Error: %v", err)
	}

	// remove test directory
	if err := os.RemoveAll(CONVERSATIONS_TEST_DIR); err != nil {
		t.Errorf("Error: %v", err)
	}
}
