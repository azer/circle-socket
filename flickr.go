package circle

import (
	"encoding/json"
	"github.com/azer/go-flickr"
	. "github.com/visionmedia/go-debug"
	"os"
)

var debug = Debug("flickr")
var client *flickr.Client

func SubscribeTo(username string, ch chan string) error {
	userId, err := FindUserId(username)

	if err != nil {
		return err
	}

	following, err := GetFollowing(userId)

	if err != nil {
		return err
	}

	debug("Fetching %s's circle (%d) photos.", username, len(following))

	for _, user := range following {
		favs, err := GetFavs(user.Id)

		if err != nil {
			debug("Failed to get %'s favs", user.Id)
			continue
		}

		slice, err := json.Marshal(favs)

		if err != nil {
			debug("Unable to parse %s's favs", user.Id)
			continue
		}

		ch <- string(slice)
	}

	close(ch)

	debug("Done (%s, %d)", username, len(following))

	return nil
}

func CreateFlickrClient() {
	client = &flickr.Client{
		Key: os.Getenv("FLICKR_API_KEY"),
	}
}

func FindUserId(username string) (string, error) {
	user, err := client.FindUser(username)

	if err != nil {
		return "", err
	}

	return user.Id, nil
}

func GetFollowing(userId string) ([]flickr.User, error) {
	users, err := client.Following(userId)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetFavs(userId string) ([]flickr.Fav, error) {
	favs, err := ReadFavs(userId)

	if err == nil && len(favs) > 0 {
		//debug("we got %s's favs (%d) cached already.", userId, len(favs))
		return favs, nil
	}

	favs, err = client.Favs(userId)

	if err != nil {
		return nil, err
	}

	for _, fav := range favs {
		err = SaveFav(fav)

		if err != nil {
			debug("Unable to save %s", fav.Id)
		}
	}

	return favs, nil
}
