package mdx

import "encoding/json"

// MdictAccessor provides a simplified interface for accessing Mdict data,
// suitable for serialization and remote access.
type MdictAccessor struct {
	Filepath          string `json:"filepath"`
	IsRecordEncrypted bool   `json:"is_record_encrypted"`
	IsMDD             bool   `json:"is_mdd"`
	IsUTF16           bool   `json:"is_utf_16"`
}

// NewAccessor creates a new MdictAccessor from an Mdict instance.
func NewAccessor(mdict *Mdict) *MdictAccessor {
	return &MdictAccessor{
		Filepath:          mdict.filePath,
		IsRecordEncrypted: mdict.meta.encryptType == EncryptRecordEnc,
		IsMDD:             mdict.fileType == MdictTypeMdd,
		IsUTF16:           mdict.meta.encoding == EncodingUtf16,
	}
}

// NewAccessorFromJSON creates a new MdictAccessor from a JSON byte slice.
func NewAccessorFromJSON(data []byte) (*MdictAccessor, error) {
	mdi := new(MdictAccessor)
	err := json.Unmarshal(data, mdi)
	return mdi, err
}

// Serialize converts the MdictAccessor to its JSON representation.
func (mdi *MdictAccessor) Serialize() ([]byte, error) {
	return json.Marshal(mdi)
}

// RetrieveDefByIndex retrieves a definition by its keyword index.
func (mdi *MdictAccessor) RetrieveDefByIndex(index *MDictKeywordIndex) ([]byte, error) {
	return locateDefByKWIndex(index, mdi.Filepath, mdi.IsRecordEncrypted, mdi.IsMDD, mdi.IsUTF16)
}
