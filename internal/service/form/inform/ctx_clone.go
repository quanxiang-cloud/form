package inform

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
)

type cs string

const (
	_magic   cs = "magic"
	_seq     cs = "seq"
	_version cs = "version"
	_method  cs = "method"

	version = "0.1"
	magic   = " "
)

// HeaderOpt HeaderOpt.
type HeaderOpt func(ctx context.Context)

// WithVersion WithVersion.
func WithVersion(formData *FormData) HeaderOpt {
	return func(ctx context.Context) {
		versions, ok := ctx.Value(_version).(string)
		if !ok || versions == "" {
			versions = version
		}
		formData.Version = versions
	}
}

// WithMagic WithMagic.
func WithMagic(formData *FormData) HeaderOpt {
	return func(ctx context.Context) {
		magics, ok := ctx.Value(_magic).(string)
		if !ok || magics == "" {
			magics = magic
		}

		formData.Magic = magics
	}
}

// WithMethod WithMethod.
func WithMethod(formData *FormData, method string) HeaderOpt {
	return func(ctx context.Context) {
		//methods, ok := ctx.Value(_method).(string)
		//if !ok || methods == "" {
		//	methods = method
		//}
		formData.Method = method
	}
}

// WithSeq WithSeq.
func WithSeq(formData *FormData) HeaderOpt {
	return func(ctx context.Context) {
		sequences, ok := ctx.Value(_seq).(string)
		if !ok || sequences == "" {
			sequences = md5Value(time.Now().Format("2006-01-02 15:04:05"))
		}
		formData.Seq = sequences
	}
}

// CloneHeader CloneHeader.
func CloneHeader(ctx context.Context, opts ...HeaderOpt) {
	for _, opt := range opts {
		opt(ctx)
	}
}

// CTXHeader Context.
func CTXHeader(c context.Context, ctx *gin.Context) context.Context {

	c = context.WithValue(c, _magic, ctx.Request.Header.Get(string(_magic)))
	c = context.WithValue(c, _seq, ctx.Request.Header.Get(string(_seq)))
	c = context.WithValue(c, _version, ctx.Request.Header.Get(string(_version)))
	c = context.WithValue(c, _method, ctx.Request.Header.Get(string(_method)))
	return c
}

// DefaultFormFiled DefaultFormFiled.
func DefaultFormFiled(ctx context.Context, data *FormData, method string) {
	CloneHeader(ctx, WithMethod(data, method), WithVersion(data), WithMagic(data), WithSeq(data))
}

func md5Value(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
