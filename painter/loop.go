package painter

import (
	"image"
	"time"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циелі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver
	next     screen.Texture // текстура, яка зараз формується
	prev     screen.Texture // текстура, яка була відправленя останнього разу у Receiver
	Mq       MessageQueue
	stopChan chan struct{}
}

var size = image.Pt(800, 800)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.Mq = MessageQueue{}
	l.stopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-l.stopChan:
				l.stopChan <- struct{}{} // Підтверджуємо повну зупинку циклу
				return
			default:
				if op := l.Mq.Pull(); op != nil {
					update := op.Do(l.next)
					if update {
						l.Receiver.Update(l.next)
						l.next, l.prev = l.prev, l.next
					}
				}
			}
		}
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	if op != nil {
		l.Mq.Push(op)
	}
}

// StopAndWait сигналізує
func (l *Loop) StopAndWait() {
	l.stopChan <- struct{}{}
	<-l.stopChan
}

// TODO: реалізувати власну чергу повідомлень.
type MessageQueue struct {
	Queue []Operation
}

func (Mq *MessageQueue) Push(op Operation) {
	Mq.Queue = append(Mq.Queue, op)
}

func (Mq *MessageQueue) Pull() Operation {
	if len(Mq.Queue) == 0 {
		for len(Mq.Queue) == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	op := Mq.Queue[0]
	Mq.Queue = Mq.Queue[1:]
	return op
}
