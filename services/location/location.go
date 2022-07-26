package location

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"some-api/utils/db"
)

type Location struct {
	Pk        string
	PersonID  string
	Latitude  string
	Longitude string
}

func Get(store db.DataStore, personID string) (*Location, error) {
	out, err := store.GetByPk(personID)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, nil
	}

	var location Location
	err = attributevalue.UnmarshalMap(out, &location)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func Save(store db.DataStore, personID string, lat string, long string) error {
	err := store.UpdateItem(personID, map[string]string{
		"personID":  personID,
		"latitude":  lat,
		"longitude": long,
	})

	if err != nil {
		return err
	}

	return nil
}
