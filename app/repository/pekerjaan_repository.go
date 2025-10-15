package repository

import (
	"database/sql"
	"fmt"
	"prak4/app/model"
	"time"
)

type PekerjaanRepository struct {
	DB *sql.DB
}

func (r *PekerjaanRepository) GetAllPekerjaan(search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	query := fmt.Sprintf(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, deleted_at 
		FROM pekerjaan_alumni 
		WHERE (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1) AND deleted_at IS NULL
		ORDER BY %s %s 
		LIMIT $2 OFFSET $3`, sortBy, order)

	rows, err := r.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func (r *PekerjaanRepository) GetAllPekerjaanForUser(userID int, search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	query := fmt.Sprintf(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, deleted_at 
		FROM pekerjaan_alumni 
		WHERE alumni_id = $1 AND (nama_perusahaan ILIKE $2 OR posisi_jabatan ILIKE $2 OR bidang_industri ILIKE $2) AND deleted_at IS NULL
		ORDER BY %s %s 
		LIMIT $3 OFFSET $4`, sortBy, order)

	rows, err := r.DB.Query(query, userID, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func (r *PekerjaanRepository) CountPekerjaanForUser(userID int, search string) (int, error) {
	var total int
	query := "SELECT COUNT(*) FROM pekerjaan_alumni WHERE alumni_id = $1 AND (nama_perusahaan ILIKE $2 OR posisi_jabatan ILIKE $2 OR bidang_industri ILIKE $2) AND deleted_at IS NULL"
	err := r.DB.QueryRow(query, userID, "%"+search+"%").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *PekerjaanRepository) GetPekerjaanByID(id int) (*model.PekerjaanAlumni, error) {
	var p model.PekerjaanAlumni
	query := "SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, deleted_at FROM pekerjaan_alumni WHERE id = $1 AND deleted_at IS NULL"
	
	err := r.DB.QueryRow(query, id).Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PekerjaanRepository) GetPekerjaanByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error) {
    query := `
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
        gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, 
        created_at, updated_at, deleted_at 
        FROM pekerjaan_alumni 
        WHERE alumni_id = $1 AND deleted_at IS NULL 
        ORDER BY tanggal_mulai_kerja DESC
    `

    rows, err := r.DB.Query(query, alumniID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var pekerjaanList []model.PekerjaanAlumni
    for rows.Next() {
        var p model.PekerjaanAlumni
        if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
            return nil, err
        }
        pekerjaanList = append(pekerjaanList, p)
    }
    return pekerjaanList, nil
}

func (r *PekerjaanRepository) CreatePekerjaan(p *model.PekerjaanAlumni) (int, error) {
	var id int
	err := r.DB.QueryRow(
		`INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
		p.AlumniID, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange, p.TanggalMulaiKerja, p.TanggalSelesaiKerja, p.StatusPekerjaan, p.DeskripsiPekerjaan, time.Now(), time.Now(),
	).Scan(&id)
	return id, err
}

func (r *PekerjaanRepository) UpdatePekerjaan(id int, p *model.PekerjaanAlumni) (int64, error) {
	result, err := r.DB.Exec(
		`UPDATE pekerjaan_alumni SET nama_perusahaan = $1, posisi_jabatan = $2, bidang_industri = $3, lokasi_kerja = $4, gaji_range = $5, tanggal_mulai_kerja = $6, tanggal_selesai_kerja = $7, status_pekerjaan = $8, deskripsi_pekerjaan = $9, updated_at = $10 
		 WHERE id = $11`,
		p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange, p.TanggalMulaiKerja, p.TanggalSelesaiKerja, p.StatusPekerjaan, p.DeskripsiPekerjaan, time.Now(), id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *PekerjaanRepository) SoftDeletePekerjaanByID(id int) (int64, error) {
	query := "UPDATE pekerjaan_alumni SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL"
	result, err := r.DB.Exec(query, time.Now(), id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *PekerjaanRepository) SoftDeletePekerjaanForUser(pekerjaanID int, userID int) (int64, error) {
	query := "UPDATE pekerjaan_alumni SET deleted_at = $1 WHERE id = $2 AND alumni_id = $3 AND deleted_at IS NULL"
	result, err := r.DB.Exec(query, time.Now(), pekerjaanID, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *PekerjaanRepository) DeletePekerjaan(id int) (int64, error) {
	query := "UPDATE pekerjaan_alumni SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL"
	result, err := r.DB.Exec(query, time.Now(), id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *PekerjaanRepository) CountPekerjaan(search string) (int, error) {
	var total int
	query := "SELECT COUNT(*) FROM pekerjaan_alumni WHERE (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1) AND deleted_at IS NULL"
	err := r.DB.QueryRow(query, "%"+search+"%").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}


//INI UTS
func (r *PekerjaanRepository) GetTrashedPekerjaan() ([]model.PekerjaanTrash, error) {
	query := "SELECT id, nama_perusahaan, deleted_at FROM pekerjaan_alumni WHERE deleted_at IS NOT NULL ORDER BY deleted_at DESC"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trashList []model.PekerjaanTrash 
	for rows.Next() {
		var p model.PekerjaanTrash 
		if err := rows.Scan(&p.ID, &p.NamaPerusahaan, &p.DeletedAt); err != nil {
			return nil, err
		}
		trashList = append(trashList, p)
	}
	return trashList, nil
}

func (r *PekerjaanRepository) GetTrashedPekerjaanForUser(userID int) ([]model.PekerjaanTrash, error) {
	query := `
		SELECT id, nama_perusahaan, deleted_at 
		FROM pekerjaan_alumni 
		WHERE deleted_at IS NOT NULL AND alumni_id = $1 
		ORDER BY deleted_at DESC`
		
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trashList []model.PekerjaanTrash 
	for rows.Next() {
		var p model.PekerjaanTrash 
		if err := rows.Scan(&p.ID, &p.NamaPerusahaan, &p.DeletedAt); err != nil {
			return nil, err
		}
		trashList = append(trashList, p)
	}
	return trashList, nil
}

func (r *PekerjaanRepository) RestorePekerjaanByID(id int) (int64, error) {
	query := "UPDATE pekerjaan_alumni SET deleted_at = NULL WHERE id = $1 AND deleted_at IS NOT NULL"
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *PekerjaanRepository) HardDeletePekerjaanByID(id int) (int64, error) {
	query := "DELETE FROM pekerjaan_alumni WHERE id = $1 AND deleted_at IS NOT NULL"
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *PekerjaanRepository) GetPekerjaanByIDIncludeTrashed(id int) (*model.PekerjaanAlumni, error) {
	var p model.PekerjaanAlumni
	query := `
		SELECT 
			id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
			lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
			status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, deleted_at 
		FROM pekerjaan_alumni 
		WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
		&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
		&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	
	if err != nil {
		return nil, err
	}
	return &p, nil
}