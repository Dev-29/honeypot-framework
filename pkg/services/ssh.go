package services

import (
	"honeypot-framework/pkg/database"
	"honeypot-framework/pkg/logger"
)

type SSHHoneypot struct {
	port     string
	banner   string
	logger   *logger.Logger
	database *database.Database
}

func NewSSHHoneypot(port, banner string, logger *logger.Logger, db *database.Database) *SSHHoneypot {
	return &SSHHoneypot{
		port:     port,
		banner:   banner,
		logger:   logger,
		database: db,
	}
}

// TODO
func (s *SSHHoneypot) Start() error {
	return nil
}
