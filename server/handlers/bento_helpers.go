package handlers

type bentoMetadata struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	OwnerID   string `json:"owner_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func getBentoMetadata() (*bentoMetadata, error) {
	return nil, nil
}

type bentoMetadataExtended struct {
	bentoMetadata
	IngredientCount  int      `json:"ingredient_count"`
	UsersWithAccess  []string `json:"users_with_access"`
	GroupsWithAccess []string `json:"groups_with_access"`
}

func getBentoMetadataExtended() (*bentoMetadataExtended, error) {
	return nil, nil
}
