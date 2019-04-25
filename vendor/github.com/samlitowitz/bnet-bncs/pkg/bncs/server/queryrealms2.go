package server

// QueryRealms2ResponseRealm is the realm structure of a SID_QUERYREALMS2 response
type QueryRealms2ResponseRealm struct {
	Unknown     uint32 // usually 1
	Title       string
	Description string
}

// QueryRealms2Response is the structure of a SID_QUERYREALMS2 response
type QueryRealms2 struct {
	Unknown uint32                      // Usually 0
	Count   uint32                      `bnet:"save-QR2Count"`
	Realms  []QueryRealms2ResponseRealm `bnet:"len-QR2Count"`
}
