# gen-id
一个身份证、名字、邮箱、地址、手机号码等随机生成的sdk

# Installation
`go get github.com/douguohai/gen-id`

如果网速过慢:
```
export GO111MODULE=on
export GOPROXY=https://goproxy.io
go get github.com/douguohai/gen-id@master
```

# Usage
```shell
    go build main.go
    ./main -count=10
```

# Example

```golang
package main

import (
	"fmt"
	"github.com/douguohai/gen-id"
	"github.com/douguohai/gen-id/generator"
)

func main()  {
	// 生成总的信息
	fmt.Println(generator.NewGeneratorData(nil))
	// 分个单独获取
	g:=new(generator.GeneratorData)
	fmt.Println(g.GeneratorPhone())
	fmt.Println(g.GeneratorName())
	fmt.Println(g.GeneratorIDCart(nil))
	fmt.Println(g.GeneratorEmail())
	fmt.Println(g.GeneratorBankID())
	fmt.Println(g.GeneratorAddress())
}

```

```go
package generator

import (
	"fmt"
	"github.com/douguohai/gen-id/metadata"
	"github.com/douguohai/gen-id/utils"
	"math"
	"math/rand"
	"strconv"
	"strings"
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

// LawyerData 律师数据
type LawyerData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Address  string `json:"address"` // 目前地址
	BankID   string `json:"bank_id"`
	PhoneNum string `json:"phone_num"` // 手机号码
	// TODO 身份证
	IDCard      string `json:"id_card"`      // 身份证号
	IssueOrg    string `json:"issue_org"`    // 身份证发证机关
	Birthday    string `json:"birthday"`     // 出生日期
	ValidPeriod string `json:"valid_period"` // 有效时期
	IDCardAddr  string `json:"id_card_addr"` // 身份证地址
	// other
	PreCardNo             string `json:"-"`
	VocationalCertificate string `json:"vocational_certificate"` //职业证号
}

func NewGeneratorData(isFullAge *bool) (ret *LawyerData) {
	var (
		data = new(LawyerData)
	)
	data.GeneratorBankID()
	data.GeneratorAddress()
	data.GeneratorEmail()
	data.GeneratorIDCart(isFullAge)
	data.GeneratorName()
	data.GeneratorPhone()
	data.GeneratorVocationalCertificate()
	ret = data
	return
}

// GeneratorProvinceAdnCityRand 随机获取城市和地址
func (g *LawyerData) GeneratorProvinceAdnCityRand() (ret string) {
	return metadata.ProvinceCity[utils.RandInt(0, ProvinceCityLength)]
}

// GeneratorAddress 获取地址
func (g *LawyerData) GeneratorAddress() (ret string) {
	g.Address = g.GeneratorProvinceAdnCityRand() +
		utils.GenRandomLengthChineseChars(2, 3) + "路" +
		strconv.Itoa(utils.RandInt(1, 8000)) + "号" +
		utils.GenRandomLengthChineseChars(2, 3) + "小区" +
		strconv.Itoa(utils.RandInt(1, 20)) + "单元" +
		strconv.Itoa(utils.RandInt(101, 2500)) + "室"
	return g.Address
}

// GeneratorBankID 获取银行卡号
func (g *LawyerData) GeneratorBankID() (ret string) {
	var (
		// 随机获取卡头
		bank = metadata.CardBins[utils.RandInt(0, CardBinsLength)]
	)
	// 生成 长度 bank.length-1 位卡号
	g.PreCardNo = strconv.Itoa(bank.Prefixes[utils.RandInt(0, len(bank.Prefixes))]) + fmt.Sprintf(
		"%0*d", bank.Length-7, utils.RandInt64(0, int64(math.Pow10(bank.Length-7))))
	g.processLUHN()

	return g.BankID
}

// processLUHN 合成卡号
func (g *LawyerData) processLUHN() {
	checkSum := 0
	tmpCardNo := utils.ReverseString(g.PreCardNo)
	for i, v := range tmpCardNo {
		// 数据层确保卡号正确
		tmp, _ := strconv.Atoi(string(v))
		// 由于卡号实际少了一位，所以反转后卡号第一位一定为偶数位
		// 同时 i 正好也是偶数，此时 i 将和卡号奇偶位同步
		if i%2 == 0 {
			// 偶数位 *2 是否为两位数(>9)
			if tmp*2 > 9 {
				// 如果为两位数则 -9
				checkSum += tmp*2 - 9
			} else {
				// 否则直接相加即可
				checkSum += tmp * 2
			}
		} else {
			// 奇数位直接相加
			checkSum += tmp
		}
	}
	if checkSum%10 != 0 {
		g.BankID = g.PreCardNo + strconv.Itoa(10-checkSum%10)
	} else {
		// 如果不巧生成的前 卡长度-1 位正好符合 LUHN 算法
		// 那么需要递归重新生成(需要符合 cardBind 中卡号长度)
		g.GeneratorBankID()
	}
}

// GeneratorEmail 生成邮箱
func (g *LawyerData) GeneratorEmail() (ret string) {
	g.Email = utils.RandStr(8) + "@" + utils.RandStr(5) + metadata.DomainSuffix[utils.RandInt(0, DomainSuffixLength)]
	return g.Email
}

// GeneratorIDCart 生成身份证信息
func (g *LawyerData) GeneratorIDCart(isFullAge *bool) (ret *LawyerData, err error) {
	// AreaCode
	areaCode := metadata.AreaCode[utils.RandInt(0, AreaCodeLength)]
	// 获取身份证地址
	addr := metadata.IDPrefix[areaCode] + utils.GenRandomLengthChineseChars(2, 3) + "路" +
		strconv.Itoa(utils.RandInt(1, 8000)) + "号" +
		utils.GenRandomLengthChineseChars(2, 3) + "小区" +
		strconv.Itoa(utils.RandInt(1, 20)) + "单元" +
		strconv.Itoa(utils.RandInt(101, 2500)) + "室"
	g.IDCardAddr = addr
	g.IssueOrg = metadata.IDPrefix[areaCode] + "公安局某某分局"
	// 获取随机生日
	var (
		birthday   string
		code       string
		begin, end time.Time
	)
	if isFullAge == nil {
		isFullAge = new(bool)
		*isFullAge = true
	}
	if t, _err := g.randBirthday(*isFullAge); _err != nil {
		err = _err
		return
	} else {
		g.Birthday = t.UTC().Format("2006-01-02")
		birthday = t.UTC().Format("20060102")
	}

	randomCode := fmt.Sprintf("%0*d", 3, utils.RandInt(0, 999))
	// 合成身份证
	prefix := strconv.Itoa(areaCode) + birthday + randomCode
	if code, err = g.VerifyCode(prefix); err != nil {
		return
	}
	g.IDCard = prefix + code

	// 获取随机有效时间
	if begin, err = g.randDate(); err != nil {
		return
	}
	end = begin.AddDate(10, 0, 0)
	g.ValidPeriod = begin.Format("2006.01.02") + "-" + end.Format("2006.01.02")

	//
	ret = g
	return
}

// randBirthday isFullAge: true 年满18岁
func (g *LawyerData) randBirthday(isFullAge bool) (ret time.Time, err error) {
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

// VerifyCode 获取 VerifyCode
func (g *LawyerData) VerifyCode(cardId string) (ret string, err error) {
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

// randDate 身份证有效期随机时间 有效期最低限制now-10年
func (g *LawyerData) randDate() (ret time.Time, err error) {
	var (
		begin, end time.Time
	)
	if begin, err = time.Parse("2006-01-02 15:04:05", time.Now().AddDate(-10, 0, 0).Format("2006-01-02 15:04:05")); err != nil {
		return
	}
	if end, err = time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05")); err != nil {
		return
	}
	return time.Unix(utils.RandInt64(begin.Unix(), end.Unix()), 0), err
}

// GeneratorPhone 生成手机号码
func (g *LawyerData) GeneratorPhone() (ret string) {
	g.PhoneNum = metadata.MobilePrefix[utils.RandInt(0, MobilePrefix)] + fmt.Sprintf("%0*d", 8, utils.RandInt(0, 100000000))
	return g.PhoneNum
}

// GeneratorName 生成姓名
func (g *LawyerData) GeneratorName() (ret string) {
	rand.Seed(time.Now().UnixNano())
	if rand.Int63()%3 == 0 {
		g.Name = metadata.LastName[utils.RandInt(0, len(metadata.LastName))] + metadata.FirstName[utils.RandInt(
			0, len(metadata.LastName))]
	} else {
		arr := strings.Split(metadata.NameStr, "\n")
		rand.Seed(time.Now().UnixNano())
		n := rand.Int63n(int64(len(arr)))
		g.Name = arr[n]
	}

	rand.Seed(time.Now().UnixNano())
	if rand.Int63()%5 == 0 {
		g.Name += metadata.FirstName[utils.RandInt(0, len(metadata.LastName))]
	}
	return g.Name
}

// GeneratorVocationalCertificate 生成律师职业证件号 17为数字
func (g *LawyerData) GeneratorVocationalCertificate() (ret string) {
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
	g.VocationalCertificate = tempAreaCode + tempZipcode + tempTelAreaCode + last
	return g.VocationalCertificate
}

```

# Statement
本项目用于开发环境,涉及商业用途用本人无关
