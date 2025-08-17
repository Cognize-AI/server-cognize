package config

func SyncDB() {
	//DB.Exec(`DROP INDEX IF EXISTS idx_view_events_shared_collection_id;`)
	err := DB.AutoMigrate()
	if err != nil {
		return
	}
}
