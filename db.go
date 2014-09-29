package circle

import (
	"encoding/json"
	"fmt"
	"github.com/azer/go-flickr"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	DB *leveldb.DB
)

func CreateDBConn(path string) {
	var err error
	DB, err = leveldb.OpenFile(path, nil)

	if err != nil {
		panic(err)
	}
}

func ReadFavs(ownerId string) ([]flickr.Fav, error) {
	iter := DB.NewIterator(util.BytesPrefix([]byte(fmt.Sprintf("favs:%s:", ownerId))), nil)
	result := []flickr.Fav{}

	for iter.Next() {
		fav := flickr.Fav{}
		err := json.Unmarshal(iter.Value(), &fav)

		if err != nil {
			continue
		}

		result = append(result, fav)
	}

	iter.Release()
	err := iter.Error()
	return result, err
}

func SaveFav(fav flickr.Fav) error {
	value, err := json.Marshal(&fav)

	if err != nil {
		return err
	}

	key := fmt.Sprintf("favs:%s:%s", fav.FavedBy, fav.DateFaved)

	err = DB.Put([]byte(key), value, nil)

	if err != nil {
		return err
	}

	return nil
}
