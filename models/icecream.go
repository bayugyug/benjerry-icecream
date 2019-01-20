package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/bayugyug/benjerry-icecream/utils"
)

type Icecream struct {
	ID                    string   `json:"productId"`
	Name                  string   `json:"name"`
	Story                 string   `json:"story"`
	Description           string   `json:"description"`
	ImageOpen             string   `json:"image_open"`
	ImageClosed           string   `json:"image_closed"`
	DietaryCertifications string   `json:"dietary_certifications"`
	AllergyInfo           string   `json:"allergy_info"`
	SourcingValues        []string `json:"sourcing_values"`
	Ingredients           []string `json:"ingredients"`
	Created               string   `json:"-"`
	Modified              string   `json:"-"`
	CreatedBy             int64    `json:"-"`
	ModifiedBy            int64    `json:"-"`
	Status                string   `json:"-"`
}

func NewIcecream() *Icecream {
	return &Icecream{}

}

func (u *Icecream) Bind(r *http.Request) error {
	//sanity check
	if u == nil {
		return errors.New("Missing required parameter")
	}
	// just a post-process after a decode..
	return nil
}

func (u *Icecream) SanityCheck(data *Icecream, which string) bool {
	switch which {
	case "ADD":
		if data.Name == "" {
			return false
		}
	case "UPDATE":
		if data.Name == "" || data.ID == "" {
			return false
		}
	case "GET", "DELETE":
		if data.ID == "" {
			return false
		}
	case "ADD-INGREDIENT", "UPDATE-INGREDIENT":
		if len(data.Ingredients) <= 0 || data.ID == "" {
			return false
		}
	case "ADD-SOURCING", "UPDATE-SOURCING":
		if len(data.SourcingValues) <= 0 || data.ID == "" {
			return false
		}
	case "DELETE-INGREDIENT", "DELETE-SOURCING":
		if data.ID == "" {
			return false
		}

	}
	return true
}

func (u *Icecream) Get(ctx context.Context, db *sql.DB, pid int64) (*Icecream, error) {
	r := `SELECT 
			ifnull(ice.id,''), 
			ifnull(ice.name,''), 
			ifnull(ice.story,''), 
			ifnull(ice.description,''), 
			ifnull(ice.image_open,''), 
			ifnull(ice.image_closed,''), 
			ifnull(ice.allergy_info,''), 
			ifnull(ice.dietary_certifications,''), 
			ifnull(ice.status,''), 
			ifnull(ice.created_by,0), 
			ifnull(ice.modified_by,0), 
			ifnull(ice.created_dt,''), 
			ifnull(ice.modified_dt,''),
			ifnull(src.name,''),
			ifnull(ing.name,'')
		FROM  icecreams ice
			  LEFT JOIN sourcing_values src ON src.icecream_id = ice.id
			  LEFT JOIN ingredients     ing ON ing.icecream_id = ice.id
		WHERE ice.id = ?`

	data := &Icecream{}
	rows, err := db.QueryContext(ctx, r, pid)
	if err != nil {
		log.Println("SQL_ERR::", err)
		return nil, err

	} else {
		defer rows.Close()
		var srcName, ingName string
		for rows.Next() {
			if err := rows.Scan(
				&data.ID,
				&data.Name,
				&data.Story,
				&data.Description,
				&data.ImageOpen,
				&data.ImageClosed,
				&data.AllergyInfo,
				&data.DietaryCertifications,
				&data.Status,
				&data.CreatedBy,
				&data.ModifiedBy,
				&data.Created,
				&data.Modified,
				&srcName,
				&ingName,
			); err != nil {
				log.Println("SQL_ERR::", err)
				return nil, err
			}
			//save
			if srcName != "" {
				data.SourcingValues = append(data.SourcingValues, srcName)
			}
			if ingName != "" {
				data.Ingredients = append(data.Ingredients, ingName)
			}

		} //iterate
		if err := rows.Err(); err != nil {
			log.Println("SQL_ERR::", err)
			return nil, err
		}
	}
	data.SourcingValues = utils.UHelper.RemoveStrDuplicates(data.SourcingValues)
	data.Ingredients = utils.UHelper.RemoveStrDuplicates(data.Ingredients)
	//sounds good ;-)
	return data, nil
}

func (u *Icecream) GetByName(ctx context.Context, db *sql.DB, pid string) (*Icecream, error) {
	r := `SELECT 
			ifnull(ice.id,''), 
			ifnull(ice.name,''), 
			ifnull(ice.story,''), 
			ifnull(ice.description,''), 
			ifnull(ice.image_open,''), 
			ifnull(ice.image_closed,''), 
			ifnull(ice.allergy_info,''), 
			ifnull(ice.dietary_certifications,''), 
			ifnull(ice.status,''), 
			ifnull(ice.created_by,0), 
			ifnull(ice.modified_by,0), 
			ifnull(ice.created_dt,''), 
			ifnull(ice.modified_dt,'')
			ifnull(src.name,''),
            ifnull(ing.name,''),
		FROM  icecreams ice
			LEFT JOIN sourcing_values src ON src.icecream_id = ice.id
			LEFT JOIN ingredients     ing ON ing.icecream_id = ice.id
		WHERE ice.name = ?`

	data := &Icecream{}
	rows, err := db.QueryContext(ctx, r, pid)
	if err != nil {
		log.Println("SQL_ERR::", err)
		return nil, err

	} else {
		defer rows.Close()
		var srcName, ingName string
		for rows.Next() {
			if err := rows.Scan(
				&data.ID,
				&data.Name,
				&data.Story,
				&data.Description,
				&data.ImageOpen,
				&data.ImageClosed,
				&data.AllergyInfo,
				&data.DietaryCertifications,
				&data.Status,
				&data.CreatedBy,
				&data.ModifiedBy,
				&data.Created,
				&data.Modified,
				&srcName,
				&ingName,
			); err != nil {
				log.Println("SQL_ERR::", err)
				return nil, err
			}
			//save
			if srcName != "" {
				data.SourcingValues = append(data.SourcingValues, srcName)
			}
			if ingName != "" {
				data.Ingredients = append(data.Ingredients, ingName)
			}
		} //iterate
		if err := rows.Err(); err != nil {
			log.Println("SQL_ERR::", err)
			return nil, err
		}
	}
	data.SourcingValues = utils.UHelper.RemoveStrDuplicates(data.SourcingValues)
	data.Ingredients = utils.UHelper.RemoveStrDuplicates(data.Ingredients)
	//sounds good ;-)
	return data, nil
}

func (u *Icecream) Exists(ctx context.Context, db *sql.DB, data *Icecream) int {
	r := `SELECT count(id)
                FROM  icecreams WHERE id = ?`

	stmt, err := db.PrepareContext(ctx, r)
	if err != nil {
		log.Println("SQL_ERR::", err)
		return -1
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRowContext(ctx, data.ID).Scan(&id)
	if err != nil {
		log.Println("SQL_ERR::", err)
		return -2
	}
	//sounds good ;-)
	return id
}

func (u *Icecream) Create(ctx context.Context, db *sql.DB, data *Icecream) int64 {
	//fmt
	r := `INSERT INTO icecreams (
                name,
                description,
                story,
                image_open,
                image_closed,
                allergy_info,
                dietary_certifications,
                status,
                created_by,
                created_dt)
              VALUES (?, ?, ?, ?, ?, ?, ?, 'active', ?, Now())
			  ON DUPLICATE KEY UPDATE
                modified_by = ?,
                modified_dt = Now()
                `
	//exec
	result, err := db.ExecContext(ctx, r,
		data.Name,
		data.Description,
		data.Story,
		data.ImageOpen,
		data.ImageClosed,
		data.AllergyInfo,
		data.DietaryCertifications,
		data.CreatedBy,
		data.CreatedBy,
	)
	if err != nil {
		log.Println("SQL_ERR::CREATE", err)
		return -1
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("SQL_ERR::CREATE", err, id)
		return -2
	}
	if id <= 0 {
		log.Println("SQL_ERR::CREATE", err, id)
		return -3
	}
	//sounds good ;-)
	return int64(id)

}

func (u *Icecream) Update(ctx context.Context, db *sql.DB, data *Icecream) (bool, error) {
	//fmt
	r := `UPDATE icecreams
                SET
			    description = ? ,
                story       = ? ,
                image_open  = ? ,
                image_closed= ? ,
                allergy_info= ? ,
                dietary_certifications = ? ,
                modified_by = ? ,
                status      = 'active',
                modified_dt = Now()
              WHERE  id     = ?`
	//exec
	result, err := db.ExecContext(ctx, r,
		data.Description,
		data.Story,
		data.ImageOpen,
		data.ImageClosed,
		data.AllergyInfo,
		data.DietaryCertifications,
		data.ModifiedBy,
		data.ID,
	)
	if err != nil {
		log.Println("SQL_ERR::UPDATE", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR::UPDATE", err)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}

func (u *Icecream) Delete(ctx context.Context, db *sql.DB, data *Icecream) (bool, error) {
	//fmt
	r := `UPDATE icecreams
                SET
                status      = 'deleted',
                modified_by = ? ,
                modified_dt = Now()
              WHERE  id     = ?`
	//exec
	result, err := db.ExecContext(ctx, r, data.ModifiedBy, data.ID)
	if err != nil {
		log.Println("SQL_ERR::DELETE", err)
		return false, errors.New("Failed to update")
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR::DELETE", err, data.ID, data.ModifiedBy)
		return false, errors.New("Failed to update")
	}
	//sounds good ;-)
	return true, nil
}

func (u *Icecream) CreateSourcingValue(ctx context.Context, db *sql.DB, name string, pid int64) int64 {
	//fmt
	r := `INSERT INTO sourcing_values (
                name,
                icecream_id,
                created_dt)
              VALUES (?, ?, Now())
              ON DUPLICATE KEY UPDATE
                 modified_dt = Now() `
	//exec
	result, err := db.ExecContext(ctx, r,
		name,
		pid,
	)
	if err != nil {
		log.Println("SQL_ERR::SOURCE", err)
		return 0
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("SQL_ERR::SOURCE", err, id, pid)
		return 0
	}
	//sounds good ;-)
	return int64(id)
}

func (u *Icecream) DeleteSourcingValue(ctx context.Context, db *sql.DB, pid int64) int64 {
	//fmt
	r := `DELETE FROM sourcing_values WHERE icecream_id = ? `
	//exec
	result, err := db.ExecContext(ctx, r,
		pid,
	)
	if err != nil {
		log.Println("SQL_ERR::SOURCE", err)
		return 0
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR::SOURCE", err, pid)
		return 0
	}
	//sounds good ;-)
	return pid
}

func (u *Icecream) CreateIngredient(ctx context.Context, db *sql.DB, name string, pid int64) int64 {
	//fmt
	r := `INSERT INTO ingredients (
                name,
                icecream_id,
                created_dt)
              VALUES (?, ?, Now())
              ON DUPLICATE KEY UPDATE
                 modified_dt = Now() `
	//exec
	result, err := db.ExecContext(ctx, r,
		name,
		pid,
	)
	if err != nil {
		log.Println("SQL_ERR::INGREDIENT", err)
		return 0
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("SQL_ERR::INGREDIENT", err, id, pid)
		return 0
	}
	//sounds good ;-)
	return int64(id)
}

func (u *Icecream) DeleteIngredient(ctx context.Context, db *sql.DB, pid int64) int64 {
	//fmt
	r := `DELETE FROM ingredients WHERE icecream_id = ? `
	//exec
	result, err := db.ExecContext(ctx, r,
		pid,
	)
	if err != nil {
		log.Println("SQL_ERR::INGREDIENT", err)
		return 0
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Println("SQL_ERR::INGREDIENT", err, pid)
		return 0
	}
	//sounds good ;-)
	return pid
}
