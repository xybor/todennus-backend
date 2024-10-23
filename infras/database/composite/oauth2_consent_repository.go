package composite

import (
	"github.com/redis/go-redis/v9"
	dbgorm "github.com/xybor/todennus-backend/infras/database/gorm"
	dbredis "github.com/xybor/todennus-backend/infras/database/redis"
	"gorm.io/gorm"
)

type OAuth2ConsentRepository struct {
	*dbgorm.OAuth2ConsentRepository
	*dbredis.OAuth2ConsentResultRepository
}

func NewOAuth2ConsentRepository(db *gorm.DB, redis *redis.Client) *OAuth2ConsentRepository {
	return &OAuth2ConsentRepository{
		OAuth2ConsentRepository:       dbgorm.NewOAuth2ConsentRepository(db),
		OAuth2ConsentResultRepository: dbredis.NewOAuth2ConsentResultRepository(redis),
	}
}
