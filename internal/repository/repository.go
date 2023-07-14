package repository

import "github.com/OoThan/usermanagement/internal/ds"

type Repository struct {
	DS   *ds.DataSource
	User *userRepository
}

type RepoConfig struct {
	DS *ds.DataSource
}

func NewRepository(rConfig *RepoConfig) *Repository {
	return &Repository{
		DS:   rConfig.DS,
		User: newUserRespository(rConfig),
	}
}
