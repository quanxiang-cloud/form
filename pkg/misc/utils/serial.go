package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	error2 "git.internal.yunify.com/qxp/misc/error2"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	"math/big"
	"regexp"
	"strings"
	"text/template"
	"time"
)

const (
	date = "date"
	incr = "incr"
	step = "step"
)
const (
	formate = "{{.%s}}"
)

// ParseTemplate 解析模版
func ParseTemplate(str string) (models.SerialScheme, string) {
	temp := strings.Split(str, ".")
	templates := temp[1 : len(temp)-1]
	var res = temp[0]
	var serialScheme models.SerialScheme
	for _, template := range templates {
		switch {
		case strings.HasPrefix(template, date):
			serialScheme.Date = getDate(template)
			res += fmt.Sprintf(formate, "Date")
		case strings.HasPrefix(template, incr):
			bit, value := getIncr(template)
			serialScheme.Bit = bit
			serialScheme.Value = value
			res += fmt.Sprintf(formate, "Incr")
		case strings.HasPrefix(template, step):
			serialScheme.Step = getStep(template)
		}
	}
	res += temp[len(temp)-1]

	return serialScheme, res
}

func getIncr(value string) (string, string) {
	result := strings.Split(getVal(value), ",")
	return result[0], result[1]
}

func getStep(value string) string {
	return getVal(value)
}

func getDate(value string) string {
	return getVal(value)
}

func getVal(val string) string {
	rex := regexp.MustCompile(`\{(.*?)\}`)
	out := rex.FindAllStringSubmatch(val, -1)
	return out[0][1]
}

// CheckSerial CheckSerial
func CheckSerial(serial *models.SerialScheme, oldSerialStr string) error {
	var oldSerial models.SerialScheme
	err := json.Unmarshal([]byte(oldSerialStr), &oldSerial)
	if err != nil {
		return err
	}

	oldBit, ok := big.NewInt(0).SetString(oldSerial.Bit, 10)
	if !ok {
		return error2.NewError(code.ErrParameter)
	}

	newBit, ok := big.NewInt(0).SetString(serial.Bit, 10)
	if !ok {
		return error2.NewError(code.ErrParameter)
	}

	_, ok = big.NewInt(0).SetString(serial.Value, 10)
	if !ok {
		return error2.NewError(code.ErrParameter)
	}

	_, ok = big.NewInt(0).SetString(serial.Step, 10)
	if !ok {
		return error2.NewError(code.ErrParameter)
	}

	if newBit.Cmp(oldBit) <= 0 {
		serial.Bit = oldSerial.Bit
		serial.Value = oldSerial.Value
	}
	return nil
}

// ExecuteTemplate ExecuteTemplate
func ExecuteTemplate(serialMap map[string]string) (*models.SerialScheme, string, error) {
	var serialScheme *models.SerialScheme
	ser := serialMap[redis.Serials]
	if err := json.Unmarshal([]byte(ser), &serialScheme); err != nil {
		return nil, "", err
	}

	serTemplate := serialMap[redis.Template]
	t, err := template.New("").Parse(serTemplate)
	if err != nil {
		return nil, "", err
	}

	var serialExcute models.SerialExcute
	serialExcute.Date = executeDate(serialScheme)
	incr, err := executeInCr(serialScheme)
	if err != nil {
		return nil, "", err
	}
	serialExcute.Incr = incr

	var buf bytes.Buffer
	if err = t.Execute(&buf, serialExcute); err != nil {
		return nil, "", err
	}
	return serialScheme, buf.String(), nil
}

var formats = map[string]string{
	"yyyy":           "2006",
	"yyyyMM":         "200601",
	"yyyyMMdd":       "20060102",
	"yyyyMMddHH":     "2006010215",
	"yyyyMMddHHmm":   "200601021504",
	"yyyyMMddHHmmss": "20060102150405",
}

func executeDate(serialScheme *models.SerialScheme) string {
	format := formats[serialScheme.Date]
	return time.Now().Local().UTC().Format(format)
}

func executeInCr(serialScheme *models.SerialScheme) (string, error) {
	value, _ := big.NewInt(0).SetString(serialScheme.Value, 10)
	bit, _ := big.NewInt(0).SetString(serialScheme.Bit, 10)
	if value.Cmp(maxSize(serialScheme.Bit)) >= 0 {
		serialScheme.Bit = bit.Add(bit, big.NewInt(1)).String()
	}

	format := "%0" + serialScheme.Bit + "s"
	res := fmt.Sprintf(format, serialScheme.Value)
	step, _ := big.NewInt(0).SetString(serialScheme.Step, 10)
	serialScheme.Value = value.Add(value, step).String()

	return res, nil
}

func maxSize(bitStr string) *big.Int {
	var res = big.NewInt(0)
	i, _ := big.NewInt(0).SetString(bitStr, 10)
	bit := i.Int64()
	for bit > 0 {
		bit--
		powNum := big.NewInt(10)
		powNum.Exp(powNum, big.NewInt(bit), nil)

		var mulNum = big.NewInt(9)
		mulNum.Mul(mulNum, powNum)
		res.Add(res, mulNum)
	}
	return res
}
