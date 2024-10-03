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

func easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx(in *jlexer.Lexer, out *RequestOrderHistory) {
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
		case "InstrumentType":
			out.InstrumentType = string(in.String())
		case "StartTime":
			out.StartTime = string(in.String())
		case "EndTime":
			out.EndTime = string(in.String())
		case "Limit":
			out.Limit = string(in.String())
		case "After":
			out.After = string(in.String())
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
func easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx(out *jwriter.Writer, in RequestOrderHistory) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"InstrumentType\":"
		out.RawString(prefix[1:])
		out.String(string(in.InstrumentType))
	}
	{
		const prefix string = ",\"StartTime\":"
		out.RawString(prefix)
		out.String(string(in.StartTime))
	}
	{
		const prefix string = ",\"EndTime\":"
		out.RawString(prefix)
		out.String(string(in.EndTime))
	}
	{
		const prefix string = ",\"Limit\":"
		out.RawString(prefix)
		out.String(string(in.Limit))
	}
	{
		const prefix string = ",\"After\":"
		out.RawString(prefix)
		out.String(string(in.After))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestOrderHistory) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestOrderHistory) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestOrderHistory) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestOrderHistory) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx(l, v)
}
func easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx1(in *jlexer.Lexer, out *RequestOrderBook) {
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
		case "Instrument":
			out.Instrument = string(in.String())
		case "Limit":
			out.Limit = int(in.Int())
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
func easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx1(out *jwriter.Writer, in RequestOrderBook) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Instrument\":"
		out.RawString(prefix[1:])
		out.String(string(in.Instrument))
	}
	{
		const prefix string = ",\"Limit\":"
		out.RawString(prefix)
		out.Int(int(in.Limit))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestOrderBook) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestOrderBook) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestOrderBook) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestOrderBook) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx1(l, v)
}
func easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx2(in *jlexer.Lexer, out *RequestOpenOrders) {
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
		case "InstrumentType":
			out.InstrumentType = string(in.String())
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
func easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx2(out *jwriter.Writer, in RequestOpenOrders) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"InstrumentType\":"
		out.RawString(prefix[1:])
		out.String(string(in.InstrumentType))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestOpenOrders) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestOpenOrders) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestOpenOrders) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestOpenOrders) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx2(l, v)
}
func easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx3(in *jlexer.Lexer, out *RequestGetOrder) {
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
		case "Instrument":
			out.Instrument = string(in.String())
		case "OrderID":
			out.OrderID = string(in.String())
		case "ClientOrderID":
			out.ClientOrderID = string(in.String())
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
func easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx3(out *jwriter.Writer, in RequestGetOrder) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Instrument\":"
		out.RawString(prefix[1:])
		out.String(string(in.Instrument))
	}
	{
		const prefix string = ",\"OrderID\":"
		out.RawString(prefix)
		out.String(string(in.OrderID))
	}
	{
		const prefix string = ",\"ClientOrderID\":"
		out.RawString(prefix)
		out.String(string(in.ClientOrderID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestGetOrder) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestGetOrder) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestGetOrder) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestGetOrder) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx3(l, v)
}
func easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx4(in *jlexer.Lexer, out *RequestFillsHistory) {
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
		case "InstrumentType":
			out.InstrumentType = string(in.String())
		case "StartTime":
			out.StartTime = string(in.String())
		case "EndTime":
			out.EndTime = string(in.String())
		case "Limit":
			out.Limit = string(in.String())
		case "After":
			out.After = string(in.String())
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
func easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx4(out *jwriter.Writer, in RequestFillsHistory) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"InstrumentType\":"
		out.RawString(prefix[1:])
		out.String(string(in.InstrumentType))
	}
	{
		const prefix string = ",\"StartTime\":"
		out.RawString(prefix)
		out.String(string(in.StartTime))
	}
	{
		const prefix string = ",\"EndTime\":"
		out.RawString(prefix)
		out.String(string(in.EndTime))
	}
	{
		const prefix string = ",\"Limit\":"
		out.RawString(prefix)
		out.String(string(in.Limit))
	}
	{
		const prefix string = ",\"After\":"
		out.RawString(prefix)
		out.String(string(in.After))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestFillsHistory) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestFillsHistory) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestFillsHistory) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestFillsHistory) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx4(l, v)
}
func easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx5(in *jlexer.Lexer, out *RequestCreateOrder) {
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
		case "instId":
			out.Instrument = string(in.String())
		case "clOrdId":
			out.ClientOrderID = string(in.String())
		case "side":
			out.Side = string(in.String())
		case "ordType":
			out.Type = string(in.String())
		case "sz":
			out.Quantity = string(in.String())
		case "px":
			out.Price = string(in.String())
		case "tdMode":
			out.TdMode = string(in.String())
		case "tgtCcy":
			out.TargetCurrency = string(in.String())
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
func easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx5(out *jwriter.Writer, in RequestCreateOrder) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"instId\":"
		out.RawString(prefix[1:])
		out.String(string(in.Instrument))
	}
	{
		const prefix string = ",\"clOrdId\":"
		out.RawString(prefix)
		out.String(string(in.ClientOrderID))
	}
	{
		const prefix string = ",\"side\":"
		out.RawString(prefix)
		out.String(string(in.Side))
	}
	{
		const prefix string = ",\"ordType\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"sz\":"
		out.RawString(prefix)
		out.String(string(in.Quantity))
	}
	if in.Price != "" {
		const prefix string = ",\"px\":"
		out.RawString(prefix)
		out.String(string(in.Price))
	}
	{
		const prefix string = ",\"tdMode\":"
		out.RawString(prefix)
		out.String(string(in.TdMode))
	}
	if in.TargetCurrency != "" {
		const prefix string = ",\"tgtCcy\":"
		out.RawString(prefix)
		out.String(string(in.TargetCurrency))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestCreateOrder) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestCreateOrder) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestCreateOrder) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestCreateOrder) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx5(l, v)
}
func easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx6(in *jlexer.Lexer, out *RequestCancelOrder) {
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
		case "instId":
			out.Instrument = string(in.String())
		case "ordId":
			out.OrderID = string(in.String())
		case "clOrdId":
			out.ClientOrderID = string(in.String())
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
func easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx6(out *jwriter.Writer, in RequestCancelOrder) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"instId\":"
		out.RawString(prefix[1:])
		out.String(string(in.Instrument))
	}
	if in.OrderID != "" {
		const prefix string = ",\"ordId\":"
		out.RawString(prefix)
		out.String(string(in.OrderID))
	}
	if in.ClientOrderID != "" {
		const prefix string = ",\"clOrdId\":"
		out.RawString(prefix)
		out.String(string(in.ClientOrderID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestCancelOrder) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestCancelOrder) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestCancelOrder) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestCancelOrder) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx6(l, v)
}
func easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx7(in *jlexer.Lexer, out *RequestCancelAllAfter) {
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
		case "timeOut":
			out.Timeout = string(in.String())
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
func easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx7(out *jwriter.Writer, in RequestCancelAllAfter) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"timeOut\":"
		out.RawString(prefix[1:])
		out.String(string(in.Timeout))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestCancelAllAfter) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestCancelAllAfter) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestCancelAllAfter) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestCancelAllAfter) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx7(l, v)
}
func easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx8(in *jlexer.Lexer, out *RequestAmendOrder) {
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
		case "instId":
			out.Instrument = string(in.String())
		case "ordId":
			out.OrderID = string(in.String())
		case "clOrdId":
			out.ClientOrderID = string(in.String())
		case "newPx":
			out.NewPrice = string(in.String())
		case "newSz":
			out.NewSize = string(in.String())
		case "reqId":
			out.RequestID = string(in.String())
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
func easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx8(out *jwriter.Writer, in RequestAmendOrder) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"instId\":"
		out.RawString(prefix[1:])
		out.String(string(in.Instrument))
	}
	if in.OrderID != "" {
		const prefix string = ",\"ordId\":"
		out.RawString(prefix)
		out.String(string(in.OrderID))
	}
	if in.ClientOrderID != "" {
		const prefix string = ",\"clOrdId\":"
		out.RawString(prefix)
		out.String(string(in.ClientOrderID))
	}
	{
		const prefix string = ",\"newPx\":"
		out.RawString(prefix)
		out.String(string(in.NewPrice))
	}
	{
		const prefix string = ",\"newSz\":"
		out.RawString(prefix)
		out.String(string(in.NewSize))
	}
	{
		const prefix string = ",\"reqId\":"
		out.RawString(prefix)
		out.String(string(in.RequestID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RequestAmendOrder) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RequestAmendOrder) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson11d1a9baEncodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RequestAmendOrder) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RequestAmendOrder) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson11d1a9baDecodeStudentgitKataAcademyQuantTorqueInternalAdaptersClientsTypesOkx8(l, v)
}