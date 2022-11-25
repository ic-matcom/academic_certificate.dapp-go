package mapper

import (
	"dapp/schema"
	"dapp/schema/dto"
)

func MapCreateAsset2Asset(asset dto.CreateAsset) dto.Asset {
	return dto.Asset{
		DocType:               schema.DocType,
		Certification:         asset.Certification,
		GoldCertificate:       asset.GoldCertificate,
		Emitter:               asset.Emitter,
		Accredited:            asset.Accredited,
		Date:                  asset.Date,
		CreatedBy:             asset.CreatedBy,
		FacultyVolumeFolio:    asset.FacultyVolumeFolio,
		UniversityVolumeFolio: asset.UniversityVolumeFolio,
	}
}
