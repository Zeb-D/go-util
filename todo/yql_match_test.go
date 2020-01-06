package todo

import (
	"fmt"
	"github.com/caibirdme/yql"
	"testing"
	"time"
)

//Though the to be matched data is the type of map[string]interface{}, there're only 5 types supported:
//
//int
//int64
//float64
//string
//bool

func TestMatch(t *testing.T) {
	rawYQL := `name='deen' and age>=23 and (hobby in ('soccer', 'swim') or score>90))`
	result, _ := yql.Match(rawYQL, map[string]interface{}{
		"name":  "deen",
		"age":   int64(23),
		"hobby": "basketball",
		"score": int64(100),
	})
	fmt.Println(result)
	rawYQL = `score ∩ (7,1,9,5,3)`
	result, _ = yql.Match(rawYQL, map[string]interface{}{
		"score": []int64{3, 100, 200},
	})
	fmt.Println(result)
	rawYQL = `score in (7,1,9,5,3)`
	result, _ = yql.Match(rawYQL, map[string]interface{}{
		"score": []int64{3, 5, 2},
	})
	fmt.Println(result)
	rawYQL = `score.sum() > 10`
	result, _ = yql.Match(rawYQL, map[string]interface{}{
		"score": []int{1, 2, 3, 4, 5},
	})
	fmt.Println(result)
}

func TestRuleMatch(t *testing.T) {
	rawYQL := `name='deen' and age>=23 and (hobby in ('soccer', 'swim') or score>90)`
	ruler, _ := yql.Rule(rawYQL)

	result, _ := ruler.Match(map[string]interface{}{
		"name":  "deen",
		"age":   23,
		"hobby": "basketball",
		"score": int64(100),
	})
	fmt.Println(result)
	result, _ = ruler.Match(map[string]interface{}{
		"name":  "deen",
		"age":   23,
		"hobby": "basketball",
		"score": int64(90),
	})
	fmt.Println(result)
	//Output:
	//true
	//false
}

func TestMatchFunction(t *testing.T) {
	user := User{
		vipFlag:     2,
		createdTime: time.Now().UnixNano(),
		ExpireTime:  time.Now().UnixNano() + 3*24*60*60*1000*1000*1000,
		testFlag:    2,
		eternalVip:  2,
	}
	fmt.Println(isVIP(user))
}

func isVIP(user User) bool {
	rule := fmt.Sprintf("monthly_vip=true and now<%d or eternal_vip=1 or ab_test!=false", user.ExpireTime)
	ok, _ := yql.Match(rule, map[string]interface{}{
		"monthly_vip": user.IsMonthlyVIP,
		"now":         time.Now().Unix(),
		"eternal_vip": user.EternalFlag,
		"ab_test":     user.isABTestMatched,
	})
	return ok
}

type User struct {
	vipFlag     int
	createdTime int64
	ExpireTime  int64
	testFlag    byte
	eternalVip  int
}

func (u User) IsMonthlyVIP() bool {
	return u.vipFlag == 2
}

func (u User) EternalFlag() int {
	return u.eternalVip
}

func (u User) isABTestMatched() bool {
	return u.testFlag == 2
}
