package events

import "fmt"
import "testing"

func TestNewEventEmitter(t *testing.T){
	emitter := NewEventEmitter()
	if emitter == nil {
		t.Fatal("new emitter but it is nil")
	}
}

func test() {
	emitter := NewEventEmitter()
	fn1 := func(args ...interface{}) {
		fmt.Println("1", args)
	}
	emitter.AddListener("hello", fn1)
	emitter.AddListener("hello", fn1)
	emitter.On("hello", func(args ...interface{}) {
		fmt.Println("2", args)
	})
	emitter.PrependListener("hello", func(args ...interface{}) {
		fmt.Println("3", args)
	})
	fn4 := func(args ...interface{}) {
		fmt.Println("4", args)
	}
	emitter.PrependListener("hello", fn4)
	emitter.Emit("hello", "Jax", ",hi Lily")
	fmt.Println("--------------------")
	//emitter.RemoveAllListener("hello")
	emitter.RemoveListener("hello", fn1)
	emitter.Emit("hello", "Jax", ",hi Lily")
	fmt.Println("--------------------")
	emitter.RemoveListener("hello", fn4)
	emitter.Emit("hello", "Jax", ",hi Lily")

	emitter.Once("eat", func(args ...interface{}) {
		fmt.Println("eat", args)
	})
	emitter.Emit("eat", "Lily")
	emitter.Emit("eat", "Lily")
	fmt.Println(emitter.Listeners("hello"))
	fmt.Println(emitter.ListenerCount("hello"))
	fmt.Println(emitter.EventNames())
	emitter.MaxListeners = 5
	for i := 0; i < 11; i++ {
		ok := emitter.AddListener("greet", func(args ...interface{}) {
			fmt.Println("greet")
		})
		fmt.Println(i, ok)
	}
	emitter.Emit("greet")
}
