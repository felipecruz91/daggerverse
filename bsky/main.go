// A generated module for Bsky functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"bytes"
	"context"
	"dagger/bsky/internal/dagger"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Bsky struct{}

func (m *Bsky) Post(
	ctx context.Context,
	// pdsURL is the URL of the PDS (Personal Data Server)
	// +default="https://bsky.social"
	pdsURL string,
	// handle is the user handle
	// +required
	handle string,
	// password is the user password
	// +required
	password *dagger.Secret,
	// text is the post's content
	// +required
	text string,
) error {
	plainPasswd, err := password.Plaintext(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("error reading password from secret: %w", err))
	}
	s, err := createSession(pdsURL, handle, plainPasswd)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating session: %w", err))
	}

	log.Printf("Session created successfully, user id: %s\n", s.UserID)

	post := Post{
		Type:      "app.bsky.feed.post",
		Text:      text,
		CreatedAt: time.Now().Format(time.RFC3339),
		Langs:     []string{}, // TODO
	}
	if err := publishPost(pdsURL, s, &post); err != nil {
		log.Fatal(fmt.Errorf("failed to publish post: %w", err))
	}

	return nil
}

// SessionResponse holds authentication session information after a successful login.
type SessionResponse struct {
	AccessToken string `json:"accessJwt"`
	UserID      string `json:"did"`
}

type Post struct {
	Type      string   `json:"$type"`
	Text      string   `json:"text"`
	CreatedAt string   `json:"createdAt"`
	Langs     []string `json:"langs,omitempty"`
}

func createSession(pdsURL, handle, password string) (*SessionResponse, error) {
	loginURL := fmt.Sprintf("%s/xrpc/com.atproto.server.createSession", pdsURL)
	requestBody, err := json.Marshal(map[string]string{
		"identifier": handle,
		"password":   password,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create session, status code: %d", resp.StatusCode)
	}

	var sessionResponse SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse, nil
}

func publishPost(pdsURL string, session *SessionResponse, post *Post) error {
	postURL := fmt.Sprintf("%s/xrpc/com.atproto.repo.createRecord", pdsURL)
	postData, err := json.Marshal(map[string]interface{}{
		"repo":       session.UserID,
		"collection": "app.bsky.feed.post",
		"record":     post,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal post data: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, postURL, bytes.NewBuffer(postData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+session.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to publish post, status code: %d, resp. body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
