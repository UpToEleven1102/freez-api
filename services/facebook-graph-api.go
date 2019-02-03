package services

import (
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	fb "github.com/huandu/facebook"
)

func GetFaceBookUserInfo(data models.FacebookTokenData) (interface{}, error) {
	result, err := fb.Get("/me", fb.Params{
		"fields": "name, email, picture.type(large)",
		"access_token": data.AccessToken,
	})

	if err != nil {
		return nil, err
	}

	var res models.FacebookUserInfo
	res.ID = result.Get("id").(string)
	res.Email = result.Get("email").(string)
	res.Name = result.Get("name").(string)
	res.Picture = result.Get("picture").(map[string]interface{})["data"].(map[string]interface{})["url"].(string)

	return res, nil
}
