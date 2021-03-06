package model

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
	"qilin-api/pkg/model/game"
	"qilin-api/pkg/model/utils"
	"time"
)

type (
	GameTag struct {
		ID    int64                 `gorm:"primary_key"`
		Title utils.LocalizedString `gorm:"type:jsonb; not null"`
	}

	GameGenre struct {
		GameTag
	}

	Descriptor struct {
		ID     uint                  `gorm:"primary_key"`
		Title  utils.LocalizedString `gorm:"type:jsonb; not null; default:'{}'"`
		System string                `gorm:"not null"`
	}

	Game struct {
		ID                   uuid.UUID             `gorm:"type:uuid; primary_key; default:gen_random_uuid()"`
		InternalName         string                `gorm:"unique; not null"`
		Title                string                `gorm:"type:text; not null"`
		Developers           string                `gorm:"type:text; not null"`
		Publishers           string                `gorm:"type:text; not null"`
		ReleaseDate          time.Time             `gorm:"type:timestamp; not null; default:now()"`
		DisplayRemainingTime bool                  `gorm:"type:boolean; not null"`
		AchievementOnProd    bool                  `gorm:"type:boolean; not null"`
		FeaturesCommon       pq.StringArray        `gorm:"type:text[]; not null; default:array[]::text[]"`
		FeaturesCtrl         string                `gorm:"type:text; not null"`
		Platforms            game.Platforms        `gorm:"type:jsonb; not null; default:'{}'"`
		Requirements         game.GameRequirements `gorm:"type:jsonb; not null; default:'{}'"`
		Languages            game.GameLangs        `gorm:"type:jsonb; not null; default:'{}'"`
		GenreMain            int64                 `gorm:"type:integer"`
		GenreAddition        pq.Int64Array         `gorm:"type:integer[]; not null; default:array[]::integer[]"`
		Tags                 pq.Int64Array         `gorm:"type:integer[]; not null; default:array[]::integer[]"`

		Vendor    *Vendor   /// VendorID is foreignKey for Vendor
		VendorID  uuid.UUID `gorm:"type:uuid"`
		Creator   *User     /// CreatorID is foreignKey for Creator
		CreatorID string

		CreatedAt time.Time  `gorm:"default:now()"`
		UpdatedAt time.Time  `gorm:"default:now()"`
		DeletedAt *time.Time `sql:"index"`

		DefaultPackageID uuid.UUID

		Product ProductEntry `gorm:"polymorphic:Entry;"`
	}

	GameDescr struct {
		gorm.Model
		Tagline               utils.LocalizedString `gorm:"type:jsonb; not null; default:'{}'"`
		Description           utils.LocalizedString `gorm:"type:jsonb; not null; default:'{}'"`
		Reviews               game.GameReviews      `gorm:"type:jsonb; not null; default:'[]'"`
		AdditionalDescription string                `gorm:"type:text; not null"`
		GameSite              string                `gorm:"type:text; not null"`
		Socials               game.Socials          `gorm:"type:jsonb; not null; default:'{}'"`

		Game   *Game
		GameID uuid.UUID
	}

	ShortGameInfo struct {
		Game
		Price
	}

	ProductGameImpl struct {
		Game
		Media
	}

	// GameService is a helper service class to interact with Game object.
	GameService interface {
		CreateTags([]GameTag) error
		CreateGenres([]GameGenre) error

		GetTags([]string) ([]GameTag, error)
		GetGenres([]string) ([]GameGenre, error)
		GetRatingDescriptors(system string) ([]Descriptor, error)
		FindTags(userId string, title string, limit, offset int) ([]GameTag, error)
		FindGenres(userId string, title string, limit, offset int) ([]GameGenre, error)

		Create(userId string, vendorId uuid.UUID, internalName string) (*Game, error)
		Delete(userId string, gameId uuid.UUID) error
		GetList(userId string, vendorId uuid.UUID, offset, limit int, internalName, genre, releaseDate, sort string) ([]*ShortGameInfo, error)
		GetInfo(gameId uuid.UUID) (*Game, error)
		UpdateInfo(game *Game) error
		GetDescr(gameId uuid.UUID) (*GameDescr, error)
		UpdateDescr(descr *GameDescr) error
		GetProduct(gameId uuid.UUID) (Product, error)
	}
)

func (ShortGameInfo) TableName() string {
	return "games"
}

func (ProductGameImpl) TableName() string {
	return "games"
}

func (p *ProductGameImpl) GetID() uuid.UUID {
	return p.Game.ID
}

func (p *ProductGameImpl) GetName() string {
	return p.Game.InternalName
}

func (p *ProductGameImpl) GetType() ProductType {
	return ProductGame
}

func (p *ProductGameImpl) GetImage() (res *utils.LocalizedString) {
	return &p.Media.CoverImage
}
