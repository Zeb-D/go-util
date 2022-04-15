package pulsar

type AuthProvider interface {
	AuthMethod() string
	AuthData() []byte
}

//type authProvider struct {
//	AccessID  string
//	AccessKey string
//}
//
//func NewAuthProvider(accessID, accessKey string) *authProvider {
//	return &authProvider{
//		AccessID:  accessID,
//		AccessKey: accessKey,
//	}
//}
//
//func (a *authProvider) AuthMethod() string {
//	return "auth1"
//}
//
//func (a *authProvider) AuthData() []byte {
//	key := md5Hex(a.AccessID + md5Hex(a.AccessKey))
//	key = key[8:24]
//	return []byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, a.AccessID, key))
//}
//
//func md5Hex(s string) string {
//	h := md5.New()
//	h.Write([]byte(s))
//	return hex.EncodeToString(h.Sum(nil))
//}

type NoAuth struct {
}

func (a *NoAuth) AuthMethod() string {
	return ""
}

func (a *NoAuth) AuthData() []byte {
	return nil
}
