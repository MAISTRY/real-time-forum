package DB

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	enforcementOfFKs = `PRAGMA FOREIGN_KEYS = 1;`

	CreateUserTableQuery = `CREATE TABLE IF NOT EXISTS User(
		UserID INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		firstname TEXT NOT NULL,
		lastname TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE CHECK(email LIKE '%@%.%'),
        password TEXT NOT NULL,
		gender TEXT NOT NULL CHECK(gender IN ('M', 'F')),
		age INTEGER NOT NULL CHECK(age >= 18 AND age <= 100),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        privilege INTEGER NOT NULL CHECK(privilege >= 1 AND privilege <= 3) DEFAULT 1
	);`
	CreatePostTableQuery = `CREATE TABLE IF NOT EXISTS Post(
        PostID INTEGER PRIMARY KEY AUTOINCREMENT,
        UserID INTEGER NOT NULL,
		PostDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		title TEXT NOT NULL,
        content TEXT NOT NULL,
		ImagePath TEXT,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
	);`
	CreateCategoryTableQuery = `CREATE TABLE IF NOT EXISTS Category(
		CategoryID INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
        description TEXT NOT NULL,
        UserID INTEGER NOT NULL,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
	);`
	CreatePostCategoryTableQuery = `CREATE TABLE IF NOT EXISTS PostCategory(
		PostID INTEGER NOT NULL,
        CategoryID INTEGER NOT NULL,
        PRIMARY KEY (PostID, CategoryID),
        FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE CASCADE,
        FOREIGN KEY (CategoryID) REFERENCES Category(CategoryID) ON DELETE CASCADE
	);`
	CreateCommentTableQuery = `CREATE TABLE IF NOT EXISTS Comment(
		CommentID INTEGER PRIMARY KEY AUTOINCREMENT,
        PostID INTEGER NOT NULL,
        UserID INTEGER NOT NULL,
        content TEXT NOT NULL,
		CmtDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE CASCADE,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
	);`
	CreatePostLikeTableQuery = `CREATE TABLE IF NOT EXISTS PostLike(
        PostID INTEGER NOT NULL,
        UserID INTEGER NOT NULL,
		PRIMARY KEY (PostID, UserID),
		FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE CASCADE,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
    );`
	CreatePostDislikeTableQuery = `CREATE TABLE IF NOT EXISTS PostDislike(
        PostID INTEGER NOT NULL,
        UserID INTEGER NOT NULL,
		PRIMARY KEY (PostID, UserID),
		FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE CASCADE,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
    );`
	///// TODO: review the tables and check how fesiable the code will be...
	CreateCommentLikeTableQuery = `CREATE TABLE IF NOT EXISTS CommentLike(
		CommentID INTEGER NOT NULL,
		UserID INTEGER NOT NULL,
		PRIMARY KEY (CommentID, UserID),
		FOREIGN KEY (CommentID) REFERENCES Comment(CommentID) ON DELETE CASCADE,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
	);`
	CreateCommentDislikeTableQuery = `CREATE TABLE IF NOT EXISTS CommentDislike(
		CommentID INTEGER NOT NULL,
		UserID INTEGER NOT NULL,
		PRIMARY KEY (CommentID, UserID),
		FOREIGN KEY (CommentID) REFERENCES Comment(CommentID) ON DELETE CASCADE,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
	);`
	// * Messages Table
	CreateMessagesTableQuery = `CREATE TABLE IF NOT EXISTS Messages(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    	sender_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		message TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_read BOOLEAN DEFAULT FALSE,
		FOREIGN KEY (sender_id) REFERENCES User (UserID) ON DELETE CASCADE,
		FOREIGN KEY (receiver_id) REFERENCES User (UserID) ON DELETE CASCADE
	);`
	// * added UserToNotify (to know who's the user to get the notification)
	CreateNotificationTableQuery = `CREATE TABLE IF NOT EXISTS Notification (
		NotificationID INTEGER PRIMARY KEY AUTOINCREMENT,
		UserID INTEGER NOT NULL,  -- User receiving the notification
		UserToNotify INTEGER NOT NULL,  -- User who is getting the notification (null if system notification)
		PostID INTEGER,           -- Post related to the notification (nullable if comment only)
		CommentID INTEGER,        -- Comment related to the notification (nullable if only a like)
		NotificationType TEXT NOT NULL CHECK(NotificationType IN ('PostLike', 'PostDislike', 'Comment', 'CommentLike', 'CommentDislike')),
		CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		IsRead BOOLEAN NOT NULL DEFAULT FALSE,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE,
		FOREIGN KEY (UserToNotify) REFERENCES User(UserID) ON DELETE CASCADE,
		FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE SET NULL,
		FOREIGN KEY (CommentID) REFERENCES Comment(CommentID) ON DELETE SET NULL
	);`
	sessionTableQuery = `CREATE TABLE IF NOT EXISTS Session(
		session_id TEXT PRIMARY KEY,
		user_id INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        expiry_date TIMESTAMP,
		ip_address TEXT,
		FOREIGN KEY (user_id) REFERENCES User(UserID) ON DELETE CASCADE
	);`
)

var predefinedCategories = []string{"Technology", "Education", "Entertainment", "Travel", "Cars", "Sports", "Lifestyle", "Science", "Business"}

func InsertDefaultUsers(db *sql.DB) {
	defaultUsers := []struct {
		username, firstname, lastname, email, password, gender string
		age, privilege                                         int
	}{
		{"admin", "Admin", "User", "admin@gmail.com", "$2a$10$2COY2pQOxsPFA6.LrOsoj.0b7cEOmiD2q4pmHgdUI3Wf1fTBX5L86", "M", 21, 3},       // * password: adminadmin
		{"maistry", "Mujtaba", "User", "mujtaba@gmail.com", "$2a$10$SsAxMwWXMMbfT9ziRrpTU.2datBjmkVIoQKMj7.PLkh3daKSyg0sO", "M", 20, 2}, // * password: mujtaba123
		{"meow", "Mahmood", "User", "mahmood@gmail.com", "$2a$10$XDHVr9yLMQbdZ72S0Nig/e71zh8nYy1.FnY82kP4Ng16wAppryx4m", "M", 20, 2},    // * password: mahmood123
	}

	for _, user := range defaultUsers {
		_, err := db.Exec(`INSERT INTO User (username, firstname, lastname, email, password, gender, age, privilege) 
			SELECT ?, ?, ?, ?, ?, ?, ?,?
			WHERE NOT EXISTS (SELECT 1 FROM User WHERE username = ?)`,
			user.username, user.firstname, user.lastname, user.email, user.password, user.gender, user.age, user.privilege, user.username)
		if err != nil {
			log.Printf("error inserting user %s: %v", user.username, err)
		}
	}
	log.Println("Users Inserted successfully...")
}
func InsertDefaultCategories(db *sql.DB) {

	for _, category := range predefinedCategories {
		_, err := db.Exec(`INSERT INTO Category (title, description, UserID) 
			SELECT ?, ?, ? 
			WHERE NOT EXISTS (SELECT 1 FROM Category WHERE title = ?)`,
			category, category+" description", 1, category)
		if err != nil {
			log.Printf("error inserting category %s: %v", category, err)
		}
	}
	log.Println("Categorys Inserted successfully...")
}

func InsertDefaultPosts(db *sql.DB) {
	defaultPosts := []struct {
		UserID                    int
		title, content, ImagePath string
	}{
		{1, "Welcome to Penguinity!",
			`We’re thrilled to have you here at Penguinity,
			Thank you for being part of Penguinity. 
			We can’t wait to see what you bring to the table!`,
			"../uploads/Penguinity.png",
		},
	}

	for _, post := range defaultPosts {
		_, err := db.Exec(`INSERT INTO Post (UserID, title, content, ImagePath) 
			SELECT ?, ?, ?, ?
			WHERE NOT EXISTS (SELECT 1 FROM Post WHERE title = ?)`,
			post.UserID, post.title, post.content, post.ImagePath, post.title)
		if err != nil {
			log.Printf("error inserting post %s: %v", post.title, err)
		}
	}
	for _, category := range predefinedCategories {
		var categoryID int
		err := db.QueryRow(`SELECT CategoryID FROM Category WHERE title = ?`, category).Scan(&categoryID)
		if err != nil {
			log.Printf("error fetching CategoryID for %s: %v", category, err)
			continue
		}

		_, err = db.Exec(`INSERT INTO PostCategory (PostID, CategoryID) 
			SELECT ?, ?
			WHERE NOT EXISTS (SELECT 1 FROM PostCategory WHERE PostID = ? AND CategoryID = ?)`,
			1, categoryID, 1, categoryID)
		if err != nil {
			log.Printf("error inserting PostCategory %s: %v", category, err)
		}
	}

	log.Println("Posts Inserted successfully...")
}

func CreateTables(db *sql.DB) {
	if _, err := db.Exec(enforcementOfFKs); err != nil {
		log.Fatalf("error enabling foreign keys: %v", err)
	}
	if _, err := db.Exec(CreateUserTableQuery); err != nil {
		log.Fatalf("error creating the user table: %v", err)
	}
	if _, err := db.Exec(CreatePostTableQuery); err != nil {
		log.Fatalf("error creating the post table: %v", err)
	}
	if _, err := db.Exec(CreateCategoryTableQuery); err != nil {
		log.Fatalf("error creating the category table: %v", err)
	}
	if _, err := db.Exec(CreatePostCategoryTableQuery); err != nil {
		log.Fatalf("error creating the post_category table: %v", err)
	}
	if _, err := db.Exec(CreateCommentTableQuery); err != nil {
		log.Fatalf("error creating the comment table: %v", err)
	}
	if _, err := db.Exec(CreatePostLikeTableQuery); err != nil {
		log.Fatalf("error creating the like table: %v", err)
	}
	if _, err := db.Exec(CreatePostDislikeTableQuery); err != nil {
		log.Fatalf("error creating the dislike table: %v", err)
	}
	if _, err := db.Exec(CreateCommentLikeTableQuery); err != nil {
		log.Fatalf("error creating the comment_like table: %v", err)
	}
	if _, err := db.Exec(CreateCommentDislikeTableQuery); err != nil {
		log.Fatalf("error creating the comment_dislike table: %v", err)
	}
	if _, err := db.Exec(CreateMessagesTableQuery); err != nil {
		log.Fatalf("error creating the Messages table: %v", err)
	}
	if _, err := db.Exec(CreateNotificationTableQuery); err != nil {
		log.Fatalf("error creating the Notification table: %v", err)
	}
	if _, err := db.Exec(sessionTableQuery); err != nil {
		log.Fatalf("error creating the session table: %v", err)
	}

	log.Println("Tables created successfully...")
}

func InitDB() {
	db, err := sql.Open("sqlite3", "./meow.db")
	if err != nil {
		log.Fatalf("error creating database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	CreateTables(db)
	InsertDefaultUsers(db)
	InsertDefaultCategories(db)
	InsertDefaultPosts(db)
}
