package models

type Category struct {
	ID           string                `json:"id" db:"id"`
	Name         string                `json:"name" db:"name"`
	Image        *string               `json:"image" db:"image"`
	ParentID     *string               `json:"parent_id,omitempty" db:"parent_id"`
	Sub          []*Category           `json:"subcategories,omitempty"`
	Translations []CategoryTranslation `json:"translations"`
}

type CategoryTranslation struct {
	CategoryID   string `json:"category_id" db:"category_id"`
	Name         string `json:"name" db:"name"`
	LanguageCode string `json:"language_code" db:"language"`
}

type CategorySwagger struct {
	ID           string                `json:"id" example:"a2d70daf-b4ac-4198-a6a0-999447483c18"`
	Name         string                `json:"name" example:"electronic"`
	Image        *string               `json:"image" example:"https://firebasestorage.googleapis.com/v0/b/ulab-market.appspot.com/o/test%2Fcategory%2Fa2d70daf-b4ac-4198-a6a0-999447483c18?alt=media&token=test%2Fcategory%2Fa2d70daf-b4ac-4198-a6a0-999447483c18"`
	ParentID     *string               `json:"parent_id,omitempty" example:""`
	Sub          []*SubCategorySwagger `json:"subcategories,omitempty"`
	Translations []CategoryTranslation `json:"translations"`
}

type SubCategorySwagger struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Image        *string               `json:"image"`
	ParentID     *string               `json:"parent_id,omitempty"`
	Translations []CategoryTranslation `json:"translations"`
}
