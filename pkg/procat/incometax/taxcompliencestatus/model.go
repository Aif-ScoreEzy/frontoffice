package taxcompliancestatus

type taxComplianceStatusRequest struct {
	Npwp string `json:"npwp" validate:"required~NPWP tidak boleh kosong., numeric~NPWP hanya berupa angka., length(16)~NPWP 15 digit tidak berlaku. Untuk pribadi gunakan NIK. Bila badan atau perusahaan tambahkan angka 0 di depan."`
}

type taxComplianceDataResponse struct {
	Nama   string `json:"nama"`
	Alamat string `json:"alamat"`
	Status string `json:"status"`
}
