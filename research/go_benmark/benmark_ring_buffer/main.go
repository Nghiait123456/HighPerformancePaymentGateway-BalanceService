package benmark_ring_buffer

//
//import (
//	"bytes"
//	"encoding/gob"
//	"fmt"
//	"github.com/smallnest/ringbuffer"
//	"log"
//	"sync"
//	"sync/atomic"
//	"time"
//)
//
//type Person struct {
//	Name    string
//	Address Address
//}
//
//type Address struct {
//	House    int
//	Street1  string
//	Town     string
//	PostCode PostCode
//}
//
//type PostCode struct {
//	Value string
//}
//
//func EncodeToBytes(p interface{}) []byte {
//
//	buf := bytes.Buffer{}
//	enc := gob.NewEncoder(&buf)
//	err := enc.Encode(p)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("uncompressed size (bytes): ", len(buf.Bytes()))
//	return buf.Bytes()
//}
//
//func DecodeToPerson(s []byte) (Person, error) {
//
//	p := Person{}
//	dec := gob.NewDecoder(bytes.NewReader(s))
//	err := dec.Decode(&p)
//	if err != nil {
//		return p, err
//	}
//	return p, nil
//}
//
//func BenmarkRingBuffer(r *ringbuffer.RingBuffer) {
//	p := Person{
//		Name: "Joe Bloggs",
//		Address: Address{
//			House:   1,
//			Street1: "The Lane",
//			Town:    "Blackburn",
//			PostCode: PostCode{
//				Value: "BB2 5LB",
//			},
//		},
//	}
//
//	dataOut := EncodeToBytes(p)
//	var numberErrRead uint64
//	var numberErrWrite uint64
//	var numberReadSucess uint64
//
//	oneTest := 1000
//	numberRoutine := 10
//	buf := make([]byte, 800)
//	starTime := time.Now().Unix()
//
//	var wg sync.WaitGroup
//
//	var k uint64
//	for k = 0; k < uint64(numberRoutine); k++ {
//		fmt.Println("k= ", k)
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			var i uint64
//			for i = 0; i < uint64(oneTest); {
//				fmt.Println("i= ", i)
//				_, err := r.Write(dataOut)
//				if err != nil {
//					fmt.Println(fmt.Sprintf("write fail, err = %s,", err.Error()))
//					atomic.AddUint64(&numberErrWrite, 1)
//				} else {
//					atomic.AddUint64(&numberReadSucess, 1)
//					atomic.AddUint64(&i, 1)
//				}
//			}
//		}()
//	}
//
//	for j := 0; j < numberRoutine; j++ {
//		fmt.Println("j= ", j)
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			var i uint64
//			for i = 0; i < uint64(oneTest); {
//				fmt.Println("i= ", i)
//				n, err := r.Read(buf)
//				if n <= 0 {
//					fmt.Println(fmt.Sprintf("read faild, n = %d,", n))
//					atomic.AddUint64(&numberErrRead, 1)
//				} else if err != nil {
//					fmt.Println(fmt.Sprintf("read faild, err = %s,", err.Error()))
//					atomic.AddUint64(&numberErrRead, 1)
//				} else {
//					fmt.Println("read success")
//					atomic.AddUint64(&numberReadSucess, 1)
//					atomic.AddUint64(&i, 1)
//				}
//				buf = nil
//			}
//		}()
//	}
//
//	wg.Wait()
//
//	endTime := time.Now().Unix()
//	fmt.Println(fmt.Sprintf("starttime %d, endtime %d, rangetime %d errorRead %d  errorWrite %d  readSuccess %d totalTest %d", starTime, endTime, endTime-starTime, numberErrRead, numberErrWrite, numberReadSucess, numberRoutine*oneTest))
//
//}
//
//func main() {
//	rb := ringbuffer.New(1000240)
//
//	BenmarkRingBuffer(rb)
//
//	//p := Person{
//	//	Name: "Joe Bloggs",
//	//	Address: Address{
//	//		House:   1,
//	//		Street1: "The Lane",
//	//		Town:    "Blackburn",
//	//		PostCode: PostCode{
//	//			Value: "BB2 5LB",
//	//		},
//	//	},
//	//}
//	//
//	//dataOut := EncodeToBytes(p)
//	//
//	//// write
//	//fmt.Println("start write")
//	//rb.Write(dataOut)
//	//fmt.Println(rb.Length())
//	//fmt.Println(rb.Free())
//	//
//	//// read
//	//
//	//buf := make([]byte, 400)
//	//fmt.Println("start read, len buf")
//	//fmt.Println(len(buf))
//	//n, err := rb.Read(buf)
//	//if n <= 0 || err != nil {
//	//	fmt.Println(fmt.Sprintf("have error, %s, n : %d", err.Error(), n))
//	//}
//	//
//	////fmt.Println(buf)
//	//p2, err := DecodeToPerson(buf)
//	//if err != nil {
//	//	fmt.Println(fmt.Sprintf("have error : %s", err.Error()))
//	//} else {
//	//	fmt.Println("success")
//	//	fmt.Println(p2)
//	//}
//
//}
