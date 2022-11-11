package generator

import (
	"fmt"
	"github.com/douguohai/gen-id/metadata"
	"github.com/douguohai/gen-id/utils"
	"math/rand"
	"strconv"
	"time"
)

const (
	ProvinceCityLength = len(metadata.ProvinceCity)
	CardBinsLength     = len(metadata.CardBins)
	DomainSuffixLength = len(metadata.DomainSuffix)
	AreaCodeLength     = len(metadata.AreaCode)
	CityNameLength     = len(metadata.CityName)
	MobilePrefix       = len(metadata.MobilePrefix)
	ZipCodeLength      = len(metadata.ZipCode)
	TelAreaCodeLength  = len(metadata.TelAreaCode)
)

// GeneratorPhone 生成手机号码
func GeneratorPhone() (ret string) {
	return metadata.MobilePrefix[utils.RandInt(0, MobilePrefix)] + fmt.Sprintf("%0*d", 8, utils.RandInt(0, 100000000))
}

// GeneratorIDCart 生成身份证信息
func GeneratorIDCart(isFullAge *bool) (IDCard string) {
	// AreaCode
	areaCode := metadata.AreaCode[utils.RandInt(0, AreaCodeLength)]
	// 获取身份证地址

	// 获取随机生日
	var (
		birthday string
		code     string
	)
	if isFullAge == nil {
		isFullAge = new(bool)
		*isFullAge = true
	}
	if t, _err := randBirthday(*isFullAge); _err != nil {
		return
	} else {
		birthday = t.UTC().Format("20060102")
	}

	randomCode := fmt.Sprintf("%0*d", 3, utils.RandInt(0, 999))
	// 合成身份证
	prefix := strconv.Itoa(areaCode) + birthday + randomCode
	code, _ = VerifyCode(prefix)
	return prefix + code
}

// VerifyCode 获取 VerifyCode
func VerifyCode(cardId string) (ret string, err error) {
	tmp := 0
	for i, v := range metadata.Wi {
		if t, _err := strconv.Atoi(string(cardId[i])); _err == nil {
			tmp += t * v
		} else {
			err = _err
			return
		}
	}
	return metadata.ValCodeArr[tmp%11], nil
}

// randBirthday isFullAge: true 年满18岁
func randBirthday(isFullAge bool) (ret time.Time, err error) {
	var (
		begin, end time.Time
	)
	if isFullAge {
		if begin, err = time.Parse("2006-01-02 15:04:05", time.Now().AddDate(-70, 0, 0).Format("2006-01-02 15:04:05")); err != nil {
			return
		}
		if end, err = time.Parse("2006-01-02 15:04:05", time.Now().AddDate(-18, 0, 0).Format("2006-01-02 15:04:05")); err != nil {
			return
		}
		ret = time.Unix(utils.RandInt64(begin.UTC().Unix(), end.UTC().Unix()), 0)
	} else {
		if begin, err = time.Parse("2006-01-02 15:04:05", "1970-01-01 00:00:00"); err != nil {
			return
		}
		if end, err = time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05")); err != nil {
			return
		}
		ret = time.Unix(utils.RandInt64(begin.UTC().Unix(), end.UTC().Unix()), 0)
	}

	return
}

// GeneratorName 生成姓名
func GeneratorName() (ret string) {
	var name string
	rand.Seed(time.Now().UnixNano())
	if rand.Int63()%3 == 0 {
		name = metadata.LastName[utils.RandInt(0, len(metadata.LastName))] + metadata.FirstName[utils.RandInt(
			0, len(metadata.LastName))]
	} else {
		rand.Seed(time.Now().UnixNano())
		n := rand.Int63n(int64(len(metadata.NameStr)))
		name = metadata.NameStr[n]
	}

	rand.Seed(time.Now().UnixNano())
	if rand.Int63()%5 == 0 {
		name += metadata.FirstName[utils.RandInt(0, len(metadata.LastName))]
	}
	return name
}

// GeneratorVocationalCertificate 生成律师职业证件号 17为数字
func GeneratorVocationalCertificate() (ret string) {
	rand.Seed(time.Now().UnixNano())
	//6位城市编号
	areaCode := metadata.AreaCode[utils.RandInt(0, AreaCodeLength)]
	tempAreaCode := strconv.Itoa(areaCode)
	//6位邮政编号
	zipcode := metadata.ZipCode[utils.RandInt(0, ZipCodeLength)]
	tempZipcode := utils.PaddingZeroForNumberStart(6, strconv.Itoa(zipcode))
	//4位电话区号
	telAreaCode := metadata.TelAreaCode[utils.RandInt(0, TelAreaCodeLength)]
	tempTelAreaCode := utils.PaddingZeroForNumberStart(4, strconv.Itoa(telAreaCode))
	last := strconv.Itoa(utils.RandInt(0, 9))
	return tempAreaCode + tempZipcode + tempTelAreaCode + last
}
