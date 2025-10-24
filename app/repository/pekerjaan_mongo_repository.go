package repository

import (
	"context"
	"errors"
	"prak4/app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PekerjaanRepositoryMongo struct {
	collection *mongo.Collection
}

func NewPekerjaanRepositoryMongo(db *mongo.Database) *PekerjaanRepositoryMongo {
	return &PekerjaanRepositoryMongo{
		collection: db.Collection("pekerjaan_alumni"),
	}
}

func (r *PekerjaanRepositoryMongo) Create(ctx context.Context, pekerjaan *model.PekerjaanAlumniMongo) (*model.PekerjaanAlumniMongo, error) {
	pekerjaan.ID = primitive.NewObjectID()
	pekerjaan.CreatedAt = time.Now()
	pekerjaan.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, pekerjaan)
	if err != nil {
		return nil, err
	}
	pekerjaan.ID = result.InsertedID.(primitive.ObjectID)
	return pekerjaan, nil
}

func (r *PekerjaanRepositoryMongo) FindAll(ctx context.Context, userID int, isAdmin bool, search, sortBy string, order int, limit, offset int) ([]model.PekerjaanAlumniMongo, int64, error) {
	filter := bson.M{"deleted_at": nil}
	if !isAdmin {
		filter["alumni_id"] = userID
	}
	if search != "" {
		filter["$or"] = []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)).SetSort(bson.D{{Key: sortBy, Value: order}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var pekerjaanList []model.PekerjaanAlumniMongo
	if err = cursor.All(ctx, &pekerjaanList); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return pekerjaanList, total, nil
}

func (r *PekerjaanRepositoryMongo) FindByID(ctx context.Context, id string, includeTrashed bool) (*model.PekerjaanAlumniMongo, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID tidak valid")
	}

	var pekerjaan model.PekerjaanAlumniMongo
	filter := bson.M{"_id": objID}
	if !includeTrashed {
		filter["deleted_at"] = nil
	}

	err = r.collection.FindOne(ctx, filter).Decode(&pekerjaan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &pekerjaan, nil
}

func (r *PekerjaanRepositoryMongo) FindAllByAlumniID(ctx context.Context, alumniID int) ([]model.PekerjaanAlumniMongo, error) {
	filter := bson.M{"alumni_id": alumniID, "deleted_at": nil}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "tanggal_mulai_kerja", Value: -1}}))
	if err != nil { return nil, err }
	defer cursor.Close(ctx)

	var pekerjaanList []model.PekerjaanAlumniMongo
	if err = cursor.All(ctx, &pekerjaanList); err != nil {
		return nil, err
	}
	return pekerjaanList, nil
}

func (r *PekerjaanRepositoryMongo) Update(ctx context.Context, id string, p *model.PekerjaanAlumniMongo) (*mongo.UpdateResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	update := bson.M{
		"$set": bson.M{
			"nama_perusahaan":      p.NamaPerusahaan,
			"posisi_jabatan":       p.PosisiJabatan,
			"bidang_industri":      p.BidangIndustri,
			"lokasi_kerja":         p.LokasiKerja,
			"gaji_range":           p.GajiRange,
			"tanggal_mulai_kerja":  p.TanggalMulaiKerja,
			"tanggal_selesai_kerja": p.TanggalSelesaiKerja,
			"status_pekerjaan":     p.StatusPekerjaan,
			"deskripsi_pekerjaan":  p.DeskripsiPekerjaan,
			"updated_at":           time.Now(),
		},
	}
	return r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
}

func (r *PekerjaanRepositoryMongo) SoftDelete(ctx context.Context, id string) (*mongo.UpdateResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	return r.collection.UpdateOne(ctx, filter, update)
}

func (r *PekerjaanRepositoryMongo) Restore(ctx context.Context, id string) (*mongo.UpdateResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID, "deleted_at": bson.M{"$ne": nil}}
	update := bson.M{"$set": bson.M{"deleted_at": nil}}
	return r.collection.UpdateOne(ctx, filter, update)
}

func (r *PekerjaanRepositoryMongo) HardDelete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID, "deleted_at": bson.M{"$ne": nil}}
	return r.collection.DeleteOne(ctx, filter)
}

func (r *PekerjaanRepositoryMongo) FindTrashed(ctx context.Context, userID int, isAdmin bool) ([]model.PekerjaanTrashMongo, error) {
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if !isAdmin {
		filter["alumni_id"] = userID
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "deleted_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trashList []model.PekerjaanTrashMongo
	if err = cursor.All(ctx, &trashList); err != nil {
		return nil, err
	}
	return trashList, nil
}