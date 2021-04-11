package gofsm

import (
    "errors"
    "strings"
    "testing"
)

type Order struct {
    Id      int
    Address string

    StatedValue
}

func TestSimpleTransition(t *testing.T) {
    order := &Order{}

    record, err := getFSM().Trigger("checkout", order)

    if err != nil {
        t.Errorf("should not raise any error when trigger event checkout")
    }

    if order.GetState() != "checkout" {
        t.Errorf("state doesn't changed to checkout")
    }

    if record.From != "draft" {
        t.Errorf("state from not set")
    }

    if record.To != "checkout" {
        t.Errorf("state to not set")
    }
}

func TestMultipleTransitions(t *testing.T) {
    order := &Order{}

    record, err := getFSM().Trigger("checkout", order, "checkout note")
    if err != nil {
        t.Errorf("should not raise any error when trigger event checkout")
    }

    record, err = getFSM().Trigger("pay", order, "pay note")
    if err != nil {
        t.Errorf("should not raise any error when trigger event checkout")
    }

    if order.GetState() != "paid" {
        t.Errorf("state doesn't changed to paid")
    }

    if record.To != "paid" {
        t.Errorf("state to not set")
    } else {
        if record.From != "checkout" {
            t.Errorf("state from not set")
        }

        if record.Note != "pay note" {
            t.Errorf("state note not set")
        }
    }
}

func TestMultipleTransitionsWithSameEventOnDifferentStatedValues(t *testing.T) {
    orderFSM := getFSM()
    cancelEvent := orderFSM.Event("cancel")
    cancelEvent.To("cancelled").From("draft", "checkout")
    cancelEvent.To("paid_cancelled").From("paid", "processed")

    unpaidOrder1 := &Order{}
    if _, err := orderFSM.Trigger("cancel", unpaidOrder1); err != nil {
        t.Errorf("should not raise any error when trigger event cancel")
    }

    if unpaidOrder1.State != "cancelled" {
        t.Errorf("order status doesn't transitioned correctly")
    }

    unpaidOrder2 := &Order{}
    unpaidOrder2.State = "checkout"
    if _, err := orderFSM.Trigger("cancel", unpaidOrder2); err != nil {
        t.Errorf("should not raise any error when trigger event cancel")
    }

    if unpaidOrder2.State != "cancelled" {
        t.Errorf("order status doesn't transitioned correctly")
    }

    paidOrder := &Order{}
    paidOrder.State = "paid"
    if _, err := orderFSM.Trigger("cancel", paidOrder); err != nil {
        t.Errorf("should not raise any error when trigger event cancel")
    }

    if paidOrder.State != "paid_cancelled" {
        t.Errorf("order status doesn't transitioned correctly")
    }
}

func TestStateCallbacks(t *testing.T) {
    orderFSM := getFSM()
    order := &Order{}

    address1 := "I'm an address should be set when enter checkout"
    address2 := "I'm an address should be set when exit checkout"
    orderFSM.State("checkout").Enter(func(order interface{}) error {
        order.(*Order).Address = address1
        return nil
    }).Exit(func(order interface{}) error {
        order.(*Order).Address = address2
        return nil
    })

    if _, err := orderFSM.Trigger("checkout", order); err != nil {
        t.Errorf("should not raise any error when trigger event checkout")
    }

    if order.Address != address1 {
        t.Errorf("enter callback not triggered")
    }

    if _, err := orderFSM.Trigger("pay", order); err != nil {
        t.Errorf("should not raise any error when trigger event pay")
    }

    if order.Address != address2 {
        t.Errorf("exit callback not triggered")
    }
}

func TestEventCallbacks(t *testing.T) {
    var (
        order = &Order{}
        orderFSM = getFSM()
        prevState, afterState string
    )

    orderFSM.Event("checkout").To("checkout").From("draft").Before(func(order interface{}) error {
        prevState = order.(*Order).State
        return nil
    }).After(func(order interface{}) error {
        afterState = order.(*Order).State
        return nil
    })

    order.State = "draft"
    if _, err := orderFSM.Trigger("checkout", order); err != nil {
        t.Errorf("should not raise any error when trigger event checkout")
    }

    if prevState != "draft" {
        t.Errorf("Before callback triggered after state change")
    }

    if afterState != "checkout" {
        t.Errorf("After callback triggered after state change")
    }
}

func TestStateOnEnterCallbackError(t *testing.T) {
    var (
        order = &Order{}
        orderFSM = getFSM()
    )

    orderFSM.State("checkout").Enter(func(order interface{}) (err error) {
        return errors.New("intentional error")
    })

    if _, err := orderFSM.Trigger("checkout", order); err == nil {
        t.Errorf("should raise an intentional error")
    }

    if order.State != "draft" {
        t.Errorf("state transitioned on Enter callback error")
    }
}

func TestStateOnExitCallbackError(t *testing.T) {
    var (
        order = &Order{}
        orderFSM = getFSM()
    )

    orderFSM.State("checkout").Exit(func(order interface{}) (err error) {
        return errors.New("intentional error")
    })

    if _, err := orderFSM.Trigger("checkout", order); err != nil {
        t.Errorf("should not raise error when checkout")
    }

    if _, err := orderFSM.Trigger("pay", order); err == nil {
        t.Errorf("should raise an intentional error")
    }

    if order.State != "checkout" {
        t.Errorf("state transitioned on Enter callback error")
    }
}

func TestEventOnBeforeCallbackError(t *testing.T) {
    var (
        order = &Order{}
        orderFSM = getFSM()
    )

    orderFSM.Event("checkout").To("checkout").From("draft").Before(func(order interface{}) error {
        return errors.New("intentional error")
    })

    if _, err := orderFSM.Trigger("checkout", order); err == nil {
        t.Errorf("should raise an intentional error")
    }

    if order.State != "draft" {
        t.Errorf("state transitioned on Enter callback error")
    }
}

func TestEventOnAfterCallbackError(t *testing.T) {
    var (
        order = &Order{}
        orderFSM = getFSM()
    )

    orderFSM.Event("checkout").To("checkout").From("draft").After(func(order interface{}) error {
        return errors.New("intentional error")
    })

    if _, err := orderFSM.Trigger("checkout", order); err == nil {
        t.Errorf("should raise an intentional error")
    }

    if order.State != "draft" {
        t.Errorf("state transitioned on Enter callback error")
    }
}

func TestWrongEvent(t *testing.T) {
    order := &Order{}

    if _, err := getFSM().Trigger("unknown_event", order); err == nil {
        t.Errorf("should raise an unknown event error")
    }
}

func TestBadFSMStructure(t *testing.T) {
    order := &Order{}
    sm := NewFSM()

    sm.Initial("state1")
    sm.State("state2")
    e := sm.Event("event1")
    e.To("state2").From("state1")
    e.To("state1").From("state1")

    if _, err := sm.Trigger("event1", order); err == nil {
        t.Errorf("should raise invalid transtions error")
    }
}

func TestFSMGraphvizOutput(t *testing.T) {

	got := Visualize(getFSM())

    want := `
digraph fsm {
    "[draft]" -> "checkout" [ label = "checkout" ];
    "[checkout]" -> "paid" [ label = "pay" ];

    "draft (initial);"
    "cancelled";
    "checkout";
    "delivered";
    "paid";
    "paid_cancelled";
    "processed";
}`

	normGot := strings.ReplaceAll(got, "\n", "")
	normWant := strings.ReplaceAll(want, "\n", "")
	if normGot != normWant {
		t.Errorf("build graphivz graph failed. \nwant \n%s\nand got \n%s\n", want, got)
	}
}

func getFSM() *FSM {
    sm := NewFSM()

    sm.Initial("draft")

    sm.State("checkout")
    sm.State("paid")
    sm.State("processed")
    sm.State("delivered")
    sm.State("cancelled")
    sm.State("paid_cancelled")

    sm.Event("checkout").To("checkout").From("draft")
    sm.Event("pay").To("paid").From("checkout")

    return sm
}

