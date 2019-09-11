package main

import (
    "encoding/binary"
    "crypto/aes"
    "crypto/cipher"
    "fmt"
    "crypto/rand"
    rand2 "math/rand"
    "io"
    "encoding/base64"
    "encoding/hex"
    "crypto/md5"
    "time"
    "bytes"
)

func Pad(src []byte) []byte {
    padding := aes.BlockSize - len(src)%aes.BlockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(src, padtext...)
}

func Unpad(src []byte) []byte {
    length := len(src)
    unpadding := int(src[length-1])
    return src[:(length - unpadding)]
}

// CBC
func encryptCBC(key, plaintext []byte) (ciphertext []byte, err error) {
    if len(plaintext)%aes.BlockSize != 0 {
        panic("plaintext is not a multiple of the block size")
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    ciphertext = make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err)
    }
    //iv, _ := hex.DecodeString("acfa7a047800b2f221f2c4f7d626eafb")
    //copy(ciphertext[:aes.BlockSize], iv)

    fmt.Printf("CBC Key: %s\n", hex.EncodeToString(key))
    fmt.Printf("CBC IV: %s\n", hex.EncodeToString(iv))

    cbc := cipher.NewCBCEncrypter(block, iv)
    cbc.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

    return
}

func decryptCBC(key, ciphertext []byte) (plaintext []byte, err error) {
    var block cipher.Block

    if block, err = aes.NewCipher(key); err != nil {
        return
    }

    if len(ciphertext) < aes.BlockSize {
        fmt.Printf("ciphertext too short")
        return
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    cbc := cipher.NewCBCDecrypter(block, iv)
    cbc.CryptBlocks(ciphertext, ciphertext)

    plaintext = ciphertext

    return
}


func createReq(req[]byte) (message string) {
// The key length can be 32, 24, 16  bytes (OR in bits: 128, 192 or 256)
    var cipher[]byte
    key := []byte("11111111111111111111111111111111")
    var err error

    ts := make([]byte, 8)
    fmt.Printf("timestamp %d\n", time.Now().Unix());
    binary.LittleEndian.PutUint64(ts, uint64(time.Now().Unix()))

    fmt.Printf("before md5 %s %d\n", string(req), len(req))
    sum := md5.Sum(req)
    checksum  := sum[:]
    fmt.Printf("after md5 %s\n", base64.StdEncoding.EncodeToString(checksum))
    plaintext := Pad(append(append(checksum, ts...), req...))

    fmt.Printf(string(plaintext) + "\n")

    if cipher, err = encryptCBC(key, plaintext); err != nil {
        panic(err)
    }
    //message = base64.StdEncoding.EncodeToString(cipher[aes.BlockSize:])
    message = base64.StdEncoding.EncodeToString(cipher)
    fmt.Printf("CBC: %s\n", message)
    return
}

func createReq2(req []byte) (message string) {
    //key := []byte("11111111111111111111111111111111")

    plaintext := make([]byte, 8 + 16 + 8 + len(req))
    random := rand2.Uint64()
    binary.LittleEndian.PutUint64(plaintext, random)
    binary.LittleEndian.PutUint64(plaintext[8+16:], uint64(time.Now().Unix()))
    copy(plaintext[8+16+8:], req)




    return
}

func checkMessage(message string) {
    key := []byte("11111111111111111111111111111111")
    var err error
    var cipher, plaintext []byte

    if cipher, err = base64.StdEncoding.DecodeString(message); err != nil {
        panic(err)
    }

    if plaintext, err = decryptCBC(key, cipher); err != nil {
        panic(err)
    }

    text := Unpad(plaintext);

    fmt.Printf(string(text) + "\n")

    checksum := text[:16]
    //ts := text[16:16+8]
    //i := int64(binary.LittleEndian.Uint64(ts))
    req := text[16+8:]
    sum := md5.Sum(req)

    if bytes.Compare(checksum, sum[:]) != 0 {
        panic("not equal")
    }
    return
}

func main() {

    message := createReq([]byte("longer MEANS more POSSIBLE keys"))

    checkMessage(message)

    fmt.Printf("Clear from CBC: %s\n", message)

    createReq2([]byte("ddddddd"))
}

