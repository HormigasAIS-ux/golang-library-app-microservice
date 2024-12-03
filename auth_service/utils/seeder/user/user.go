package seeder_util

import (
	"auth_service/config"
	"auth_service/domain/model"
	"auth_service/repository"
	"fmt"

	"github.com/google/uuid"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("main")

func SeedUser(userRepo repository.IUserRepo) error {
	users := []model.User{}

	if config.Envs.INITIAL_ADMIN_USERNAME != "" && config.Envs.INITIAL_ADMIN_PASSWORD != "" {
		users = append(users, model.User{
			UUID:     uuid.New(),
			Username: config.Envs.INITIAL_ADMIN_USERNAME,
			Password: config.Envs.INITIAL_ADMIN_PASSWORD,
			Fullname: config.Envs.INITIAL_ADMIN_USERNAME,
			Email:    fmt.Sprint(config.Envs.INITIAL_ADMIN_USERNAME, "@gmail.com"),
			Role:     "admin",
		})
	} else {
		logger.Warningf("initial admin username and password not set")
	}

	if config.Envs.INITIAL_USER_USERNAME != "" && config.Envs.INITIAL_USER_PASSWORD != "" {
		users = append(users, model.User{
			UUID:     uuid.New(),
			Username: config.Envs.INITIAL_USER_USERNAME,
			Password: config.Envs.INITIAL_USER_PASSWORD,
			Fullname: config.Envs.INITIAL_USER_USERNAME,
			Email:    fmt.Sprint(config.Envs.INITIAL_USER_USERNAME, "@gmail.com"),
			Role:     "user",
		})
	} else {
		logger.Warningf("initial user username and password not set")
	}

	for _, user := range users {
		logger.Infof("seeding user: %s", user.Username)

		existing, _ := userRepo.GetByUsername(user.Username)
		if existing != nil {
			logger.Infof("user already exists: %s", user.Username)
			continue
		}

		err := userRepo.Create(&user)
		if err != nil {
			return err
		}

		logger.Infof("user seeded: %s", user.Username)
	}

	return nil
}
