package model

// MetaInfo berisi informasi pagination, sorting, dan search
type MetaInfo struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Total  int    `json:"total"`
	Pages  int    `json:"pages"`
	SortBy string `json:"sortBy"`
	Order  string `json:"order"`
	Search string `json:"search"`
}

// AlumniResponse adalah format respons akhir untuk endpoint /alumni
type AlumniResponse struct {
	Data []Alumni `json:"data"`
	Meta MetaInfo `json:"meta"`
}

// PekerjaanResponse adalah format respons akhir untuk endpoint /pekerjaan
type PekerjaanResponse struct {
	Data []PekerjaanAlumni `json:"data"`
	Meta MetaInfo          `json:"meta"`
}