package taxverificationdetail

type taxVerificationRequest struct {
	NpwpOrNik string `json:"npwp_or_nik" validate:"required~NPWP or NIK cannot be empty., numeric~NPWP is only number., length(16)~NPWP 15 digit tidak berlaku. Untuk pribadi gunakan NIK. Bila badan atau perusahaan tambahkan angka 0 di depan."`
}

type taxVerificationRespData struct {
	Nama             string `json:"nama"`
	Alamat           string `json:"alamat"`
	NPWP             string `json:"npwp"`
	NPWPVerification string `json:"npwp_verification"`
	TaxCompliance    string `json:"tax_compliance"`
	Status           string `json:"status"`
}
