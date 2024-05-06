package maatidb

import (
	"database/sql"
	"time"
)

type CommentInfos struct {
	PostId    int
	CommentId int
	UserImage string
	UserName  string
	Date      time.Time
	LoveNumb  int
	HateNumb  int
	Content   string
}

func GetCommentsInfos() ([]CommentInfos, error) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT postId, commentId, userImage, userName, date, loveNumb, hateNumb, content FROM CommentsInfosView")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []CommentInfos{}
	for rows.Next() {
		var c CommentInfos

		err := rows.Scan(&c.PostId, &c.CommentId, &c.UserImage, &c.UserName, &c.Date, &c.LoveNumb, &c.HateNumb, &c.Content)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
