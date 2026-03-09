package repositories

import "gosvelte/app/models"

type WishlistRepository struct {
	BaseRepository[models.Wishlist]
}
