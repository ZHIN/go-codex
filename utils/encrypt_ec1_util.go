package utils

import (
	"bytes"
)

func (e *Encrypt_Util) EC1Encrypt(origData, key []byte) (*bytes.Buffer, error) {

	return bytes.NewBuffer(origData), nil
	// var err error
	// timestamp := time.Now().Unix()
	// encryptByteBuff := new(bytes.Buffer)
	// bs := make([]byte, 8)
	// binary.LittleEndian.PutUint64(bs, uint64(timestamp))
	// encryptByteBuff.Write(bs)
	// buff2, err = e.DesEncrypt(origData, bs)
	// if err != nil {
	// 	return nil, err
	// }
	// sign := StringUtil.MD5Byte(buff2)
	// encryptByteBuff.Write(sign)
	// encryptByteBuff.Write(buff2)

	// return encryptByteBuff, nil

}

func (*Encrypt_Util) EC1Decrypt(crypted, key []byte) ([]byte, error) {
	// 校验内容是否有篡改
	// 解密内容
	// key4 := crypted[:4]
	// timestamp := crypted[4:12]
	// buff2 := append(append(origData[12:], bs...), key...)
	// sign := StringUtil.MD5Byte(buff2)
	return crypted, nil
}
