package repository

import (
	"database/sql"
	"fmt"
	"prak4/app/model"
	"time"
)

type AlumniRepository struct {
	DB *sql.DB
}

func (r *AlumniRepository) GetAllAlumni(search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	query := fmt.Sprintf(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at 
		FROM alumni 
		WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1
		ORDER BY %s %s 
		LIMIT $2 OFFSET $3`, sortBy, order)
	
	rows, err := r.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []model.Alumni
	for rows.Next() {
		var a model.Alumni
		if err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		alumniList = append(alumniList, a)
	}
	return alumniList, nil
}
func (r *AlumniRepository) GetAlumniByID(id int) (*model.Alumni, error) {
	var a model.Alumni
	err := r.DB.QueryRow("SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni WHERE id = $1", id).Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AlumniRepository) GetAlumniByAngkatan(angkatan int) (*model.AlumniAngkatan, error) {
    jumlahalumni := &model.AlumniAngkatan{Angkatan: angkatan}
    err := r.DB.QueryRow("SELECT COUNT(*) FROM alumni WHERE angkatan = $1", angkatan).Scan(&jumlahalumni.Jumlah)
    if err != nil {
        return nil, err
    }
    return jumlahalumni, nil
}

func (r *AlumniRepository) CreateAlumni(alumni *model.Alumni) (int, error) {
	var id int
	err := r.DB.QueryRow(
		"INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		alumni.NIM, alumni.Nama, alumni.Jurusan, alumni.Angkatan, alumni.TahunLulus, alumni.Email, alumni.NoTelepon, alumni.Alamat, time.Now(), time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AlumniRepository) UpdateAlumni(id int, alumni *model.Alumni) (int64, error) {
	result, err := r.DB.Exec(
		"UPDATE alumni SET nama = $1, jurusan = $2, angkatan = $3, tahun_lulus = $4, email = $5, no_telepon = $6, alamat = $7, updated_at = $8 WHERE id = $9",
		alumni.Nama, alumni.Jurusan, alumni.Angkatan, alumni.TahunLulus, alumni.Email, alumni.NoTelepon, alumni.Alamat, time.Now(), id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *AlumniRepository) DeleteAlumni(id int) (int64, error) {
	result, err := r.DB.Exec("DELETE FROM alumni WHERE id = $1", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *AlumniRepository) CountAlumni(search string) (int, error) {
	var total int
	query := "SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1"
	err := r.DB.QueryRow(query, "%"+search+"%").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}