package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strings"
	"testing"

	"github.com/bailu1901/lockstepserver/proto"
	"github.com/golang/protobuf/proto"
)

func Test_SCPacket(t *testing.T) {

	var sID int32 = 19234333
	msg := &pb.C2S_InputSkillMsg{
		Sid: proto.Int32(sID),
		X:   proto.Int32(10),
		Y:   proto.Int32(20),
	}
	raw, _ := proto.Marshal(msg)
	p := NewPacket(uint8(pb.ID_C2S_InputSkill), msg)
	if nil == p {
		t.Fail()
	}

	buff := p.Serialize()

	dataLen := binary.BigEndian.Uint16(buff[0:])
	if dataLen != uint16(len(raw)) {
		t.Error("dataLen != uint16(len(raw))")
	}

	if MinPacketLen+dataLen != MinPacketLen+uint16(len(raw)) {
		t.Error("MinPacketLen+dataLen != MinPacketLen+uint16(len(raw))")
	}

	id := buff[DataLen]
	if p.id != id {
		t.Error("uint8(ID_C2S_Connect) != id")
	}

	msg1 := &pb.C2S_InputSkillMsg{}
	if err := proto.Unmarshal(buff[MinPacketLen:], msg1); nil != err {
		t.Error(err)
	}

	if msg.GetSid() != msg1.GetSid() || msg.GetX() != msg1.GetX() || msg.GetY() != msg1.GetY() {
		t.Error("msg.Sid != data1.Sid || msg.X != data1.X || msg.Y != data1.Y")
	}
}

func Benchmark_SCPacket(b *testing.B) {

	var sID int32 = 19234333
	msg := &pb.C2S_InputSkillMsg{
		Sid: proto.Int32(sID),
		X:   proto.Int32(10),
		Y:   proto.Int32(20),
	}

	for i := 0; i < b.N; i++ {
		NewPacket(uint8(pb.ID_C2S_InputSkill), msg)
	}

}

func Test_Packet(t *testing.T) {
	var sID int32 = 19234333
	msg := &pb.C2S_InputSkillMsg{
		Sid: proto.Int32(sID),
		X:   proto.Int32(10),
		Y:   proto.Int32(20000),
	}

	temp, _ := proto.Marshal(msg)

	p := &Packet{
		id:   uint8(pb.ID_C2S_InputSkill),
		data: temp,
	}

	b := p.Serialize()

	r := strings.NewReader(string(b))

	proto := &MsgProtocol{}

	ret, err := proto.ReadPacket(r)
	if nil != err {
		t.Error(err)
	}

	packet, _ := ret.(*Packet)
	if packet.GetMessageID() != p.id {
		t.Error("packet.GetMessageID() !=  uint8(ID_C2S_InputSkill)")
	}

	if len(packet.data) != len(p.data) {
		t.Error("len(packet.data)!=len(p.data)")
	}

	msg1 := &pb.C2S_InputSkillMsg{}
	err = packet.UnmarshalPB(msg1)
	if nil != err {
		t.Error(err)
	}
	if msg.GetSid() != msg1.GetSid() || msg.GetX() != msg1.GetX() || msg.GetY() != msg1.GetY() {
		t.Error("msg.Sid != data1.Sid || msg.X != data1.X || msg.Y != data1.Y")
	}
}

func Benchmark_Packet(b *testing.B) {
	var sID int32 = 19234333
	msg := &pb.C2S_InputSkillMsg{
		Sid: proto.Int32(sID),
		X:   proto.Int32(10),
		Y:   proto.Int32(20000),
	}

	temp, _ := json.Marshal(msg)

	p := &Packet{
		id:   uint8(pb.ID_C2S_InputSkill),
		data: temp,
	}

	buf := p.Serialize()

	//strings.NewReader(string(b))

	proto := &MsgProtocol{}

	r := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		r.Write(buf)
		if _, err := proto.ReadPacket(r); nil != err {
			b.Error(err)
		}
	}

}
