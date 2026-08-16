package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

type Student struct {
	Name    string
	Age     uint8
	Address string
}

func encode(v interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(v)
	return buffer.Bytes(), err
}

func decode(b []byte, v interface{}) error {
	decoder := gob.NewDecoder(bytes.NewReader(b)) //创建解密器
	return decoder.Decode(v)
}

func TestGOB(t *testing.T) {
	//序列化
	s1 := Student{
		"张三",
		18,
		"江苏省",
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer) //创建编码器
	err1 := encoder.Encode(&s1)        //编码
	if err1 != nil {
		t.Fatal(err1)
	}

	fmt.Printf("序列化后：%x\n", buffer.Bytes())

	//反序列化
	byteEn := buffer.Bytes()
	decoder := gob.NewDecoder(bytes.NewReader(byteEn)) //创建解密器
	var s2 Student
	err2 := decoder.Decode(&s2) //解密
	if err2 != nil {
		t.Fatal(err2)
	}
	fmt.Println("反序列化后：", s2)
	if s2 != s1 {
		t.Fatalf("round trip = %#v, want %#v", s2, s1)
	}
}

type messageSnapshot struct {
	MethodID uint32
	Target   string
	Payload  Student
}

func TestMessage(t *testing.T) {
	want := messageSnapshot{
		MethodID: 1001,
		Target:   "0.0.5.1.student",
		Payload:  Student{Name: "李四", Age: 20, Address: "深圳"},
	}

	data, err := encode(&want)
	if err != nil {
		t.Fatal(err)
	}
	var got messageSnapshot
	if err = decode(data, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip = %#v, want %#v", got, want)
	}
}

type FileAll struct {
	Name string
	Cxt  []byte
}

func Test111(t *testing.T) {
	fileName := filepath.Join(t.TempDir(), "fixture.txt")
	want := FileAll{Name: fileName, Cxt: []byte("fixture")}

	data, err := encode(&want)
	if err != nil {
		t.Fatal(err)
	}
	var got FileAll
	if err = decode(data, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip = %#v, want %#v", got, want)
	}
}
