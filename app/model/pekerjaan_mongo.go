package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PekerjaanAlumniMongo struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AlumniID            int                `json:"alumni_id" bson:"alumni_id"`
	NamaPerusahaan      string             `json:"nama_perusahaan" bson:"nama_perusahaan"`
	PosisiJabatan       string             `json:"posisi_jabatan" bson:"posisi_jabatan"`
	BidangIndustri      string             `json:"bidang_industri" bson:"bidang_industri"`
	LokasiKerja         string             `json:"lokasi_kerja" bson:"lokasi_kerja"`
	GajiRange           *string            `json:"gaji_range,omitempty" bson:"gaji_range,omitempty"`
	TanggalMulaiKerja   time.Time          `json:"tanggal_mulai_kerja" bson:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time         `json:"tanggal_selesai_kerja,omitempty" bson:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string             `json:"status_pekerjaan" bson:"status_pekerjaan"`
	DeskripsiPekerjaan  *string            `json:"deskripsi_pekerjaan,omitempty" bson:"deskripsi_pekerjaan,omitempty"`
	CreatedAt           time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt           *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type PekerjaanTrashMongo struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	NamaPerusahaan string             `json:"nama_perusahaan" bson:"nama_perusahaan"`
	DeletedAt      *time.Time         `json:"deleted_at" bson:"deleted_at"`
}