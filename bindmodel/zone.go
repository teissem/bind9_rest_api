package bindmodel

type Zone struct {
	Name         string `json:"name" binding:"required"`
	FileLocation string `json:"file_location" binding:"required"`
}
