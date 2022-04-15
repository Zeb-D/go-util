package crypto

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestEcbDecrypter_BlockSize(t *testing.T) {
	data := []byte("296ed36974cc162c276d1837bd8d92dd8f6a1f32e07bb419deb0765186e336171c44062c93ad0754250f307404561567920e0b387be8d0e9b39b831524ce21a90739e7082384a1a72a46db97bb9bb8b077eebc8986430b5e698892ceb54e68cc1418e3bf29fb31f20167e299a55a5afd")
	key := []byte("nq2q5pLmdy7ViEEAeAVVltgS9pCxDpx8")

	d, err := AesECBDecrypt(data, key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(d))

	m := make(map[string]interface{})

	err = json.Unmarshal(d, m)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(m)
}
