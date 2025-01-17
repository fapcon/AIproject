// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package bybit

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6ff3ac1dDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesBybit(in *jlexer.Lexer, out *CommonResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "retCode":
			out.Code = int(in.Int())
		case "retMsg":
			out.Msg = string(in.String())
		case "retExtInfo":
			easyjson6ff3ac1dDecode(in, &out.Info)
		case "time":
			out.Time = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesBybit(out *jwriter.Writer, in CommonResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"retCode\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Code))
	}
	{
		const prefix string = ",\"retMsg\":"
		out.RawString(prefix)
		out.String(string(in.Msg))
	}
	{
		const prefix string = ",\"retExtInfo\":"
		out.RawString(prefix)
		easyjson6ff3ac1dEncode(out, in.Info)
	}
	{
		const prefix string = ",\"time\":"
		out.RawString(prefix)
		out.Int64(int64(in.Time))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CommonResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesBybit(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CommonResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesBybit(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CommonResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesBybit(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CommonResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesBybit(l, v)
}
func easyjson6ff3ac1dDecode(in *jlexer.Lexer, out *struct{}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncode(out *jwriter.Writer, in struct{}) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}
