package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"

	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"

	"mime/multipart"

	"github.com/larspensjo/config"
	"github.com/tealeg/xlsx"
)

var MYSQL_DB_CONNECT = "train:train123@tcp(192.168.x.xxx:3306)/ods_train?parseTime=true"
var PLAN_CONNECT = "plan:plan123@tcp(192.168.x.xxx:3306)/ods_sales_plan?parseTime=true"
var MONGODB = "mongodb://xxx:xxxx@192.168.x.xxx:27017/"
var PORT = 8766

const (
	HANA_DRIVER = "hdb"
	HANA_DNS    = "hdb://xxx:xxx@192.168.x.xxx:30015"
)

const (
	ORA_DRIVER  = "oci8"
	ORA_URI     = "mfdev/grace1996erp@M4PRO"
	ORA_URIbank = "system/information@MUTIBANK"
	//ods_bak/odsbak@110.1.5.54:1521/ODS_DEV
)

var CurrentMode = "DEV"

func RunMode() {
	conf, _ := config.ReadDefault("conf/app.conf")
	switch CurrentMode {
	case "DEV":
		MYSQL_DB_CONNECT, _ = conf.String("DEV", "MYSQL_CONN")
		PLAN_CONNECT, _ = conf.String("DEV", "PLAN_CONN")
		MONGODB, _ = conf.String("DEV", "MONGO_CONN")
		_PORT, _ := conf.String("DEV", "PORT")
		PORT, _ = strconv.Atoi(_PORT)
	case "TEST":
		MYSQL_DB_CONNECT, _ = conf.String("TEST", "MYSQL_CONN")
		PLAN_CONNECT, _ = conf.String("TEST", "PLAN_CONN")
		MONGODB, _ = conf.String("TEST", "MONGO_CONN")
		_PORT, _ := conf.String("TEST", "PORT")
		PORT, _ = strconv.Atoi(_PORT)
	case "PRD":
		MYSQL_DB_CONNECT, _ = conf.String("PRD", "MYSQL_CONN")
		PLAN_CONNECT, _ = conf.String("PRD", "PLAN_CONN")
		MONGODB, _ = conf.String("PRD", "MONGO_CONN")
		_PORT, _ := conf.String("PRD", "PORT")
		PORT, _ = strconv.Atoi(_PORT)
	}

}

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		//中断 panic(err)
	}

}

func ToFloat64(ori []byte) (re float64) {
	var bi big.Int
	var neg bool
	var i int

	neg, i = decodeDecimal(ori, &bi)
	re = bigIntToFloat(neg, &bi, i)
	return re
}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//获取source的子串,如果start小于0或者end大于source长度则返回""
//start:开始index，从0开始，包括0
//end:结束index，以end结束，但不包括end
func SubString(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

func decodeDecimal(b []byte, m *big.Int) (bool, int) {

	//bigint word size (*--> src/pkg/math/big/arith.go)
	const (
		dec128Bias = 6176
		// Compute the size _S of a Word in bytes.
		_m    = ^big.Word(0)
		_logS = _m>>8&1 + _m>>16&1 + _m>>32&1
		_S    = 1 << _logS
	)

	neg := (b[15] & 0x80) != 0
	exp := int((((uint16(b[15])<<8)|uint16(b[14]))<<1)>>2) - dec128Bias

	b14 := b[14]  // save b[14]
	b[14] &= 0x01 // keep the mantissa bit (rest: sign and exp)

	//most significand byte
	msb := 14
	for msb > 0 {
		if b[msb] != 0 {
			break
		}
		msb--
	}

	//calc number of words
	numWords := (msb / _S) + 1
	w := make([]big.Word, numWords)

	k := numWords - 1
	d := big.Word(0)
	for i := msb; i >= 0; i-- {
		d |= big.Word(b[i])
		if k*_S == i {
			w[k] = d
			k--
			d = 0
		}
		d <<= 8
	}
	b[14] = b14 // restore b[14]
	m.SetBits(w)
	return neg, exp
}

func bigIntToFloat(sign bool, m *big.Int, exp int) float64 {
	var neg int64
	if sign {
		neg = -1
	} else {
		neg = 1
	}

	return float64(neg*m.Int64()) * math.Pow10(exp)
}

func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

func Encrypt(password, deskey string) (reuslt string, err error) {
	origData := []byte(password)
	key := []byte(deskey)
	block, err := des.NewCipher(key)
	if err != nil {
		return "error", err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	reuslt = base64.StdEncoding.EncodeToString(crypted) //一般64位编码处理
	reuslt = fmt.Sprintf("%X", crypted)                 //ODS的处理，16进制编码
	return reuslt, nil
}

func Decrypt(password, deskey string) ([]byte, error) {
	//针对ODS的反解处理，16进制字符串编码变为二进制
	crypted, _ := hex.DecodeString(password)
	// crypted := []byte(password)

	key := []byte(deskey)
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// get 网络请求
func HttpGet(apiURL string, params url.Values) (rs []byte, err error) {
	var Url *url.URL
	Url, err = url.Parse(apiURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return nil, err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlstr := Url.String()
	resp, err := http.Get(urlstr)
	fmt.Println(urlstr)
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// post 网络请求 ,params 是url.Values类型
func HttpPost(apiURL string, params url.Values) (rs []byte, err error) {
	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//网络文件下载
func DownloadFile(fileName string, url string) (err error) {
	// Create the file
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

//移除重复数据
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	sort.Strings(arr)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// 获取文件大小的接口
type Size interface {
	Size() int64
}

func XlsxFileReader(mimeFile multipart.File) (*xlsx.File, error) {

	defer mimeFile.Close()
	var size int64
	if sizeInterface, ok := mimeFile.(Size); ok {
		size = sizeInterface.Size()
	}

	xlFile, err := xlsx.OpenReaderAt(mimeFile, size)
	return xlFile, err
}

// 随机生成大写字母
func GetRandomString(l int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

const (
	LongDateFormat  = "2006-01-02 15:04:05"
	ShortDateFormat = "2006-01-02"
)

//获取日期格式
func GetLongDateString(date string, Hours int64) (dateString string, err error) {
	if len(date) <= 0 {
		return "", errors.New("时间为空")
	}
	inputDate, err := time.Parse(ShortDateFormat, date)
	if err == nil {
		h, _ := time.ParseDuration("1h")
		d := inputDate.Add(time.Duration(Hours) * h)
		return d.Format(LongDateFormat), err
	} else {
		return "", errors.New("时间格式不对")
	}
}

//获取相差时间
func GetMinuteDiffer(start_time, end_time string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff / 60
		return hour
	} else {
		return hour
	}
}
