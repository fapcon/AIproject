// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package okx

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

func easyjsonE0009855DecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx(in *jlexer.Lexer, out *OrderBookResponse) {
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
		case "code":
			out.Code = string(in.String())
		case "msg":
			out.Msg = string(in.String())
		case "data":
			if in.IsNull() {
				in.Skip()
				out.Data = nil
			} else {
				in.Delim('[')
				if out.Data == nil {
					if !in.IsDelim(']') {
						out.Data = make([]struct {
							Asks [][]string `json:"asks"`
							Bids [][]string `json:"bids"`
							TS   string     `json:"ts"`
						}, 0, 1)
					} else {
						out.Data = []struct {
							Asks [][]string `json:"asks"`
							Bids [][]string `json:"bids"`
							TS   string     `json:"ts"`
						}{}
					}
				} else {
					out.Data = (out.Data)[:0]
				}
				for !in.IsDelim(']') {
					var v1 struct {
						Asks [][]string `json:"asks"`
						Bids [][]string `json:"bids"`
						TS   string     `json:"ts"`
					}
					easyjsonE0009855Decode(in, &v1)
					out.Data = append(out.Data, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjsonE0009855EncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx(out *jwriter.Writer, in OrderBookResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"code\":"
		out.RawString(prefix[1:])
		out.String(string(in.Code))
	}
	{
		const prefix string = ",\"msg\":"
		out.RawString(prefix)
		out.String(string(in.Msg))
	}
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		if in.Data == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Data {
				if v2 > 0 {
					out.RawByte(',')
				}
				easyjsonE0009855Encode(out, v3)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v OrderBookResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE0009855EncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v OrderBookResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE0009855EncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *OrderBookResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE0009855DecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *OrderBookResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE0009855DecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx(l, v)
}
func easyjsonE0009855Decode(in *jlexer.Lexer, out *struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
	TS   string     `json:"ts"`
}) {
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
		case "asks":
			if in.IsNull() {
				in.Skip()
				out.Asks = nil
			} else {
				in.Delim('[')
				if out.Asks == nil {
					if !in.IsDelim(']') {
						out.Asks = make([][]string, 0, 2)
					} else {
						out.Asks = [][]string{}
					}
				} else {
					out.Asks = (out.Asks)[:0]
				}
				for !in.IsDelim(']') {
					var v4 []string
					if in.IsNull() {
						in.Skip()
						v4 = nil
					} else {
						in.Delim('[')
						if v4 == nil {
							if !in.IsDelim(']') {
								v4 = make([]string, 0, 4)
							} else {
								v4 = []string{}
							}
						} else {
							v4 = (v4)[:0]
						}
						for !in.IsDelim(']') {
							var v5 string
							v5 = string(in.String())
							v4 = append(v4, v5)
							in.WantComma()
						}
						in.Delim(']')
					}
					out.Asks = append(out.Asks, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "bids":
			if in.IsNull() {
				in.Skip()
				out.Bids = nil
			} else {
				in.Delim('[')
				if out.Bids == nil {
					if !in.IsDelim(']') {
						out.Bids = make([][]string, 0, 2)
					} else {
						out.Bids = [][]string{}
					}
				} else {
					out.Bids = (out.Bids)[:0]
				}
				for !in.IsDelim(']') {
					var v6 []string
					if in.IsNull() {
						in.Skip()
						v6 = nil
					} else {
						in.Delim('[')
						if v6 == nil {
							if !in.IsDelim(']') {
								v6 = make([]string, 0, 4)
							} else {
								v6 = []string{}
							}
						} else {
							v6 = (v6)[:0]
						}
						for !in.IsDelim(']') {
							var v7 string
							v7 = string(in.String())
							v6 = append(v6, v7)
							in.WantComma()
						}
						in.Delim(']')
					}
					out.Bids = append(out.Bids, v6)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "ts":
			out.TS = string(in.String())
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
func easyjsonE0009855Encode(out *jwriter.Writer, in struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
	TS   string     `json:"ts"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"asks\":"
		out.RawString(prefix[1:])
		if in.Asks == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Asks {
				if v8 > 0 {
					out.RawByte(',')
				}
				if v9 == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
					out.RawString("null")
				} else {
					out.RawByte('[')
					for v10, v11 := range v9 {
						if v10 > 0 {
							out.RawByte(',')
						}
						out.String(string(v11))
					}
					out.RawByte(']')
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"bids\":"
		out.RawString(prefix)
		if in.Bids == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v12, v13 := range in.Bids {
				if v12 > 0 {
					out.RawByte(',')
				}
				if v13 == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
					out.RawString("null")
				} else {
					out.RawByte('[')
					for v14, v15 := range v13 {
						if v14 > 0 {
							out.RawByte(',')
						}
						out.String(string(v15))
					}
					out.RawByte(']')
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"ts\":"
		out.RawString(prefix)
		out.String(string(in.TS))
	}
	out.RawByte('}')
}