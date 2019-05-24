package en_decrypt

/*
提供简单的明文字符串形式加密，使得保存的token.json不被窃取
Simple string encrypt and decrypt for safety of your token stored in json file
*/

import (
	"crypto/aes"
	"crypto/cipher"
)

//16 bytes key
const KEY = "1234567890abcdef"

//16 bytes cipher text
var iv = []byte("abcdef1234567890")
//encrypt block
var block cipher.Block

func init() {
	var err error
	//initialize block
	block, err = aes.NewCipher([]byte(KEY))
	if err != nil {
		panic(err)
	}
}

//加密字符串，返回加密结果
//Encrypt -> string
func EncryptText(what string) (string, error) {

	str := []byte(what)

	encryptIt := cipher.NewCFBEncrypter(block, iv)
	encrypted := make([]byte, len(str))
	encryptIt.XORKeyStream(encrypted, str)

	return string(encrypted), nil
}

//解密字符串，返回解密结果
//Decrypt -> string
func DecryptText(what string) string {
	decryptIt := cipher.NewCFBDecrypter(block, iv)

	decrypted := make([]byte, len(what))
	decryptIt.XORKeyStream(decrypted, []byte(what))

	return string(decrypted)
}
