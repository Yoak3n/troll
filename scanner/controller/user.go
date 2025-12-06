package controller

import "github.com/Yoak3n/troll/scanner/model"

func (d *Database) QueryTopNUserInTopic(topic string, n int) ([]model.UserQuery, error) {
	var ret []model.UserQuery
	query := `
	SELECT
		u.*,
		COUNT(c.comment_id) AS count
	FROM
		user_tables u
		INNER JOIN
		comment_tables c ON u.uid = c.owner
		INNER JOIN
		video_tables v ON c.video_avid = v.avid
	WHERE
		v.topic = ?
	GROUP BY
		u.uid
	ORDER BY
	    count DESC
	LIMIT ?;`
	err := d.db.Raw(query, topic, n).Scan(&ret).Error
	return ret, err
}
