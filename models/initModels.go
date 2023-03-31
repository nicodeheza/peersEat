package models

func InitModels(databaseName string) {
	InitPeerModel(databaseName)
	InitRestaurantModel(databaseName)
}
