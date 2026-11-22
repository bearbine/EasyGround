package hosts

import (
	"encoding/hex"
	"strings"
)

//======================
func NormalizeID(id string)(string,error){
	if len(id) != 36{return "",ErrInvalidID}
	//короче говоря дефисы должны разделять группы
	//  8-4-4-4-12
	if id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-'{return "",ErrInvalidID}

	//без дефисов UUID должен содержать 16 байт в hex-записи
	decoded,err := hex.DecodeString(strings.ReplaceAll(id,"-",""))
	if err != nil || len(decoded) != 16 {return "", ErrInvalidID}
	
	return strings.ToLower(id),nil
}