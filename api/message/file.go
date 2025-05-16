package message

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func (s *MessageService) sendFile(ctx context.Context, msg *FileMessage, path string) (msgID string, fileID string, err error) {
	params := buildParams(&msg.Message)
	params.Set("caption", msg.Text)
	if msg.FileID != "" {
		params.Set("fileId", msg.FileID)
	}
	response := struct {
		Id          string `json:"msgId"`
		FileID      string `json:"fileId"`
		Ok          bool   `json:"ok"`
		Description string `json:"description"`
	}{}

	// Send POST; Upload file
	if msg.Contents != nil {
		buffer := &bytes.Buffer{}
		part := multipart.NewWriter(buffer)
		fileWriter, err := part.CreateFormFile("file", msg.Filename)
		if err != nil {
			return "", "", fmt.Errorf("unable to create file writer: %w", err)
		}
		_, err = io.Copy(fileWriter, msg.Contents)
		if err != nil {
			return "", "", fmt.Errorf("unable to copy file contents: %w", err)
		}
		part.Close()
		req, err := s.client.PerformRequest(ctx, http.MethodPost, path, params, buffer)
		if err != nil {
			return "", "", err
		}
		req.Header.Set("Content-Type", part.FormDataContentType())
		resp, err := s.client.Do(req)
		if err != nil {
			return "", "", fmt.Errorf("unable to upload file: %w", err)
		}
		defer resp.Body.Close()
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return "", "", fmt.Errorf("unable to decode response: %w", err)
		}
		return response.Id, response.FileID, nil
	}

	req, err := s.client.PerformRequest(ctx, http.MethodPost, path, params, nil)
	if err != nil {
		return "", "", err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("unable to upload file: %w", err)
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", "", fmt.Errorf("unable to decode response: %w", err)
	}
	if !response.Ok {
		return "", "", fmt.Errorf("%w: %s", ErrNotOk, response.Description)
	}
	return response.Id, msg.FileID, nil
}
