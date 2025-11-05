package controller

import "github.com/Yoak3n/troll/scanner/model"

type SearchOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type SearchOptionResponse struct {
	Typ     string         `json:"type"`
	Options []SearchOption `json:"options"`
}

func (d *Database) GetSearchOptionsByUid(uid uint, typ string) SearchOptionResponse {
	videos := make([]model.VideoTable, 0)
	d.db.Model(&videos)

	if err := d.db.Joins("JOIN comment_tables ON video_tables.avid = comment_tables.video_avid").
		Joins("JOIN user_tables ON comment_tables.owner = user_tables.uid").
		Where("user_tables.uid = ?", uid).
		Distinct().
		Find(&videos).Error; err != nil {
		return SearchOptionResponse{}
	}
	return videoToSearchOption(videos, typ)
}

func (d *Database) GetSearchOptionsByUserName(name string, typ string) SearchOptionResponse {
	videos := make([]model.VideoTable, 0)
	if err := d.db.Joins("JOIN comment_tables ON video_tables.avid = comment_tables.video_avid").
		Joins("JOIN user_tables ON comment_tables.owner = user_tables.uid").
		Where("user_tables.username = ?", name).
		Distinct().
		Find(&videos).Error; err != nil {
		return SearchOptionResponse{}
	}
	return videoToSearchOption(videos, typ)
}

func videoToSearchOption(videos []model.VideoTable, typ string) SearchOptionResponse {
	ret := SearchOptionResponse{
		Typ:     typ,
		Options: make([]SearchOption, 0),
	}
	switch typ {
	case "topic":
		topic2Videos := make(map[string]bool)
		for _, v := range videos {
			_, ok := topic2Videos[v.Topic]
			if !ok {
				topic2Videos[v.Topic] = true
				option := SearchOption{
					Label: v.Topic,
					Value: v.Topic,
				}
				ret.Options = append(ret.Options, option)
			}
		}
	case "video":
		for _, v := range videos {
			option := SearchOption{
				Label: v.Title,
				Value: v.Bvid,
			}
			ret.Options = append(ret.Options, option)
		}
	}
	return ret
}

func (d *Database) GetSearchOptionsFromAllVideos(typ string) SearchOptionResponse {
	videos := make([]model.VideoTable, 0)
	d.db.Find(&videos)
	return videoToSearchOption(videos, typ)
}

func (d *Database) GetUserCommentsByRange(uid uint, name string, rangeType string, rangeData []string) {

}
