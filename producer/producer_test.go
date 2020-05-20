// package producer

// import "testing"

// func TestProducerInit(t *testing.T) {
// 	p, _ := Init("")
// 	if p != nil {
// 		t.Errorf("Init(\"\") FAILED, expected some error got %v", p)
// 	} else {
// 		t.Logf("Init(\"\") PASSED,expected <nil> got %v", p)
// 	}

// 	p, _ = Init("192.168.99.100:19092")
// 	if p == nil {
// 		t.Errorf("Init(\"192.168.99.100:19092\") FAILED, expected some error got %v", p)
// 	} else {
// 		t.Logf("Init(\"192.168.99.100:19092\") PASSED,expected <nil> got %v", p)
// 	}
// }
