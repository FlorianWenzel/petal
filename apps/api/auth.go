package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5"
)

const (
	sessionCookie = "petal_session"
	sessionTTL    = 30 * 24 * time.Hour
)

var usernameRegex = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{1,22}[a-z0-9])?$`)

type user struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type credentialsReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type renameReq struct {
	Username string `json:"username"`
}

type ctxKey int

const userCtxKey ctxKey = 1

func userFromCtx(ctx context.Context) *user {
	u, _ := ctx.Value(userCtxKey).(*user)
	return u
}

func (s *server) signup(w http.ResponseWriter, r *http.Request) {
	var body credentialsReq
	if err := decodeJSON(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body")
		return
	}
	if !validEmail(body.Email) || len(body.Password) < 8 {
		writeErr(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	ctx := r.Context()
	var existing string
	err := s.pool.QueryRow(ctx, `select id from users where email = $1`, body.Email).Scan(&existing)
	if err == nil {
		writeErr(w, http.StatusConflict, "email already registered")
		return
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		writeErr(w, http.StatusInternalServerError, "db")
		return
	}

	hash, err := argon2id.CreateHash(body.Password, argon2id.DefaultParams)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "hash")
		return
	}

	uname, err := s.allocateUsername(ctx)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "username")
		return
	}

	var u user
	err = s.pool.QueryRow(ctx,
		`insert into users (email, password_hash, username)
		 values ($1, $2, $3)
		 returning id, email, username`,
		body.Email, hash, uname,
	).Scan(&u.ID, &u.Email, &u.Username)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "create")
		return
	}

	if err := s.openSession(ctx, w, u.ID); err != nil {
		writeErr(w, http.StatusInternalServerError, "session")
		return
	}
	writeJSON(w, http.StatusCreated, u)
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	var body credentialsReq
	if err := decodeJSON(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body")
		return
	}

	ctx := r.Context()
	var u user
	var hash string
	err := s.pool.QueryRow(ctx,
		`select id, email, username, password_hash from users where email = $1`,
		body.Email,
	).Scan(&u.ID, &u.Email, &u.Username, &hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeErr(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		writeErr(w, http.StatusInternalServerError, "db")
		return
	}

	ok, err := argon2id.ComparePasswordAndHash(body.Password, hash)
	if err != nil || !ok {
		writeErr(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := s.openSession(ctx, w, u.ID); err != nil {
		writeErr(w, http.StatusInternalServerError, "session")
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (s *server) logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookie); err == nil && c.Value != "" {
		_, _ = s.pool.Exec(r.Context(), `delete from sessions where id = $1`, c.Value)
	}
	clearCookie(w, s.cfg.Production)
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, userFromCtx(r.Context()))
}

func (s *server) rename(w http.ResponseWriter, r *http.Request) {
	var body renameReq
	if err := decodeJSON(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body")
		return
	}
	uname := strings.ToLower(strings.TrimSpace(body.Username))
	if !usernameRegex.MatchString(uname) {
		writeErr(w, http.StatusBadRequest, "username must be 3-24 chars, lowercase letters, digits, hyphens")
		return
	}
	cur := userFromCtx(r.Context())

	var taken string
	err := s.pool.QueryRow(r.Context(),
		`select id from users where username = $1 and id != $2`,
		uname, cur.ID,
	).Scan(&taken)
	if err == nil {
		writeErr(w, http.StatusConflict, "username taken")
		return
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		writeErr(w, http.StatusInternalServerError, "db")
		return
	}

	var u user
	err = s.pool.QueryRow(r.Context(),
		`update users set username = $1 where id = $2 returning id, email, username`,
		uname, cur.ID,
	).Scan(&u.ID, &u.Email, &u.Username)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "update")
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (s *server) openSession(ctx context.Context, w http.ResponseWriter, userID string) error {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return err
	}
	sid := base64.RawURLEncoding.EncodeToString(buf)
	expires := time.Now().Add(sessionTTL)
	if _, err := s.pool.Exec(ctx,
		`insert into sessions (id, user_id, expires_at) values ($1, $2, $3)`,
		sid, userID, expires,
	); err != nil {
		return err
	}
	setCookie(w, sid, expires, s.cfg.Production)
	return nil
}

func setCookie(w http.ResponseWriter, sid string, expires time.Time, prod bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    sid,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   prod,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearCookie(w http.ResponseWriter, prod bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   prod,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionCookie)
		if err != nil || c.Value == "" {
			writeErr(w, http.StatusUnauthorized, "no session")
			return
		}
		var u user
		err = s.pool.QueryRow(r.Context(),
			`select users.id, users.email, users.username
			   from sessions
			   join users on users.id = sessions.user_id
			  where sessions.id = $1 and sessions.expires_at > now()`,
			c.Value,
		).Scan(&u.ID, &u.Email, &u.Username)
		if err != nil {
			writeErr(w, http.StatusUnauthorized, "invalid session")
			return
		}
		ctx := context.WithValue(r.Context(), userCtxKey, &u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}
