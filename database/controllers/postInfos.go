package maatidb

import (
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type PostInfo struct {
	PostId        int
	UserImage     string
	UserName      string
	PostDate      time.Time
	LoveNumb      int
	HateNumb      int
	Title         string
	Description   string
	Categories    string
	CategoriesTab []string
}

// function to get all post infos for the view
// ordered by date
func GetPostInfos() ([]PostInfo, error) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT postId, userImage, userName, postDate, loveNumb, hateNumb, title, description, categoryNames FROM ExtendedPostsInfosView")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []PostInfo{}
	for rows.Next() {
		var p PostInfo
		err := rows.Scan(&p.PostId, &p.UserImage, &p.UserName, &p.PostDate, &p.LoveNumb, &p.HateNumb, &p.Title, &p.Description, &p.Categories)
		if err != nil {
			return nil, err
		}
		p.CategoriesTab = strings.Split(p.Categories, ",")

		posts = append(posts, p)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return posts, nil
}
