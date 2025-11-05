package controller

import "github.com/Yoak3n/troll/scanner/model"

type VideoDataQuery struct {
	model.VideoTable
	model.UserTable
	Count int `json:"count"`
}

type VideoData struct {
	Avid        uint   `json:"avid"`
	Bvid        string `json:"bvid"`
	Title       string `json:"title"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Author      `json:"author"`
}

type VideoDataWithCommentsCount struct {
	Count int `json:"count"`
	VideoData
}

type Author struct {
	Uid      uint   `json:"uid"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Location string `json:"location"`
}

type CommentWithUser struct {
	model.CommentTable
	Uid      uint   `json:"uid"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Location string `json:"location"`
}

func (d *Database) GetVideoByAvid(avid uint) VideoData {
	videoDataQuery := VideoDataQuery{}

	query := `
	SELECT v.*,u.*
	FROM video_tables v 
	INNER JOIN user_tables u ON v.owner = u.uid
	WHERE v.avid = ?
	`
	if err := d.db.Raw(query, avid).Scan(&videoDataQuery).Error; err != nil {
		return VideoData{}
	}

	return VideoData{
		Avid:        videoDataQuery.VideoTable.Avid,
		Bvid:        videoDataQuery.VideoTable.Bvid,
		Cover:       videoDataQuery.VideoTable.Cover,
		Title:       videoDataQuery.VideoTable.Title,
		Topic:       videoDataQuery.VideoTable.Topic,
		Description: videoDataQuery.VideoTable.Description,
		Author: Author{
			Uid:      videoDataQuery.UserTable.Uid,
			Name:     videoDataQuery.UserTable.Username,
			Location: videoDataQuery.UserTable.Location,
			Avatar:   videoDataQuery.UserTable.Avatar,
		},
	}
}

func (d *Database) GetVideosByTopic(topicName string) []VideoDataWithCommentsCount {
	result := make([]VideoDataQuery, 0)
	query := `
	SELECT v.*,u.*, Count(c.comment_id) AS count
	FROM comment_tables c INNER JOIN video_tables v ON c.video_avid = v.avid LEFT  JOIN user_tables u ON v.owner = u.uid
	WHERE v.topic = ?
	GROUP BY v.avid, u.uid;`
	if err := d.db.Raw(query, topicName).Scan(&result).Error; err != nil {
		return nil
	}
	videosMap := make(map[uint]int)

	//
	// query2 := `
	// SELECT
	// 	c.*,
	// 	u.uid,
	// 	u.username,
	// 	u.avatar,
	// 	u.location
	// FROM comment_tables c
	// INNER JOIN video_tables v ON c.video_avid = v.avid
	// INNER JOIN user_tables u ON c.owner = u.uid
	// WHERE v.topic = ?
	// `
	// result2 := make([]CommentWithUser, 0)
	// if err := d.db.Raw(query2, topicName).Scan(&result2).Error; err != nil {
	// 	return nil
	// }
	// for _, ret := range result1 {
	// 	if _, exists := videosMap[ret.Avid]; !exists {
	// 		vdwc := &VideoDataWithComments{
	// 			Avid:        ret.Avid,
	// 			Bvid:        ret.Bvid,
	// 			Title:       ret.Title,
	// 			Topic:       ret.Topic,
	// 			Description: ret.Description,
	// 			Count:       0,
	// 			Comments:    make([]CommentData, 0),
	// 			Author: Author{
	// 				Uid:      ret.UserTable.Uid,
	// 				Name:     ret.UserTable.Username,
	// 				Avatar:   ret.UserTable.Avatar,
	// 				Location: ret.UserTable.Location,
	// 			},
	// 		}
	// 		videosMap[ret.Avid] = vdwc
	// 	}
	// 	cd := CommentData{
	// 		CommentId: ret.CommentId,
	// 		Content:   ret.Text,
	// 		Owner: Author{
	// 			Uid: ret.UserTable.Uid,
	// 		},
	// 	}
	// 	videosMap[ret.Avid].Comments = append(videosMap[ret.Avid].Comments, cd)
	// 	videosMap[ret.Avid].Count++
	// }

	for _, ret := range result {
		if _, exists := videosMap[ret.Avid]; !exists {
			videosMap[ret.Avid] = 1
		} else {
			videosMap[ret.Avid]++
		}
	}
	ret := make([]VideoDataWithCommentsCount, 0)
	for _, video := range result {
		vdwc := VideoDataWithCommentsCount{
			VideoData: VideoData{
				Avid:        video.Avid,
				Bvid:        video.Bvid,
				Title:       video.Title,
				Topic:       video.Topic,
				Description: video.Description,
				Author: Author{
					Uid:      video.UserTable.Uid,
					Name:     video.UserTable.Username,
					Avatar:   video.UserTable.Avatar,
					Location: video.UserTable.Location,
				},
			},
			Count: video.Count,
		}
		ret = append(ret, vdwc)
	}

	return ret
}
