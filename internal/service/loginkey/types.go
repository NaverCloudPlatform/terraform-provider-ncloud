package loginkey

type LoginKey struct {
	KeyName     *string `json:"key_name,omitempty"`
	Fingerprint *string `json:"fingerprint,omitempty"`
}
