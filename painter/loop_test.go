package painter

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"
)

func (om *op3Mock) Do(t screen.Texture) bool {
	args := om.Called(t)
	return args.Bool(0)
}

func TestLoop_Post(t *testing.T) {
	screen := new(screenMock)
	texture := new(textureMock)
	receiver := new(receiverMock)
	tx := image.Pt(800, 800)
	l := Loop{
		Receiver: receiver,
	}

	screen.On("NewTexture", tx).Return(texture, nil)
	receiver.On("Update", texture).Return()

	op1 := new(op1Mock)
	op2 := new(op2Mock)
	op3 := new(op3Mock)

	l.Start(screen)

	op1.On("Do", texture).Return(false).Once()
	op2.On("Do", texture).Return(true).Once()
	op3.On("Do", texture).Return(true).Once()

	l.Post(op1)
	l.Post(op2)
	l.Post(op3)

	time.Sleep(1 * time.Second) // Подождите, пока операции выполнены

	op1.AssertExpectations(t)
	op2.AssertExpectations(t)
	op3.AssertExpectations(t)

	receiver.AssertCalled(t, "Update", texture)
	screen.AssertCalled(t, "NewTexture", image.Pt(800, 800))
}

func TestMessageQueue_Push(t *testing.T) {
	Mq := &MessageQueue{}

	op1 := &operationQueueMock{}
	op2 := &operationQueueMock{}
	op3 := &operationQueueMock{}

	Mq.Push(op3)
	Mq.Push(op2)
	Mq.Push(op1)

	assert.Equal(t, 3, len(Mq.Queue))
	assert.Equal(t, op3, Mq.Queue[0])
	assert.Equal(t, op2, Mq.Queue[1])
	assert.Equal(t, op1, Mq.Queue[2])
}

func TestMessageQueue_Pull(t *testing.T) {
	mq := &MessageQueue{}

	op1 := op1Mock{}
	op2 := op2Mock{}
	op3 := op3Mock{}

	mq.Push(&op3)
	mq.Push(&op2)
	mq.Push(&op1)

	op := mq.Pull()
	assert.Equal(t, &op3, op)

	op = mq.Pull()
	assert.Equal(t, &op2, op)

	op = mq.Pull()
	assert.Equal(t, &op1, op)

	assert.Empty(t, mq.Queue)
}

type op1Mock struct {
	mock.Mock
}

func (om *op1Mock) Do(t screen.Texture) bool {
	args := om.Called(t)
	return args.Bool(0)
}

type op2Mock struct {
	mock.Mock
}

func (om *op2Mock) Do(t screen.Texture) bool {
	args := om.Called(t)
	return args.Bool(0)
}

type op3Mock struct {
	mock.Mock
}

type receiverMock struct {
	mock.Mock
}

func (rm *receiverMock) Update(t screen.Texture) {
	rm.Called(t)
}

type screenMock struct {
	mock.Mock
}

func (sm *screenMock) NewBuffer(size image.Point) (screen.Buffer, error) {
	return nil, nil
}

func (sm *screenMock) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return nil, nil
}

func (sm *screenMock) NewTexture(size image.Point) (screen.Texture, error) {
	args := sm.Called(size)
	return args.Get(0).(screen.Texture), args.Error(1)
}

type textureMock struct {
	mock.Mock
}

func (tm *textureMock) Release() {
	tm.Called()
}

func (tm *textureMock) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	tm.Called(dp, src, sr)
}

func (tm *textureMock) Bounds() image.Rectangle {
	args := tm.Called()
	return args.Get(0).(image.Rectangle)
}

func (tm *textureMock) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	tm.Called(dr, src, op)
}

func (tm *textureMock) Size() image.Point {
	args := tm.Called()
	return args.Get(0).(image.Point)
}

type operationMock struct {
	mock.Mock
}

func (om *operationMock) Do(t screen.Texture) bool {
	args := om.Called(t)
	return args.Bool(0)
}

type operationQueueMock struct{}

func (m *operationQueueMock) Do(t screen.Texture) (ready bool) {
	return false
}
