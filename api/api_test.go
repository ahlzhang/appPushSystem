/** 
 * 测试.
 *
 * User: zhangbob 
 * Date: 2018/5/29 
 * Time: 下午2:59 
 */
package api

import "testing"

func TestAddPushMessageApi(t *testing.T) {
	a := `{"channel":1,"content":"13213","img":"","pushRange":2,"title":"123","typeId":1,"uerIdList":[21],"url":""}`
	code, msg, _ := addPushMessageApi(a)
	t.Log("==>", code)
	t.Log("==>", msg)
}
