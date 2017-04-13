package postgres

import (
	"time"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/dbx"
	"github.com/WeCanHearYou/wechy/app/models"
)

// IdeaService contains read and write operations for ideas
type IdeaService struct {
	DB *dbx.Database
}

// GetAll returns all tenant ideas
func (svc *IdeaService) GetAll(tenantID int) ([]*models.Idea, error) {
	rows, err := svc.DB.Query(`SELECT i.id, i.number, i.title, i.description, i.created_on, u.id, u.name, u.email
								FROM ideas i
								INNER JOIN users u
								ON u.id = i.user_id
								WHERE i.tenant_id = $1
								ORDER BY i.created_on DESC`, tenantID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var ideas []*models.Idea
	for rows.Next() {
		idea := &models.Idea{}
		rows.Scan(&idea.ID, &idea.Number, &idea.Title, &idea.Description, &idea.CreatedOn, &idea.User.ID, &idea.User.Name, &idea.User.Email)
		ideas = append(ideas, idea)
	}

	return ideas, nil
}

// GetByID returns idea by given id
func (svc *IdeaService) GetByID(tenantID, ideaID int) (*models.Idea, error) {
	rows, err := svc.DB.Query(`SELECT i.id, i.number, i.title, i.description, i.created_on, u.id, u.name, u.email
								FROM ideas i
								INNER JOIN users u
								ON u.id = i.user_id
								WHERE i.tenant_id = $1
								AND i.id = $2
								ORDER BY i.created_on DESC`, tenantID, ideaID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		idea := &models.Idea{}
		rows.Scan(&idea.ID, &idea.Number, &idea.Title, &idea.Description, &idea.CreatedOn, &idea.User.ID, &idea.User.Name, &idea.User.Email)
		return idea, nil
	}
	return nil, app.ErrNotFound
}

// GetByNumber returns idea by tenant and number
func (svc *IdeaService) GetByNumber(tenantID, number int) (*models.Idea, error) {
	rows, err := svc.DB.Query(`SELECT i.id, i.number, i.title, i.description, i.created_on, u.id, u.name, u.email
								FROM ideas i
								INNER JOIN users u
								ON u.id = i.user_id
								WHERE i.tenant_id = $1
								AND i.number = $2
								ORDER BY i.created_on DESC`, tenantID, number)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		idea := &models.Idea{}
		rows.Scan(&idea.ID, &idea.Number, &idea.Title, &idea.Description, &idea.CreatedOn, &idea.User.ID, &idea.User.Name, &idea.User.Email)
		return idea, nil
	}
	return nil, app.ErrNotFound
}

// GetCommentsByIdeaID returns all coments from given idea
func (svc *IdeaService) GetCommentsByIdeaID(tenantID, ideaID int) ([]*models.Comment, error) {
	rows, err := svc.DB.Query(`SELECT c.id, c.content, c.created_on, u.id, u.name, u.email
								FROM comments c
								INNER JOIN ideas i
								ON i.id = c.idea_id
								INNER JOIN users u
								ON u.id = c.user_id
								WHERE i.id = $1
								AND i.tenant_id = $2
								ORDER BY c.created_on DESC`, ideaID, tenantID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var comments []*models.Comment
	for rows.Next() {
		c := &models.Comment{}
		rows.Scan(&c.ID, &c.Content, &c.CreatedOn, &c.User.ID, &c.User.Name, &c.User.Email)
		comments = append(comments, c)
	}

	return comments, nil
}

// Save a new idea in the database
func (svc *IdeaService) Save(tenantID, userID int, title, description string) (*models.Idea, error) {
	tx, err := svc.DB.Begin()
	if err != nil {
		return nil, err
	}

	idea := new(models.Idea)
	idea.Title = title
	idea.Description = description

	row := tx.QueryRow(`INSERT INTO ideas (title, number, description, tenant_id, user_id, created_on) 
						VALUES ($1, (SELECT COALESCE(MAX(number), 0) + 1 FROM ideas i WHERE i.tenant_id = $3), $2, $3, $4, $5) 
						RETURNING id`, title, description, tenantID, userID, time.Now())
	if err = row.Scan(&idea.ID); err != nil {
		tx.Rollback()
		return nil, err
	}

	return idea, tx.Commit()
}

// AddComment places a new comment on an idea
func (svc *IdeaService) AddComment(userID, ideaID int, content string) (int, error) {
	tx, err := svc.DB.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	if err = tx.QueryRow("INSERT INTO comments (idea_id, content, user_id, created_on) VALUES ($1, $2, $3, $4) RETURNING id", ideaID, content, userID, time.Now()).Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}
