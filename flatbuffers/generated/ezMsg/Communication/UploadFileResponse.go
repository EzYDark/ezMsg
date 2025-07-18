// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package Communication

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type UploadFileResponse struct {
	_tab flatbuffers.Table
}

func GetRootAsUploadFileResponse(buf []byte, offset flatbuffers.UOffsetT) *UploadFileResponse {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &UploadFileResponse{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *UploadFileResponse) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *UploadFileResponse) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *UploadFileResponse) RequestNonce() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UploadFileResponse) MutateRequestNonce(n uint64) bool {
	return rcv._tab.MutateUint64Slot(4, n)
}

func (rcv *UploadFileResponse) Success() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *UploadFileResponse) MutateSuccess(n bool) bool {
	return rcv._tab.MutateBoolSlot(6, n)
}

func (rcv *UploadFileResponse) FileUid() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UploadFileResponse) MutateFileUid(n uint64) bool {
	return rcv._tab.MutateUint64Slot(8, n)
}

func (rcv *UploadFileResponse) ErrorMessage() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func UploadFileResponseStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func UploadFileResponseAddRequestNonce(builder *flatbuffers.Builder, requestNonce uint64) {
	builder.PrependUint64Slot(0, requestNonce, 0)
}
func UploadFileResponseAddSuccess(builder *flatbuffers.Builder, success bool) {
	builder.PrependBoolSlot(1, success, false)
}
func UploadFileResponseAddFileUid(builder *flatbuffers.Builder, fileUid uint64) {
	builder.PrependUint64Slot(2, fileUid, 0)
}
func UploadFileResponseAddErrorMessage(builder *flatbuffers.Builder, errorMessage flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(errorMessage), 0)
}
func UploadFileResponseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
