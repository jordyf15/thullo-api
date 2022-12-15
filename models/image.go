package models

import (
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

type Image struct {
	Width uint   `json:"width" bson:"width"`
	ID    string `json:"id,omitempty" bson:"id"`
	URL   string `json:"url,omitempty" bson:"-"`
}

type Images []*Image

func (i *Images) UnmarshalBSON(value interface{}) error {
	type ImagesAlias Images
	b, ok := value.([]byte)
	if b == nil {
		return nil
	}

	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return bson.Unmarshal(b, (*ImagesAlias)(i))
}

func (i *Images) MarshalBSON() ([]byte, error) {
	if len(*i) == 0 {
		return nil, nil
	}

	type ImagesAlias Images
	return bson.Marshal((*ImagesAlias)(i))
}

func (image *Image) MarshalJSON() ([]byte, error) {
	type Alias Image
	newStruct := &struct {
		*Alias
	}{
		Alias: (*Alias)(image),
	}

	if len(newStruct.URL) > 0 {
		newStruct.ID = ""
	}

	return json.Marshal(newStruct)
}
