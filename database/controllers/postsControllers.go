package maatidb

import (
	"database/sql"
	"log"
	"strings"
)


// GetPostByCategories récupère les informations des posts filtrés par catégories.
func GetPostByCategories(categories []string) ([]PostInfo, error) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prépare la déclaration SQL pour interroger les données en incluant les filtres de catégorie
	query := `
        SELECT postId, userImage, userName, postDate, loveNumb, hateNumb, title, description, categoryNames
        FROM ExtendedPostsInfosView
        WHERE categoryNames LIKE ?
    `
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var posts []PostInfo
	// Pour chaque catégorie, exécutez la requête et ajoutez les résultats
	for _, cat := range categories {
		catPattern := "%" + cat + "%"
		rows, err := stmt.Query(catPattern)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var p PostInfo
			err := rows.Scan(&p.PostId, &p.UserImage, &p.UserName, &p.PostDate, &p.LoveNumb, &p.HateNumb, &p.Title, &p.Description, &p.Categories)
			if err != nil {
				rows.Close() // Assurez-vous de fermer rows avant de retourner
				return nil, err
			}
			p.CategoriesTab = strings.Split(p.Categories, ",")
			posts = append(posts, p)
		}
		rows.Close()
	}

	return posts, nil
}

// GetPostByUserName récupère les informations des posts filtrés par le nom d'utilisateur.
func GetPostByUserName(userName string) ([]PostInfo, error) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prépare la déclaration SQL pour interroger les données en incluant les filtres d'utilisateur
	query := `
	    SELECT postId, userImage, userName, postDate, loveNumb, hateNumb, title, description, categoryNames
	    FROM ExtendedPostsInfosView
	    WHERE userName = ?
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Exécute la requête avec le nom d'utilisateur spécifié
	rows, err := stmt.Query(userName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostInfo
	for rows.Next() {
		var p PostInfo
		err := rows.Scan(&p.PostId, &p.UserImage, &p.UserName, &p.PostDate, &p.LoveNumb, &p.HateNumb, &p.Title, &p.Description, &p.Categories)
		if err != nil {
			return nil, err
		}
		p.CategoriesTab = strings.Split(p.Categories, ",")
		posts = append(posts, p)
	}

	// Vérifie s'il y a des erreurs rencontrées pendant l'itération
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostById(id int) ([]PostInfo, error) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
	    SELECT postId, userImage, userName, postDate, loveNumb, hateNumb, title, description, categoryNames
	    FROM ExtendedPostsInfosView
	    WHERE postId = ?
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Exécute la requête avec le nom d'utilisateur spécifié
	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []PostInfo
	for rows.Next() {
		var p PostInfo
		err := rows.Scan(&p.PostId, &p.UserImage, &p.UserName, &p.PostDate, &p.LoveNumb, &p.HateNumb, &p.Title, &p.Description, &p.Categories)
		if err != nil {
			return nil, err
		}
		p.CategoriesTab = strings.Split(p.Categories, ",")
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
