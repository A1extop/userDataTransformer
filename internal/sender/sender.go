// Отвечает за отправку данных в другие сервисы
package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"userDataTransformer/internal/models"
)

type RemoteSender struct {
	endpoint string
	client   *http.Client
}

func NewRemoteSender(endpoint string) IRemoteSender {
	return &RemoteSender{
		endpoint: endpoint,
		client:   &http.Client{},
	}
}

// Захотелось немного защитить от падения другого сервера и потери данных сразу, поэтому 500 в последствии снова отправляется
func (s *RemoteSender) SendUser(ctx context.Context, user models.JSONUser) error {
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusInternalServerError {
		return fmt.Errorf("bad response from remote: %s", resp.Status)
	}

	return nil
}
