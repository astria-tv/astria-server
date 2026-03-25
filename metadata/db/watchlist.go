package db

import (
	"github.com/jinzhu/gorm"
)

// WatchlistItem tracks a media item (movie or series) that a user has added to their watchlist.
type WatchlistItem struct {
	gorm.Model
	UUIDable
	UserID    uint   `gorm:"unique_index:idx_unique_watchlist_per_media"`
	MediaUUID string `gorm:"unique_index:idx_unique_watchlist_per_media"`
	MediaType string
}

// AddWatchlistItem adds a media item to a user's watchlist. If the item already exists, it is a no-op.
func AddWatchlistItem(item *WatchlistItem) error {
	existing := WatchlistItem{}
	if db.Where("user_id = ? AND media_uuid = ?", item.UserID, item.MediaUUID).First(&existing).RecordNotFound() {
		return db.Create(item).Error
	}
	return nil
}

// DeleteWatchlistItem removes a media item from a user's watchlist.
func DeleteWatchlistItem(mediaUUID string, userID uint) error {
	return db.Unscoped().Where("user_id = ? AND media_uuid = ?", userID, mediaUUID).Delete(&WatchlistItem{}).Error
}

// GetWatchlistItems returns all watchlist items for a given user.
func GetWatchlistItems(userID uint) ([]WatchlistItem, error) {
	var items []WatchlistItem
	err := db.Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error
	return items, err
}

// IsOnWatchlistByUUIDs returns a map indicating which of the given media UUIDs are on the user's watchlist.
func IsOnWatchlistByUUIDs(mediaUUIDs []string, userID uint) map[string]bool {
	result := make(map[string]bool)
	if len(mediaUUIDs) == 0 {
		return result
	}

	var items []WatchlistItem
	db.Where("user_id = ? AND media_uuid IN (?)", userID, mediaUUIDs).Find(&items)
	for _, item := range items {
		result[item.MediaUUID] = true
	}
	return result
}
