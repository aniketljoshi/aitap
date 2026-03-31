package export

import (
	"encoding/json"
	"os"

	"github.com/aniketljoshi/aitap/internal/model"
	"github.com/aniketljoshi/aitap/internal/redact"
)

// ToJSONL writes all session calls to a JSONL file.
func ToJSONL(session *model.Session, path string, redactSecrets bool) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, call := range session.Calls {
		c := *call // copy
		if redactSecrets {
			c.RequestBody = redact.Redact(c.RequestBody)
			c.ResponseBody = redact.Redact(c.ResponseBody)
			c.ResponseText = redact.Redact(c.ResponseText)
			c.SystemPrompt = redact.Redact(c.SystemPrompt)
			for i := range c.Messages {
				c.Messages[i].Content = redact.Redact(c.Messages[i].Content)
			}
		}
		if err := enc.Encode(c); err != nil {
			return err
		}
	}
	return nil
}
