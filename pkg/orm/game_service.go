package orm

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"net/http"
	"qilin-api/pkg/model"
	bto "qilin-api/pkg/model/game"
	"strings"
	"time"
)

// GameService is service to interact with database and Game object.
type GameService struct {
	db *gorm.DB
}

// NewGameService initialize this service.
func NewGameService(db *Database) (*GameService, error) {
	return &GameService{db.database}, nil
}

func (p *GameService) verify_UserAndVendor(userId, vendorId uuid.UUID) (err error) {
	foundVendor := -1
	err = p.db.Table("vendor_users").Where("user_id = ? and vendor_id = ?", userId, vendorId).Count(&foundVendor).Error
	if err != nil {
		return errors.Wrap(err, "Verify vendor")
	}
	if foundVendor == 0 {
		return NewServiceError(404, "Vendor not found")
	}
	return
}

func (p *GameService) GetTags(ids []string) (tags []model.GameTag, err error) {
	stmt := p.db
	if ids != nil && len(ids) > 0 {
		stmt = stmt.Where("ID in (?)", ids)
	}
	err = stmt.Order("id").Find(&tags).Error
	if err != nil {
		return nil, errors.Wrap(err, "While fetch tags")
	}
	return
}

func (p *GameService) GetGenres(ids []string) (genres []model.GameGenre, err error) {
	stmt := p.db
	if ids != nil && len(ids) > 0 {
		stmt = stmt.Where("ID in (?)", ids)
	}
	err = stmt.Order("id").Find(&genres).Error
	if err != nil {
		return nil, errors.Wrap(err, "While fetch genres")
	}
	return
}

func (p *GameService) GetRatingDescriptors(system string) (items []model.Descriptor, err error) {
	query := p.db.Order("title ->> 'en'")

	if system != "" {
		query = query.Where("system = ?", system)
	}

	err = query.Find(&items).Error
	if err != nil {
		return nil, errors.Wrap(err, "Fetch rating descriptors")
	}
	return
}

func (p *GameService) FindTags(userId uuid.UUID, title string, limit, offset int) (tags []model.GameTag, err error) {
	stmt := p.db
	if title != "" {
		user := model.User{}
		err = p.db.Select("id, lang").Where("id = ?", userId).First(&user).Error
		if err != nil {
			return nil, errors.Wrap(err, "while fetch user")
		}
		stmt = stmt.Where("title ->> ? ilike ?", user.Lang, title)
	}
	if limit > 0 {
		stmt = stmt.Limit(limit).Offset(offset)
	}
	err = stmt.Order("id").Find(&tags).Error
	if err != nil {
		return nil, errors.Wrap(err, "While fetch tags")
	}
	return
}
func (p *GameService) FindGenres(userId uuid.UUID, title string, limit, offset int) (genres []model.GameGenre, err error) {
	stmt := p.db
	if title != "" {
		user := model.User{}
		err = p.db.Select("lang").Where("id = ?", userId).First(&user).Error
		if err != nil {
			return nil, errors.Wrap(err, "while fetch user")
		}
		stmt = stmt.Where("title ->> ? ilike ?", user.Lang, title)
	}
	if limit > 0 {
		stmt = stmt.Limit(limit).Offset(offset)
	}
	err = stmt.Order("id").Find(&genres).Error
	if err != nil {
		return nil, errors.Wrap(err, "While fetch genres")
	}
	return
}

// Creates new Game object in database
func (p *GameService) Create(userId uuid.UUID, vendorId uuid.UUID, internalName string) (item *model.Game, err error) {

	if err := p.verify_UserAndVendor(userId, vendorId); err != nil {
		return nil, err
	}

	item = &model.Game{}
	errE := p.db.First(item, `internal_name ilike ?`, internalName).Error
	if errE == nil {
		return nil, NewServiceError(400, "Name already in use")
	}

	item.ID = uuid.NewV4()
	item.InternalName = internalName
	item.FeaturesCtrl = ""
	item.FeaturesCommon = []string{}
	item.Platforms = bto.Platforms{}
	item.Requirements = bto.GameRequirements{}
	item.Languages = bto.GameLangs{}
	item.FeaturesCommon = []string{}
	item.Genre = []string{}
	item.Tags = []string{}
	item.VendorID = vendorId
	item.CreatorID = userId

	err = p.db.Create(item).Error
	if err != nil {
		return nil, errors.Wrap(err, "While create new game")
	}

	err = p.db.Create(&model.GameDescr{
		Game:    item,
		Reviews: []bto.GameReview{},
	}).Error
	if err != nil {
		return nil, errors.Wrap(err, "Create descriptions for game")
	}

	return
}

func (p *GameService) GetList(userId uuid.UUID, vendorId uuid.UUID,
	offset, limit int, internalName, genre, releaseDate, sort string, price float64) (list []*model.Game, err error) {

	if err := p.verify_UserAndVendor(userId, vendorId); err != nil {
		return nil, err
	}

	user := model.User{}
	err = p.db.Select("lang, currency").Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, errors.Wrap(err, "while fetch user")
	}

	conds := []string{}
	vals := []interface{}{}

	if internalName != "" {
		conds = append(conds, `internal_name ilike ?`)
		vals = append(vals, internalName)
	}

	if genre != "" {
		genres := []model.GameGenre{}
		err = p.db.Where("title ->> ? ilike ?", user.Lang, genre).Limit(1).Find(&genres).Error
		if err != nil {
			return nil, errors.Wrap(err, "while fetch genres")
		}
		if len(genres) == 0 {
			return // 200: No any genre found
		}
		conds = append(conds, "? = ANY(genre)")
		vals = append(vals, genres[0].ID)
	}

	if releaseDate != "" {
		rdate, err := time.Parse("2006-01-02", releaseDate)
		if err != nil {
			return nil, NewServiceError(400, "Invalid date")
		}
		conds = append(conds, `date(release_date) = ?`)
		vals = append(vals, rdate)
	}

	if price > 0 {
		conds = append(conds, `game_prices.value = ?`)
		vals = append(vals, price)
	}

	conds = append(conds, `vendor_id = ?`)
	vals = append(vals, vendorId)

	var orderBy interface{}
	orderBy = "created_at ASC"
	if sort != "" {
		switch sort {
		case "-genre":
			orderBy = "created_at DESC"
		case "+genre":
			orderBy = "genre ASC"
		case "-releaseDate":
			orderBy = "release_date DESC"
		case "+releaseDate":
			orderBy = "release_date ASC"
		case "-price":
			orderBy = "game_prices.value DESC"
		case "+price":
			orderBy = "game_prices.value ASC"
		case "-name":
			orderBy = "internal_name DESC"
		case "+name":
			orderBy = "internal_name ASC"
		}
	}

	err = p.db.
		Model(model.Game{}).
		// TODO: relate with prices
		//Joins("LEFT JOIN game_prices on game_prices.game_id = game.id and game_prices.currency = ?", vendor.Currency).
		Where(strings.Join(conds, " and "), vals...).
		Order(orderBy).
		Limit(limit).
		Offset(offset).
		Find(&list).Error
	if err != nil {
		return nil, errors.Wrap(err, "Fetch games list")
	}

	return
}

func (p *GameService) GetInfo(userId uuid.UUID, gameId uuid.UUID) (game *model.Game, err error) {

	game = &model.Game{}
	err = p.db.First(game, `id = ? and vendor_id in (select vendor_id from vendor_users where user_id = ?)`, gameId, userId).Error
	if err == gorm.ErrRecordNotFound {
		return nil, NewServiceError(404, "Game not found")
	} else if err != nil {
		return nil, errors.Wrap(err, "Fetch game info")
	}

	return game, nil
}

func (p *GameService) Delete(userId uuid.UUID, gameId uuid.UUID) (err error) {

	game, err := p.GetInfo(userId, gameId)
	if err != nil {
		return err
	}

	err = p.db.Delete(game).Error
	if err != nil {
		return errors.Wrap(err, "Delete game")
	}

	return nil
}

func (p *GameService) UpdateInfo(userId uuid.UUID, game *model.Game) (err error) {

	gameSrc, err := p.GetInfo(userId, game.ID)
	if err != nil {
		return err
	}
	game.CreatorID = gameSrc.CreatorID
	game.VendorID = gameSrc.VendorID
	game.CreatedAt = gameSrc.CreatedAt
	game.UpdatedAt = time.Now()
	err = p.db.Save(game).Error
	if err != nil && strings.Index(err.Error(), "duplicate key value") > -1 {
		return NewServiceError(http.StatusConflict, "Invalid internal_name")
	} else if err != nil {
		return errors.Wrap(err, "Update game")
	}

	return nil
}

func (p *GameService) GetDescr(userId uuid.UUID, gameId uuid.UUID) (descr *model.GameDescr, err error) {
	game, err := p.GetInfo(userId, gameId)
	if err != nil {
		return nil, err
	}
	descr = &model.GameDescr{
		Reviews: []bto.GameReview{},
	}
	err = p.db.Model(game).Related(descr).Error
	if err != nil {
		return nil, errors.Wrap(err, "Fetch game descr")
	}
	return descr, nil
}

func (p *GameService) UpdateDescr(userId uuid.UUID, descr *model.GameDescr) (err error) {
	game, err := p.GetInfo(userId, descr.GameID)
	if err != nil {
		return err
	}
	update := *descr
	if update.ID == 0 {
		found := model.GameDescr{}
		err = p.db.Model(game).Related(&found).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			update.CreatedAt = time.Now()
		} else if err != nil {
			return errors.Wrap(err, "Get game descr")
		} else {
			update.ID = found.ID
			update.CreatedAt = found.CreatedAt
		}
	}
	update.UpdatedAt = time.Now()
	err = p.db.Save(&update).Error
	if err != nil {
		return errors.Wrap(err, "Update game descr")
	}
	return
}

func (p *GameService) CreateTags(tags []model.GameTag) (err error) {
	for _, t := range tags {
		err = p.db.Create(&t).Error
		if err != nil {
			return errors.Wrap(err, "Create game tag")
		}
	}
	return
}
