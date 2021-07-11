//填充数据
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}
//去除填充数据
func PKCS7UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}
//AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    origData = PKCS7Padding(origData, blockSize)
    blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
    crypted := make([]byte, len(origData))
    blockMode.CryptBlocks(crypted, origData)
    return []byte(base64.StdEncoding.EncodeToString(crypted)), nil
}
//AES解密
func AesDecrypt(strEncrypted, key []byte) ([]byte, error) {
    b64Encrypted,err := base64.StdEncoding.DecodeString(string(strEncrypted))
    if err != nil {
        return nil, err
    }
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
    origData := make([]byte, len(b64Encrypted))
    blockMode.CryptBlocks(origData, b64Encrypted)
    origData = PKCS7UnPadding(origData)
    return origData, nil
}
func main() {
    text := "今晚打老虎"
    AesKey := []byte("0f90023fc9ae10110f90023fc9ae1011") //秘钥长度为16的倍数
    fmt.Printf("明文: %s\n", text)
    fmt.Printf("秘钥: %s\n", AesKey)
    encrypted, err := AesEncrypt([]byte(text), AesKey)
    if err != nil {
        panic(err)
    }
    fmt.Printf("加密后: %s\n", encrypted)

    origin, err := AesDecrypt(encrypted, AesKey)
    if err != nil {
        panic(err)
    }
    fmt.Printf("解密后: %s\n", origin)
}