package utils

import (
	"fmt"
	captcha "github.com/mojocn/base64Captcha"
	"math/rand"
	"strconv"
	"time"
)

/*
*
验证码生成工具 图形以及数字验证码
*/
const (
	width      = 40
	height     = 130
	codeLength = 5
	noiseCount = 4
)

func newDriver(width, height, codeLength, noiseCount int) *captcha.DriverString {
	driver := new(captcha.DriverString)
	driver.Height = width
	driver.Width = height
	driver.NoiseCount = noiseCount
	driver.ShowLineOptions = captcha.OptionShowSineLine | captcha.OptionShowSlimeLine | captcha.OptionShowHollowLine
	driver.Length = codeLength
	driver.Source = "1234567890qwertyuipkjhgfdsazxcvbnm"
	driver.Fonts = []string{"wqy-microhei.ttc"}
	return driver
}

// GenerateImageCaptcha 生成图形验证码以及答案
func GenerateImageCaptcha() (captcha.Item, string) {
	var driver = newDriver(width, height, codeLength, noiseCount).ConvertFonts()
	c := captcha.NewCaptcha(driver, captcha.DefaultMemStore)
	_, content, answer := c.Driver.GenerateIdQuestionAnswer()
	item, _ := c.Driver.DrawCaptcha(content)
	return item, answer
}

func GenerateNumberCaptcha(num int) string {
	str := "1"
	for i := 0; i < num; i++ {
		str += strconv.Itoa(0)
	}
	str10 := str
	int10, err := strconv.ParseInt(str10, 10, 32)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		j := int32(int10)
		return fmt.Sprintf("%0"+strconv.Itoa(num)+"v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(j))
	}
}
