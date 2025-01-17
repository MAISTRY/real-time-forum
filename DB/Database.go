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

func InsertDefaultUsers(db *sql.DB) {
	defaultUsers := []struct {
		username, firstname, lastname, email, password, gender string
		privilege                                              int
	}{
		{"admin", "Admin", "User", "admin@gmail.com", "$2a$10$2COY2pQOxsPFA6.LrOsoj.0b7cEOmiD2q4pmHgdUI3Wf1fTBX5L86", "M", 3},       // * password: adminadmin
		{"maistry", "Mujtaba", "User", "mujtaba@gmail.com", "$2a$10$SsAxMwWXMMbfT9ziRrpTU.2datBjmkVIoQKMj7.PLkh3daKSyg0sO", "M", 2}, // * password: mujtaba123
		{"meow", "Mahmood", "User", "mahmood@gmail.com", "$2a$10$XDHVr9yLMQbdZ72S0Nig/e71zh8nYy1.FnY82kP4Ng16wAppryx4m", "M", 2},    // * password: mahmood123
	}

	for _, user := range defaultUsers {
		_, err := db.Exec(`INSERT INTO User (username, firstname, lastname, email, password, gender, privilege) 
			SELECT ?, ?, ?, ?, ?, ?, ? 
			WHERE NOT EXISTS (SELECT 1 FROM User WHERE username = ?)`,
			user.username, user.firstname, user.lastname, user.email, user.password, user.gender, user.privilege, user.username)
		if err != nil {
			log.Printf("error inserting user %s: %v", user.username, err)
		}
	}
	log.Println("Users Inserted successfully...")
}
func InsertDefaultCategories(db *sql.DB) {

	var predefinedCategories = []string{"Technology", "Education", "Entertainment", "Travel", "Cars", "Sports", "Lifestyle", "Science", "Business"}
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

// InitDB initializes the database connection, creates the necessary tables, and performs any initial table filling.
// It opens a new SQLite database connection, checks the connection, and then calls the CreateTables function to create the required tables.
// If any errors occur during the database initialization or table creation, it logs a fatal error.
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
}

// CreateTables creates the necessary tables in the database.
// It enables foreign key constraints, creates the user, post, category, post_category, comment, like, dislike, and notification tables.
// If any errors occur during the table creation, it logs a fatal error.
// Finally, it logs a success message indicating that the tables were created successfully.
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
